package user_service

import (
	"errors"
	"github.com/google/uuid"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/ali_sms"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strconv"
	"strings"
)

func CreateAuthFactory() *UserAuth {
	return &UserAuth{
		redisClient:   redis_factory.GetOneRedisClient(),
		tokenPrefix:   "user_auth_token_",
		tokenExpire:   cur_app.GetOrganizationTokenExpire(),
		lockingPrefix: "user_auth_locking_",
		sliderPrefix:  "user_auth_slider_",
		sliderExpire:  180,
		smsPrefix:     "user_auth_sms_",
		smsExpire:     300,
	}
}

type UserAuth struct {
	redisClient   *redis_factory.RedisClient
	tokenPrefix   string
	tokenExpire   int
	lockingPrefix string
	sliderPrefix  string
	sliderExpire  int
	smsPrefix     string
	smsExpire     int
}

type UserCacheLogin struct {
	OrgId         int64   `json:"org_id"`
	UserId        int64   `json:"user_id"`
	Username      string  `json:"username"`
	IsRoot        int     `json:"is_root"`
	DeptIds       []int64 `json:"dept_ids"`        //直接管理的组织架构
	DeptWithChild []int64 `json:"dept_with_child"` //直接管理的组织架构以及子集
	LoginAt       string  `json:"login_at"`
}

func (crd *UserAuth) keyByToken(val string) string {
	return crd.tokenPrefix + val
}

func (crd *UserAuth) keyByLocking(val string) string {
	return crd.lockingPrefix + val
}

func (crd *UserAuth) keyBySlider(val string) string {
	return crd.sliderPrefix + val
}

func (crd *UserAuth) keyBySms(val string) string {
	return crd.smsPrefix + val
}

// GenerateBySlider  生成滑块码
func (crd *UserAuth) GenerateBySlider(username string) (code string) {
	code = uuid.New().String()
	code = strings.Replace(code, "-", "", -1)
	crd.redisClient.SetEx(crd.keyBySlider(username), code, crd.sliderExpire)

	return
}

// checkBySlider  验证滑块码
func (crd *UserAuth) checkBySlider(username, code string) bool {
	cacheCode := crd.redisClient.Get(crd.keyBySlider(username))
	if cacheCode == code {
		return true
	}
	return false
}

// GenerateBySms  生成短信码
func (crd *UserAuth) GenerateBySms(phone string) (code string) {
	code = int_helper.IntRandom2Str(100000, 999999)
	crd.redisClient.SetEx(crd.keyBySms(phone), code, crd.sliderExpire)

	go ali_sms.SendAuthCode(phone, code)
	return code
}

// checkBySms  验证短信码
func (crd *UserAuth) checkBySms(phone, code string) bool {
	cacheCode := crd.redisClient.Get(crd.keyBySms(phone))
	if cacheCode == code {
		return true
	}
	return false
}

// deleteBySms  删除短信码
func (crd *UserAuth) deleteBySms(phone string) {
	crd.redisClient.Del(crd.keyBySms(phone))
}

// IsLocking  是否锁定中
func (crd *UserAuth) IsLocking(username string) error {
	isOK, _ := crd.redisClient.Exists(crd.keyByLocking(username))
	if !isOK {
		return nil
	} else {
		return errors.New("用户已经被锁定，请稍后在试")
	}
}

// setLocking  设置锁定中
func (crd *UserAuth) setLocking(username string, second int) (isOk bool) {
	isOk, _ = crd.redisClient.SetEx(crd.keyByLocking(username), "veda", second)
	return isOk
}

// HasUser 存在用户
func (crd *UserAuth) HasUser(username string) (err error, model *user_models.UserModel) {
	model = CreateCurdFactory().GetItemByName(username)
	if model == nil || model.Id == 0 {
		return errors.New("用户名或密码错误"), nil
	} else {
		return nil, model
	}
}

// IsPhoneCheck 是否需要手机校验
func (crd *UserAuth) IsPhoneCheck(userPhone string, orgId int64) (err error, isPhone bool) {
	if string_helper.IsEmpty(userPhone) {
		return nil, false
	}

	//数据字典
	if smsSwitch := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.OrganizationLoginSmsSwitch); string_helper.ToInt(smsSwitch) == enums.SwitchOn {
		return nil, true
	}

	//企业开关
	orgModel, _ := organization_models.CreateOrganizationFactory().GetItemById(orgId)
	if orgModel == nil || orgModel.Id == 0 {
		return errors.New("用户名或密码错误"), false
	} else if orgModel.LoginPhoneCheck == enums.SwitchOn {
		return nil, true
	} else {
		return nil, false
	}
}

// LoginByPw 密码登录
func (crd *UserAuth) LoginByPw(username, pw, slider, sms string) (result map[string]interface{}, errCode int, err error) {
	if !crd.checkBySlider(username, slider) {
		return nil, -88, errors.New("滑块不正确")
	}
	if err = crd.IsLocking(username); err != nil {
		return nil, -1, err
	}

	err, userModel := crd.HasUser(username)
	if err != nil {
		return nil, -1, err
	}

	//短信验证
	if err, smsSwitch := crd.IsPhoneCheck(userModel.Phone, userModel.OrgId); err != nil {
		return nil, -99, err
	} else if smsSwitch {
		if string_helper.IsEmpty(sms) {
			return nil, -99, errors.New("已启用双因子认证时，短信验证码不能为空")
		}
		if !crd.checkBySms(userModel.Phone, sms) {
			return nil, -99, errors.New("短信验证码不一致")
		}
		crd.deleteBySms(userModel.Phone)
	}

	//用户名和密码不一致
	pw, _ = openssl_encrypt.CreateAesFactory().CbcEncryptToHex(pw)
	curdFact := CreateCurdFactory()
	if !curdFact.HasByUsernameAndPw(username, pw) {
		//错误次数加1
		curdFact.IncrPwErrorNum(userModel.Id)

		//数据字典
		dictIds := []string{
			enum_dictionary.OrganizationLoginErrorMax,
			enum_dictionary.OrganizationLoginLockDuration,
			enum_dictionary.OrganizationLoginErrorRemind,
		}
		mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
		//达到登录错误最大次数
		diff := string_helper.ToInt(mapDict[enum_dictionary.OrganizationLoginErrorMax]) - userModel.PwErrorNum - 1
		if diff <= 0 {
			//设置锁定
			crd.setLocking(username, string_helper.ToInt(mapDict[enum_dictionary.OrganizationLoginLockDuration])*60)
			return result, -2, errors.New("用户名或密码错误, 已被锁定" +
				mapDict[enum_dictionary.PlatformLoginLockDuration] + "分钟，请稍后在试")
		} else if diff < string_helper.ToInt(mapDict[enum_dictionary.OrganizationLoginErrorRemind]) {
			return result, -2, errors.New("用户名或密码错误, 还有" +
				strconv.Itoa(diff) + "次即将被锁定" + mapDict[enum_dictionary.OrganizationLoginErrorRemind] + "分钟")
		} else {
			return nil, -2, errors.New("用户名或密码错误")
		}
	}

	curdFact.loginSuccess(userModel.Id)

	result = crd.loginAfter(userModel)

	return result, 0, nil
}

// loginAfter
func (crd *UserAuth) loginAfter(userModel *user_models.UserModel) map[string]interface{} {
	var deptIds, deptWithChild []int64
	//不是管理员
	if userModel.IsRoot == enums.IsNot {
		deptIds = json_helper.ToSliceInt64(userModel.Depts)
		deptWithChild = department_service.CreateCurdFactory().GetIdsSelfWithChildByIds(userModel.OrgId, deptIds)
	}

	loginInfo := UserCacheLogin{
		OrgId:         userModel.OrgId,
		UserId:        userModel.Id,
		Username:      userModel.Username,
		IsRoot:        userModel.IsRoot,
		DeptIds:       deptIds,
		DeptWithChild: deptWithChild,
		LoginAt:       date_helper.CurDatetime(),
	}

	//缓存
	jsonStr := json_helper.Encode(loginInfo)
	token := md5_encrypt.MD5(jsonStr)
	crd.redisClient.SetEx(crd.keyByToken(token), jsonStr, crd.tokenExpire)

	//
	result := make(map[string]interface{})
	result["username"] = userModel.Username
	result["is_root"] = userModel.IsRoot
	encryptToken, _ := openssl_encrypt.CreateAesFactory().CbcEncryptToBase64(token)
	result["token"] = encryptToken
	return result
}

// DecryptByLogin 登录信息解密
func (crd *UserAuth) DecryptByLogin(token string, isDelay bool) (loginInfo UserCacheLogin, err error) {
	token, err = openssl_encrypt.CreateAesFactory().CbcDecryptByBase64(token)
	if err != nil {
		return UserCacheLogin{}, errors.New("token fail")
	}

	jsonStr := crd.redisClient.Get(crd.keyByToken(token))
	if string_helper.IsEmpty(jsonStr) {
		return UserCacheLogin{}, errors.New("token expired")
	}
	err = json_helper.ToStruct(jsonStr, &loginInfo)
	if err != nil {
		return UserCacheLogin{}, errors.New("token json fail")
	}

	if isDelay {
		_, _ = crd.redisClient.Expire(crd.keyByToken(token), crd.tokenExpire)
	}
	return loginInfo, nil
}
