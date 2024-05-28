package user_service

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_user"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/org_log_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateCurdFactory() *UserCurd {
	return &UserCurd{masterModel: user_models.CreateUserFactory()}
}

type UserCurd struct {
	masterModel *user_models.UserModel
}

// AddSub 添加子用户
func (crd *UserCurd) AddSub(userParam *user_models.UserModel, netIp string) error {
	//用户已存在
	userExist := crd.GetItemByName(userParam.Username)
	if userExist.Id > 0 {
		return errors.New("用户名已存在")
	}
	//
	orgModel, err := organization_models.CreateOrganizationFactory().GetItemById(userParam.OrgId)
	if err != nil || orgModel.Id == 0 {
		return errors.New("企业数据错误")
	}

	now := date_helper.CurDatetime()
	userParam.Password = CreateAuthFactory().EncryptPw(userParam.Password)
	userParam.Status = enum_user.StatusNormal
	userParam.PwErrorNum = 0
	userParam.LastLoginAt = now
	userParam.CreatedAt = now
	if !crd.masterModel.InsertItem(userParam) {
		return errors.New("添加子户失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgAddSubUser, netIp)
	orgLogParam.CreatedAt = userParam.CreatedAt
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New("保存企业日志失败")
	}
	return nil
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
		"last_login_at": date_helper.CurDatetime(),
	}
	if isOk, _ := crd.masterModel.UpdateById(id, updates); isOk {
		return true
	} else {
		return false
	}
}
