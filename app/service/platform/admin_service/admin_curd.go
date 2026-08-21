package admin_service

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/platform/admin_models"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateCurdFactory() *AdminCurd {
	return &AdminCurd{masterModel: admin_models.CreateAdminFactory()}
}

type AdminCurd struct {
	masterModel *admin_models.AdminModel
}

// GetItemByUsername 通过用户名获取
func (crd *AdminCurd) GetItemByUsername(name string) (model *admin_models.AdminModel) {
	return crd.masterModel.GetItemByUsername(name)
}

// HasByUsernameAndPw 用户名和密码一致
func (crd *AdminCurd) HasByUsernameAndPw(name, pw string) bool {
	if adminModel := crd.masterModel.GetItemByUsernameAndPw(name, pw); adminModel.Id > 0 {
		return true
	}
	return false
}

// IncrPwErrorNum 自增登录错误次数
func (crd *AdminCurd) IncrPwErrorNum(id int64) bool {
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
func (crd *AdminCurd) loginSuccess(id int64) bool {
	updates := map[string]interface{}{
		"pw_error_num":  0,
		"last_login_at": date_helper.CurDatetime(),
	}
	if isOk, _ := crd.masterModel.UpdateById(id, updates); isOk {
		return true
	} else {
		return false
	}
}

// ModifyPwById 修改密码
func (crd *AdminCurd) ModifyPwById(id int64, oldPw, newPw string) (errCode int, err error) {
	adminModel := crd.masterModel.GetItemById(id)
	if adminModel.Id == 0 {
		return -1, errors.New(my_errors.CurdUserNotExist)
	}

	//加密
	AuthFact := CreateAuthFactory()
	oldPw = AuthFact.EncryptPw(oldPw)
	newPw = AuthFact.EncryptPw(newPw)

	//
	if !crd.HasByUsernameAndPw(adminModel.Username, oldPw) {
		return -1, errors.New(my_errors.CurdOldPasswordIncorrect)
	}

	//修改密码
	updates := map[string]interface{}{
		"pw": newPw,
	}
	if isOk, _ := crd.masterModel.UpdateById(id, updates); !isOk {
		return -1, errors.New(my_errors.CurdSaveFailed)
	}
	return errCode, err
}
