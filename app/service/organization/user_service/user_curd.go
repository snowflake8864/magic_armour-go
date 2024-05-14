package user_service

import (
	"gorm.io/gorm"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/openssl_encrypt"
)

func CreateCurdFactory() *UserCurd {
	return &UserCurd{masterModel: user_models.CreateUserFactory()}
}

type UserCurd struct {
	masterModel *user_models.UserModel
}

// Init 初始化
func (crd *UserCurd) Init(param *user_models.UserModel) {
	param.RoleId = 0
	param.PwErrorNum = 0
	pw, _ := openssl_encrypt.CreateAesFactory().CbcEncryptToHex(param.Password)
	param.Password = pw
}

// Register 注册
func (crd *UserCurd) Register(param *user_models.UserModel) bool {
	crd.Init(param)
	return crd.masterModel.InsertItem(param)
}

// DeleteByOrg 通过企业编号删除
func (crd *UserCurd) DeleteByOrg(orgId int64) bool {
	conditions := user_models.SearchConditions{
		OrgId: orgId,
	}
	return crd.masterModel.DeleteByConditions(conditions)
}

// GetItemById 通过用户编号获取
func (crd *UserCurd) GetItemById(id int64) (model *user_models.UserModel) {
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// GetItemByOrgAndUser 通过企业和用户获取
func (crd *UserCurd) GetItemByOrgAndUser(orgId, userId int64) (model *user_models.UserModel) {
	model, _ = crd.masterModel.GetItemByOrgAndUser(orgId, userId)
	return model
}

// GetItemByName 通过企业名称获取
func (crd *UserCurd) GetItemByName(name string) (model *user_models.UserModel) {
	model, _ = crd.masterModel.GetItemByName(name)
	return model
}

// HasByUsernameAndPw 用户名和密码一致
func (crd *UserCurd) HasByUsernameAndPw(name, pw string) bool {
	if adminModel := crd.masterModel.GetItemByNameAndPw(name, pw); adminModel.Id > 0 {
		return true
	}
	return false
}

// IncrPwErrorNum 自增登录错误次数
func (crd *UserCurd) IncrPwErrorNum(id int64) bool {
	updates := map[string]interface{}{
		"pw_error_num": gorm.Expr("pw_error_num+1"),
	}
	if isOk, _ := crd.masterModel.UpdateById(id, updates); isOk {
		return true
	} else {
		return false
	}
}

// loginSuccess 登录成功
func (crd *UserCurd) loginSuccess(id int64) bool {
	updates := map[string]interface{}{
		"pw_error_num":  0,
		"last_login_at": helpers.DatetimeCurr(),
	}
	if isOk, _ := crd.masterModel.UpdateById(id, updates); isOk {
		return true
	} else {
		return false
	}
}
