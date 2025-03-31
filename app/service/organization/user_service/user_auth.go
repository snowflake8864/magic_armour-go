package user_service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_org_operate"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_user"
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
	"strings"
)

const (
	ActonLogin  = "1"
	ActonReg    = "2"
	ActonFindPw = "3"
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
	OrgName       string  `json:"org_name"`
	OrgCode       string  `json:"org_code"`
	UserId        int64   `json:"user_id"`
	Username      string  `json:"username"`
	RoleId        int64   `json:"role_id"`
	IsRoot        int     `json:"is_root"`
	DeptIds       []int64 `json:"dept_ids"`        //直接管理的组织架构
	DeptWithChild []int64 `json:"dept_with_child"` //直接管理的组织架构以及子集
	LoginAt       string  `json:"login_at"`
	LoginIp       string  `json:"login_ip"`
	Source        int     `json:"source"`
}

type UserBackLogin struct {
	OrgId    int64  `json:"org_id"`
	OrgName  string `json:"org_name"`
	OrgCode  string `json:"org_code"`
	OrgType  int    `json:"org_type"`
	UserId   int64  `json:"userid"`
	Username string `json:"username"`
	IsRoot   int    `json:"is_root"`
	Token    string `json:"token"`
	Position string `json:"position"`
	Source   int    `json:"source"`
}

func (crd *UserAuth) keyByToken(val string) string {
	return crd.tokenPrefix + val
}

func (crd *UserAuth) keyByLocking(val string) string {
	return crd.lockingPrefix + val
}

func (crd *UserAuth) keyBySlider(action, val string) string {
	return crd.sliderPrefix + action + "_" + val
}

func (crd *UserAuth) keyBySms(action, val string) string {
	return crd.smsPrefix + action + "_" + val
}

// GenerateBySlider  生成滑块码
func (crd *UserAuth) GenerateBySlider(action, username string) (code string) {
	code = uuid.New().String()
	code = strings.Replace(code, "-", "", -1)
	crd.redisClient.SetEx(crd.keyBySlider(action, username), code, crd.sliderExpire)

	return
}

// CheckBySlider  验证滑块码
func (crd *UserAuth) CheckBySlider(action, username, code string) bool {
	cacheCode := crd.redisClient.Get(crd.keyBySlider(action, username))
	if cacheCode == code {
		return true
	}
	return false
}

// GenerateBySms  生成短信码
func (crd *UserAuth) GenerateBySms(action, phone string) (code string) {
	code = int_helper.IntRandom2Str(100000, 999999)
	crd.redisClient.SetEx(crd.keyBySms(action, phone), code, crd.sliderExpire)

	go ali_sms.SendAuthCode(phone, code)
	return code
}

// CheckBySms  验证短信码
func (crd *UserAuth) CheckBySms(action, phone, code string) bool {
	cacheCode := crd.redisClient.Get(crd.keyBySms(action, phone))
	if cacheCode == code {
		return true
	}
	return false
}

// DeleteBySms  删除短信码
func (crd *UserAuth) DeleteBySms(action, phone string) {
	crd.redisClient.Del(crd.keyBySms(action, phone))
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
	if model = CreateCurdFactory().GetItemByName(username); model != nil && model.Id > 0 {
		return nil, model
	} else if model = CreateCurdFactory().GetItemByPhone(username); model != nil && model.Id > 0 {
		return nil, model
	} else {
		return errors.New("用户名或密码错误"), nil
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

// EncryptPw 加密密码
func (crd *UserAuth) EncryptPw(pw string) string {
	pw, _ = openssl_encrypt.CreateAesFactory().CbcEncryptToHex(pw)
	return pw
}

// LoginByPw 密码登录
func (crd *UserAuth) LoginByPw(username, pw, slider, sms string, netIp string) (backInfo *UserBackLogin, errCode int, err error) {
	if !crd.CheckBySlider(ActonLogin, username, slider) {
		return nil, -88, errors.New("滑块不正确")
	}
	if err = crd.IsLocking(username); err != nil {
		return nil, -1, err
	}

	//获取用户
	err, userModel := crd.HasUser(username)
	if err != nil {
		return nil, -1, err
	}
	if userModel.Status == enum_user.StatusLocking {
		return nil, -99, errors.New("用户已锁定，请联系管理员")
	}
	if userModel.Status == enum_user.StatusVerifying {
		return nil, -99, errors.New("用户审核中，请联系管理员")
	}

	//短信验证
	if err1, smsSwitch := crd.IsPhoneCheck(userModel.Phone, userModel.OrgId); err1 != nil {
		return nil, -99, err1
	} else if smsSwitch {
		if string_helper.IsEmpty(sms) {
			return nil, -99, errors.New("已启用双因子认证时，短信验证码不能为空")
		}
		if !crd.CheckBySms(ActonLogin, userModel.Phone, sms) {
			return nil, -99, errors.New("短信验证码不一致")
		}
		crd.DeleteBySms(ActonLogin, userModel.Phone)
	}

	//企业验证
	orgModel, err := organization_models.CreateOrganizationFactory().GetItemById(userModel.OrgId)
	if err != nil || orgModel.Id == 0 {
		return nil, -1, errors.New("企业不存在")
	}
	if orgModel.Status == enum_organization.StatusVerifying {
		return nil, -1, errors.New("账户审核中，请联系管理员")
	}
	if orgModel.Status == enum_organization.StatusLocking {
		return nil, -1, errors.New("账户已锁定，请联系管理员")
	}

	//用户名和密码不一致
	curdFact := CreateCurdFactory()
	if !curdFact.HasByUsernameAndPw(username, crd.EncryptPw(pw)) {
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
			return nil, -2,
				errors.New(fmt.Sprintf("用户名或密码错误, 已被锁定%s分钟，请稍后在试", mapDict[enum_dictionary.OrganizationLoginLockDuration]))

		} else if diff < string_helper.ToInt(mapDict[enum_dictionary.OrganizationLoginErrorRemind]) {
			return nil, -2,
				errors.New(fmt.Sprintf("用户名或密码错误, 还有%d次即将被锁定%s分钟", diff, mapDict[enum_dictionary.OrganizationLoginLockDuration]))

		} else {
			return nil, -2, errors.New("用户名或密码错误")
		}
	}

	curdFact.loginSuccess(userModel.Id)

	backInfo = crd.loginAfter(userModel, orgModel, netIp, enum_org_operate.SourceWeb)

	return backInfo, 0, nil
}

// LoginByHs 幻势登录
func (crd *UserAuth) LoginByHs(orgCode string, netIp string) (backInfo *UserBackLogin, errCode int, err error) {
	//企业验证
	orgModel, err := organization_models.CreateOrganizationFactory().GetItemByCode(orgCode)
	if err != nil || orgModel.Id == 0 {
		return nil, -1, errors.New("企业不存在")
	}
	if orgModel.Status == enum_organization.StatusVerifying {
		return nil, -1, errors.New("账户审核中，请联系管理员")
	}
	if orgModel.Status == enum_organization.StatusLocking {
		return nil, -1, errors.New("账户已锁定，请联系管理员")
	}

	//获取用户
	err, userModel := crd.HasUser(orgModel.UserName)
	if err != nil {
		return nil, -1, err
	}
	if userModel.Status == enum_user.StatusLocking {
		return nil, -99, errors.New("用户已锁定，请联系管理员")
	}
	if userModel.Status == enum_user.StatusVerifying {
		return nil, -99, errors.New("用户审核中，请联系管理员")
	}

	CreateCurdFactory().loginSuccess(userModel.Id)

	backInfo = crd.loginAfter(userModel, orgModel, netIp, enum_org_operate.SourceHs)

	return backInfo, 0, nil

}

// loginAfter
func (crd *UserAuth) loginAfter(userModel *user_models.UserModel, orgModel *organization_models.OrganizationModel, netIp string, source int) *UserBackLogin {
	var deptIds, deptWithChild []int64
	//不是管理员
	if userModel.IsRoot == enums.IsNot {
		deptIds = json_helper.ToSliceInt64(userModel.Depts)
		deptWithChild = department_service.CreateCurdFactory().GetIdsSelfWithChildByIds(userModel.OrgId, deptIds)
	}

	loginInfo := UserCacheLogin{
		OrgId:         userModel.OrgId,
		OrgName:       orgModel.Name,
		OrgCode:       orgModel.Code,
		UserId:        userModel.Id,
		Username:      userModel.Username,
		RoleId:        userModel.RoleId,
		IsRoot:        userModel.IsRoot,
		DeptIds:       deptIds,
		DeptWithChild: deptWithChild,
		LoginAt:       date_helper.CurDatetime(),
		LoginIp:       netIp,
		Source:        source,
	}

	//缓存
	jsonStr := json_helper.Encode(loginInfo)
	token := md5_encrypt.MD5(jsonStr)
	_, _ = crd.redisClient.SetEx(crd.keyByToken(token), jsonStr, crd.tokenExpire)

	//记录日志

	//返回给前端数据
	encryptToken, _ := openssl_encrypt.CreateAesFactory().CbcEncryptToBase64(token)
	backInfo := UserBackLogin{
		OrgId:    orgModel.Id,
		OrgName:  orgModel.Name,
		OrgCode:  orgModel.Code,
		OrgType:  orgModel.Type,
		UserId:   userModel.Id,
		Username: userModel.Username,
		IsRoot:   userModel.IsRoot,
		Token:    encryptToken,
		Position: userModel.Position,
		Source:   source,
	}

	return &backInfo
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

// Logout 退出
func (crd *UserAuth) Logout(token string) bool {
	token, err := openssl_encrypt.CreateAesFactory().CbcDecryptByBase64(token)
	if err != nil {
		return true
	}

	_, _ = crd.redisClient.Del(crd.keyByToken(token))

	return true
}
