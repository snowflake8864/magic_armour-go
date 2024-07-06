package organize_service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_user"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/organization/org_log_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateOrgCurdFactory() *OrgCurd {
	return &OrgCurd{masterModel: organization_models.CreateOrganizationFactory()}
}

type OrgCurd struct {
	masterModel *organization_models.OrganizationModel
}

// Add 新增企业
func (crd *OrgCurd) Add(orgParam *organization_models.OrganizationModel, userParam *user_models.UserModel, netIp string) (errCode int, err error) {
	//用户已存在
	userExist := user_service.CreateCurdFactory().GetItemByName(userParam.Username)
	if userExist.Id > 0 {
		return -1, errors.New("用户名已存在")
	}
	//企业已存在
	if orgParam.Type == enum_organization.TypePerson {
		orgParam.Name = ""
	} else {
		orgExist := crd.GetItemByName(orgParam.Name)
		if orgExist.Id > 0 {
			return -1, errors.New("企业已存在")
		}
	}

	tx := crd.masterModel.Begin()
	if tx.Error != nil {
		return -10, errors.New("保存企业失败")
	}

	//新增企业
	if orgParam.Source == enums.OperateSourcePlatform {
		orgParam.TrusteeshipSwitch = enums.IsYes
	} else {
		orgParam.TrusteeshipSwitch = enums.IsNot
	}
	orgParam.Status = enum_organization.StatusVerifying
	orgParam.IsDel = enums.IsNot
	orgParam.IsClean = enums.IsNot
	orgParam.LoginPhoneCheck = enums.IsNot
	orgParam.TrusteeshipSwitch = enums.IsYes
	orgParam.CreatedAt = date_helper.CurDatetime()
	orgParam.Code = md5_encrypt.MD5(uuid.New().String() + orgParam.CreatedAt)
	if err = tx.Create(orgParam).Error; err != nil {
		tx.Rollback()
		return -1, errors.New("保存企业失败")
	}

	//新增企业扩展
	OrgExtendParam := CreateOrgExtendFactory().InitDefault(orgParam.Id)
	if err = tx.Create(&OrgExtendParam).Error; err != nil {
		tx.Rollback()
		return -2, errors.New("保存企业配置失败")
	}

	//新增企业用户
	userParam.Password = user_service.CreateAuthFactory().EncryptPw(userParam.Password)
	userParam.OrgId = orgParam.Id
	userParam.Status = enum_user.StatusNormal
	userParam.LastLoginAt = orgParam.CreatedAt
	userParam.CreatedAt = orgParam.CreatedAt
	if err = tx.Create(&userParam).Error; err != nil {
		tx.Rollback()
		return -3, errors.New("注册企业用户失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgParam, enum_organization.LogActionPlatformAdd, netIp)
	orgLogParam.CreatedAt = orgParam.CreatedAt
	orgLogParam.OperateSource = orgParam.Source
	orgLogParam.OperateId = orgParam.CreateId
	orgLogParam.OperateName = orgParam.CreateName
	orgLogParam.Content = orgParam.Remark
	if err = tx.Create(&orgLogParam).Error; err != nil {
		tx.Rollback()
		return -3, errors.New("保存企业日志失败")
	}

	tx.Commit()
	return errCode, err
}

// Approved 审核: 批准
func (crd *OrgCurd) Approved(orgId int64, reason, netIp string, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	//验证企业
	orgExist := crd.GetItemByOrg(orgId)
	if orgExist.Id == 0 {
		return -1, errors.New("企业不存在")
	}
	if orgExist.Status != enum_organization.StatusVerifying {
		return -1, errors.New("企业状态异常，不能进行审核操作")
	}

	//修改
	updates := map[string]interface{}{
		"status": enum_organization.StatusNormal,
	}
	if !crd.masterModel.UpdateById(orgId, updates) {
		return -1, errors.New("企业账户审核失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgExist, enum_organization.LogActionPlatformApprove, netIp)
	orgLogParam.OperateSource = enums.OperateSourcePlatform
	orgLogParam.OperateId = loginInfo.AdminId
	orgLogParam.OperateName = loginInfo.Username
	orgLogParam.Content = reason
	if !logFact.Add(&orgLogParam) {
		return -3, errors.New("保存企业日志失败")
	}

	//新增

	return errCode, err
}

// Lock 管理: 锁定
func (crd *OrgCurd) Lock(orgId int64, reason, netIp string, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	//验证企业
	orgExist := crd.GetItemByOrg(orgId)
	if orgExist.Id == 0 {
		return -1, errors.New("企业不存在")
	}
	if orgExist.Status != enum_organization.StatusNormal {
		return -1, errors.New("企业状态异常，不能进行锁定操作")
	}

	//修改
	updates := map[string]interface{}{
		"status": enum_organization.StatusLocking,
	}
	if !crd.masterModel.UpdateById(orgId, updates) {
		return -1, errors.New("企业账户锁定失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgExist, enum_organization.LogActionPlatformLock, netIp)
	orgLogParam.OperateSource = enums.OperateSourcePlatform
	orgLogParam.OperateId = loginInfo.AdminId
	orgLogParam.OperateName = loginInfo.Username
	orgLogParam.Content = reason
	if !logFact.Add(&orgLogParam) {
		return -3, errors.New("保存企业日志失败")
	}

	return errCode, err
}

// Unlock 管理: 解锁
func (crd *OrgCurd) Unlock(orgId int64, reason, netIp string, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	//验证企业
	orgExist := crd.GetItemByOrg(orgId)
	if orgExist.Id == 0 {
		return -1, errors.New("企业不存在")
	}
	if orgExist.Status != enum_organization.StatusLocking {
		return -1, errors.New("企业状态异常，不能进行解锁操作")
	}

	//修改
	updates := map[string]interface{}{
		"status":   enum_organization.StatusNormal,
		"is_clean": enums.IsNot,
	}
	if !crd.masterModel.UpdateById(orgId, updates) {
		return -1, errors.New("企业账户解锁失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgExist, enum_organization.LogActionPlatformUnlock, netIp)
	orgLogParam.OperateSource = enums.OperateSourcePlatform
	orgLogParam.OperateId = loginInfo.AdminId
	orgLogParam.OperateName = loginInfo.Username
	orgLogParam.Content = reason
	if !logFact.Add(&orgLogParam) {
		return -3, errors.New("保存企业日志失败")
	}

	return errCode, err
}

// Clean 管理: 清理业务数据
func (crd *OrgCurd) Clean(orgId int64, reason, netIp string, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	//验证企业
	orgExist := crd.GetItemByOrg(orgId)
	if orgExist.Id == 0 {
		return -1, errors.New("企业不存在")
	}
	if orgExist.Status != enum_organization.StatusLocking {
		return -1, errors.New("企业状态不是锁定中，不能进行清理操作")
	}
	if orgExist.IsClean == enums.IsYes {
		return -1, errors.New("企业不能重复清理数据")
	}

	//修改
	updates := map[string]interface{}{
		"status":   enum_organization.StatusLocking,
		"is_clean": enums.IsYes,
	}
	if !crd.masterModel.UpdateById(orgId, updates) {
		return -1, errors.New("企业清理业务数据失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgExist, enum_organization.LogActionPlatformClean, netIp)
	orgLogParam.OperateSource = enums.OperateSourcePlatform
	orgLogParam.OperateId = loginInfo.AdminId
	orgLogParam.OperateName = loginInfo.Username
	orgLogParam.Content = reason
	if !logFact.Add(&orgLogParam) {
		return -3, errors.New("保存企业日志失败")
	}

	//企业级卸载终端任务
	task_service.CreateCacheFactory().AddOrgUninstall(orgId)

	//丢入队列
	msg := &kafka_go.MessageStruct{
		Type:  kafka_go.TypeCommonOrgClean,
		OrgId: orgId,
	}
	producer_kafka.CreateProducerFactory().ProduceCommon(context.Background(), msg)

	return errCode, err
}

// SetHsSwitch 管理: 设置幻势开关
func (crd *OrgCurd) SetHsSwitch(orgId int64, openSwitch int, netIp string, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	//验证企业
	orgExist := crd.GetItemByOrg(orgId)
	if orgExist.Id == 0 {
		return -1, errors.New("企业不存在")
	}

	//修改
	updates := map[string]interface{}{
		"huanshi_switch": openSwitch,
	}
	if !crd.masterModel.UpdateById(orgId, updates) {
		return -1, errors.New("保存数据失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgExist, enum_organization.LogActionPlatformHuanshiSwitch, netIp)
	orgLogParam.OperateSource = enums.OperateSourcePlatform
	orgLogParam.OperateId = loginInfo.AdminId
	orgLogParam.OperateName = loginInfo.Username
	if openSwitch == enums.IsYes {
		orgLogParam.Content = "开启"
	} else {
		orgLogParam.Content = "关闭"
	}
	if !logFact.Add(&orgLogParam) {
		return -3, errors.New("保存企业日志失败")
	}

	return errCode, err
}

// SetLoginSmsSwitch 管理: 设置企业登录短信验证开关
func (crd *OrgCurd) SetLoginSmsSwitch(orgId int64, openSwitch int, netIp string, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	//验证企业
	orgExist := crd.GetItemByOrg(orgId)
	if orgExist.Id == 0 {
		return -1, errors.New("企业不存在")
	}

	//修改
	updates := map[string]interface{}{
		"login_phone_check": openSwitch,
	}
	if !crd.masterModel.UpdateById(orgId, updates) {
		return -1, errors.New("保存数据失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgExist, enum_organization.LogActionPlatformLoginSmsSwitch, netIp)
	orgLogParam.OperateSource = enums.OperateSourcePlatform
	orgLogParam.OperateId = loginInfo.AdminId
	orgLogParam.OperateName = loginInfo.Username
	if openSwitch == enums.IsYes {
		orgLogParam.Content = "开启"
	} else {
		orgLogParam.Content = "关闭"
	}
	if !logFact.Add(&orgLogParam) {
		return -3, errors.New("保存企业日志失败")
	}

	return errCode, err
}

// Refused 审核: 回绝
func (crd *OrgCurd) Refused(orgId int64, reason, netIp string, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	//验证企业
	orgExist := crd.GetItemByOrg(orgId)
	if orgExist.Id == 0 {
		return -1, errors.New("企业不存在")
	}
	if orgExist.Status != enum_organization.StatusVerifying {
		return -1, errors.New("企业状态异常，不能进行审核操作")
	}

	//删除企业
	if !crd.masterModel.DeleteById(orgId) {
		return -1, errors.New("企业账户删除失败")
	}

	//删除企业扩展
	if !CreateOrgExtendFactory().DeleteByOrg(orgId) {
		return -1, errors.New("企业配置删除失败")
	}

	//删除企业用户
	if !user_service.CreateCurdFactory().DeleteByOrg(orgId) {
		return -1, errors.New("企业用户删除失败")
	}

	//新增企业日志
	logFact := org_log_service.CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgExist, enum_organization.LogActionPlatformRefuse, netIp)
	orgLogParam.OperateSource = enums.OperateSourcePlatform
	orgLogParam.OperateId = loginInfo.AdminId
	orgLogParam.OperateName = loginInfo.Username
	orgLogParam.Content = reason
	if !logFact.Add(&orgLogParam) {
		return -3, errors.New("保存企业日志失败")
	}

	return errCode, err
}

// GetItemByOrg 通过企业编号获取
func (crd *OrgCurd) GetItemByOrg(orgId int64) (model *organization_models.OrganizationModel) {
	model, _ = crd.masterModel.GetItemById(orgId)
	return model
}

// GetItemByName 通过企业名称获取
func (crd *OrgCurd) GetItemByName(name string) (model *organization_models.OrganizationModel) {
	model, _ = crd.masterModel.GetItemByName(name)
	return model
}
