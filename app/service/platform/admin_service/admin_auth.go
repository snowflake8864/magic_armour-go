package admin_service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/platform/admin_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/ali_sms"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/helpers/timeout_control_helper"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

func CreateAuthFactory() *AdminAuth {
	return &AdminAuth{
		redisClient:   redis_factory.GetOneRedisClient(),
		tokenPrefix:   "admin_auth_token_",
		tokenExpire:   cur_app.GetPlatformTokenExpire(),
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
	LoginIp  string `json:"login_ip"`
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
	code = int_helper.IntRandom2Str(100000, 999999)
	crd.redisClient.SetEx(crd.keyBySms(phone), code, crd.sliderExpire)

	// 调用客户端
	timeout_control_helper.TimeoutControl(func(errChan chan error) {
		// 赋值
		err := ali_sms.SendAuthCode(phone, code)

		errChan <- err
	}, fmt.Sprintf("pt action=%s,phone=%s,code=%s", "auth", phone, code))
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
func (crd *AdminAuth) IsLocking(username string) error {
	isOK, _ := crd.redisClient.Exists(crd.keyByLocking(username))
	if !isOK {
		return nil
	} else {
		return errors.New("用户已经被锁定，请稍后在试")
	}
}

// setLocking  设置锁定中
func (crd *AdminAuth) setLocking(username string, second int) (isOk bool) {
	isOk, _ = crd.redisClient.SetEx(crd.keyByLocking(username), "veda", second)
	return isOk
}

// HasUser 存在用户
func (crd *AdminAuth) HasUser(username string) (err error, model *admin_models.AdminModel) {
	model = CreateCurdFactory().GetItemByUsername(username)
	if model == nil || model.Id == 0 {
		return errors.New(my_errors.CurdUsernameOrPwFail), nil
	} else {
		return nil, model
	}
}

// EncryptPw 加密密码
func (crd *AdminAuth) EncryptPw(pw string) string {
	pw, _ = openssl_encrypt.CreateAesFactory().CbcEncryptToHex(pw)
	return pw
}

// LoginByPw 密码登录
func (crd *AdminAuth) LoginByPw(username, pw, slider, sms string, netIp, lang string) (result map[string]interface{}, errCode int, err error) {
	if !crd.checkBySlider(username, slider) {
		return nil, -88, errors.New(my_errors.CurdSliderIncorrect)
	}
	if err = crd.IsLocking(username); err != nil {
		return nil, -1, err
	}

	err, adminModel := crd.HasUser(username)
	if err != nil {
		return nil, -1, err
	}

	//短信验证
	if !string_helper.IsEmpty(adminModel.Phone) && adminModel.LoginPhoneCheck == enums.IsYes {
		if string_helper.IsEmpty(sms) {
			return nil, -99, errors.New(my_errors.CurdTwoFactorAutEnabled)
		}
		if !crd.checkBySms(adminModel.Phone, sms) {
			return nil, -99, errors.New(my_errors.CurdSmsVerifyCodeNotMatch)
		}
		crd.deleteBySms(adminModel.Phone)
	}

	//用户名和密码不一致
	curdFact := CreateCurdFactory()
	if !curdFact.HasByUsernameAndPw(username, crd.EncryptPw(pw)) {
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
		diff := string_helper.ToInt(mapDict[enum_dictionary.PlatformLoginErrorMax]) - adminModel.PwErrorNum - 1
		if diff <= 0 {
			//设置锁定
			crd.setLocking(username, string_helper.ToInt(mapDict[enum_dictionary.PlatformLoginLockDuration])*60)

			return nil, -2,
				errors.New(my_errors.GetMessagef(lang, my_errors.CurdUsernameOrPwFailLocking, mapDict[enum_dictionary.PlatformLoginLockDuration]))

		} else if diff < string_helper.ToInt(mapDict[enum_dictionary.PlatformLoginErrorRemind]) {

			return nil, -2,
				errors.New(my_errors.GetMessagef(lang, my_errors.CurdUsernameOrPwFailWillLock, diff, mapDict[enum_dictionary.PlatformLoginLockDuration]))

		} else {
			return nil, -2, errors.New(my_errors.CurdUsernameOrPwFail)
		}
	}

	curdFact.loginSuccess(adminModel.Id)

	result = crd.loginAfter(adminModel, netIp)

	return result, 0, nil
}

// loginAfter
func (crd *AdminAuth) loginAfter(adminModel *admin_models.AdminModel, netIp string) map[string]interface{} {
	loginInfo := AdminCacheLogin{
		AdminId:  adminModel.Id,
		Username: adminModel.Username,
		IsRoot:   adminModel.IsRoot,
		LoginAt:  date_helper.CurDatetime(),
		LoginIp:  netIp,
	}

	//缓存
	jsonStr := json_helper.Encode(loginInfo)
	token := md5_encrypt.MD5(jsonStr)
	crd.redisClient.SetEx(crd.keyByToken(token), jsonStr, crd.tokenExpire)

	//
	result := make(map[string]interface{})
	result["userid"] = adminModel.Id
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
	if string_helper.IsEmpty(jsonStr) {
		return AdminCacheLogin{}, errors.New("token expired")
	}
	err = json_helper.ToStruct(jsonStr, &loginInfo)
	if err != nil {
		return AdminCacheLogin{}, errors.New("token json fail")
	}

	if isDelay {
		_, _ = crd.redisClient.Expire(crd.keyByToken(token), crd.tokenExpire)
	}
	return loginInfo, nil
}

// Logout 退出
func (crd *AdminAuth) Logout(token string) bool {
	token, err := openssl_encrypt.CreateAesFactory().CbcDecryptByBase64(token)
	if err != nil {
		return true
	}

	_, _ = crd.redisClient.Del(crd.keyByToken(token))

	return true
}
