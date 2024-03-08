package admin_service

import (
	"goskeleton/app/model/model_platform"
)

func CreateCurdFactory() *AdminCurd {
	return &AdminCurd{masterModel: model_platform.CreateAdminFactory()}
}

type AdminCurd struct {
	masterModel *model_platform.AdminModel
}

// GetItemByUsername 通过用户名获取
func (crd *AdminCurd) GetItemByUsername(name string) (model *model_platform.AdminModel) {
	model, _ = crd.masterModel.GetItemByUsername(name)
	return model
}

// GetItemByUsernameAndPw 通过用户名获取
func (crd *AdminCurd) GetItemByUsernameAndPw(name, pw string) (model *model_platform.AdminModel) {
	model, _ = crd.masterModel.GetItemByUsername(name)
	return model
}
