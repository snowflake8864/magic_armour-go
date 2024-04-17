package user_service

import (
	"goskeleton/app/model/model_organization"
	"goskeleton/app/utils/openssl_encrypt"
)

func CreateCurdFactory() *UserCurd {
	return &UserCurd{masterModel: model_organization.CreateUserFactory()}
}

type UserCurd struct {
	masterModel *model_organization.UserModel
}

// Register 注册
func (crd *UserCurd) Register(param *model_organization.UserModel) bool {
	param.RoleId = 0
	param.PwErrorNum = 0
	pw, _ := openssl_encrypt.CreateAesFactory().CbcEncryptToHex(param.Password)
	param.Password = pw
	return crd.masterModel.InsertItem(param)
}

// GetItemById 通过用户编号获取
func (crd *UserCurd) GetItemById(id int64) (model *model_organization.UserModel) {
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// GetItemByOrgAndUser 通过企业和用户获取
func (crd *UserCurd) GetItemByOrgAndUser(orgId, userId int64) (model *model_organization.UserModel) {
	model, _ = crd.masterModel.GetItemByOrgAndUser(orgId, userId)
	return model
}

// GetItemByName 通过企业名称获取
func (crd *UserCurd) GetItemByName(name string) (model *model_organization.UserModel) {
	model, _ = crd.masterModel.GetItemByName(name)
	return model
}
