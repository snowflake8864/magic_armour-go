package user_service

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_user"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/validator/common/common_regexp"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/org_log_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *UserCurd {
	return &UserCurd{masterModel: user_models.CreateUserFactory()}
}

type UserCurd struct {
	masterModel *user_models.UserModel
}

func (crd *UserCurd) getOrg(OrgId int64) *organization_models.OrganizationModel {
	orgModel, err := organization_models.CreateOrganizationFactory().GetItemById(OrgId)
	if err != nil || orgModel.Id == 0 {
		return nil
	}
	return orgModel

}

// AddForSub 添加子用户
func (crd *UserCurd) AddForSub(userParam *user_models.UserModel, netIp string) error {
	if err := common_regexp.CheckUserName(userParam.Username); err != nil {
		return err
	}
	if err := common_regexp.CheckPhone(userParam.Phone); err != nil {
		return err
	}
	if err := common_regexp.CheckEmail(userParam.Email); err != nil {
		return err
	}
	if err := common_regexp.CheckPw(userParam.Password); err != nil {
		return err
	}
	//用户已存在
	if userExist := crd.GetItemByName(userParam.Username); userExist != nil && userExist.Id > 0 {
		return errors.New(my_errors.CurdUsernameAlreadyExists)
	}
	if userExist := crd.GetItemByPhone(userParam.Phone); userExist != nil && userExist.Id > 0 {
		return errors.New(my_errors.CurdPhoneAlreadyExists)
	}
	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//是否超过限制
	subUserLimit := string_helper.ToInt64(dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.OrganizationSubUserLimit))
	if subUserLimit > 0 {
		condition := user_models.SearchConditions{OrgId: userParam.OrgId}
		if crd.masterModel.Count(condition) > subUserLimit {
			return errors.New(my_errors.CurdOUserMaxLimit)
		}
	}

	now := date_helper.CurDatetime()
	userParam.Password = CreateAuthFactory().EncryptPw(userParam.Password)
	userParam.Status = enum_user.StatusNormal
	userParam.IsRoot = enums.IsNot
	userParam.PwErrorNum = 0
	userParam.CreatedAt = now
	if !crd.masterModel.InsertItem(userParam) {
		return errors.New(my_errors.CurdOrgUserSaveFailed)
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgAddSubUser, netIp)
	orgLogParam.CreatedAt = userParam.CreatedAt
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}
	return nil
}

// UpdateForSub 修改子用户
func (crd *UserCurd) UpdateForSub(userParam *user_models.UserModel, netIp string) error {
	if err := common_regexp.CheckPhone(userParam.Phone); err != nil {
		return err
	}
	if err := common_regexp.CheckEmail(userParam.Email); err != nil {
		return err
	}

	//验证用户
	if userExist := crd.GetItemByOrgAndUser(userParam.OrgId, userParam.Id); userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	} else if userExist.IsRoot == enums.IsYes {
		return errors.New(my_errors.CurdOUserAdminNotModify)
	} else if userExist.Phone != userParam.Phone {
		if userPhone := crd.GetItemByPhone(userParam.Phone); userPhone != nil && userPhone.Id > 0 {
			return errors.New(my_errors.CurdPhoneAlreadyExists)
		}
	}

	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//修改
	updates := map[string]interface{}{
		"email":    userParam.Email,
		"phone":    userParam.Phone,
		"realname": userParam.Realname,
		"position": userParam.Position,
		"intro":    userParam.Intro,
		"depts":    userParam.Depts,
		"role_id":  userParam.RoleId,
	}
	if isOk, _ := crd.masterModel.UpdateById(userParam.Id, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}
	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgEditSubUser, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}
	return nil
}

// ModifyPwForSub 子用户修改密码
func (crd *UserCurd) ModifyPwForSub(userParam *user_models.UserModel, netIp string) error {
	if err := common_regexp.CheckPw(userParam.Password); err != nil {
		return err
	}
	//验证用户
	if userExist := crd.GetItemByOrgAndUser(userParam.OrgId, userParam.Id); userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	} else if userExist.IsRoot == enums.IsYes {
		return errors.New(my_errors.CurdPhoneAlreadyExists)
	}
	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//修改
	updates := map[string]interface{}{
		"password":     CreateAuthFactory().EncryptPw(userParam.Password),
		"pw_modify_at": date_helper.CurTimestamp(),
	}
	if isOk, _ := crd.masterModel.UpdateById(userParam.Id, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}
	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgModifyPwSubUser, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}
	return nil
}

// DeleteForSub 删除子用户
func (crd *UserCurd) DeleteForSub(userParam *user_models.UserModel, netIp string) error {
	//验证用户
	if userExist := crd.GetItemByOrgAndUser(userParam.OrgId, userParam.Id); userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	} else if userExist.IsRoot == enums.IsYes {
		return errors.New(my_errors.CurdOUserAdminNotDelete)
	}
	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	conditions := user_models.SearchConditions{
		OrgId: userParam.OrgId,
		Id:    userParam.Id,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdSaveFailed)
	}
	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgDelSubUser, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
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

// GetItemByPhone 通过手机获取
func (crd *UserCurd) GetItemByPhone(phone string) (model *user_models.UserModel) {
	model, _ = crd.masterModel.GetItemByPhone(phone)
	return model
}

// HasByUsernameAndPw 用户名和密码一致
func (crd *UserCurd) HasByUsernameAndPw(name, pw string) bool {
	if userModel := crd.masterModel.GetItemByNameAndPw(name, pw); userModel.Id > 0 {
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

// ModifyPwById 修改密码
func (crd *UserCurd) ModifyPwById(userParam *user_models.UserModel, newPw string, netIp string) error {
	if err := common_regexp.CheckPw(newPw); err != nil {
		return err
	}
	//用户已存在
	userExist := crd.GetItemByOrgAndUser(userParam.OrgId, userParam.Id)
	if userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	}
	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//加密
	AuthFact := CreateAuthFactory()
	oldPw := AuthFact.EncryptPw(userParam.Password)
	newPw = AuthFact.EncryptPw(newPw)
	if !crd.HasByUsernameAndPw(userExist.Username, oldPw) {
		return errors.New(my_errors.CurdOldPasswordIncorrect)
	}

	//修改密码
	updates := map[string]interface{}{
		"password":     newPw,
		"pw_modify_at": date_helper.CurTimestamp(),
	}
	if isOk, _ := crd.masterModel.UpdateById(userExist.Id, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgSelfModifyPw, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}
	return nil
}

// ModifyPhoneById 修改手机
func (crd *UserCurd) ModifyPhoneById(userParam *user_models.UserModel, netIp string) error {
	if err := common_regexp.CheckPhone(userParam.Phone); err != nil {
		return err
	}
	//用户已存在
	userExist := crd.GetItemByOrgAndUser(userParam.OrgId, userParam.Id)
	if userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	}
	if userModel := crd.GetItemByPhone(userParam.Phone); userExist != nil && userModel.Id > 0 && userModel.Id != userExist.Id {
		return errors.New(my_errors.CurdPhoneAlreadyExists)
	}

	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//加密
	AuthFact := CreateAuthFactory()
	oldPw := AuthFact.EncryptPw(userParam.Password)
	if !crd.HasByUsernameAndPw(userExist.Username, oldPw) {
		return errors.New(my_errors.CurdOldPasswordIncorrect)
	}

	//修改手机
	updates := map[string]interface{}{
		"phone": userParam.Phone,
	}
	if isOk, _ := crd.masterModel.UpdateById(userExist.Id, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//同步修改企业
	if userExist.IsRoot == enums.IsYes {
		organization_models.CreateOrganizationFactory().UpdateById(userParam.OrgId, map[string]interface{}{
			"user_phone": userParam.Phone,
		})
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgSelfModifyPhone, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}

	return nil
}

// ModifyEmailById 修改邮箱
func (crd *UserCurd) ModifyEmailById(userParam *user_models.UserModel, netIp string) error {
	if err := common_regexp.CheckEmail(userParam.Email); err != nil {
		return err
	}
	//用户已存在
	userExist := crd.GetItemByOrgAndUser(userParam.OrgId, userParam.Id)
	if userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	}
	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//加密
	AuthFact := CreateAuthFactory()
	oldPw := AuthFact.EncryptPw(userParam.Password)
	if !crd.HasByUsernameAndPw(userExist.Username, oldPw) {
		return errors.New(my_errors.CurdOldPasswordIncorrect)
	}

	//修改邮箱失
	updates := map[string]interface{}{
		"email": userParam.Email,
	}
	if isOk, _ := crd.masterModel.UpdateById(userParam.Id, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//同步修改企业
	if userExist.IsRoot == enums.IsYes {
		organization_models.CreateOrganizationFactory().UpdateById(userParam.OrgId, map[string]interface{}{
			"user_email": userParam.Email,
		})
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgSelfModifyEmail, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}
	return nil
}

// ModifyInfoById 修改信息
func (crd *UserCurd) ModifyInfoById(userParam *user_models.UserModel, netIp string) error {
	//用户已存在
	userExist := crd.GetItemByOrgAndUser(userParam.OrgId, userParam.Id)
	if userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	}
	//
	orgModel := crd.getOrg(userParam.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//修改
	updates := map[string]interface{}{
		"position": userParam.Position,
		"intro":    userParam.Intro,
	}
	if isOk, _ := crd.masterModel.UpdateById(userParam.Id, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}
	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgSelfModifyInfo, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userParam.CreateId
	orgLogParam.OperateName = userParam.CreateName
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}
	return nil
}

// FindPwByPhone 找回密码
func (crd *UserCurd) FindPwByPhone(userPhone, newPw string, netIp string) error {
	if err := common_regexp.CheckPw(newPw); err != nil {
		return err
	}
	//用户已存在
	userExist := crd.GetItemByPhone(userPhone)
	if userExist == nil || userExist.Id == 0 {
		return errors.New(my_errors.CurdUserNotExist)
	}
	//
	orgModel := crd.getOrg(userExist.OrgId)
	if orgModel == nil {
		return errors.New(my_errors.CurdOrgNotExists)
	}

	//修改密码
	updates := map[string]interface{}{
		"password":     CreateAuthFactory().EncryptPw(newPw),
		"pw_modify_at": date_helper.CurTimestamp(),
	}
	if isOk, _ := crd.masterModel.UpdateById(userExist.Id, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgModel, enum_organization.LogActionOrgSelfFindPw, netIp)
	orgLogParam.CreatedAt = date_helper.CurDatetime()
	orgLogParam.OperateSource = enums.OperateSourceOrganize
	orgLogParam.OperateId = userExist.Id
	orgLogParam.OperateName = userExist.Username
	if !logFact.Add(&orgLogParam) {
		return errors.New(my_errors.CurdOrgLogSaveFailed)
	}
	return nil
}
