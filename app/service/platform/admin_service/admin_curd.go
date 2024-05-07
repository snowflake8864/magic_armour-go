package admin_service

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/openssl_encrypt"
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
		return true
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
		"last_login_at": helpers.DatetimeCurr(),
	}
	return crd.masterModel.UpdateById(id, updates)
}

// ModifyPwById 修改密码
func (crd *AdminCurd) ModifyPwById(id int64, oldPw, newPw string) (errCode int, err error) {
	adminModel, err := crd.masterModel.GetItemById(id)
	if err != nil || adminModel.Id == 0 {
		return -1, errors.New("用户不存在")
	}

	//加密
	AesFact := openssl_encrypt.CreateAesFactory()
	oldPw, _ = AesFact.CbcEncryptToHex(oldPw)
	newPw, _ = AesFact.CbcEncryptToHex(newPw)

	//
	if !crd.HasByUsernameAndPw(adminModel.Username, oldPw) {
		return -1, errors.New("旧密码不正确")
	}

	//修改密码
	updates := map[string]interface{}{
		"pw": newPw,
	}
	if !crd.masterModel.UpdateById(id, updates) {
		return -1, errors.New("修改密码失败")
	}
	return errCode, err
}
