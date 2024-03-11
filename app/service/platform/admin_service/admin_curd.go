package admin_service

import (
	"gorm.io/gorm"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
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

// HasByUsernameAndPw 用户名和密码一致
func (crd *AdminCurd) HasByUsernameAndPw(name, pw string) bool {
	adminModel, _ := crd.masterModel.GetItemByUsernameAndPw(name, pw)
	if adminModel.Id > 0 {
		return false
	}
	return false
}

// IncrPwErrorNum 自增登录错误次数
func (crd *AdminCurd) IncrPwErrorNum(id int64) bool {
	updates := map[string]interface{}{
		"pw_error_num": gorm.Expr("pw_error_num+1"),
	}
	return crd.masterModel.UpdateById(id, updates)
}

// loginSuccess 登录成功
func (crd *AdminCurd) loginSuccess(id int64) bool {
	updates := map[string]interface{}{
		"pw_error_num":  0,
		"last_login_at": helpers.GetDatetime(),
	}
	return crd.masterModel.UpdateById(id, updates)
}
