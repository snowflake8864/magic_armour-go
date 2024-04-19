package admin_service

import (
	"errors"
	"github.com/google/uuid"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/ali_sms"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strconv"
	"strings"
)

func CreateAuthFactory() *AdminAuth {
	return &AdminAuth{
		redisClient:   redis_factory.GetOneRedisClient(),
		tokenPrefix:   "admin_auth_token_",
		tokenExpire:   1800,
		lockingPrefix: "admin_auth_locking_",
		sliderPrefix:  "admin_auth_slider_",
		sliderExpire:  180,
		smsPrefix:     "admin_auth_sms_",
		smsExpire:     300,
	}
}

type AdminAuth struct {
	redisClient   *redis_factory.RedisClient
	tokenPrefix   string
	tokenExpire   int
	lockingPrefix string
	sliderPrefix  string
	sliderExpire  int
	smsPrefix     string
	smsExpire     int
}

type AdminCacheLogin struct {
	AdminId  int64  `json:"admin_id"`
	Username string `json:"username"`
	IsRoot   int    `json:"is_root"`
	LoginAt  string `json:"login_at"`
}

func (crd *AdminAuth) keyByToken(val string) string {
	return crd.tokenPrefix + val
}

func (crd *AdminAuth) keyByLocking(val string) string {
	return crd.lockingPrefix + val
}

func (crd *AdminAuth) keyBySlider(val string) string {
	return crd.sliderPrefix + val
}

func (crd *AdminAuth) keyBySms(val string) string {
	return crd.smsPrefix + val
}

// GenerateBySlider  生成滑块码
func (crd *AdminAuth) GenerateBySlider(username string) (code string) {
	code = uuid.New().String()
	code = strings.Replace(code, "-", "", -1)
	crd.redisClient.SetEx(crd.keyBySlider(username), code, crd.sliderExpire)

	return
}

// checkBySlider  验证滑块码
func (crd *AdminAuth) checkBySlider(username, code string) bool {
	cacheCode := crd.redisClient.Get(crd.keyBySlider(username))
	if cacheCode == code {
		return true
	}
	return false
}

// GenerateBySms  生成短信码
func (crd *AdminAuth) GenerateBySms(phone string) (code string) {
	code = helpers.IntRandom2Str(100000, 999999)
	crd.redisClient.SetEx(crd.keyBySms(phone), code, crd.sliderExpire)

	go ali_sms.SendAuthCode(phone, code)
	return code
}

// checkBySms  验证短信码
func (crd *AdminAuth) checkBySms(phone, code string) bool {
	cacheCode := crd.redisClient.Get(crd.keyBySms(phone))
	if cacheCode == code {
		return true
	}
	return false
}

// deleteBySms  删除短信码
func (crd *AdminAuth) deleteBySms(phone string) {
	crd.redisClient.Del(crd.keyBySms(phone))
}

// IsLocking  是否锁定中
func (crd *AdminAuth) IsLocking(username string) (isOk bool) {
	isOk, _ = crd.redisClient.Exists(crd.keyByLocking(username))
	return isOk
}

// setLocking  设置锁定中
func (crd *AdminAuth) setLocking(username string, second int) (isOk bool) {
	isOk, _ = crd.redisClient.SetEx(crd.keyByLocking(username), "veda", second)
	return isOk
}

// LoginByPw 密码登录
func (crd *AdminAuth) LoginByPw(username, pw, slider, sms string) (result map[string]interface{}, errCode int, err error) {
	if !crd.checkBySlider(username, slider) {
		return nil, -88, errors.New("滑块不正确")
	}
	if crd.IsLocking(username) {
		return nil, -1, errors.New("用户已经被锁定，请稍后在试")
	}

	curdFact := CreateCurdFactory()
	adminModel := curdFact.GetItemByUsername(username)
	if adminModel.Id == 0 {
		return nil, -1, errors.New("用户名或密码错误")
	}

	//短信验证
	if !helpers.StringIsEmpty(adminModel.Phone) && adminModel.LoginPhoneCheck == enums.IsYes {
		if helpers.StringIsEmpty(sms) {
			return nil, -99, errors.New("已启用双因子认证时，短信验证码不能为空")
		}
		if !crd.checkBySms(adminModel.Phone, sms) {
			return nil, -99, errors.New("短信验证码不一致")
		}
		crd.deleteBySms(adminModel.Phone)
	}

	//用户名和密码不一致
	pw, _ = openssl_encrypt.CreateAesFactory().CbcEncryptToHex(pw)
	if !curdFact.HasByUsernameAndPw(username, pw) {
		//错误次数加1
		curdFact.IncrPwErrorNum(adminModel.Id)

		//数据字典
		dictIds := []string{
			enum_dictionary.PlatformLoginErrorMax,
			enum_dictionary.PlatformLoginLockDuration,
			enum_dictionary.PlatformLoginErrorRemind,
		}
		mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
		//达到登录错误最大次数
		diff := helpers.String2Int(mapDict[enum_dictionary.PlatformLoginErrorMax]) - adminModel.PwErrorNum - 1
		if diff <= 0 {
			//设置锁定
			crd.setLocking(username, helpers.String2Int(mapDict[enum_dictionary.PlatformLoginLockDuration])*60)
			return result, -2, errors.New("用户名或密码错误, 已被锁定" +
				mapDict[enum_dictionary.PlatformLoginLockDuration] + "分钟，请稍后在试")
		} else if diff < helpers.String2Int(mapDict[enum_dictionary.PlatformLoginErrorRemind]) {
			return result, -2, errors.New("用户名或密码错误, 还有" +
				strconv.Itoa(diff) + "次即将被锁定" + mapDict[enum_dictionary.PlatformLoginLockDuration] + "分钟")
		} else {
			return nil, -2, errors.New("用户名或密码错误")
		}
	}

	curdFact.loginSuccess(adminModel.Id)

	result = crd.loginAfter(adminModel)

	return result, 0, nil
}

// loginAfter
func (crd *AdminAuth) loginAfter(adminModel *model_platform.AdminModel) map[string]interface{} {
	loginInfo := AdminCacheLogin{
		AdminId:  adminModel.Id,
		Username: adminModel.Username,
		IsRoot:   adminModel.IsRoot,
		LoginAt:  helpers.DatetimeCurr(),
	}

	//缓存
	jsonStr := helpers.JsonEncode(loginInfo)
	token := md5_encrypt.MD5(jsonStr)
	crd.redisClient.SetEx(crd.keyByToken(token), jsonStr, crd.tokenExpire)

	//
	result := make(map[string]interface{})
	result["username"] = adminModel.Username
	result["is_root"] = adminModel.IsRoot
	encryptToken, _ := openssl_encrypt.CreateAesFactory().CbcEncryptToBase64(token)
	result["token"] = encryptToken
	return result
}

// DecryptByLogin 登录信息解密
func (crd *AdminAuth) DecryptByLogin(token string, isDelay bool) (loginInfo AdminCacheLogin, err error) {
	token, err = openssl_encrypt.CreateAesFactory().CbcDecryptByBase64(token)
	if err != nil {
		return AdminCacheLogin{}, errors.New("token fail")
	}

	jsonStr := crd.redisClient.Get(crd.keyByToken(token))
	if helpers.StringIsEmpty(jsonStr) {
		return AdminCacheLogin{}, errors.New("token expired")
	}
	err = helpers.Json2Struct(jsonStr, &loginInfo)
	if err != nil {
		return AdminCacheLogin{}, errors.New("token json fail")
	}

	if isDelay {
		_, _ = crd.redisClient.Expire(crd.keyByToken(token), crd.tokenExpire)
	}
	return loginInfo, nil
}
