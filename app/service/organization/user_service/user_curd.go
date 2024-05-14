package user_service

import (
	"goskeleton/app/model/organization/user_models"
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
	conditions := map[string]interface{}{
		"org_id": orgId,
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

// PageByOrg 通过企业分页
func (crd *UserCurd) PageByOrg(orgId int64, offset, pageSize int, orderField, orderSort string) (total int64, items []*user_models.UserModel) {
	wheres := make(map[string]interface{})
	wheres["org_id"] = orgId
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(wheres, offset, pageSize, order)
	return total, items
}
