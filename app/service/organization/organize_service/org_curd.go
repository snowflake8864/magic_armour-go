package organize_service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/model/model_organization"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateOrgCurdFactory() *OrgCurd {
	return &OrgCurd{masterModel: model_organization.CreateOrganizationFactory()}
}

type OrgCurd struct {
	OrgBase
	masterModel *model_organization.OrganizationModel
}

// Add 新增企业
func (crd *OrgCurd) Add(orgParam *model_organization.OrganizationModel, userParam *model_organization.UserModel, netIp string) (errCode int, err error) {
	crd.masterModel.Begin()
	//新增企业
	if orgParam.Type == enum_organization.TypePerson {
		orgParam.Name = ""
	}
	if orgParam.Source == enum_organization.SourcePlatform {
		orgParam.Status = enum_organization.StatusNormal
		orgParam.HuanshiSwitch = enums.IsYes
	} else {
		orgParam.Status = enum_organization.StatusNormal
		orgParam.HuanshiSwitch = enums.IsNot
	}
	orgParam.IsDel = enums.IsNot
	orgParam.IsClean = enums.IsNot
	orgParam.TrusteeshipSwitch = enums.IsNot
	orgParam.LoginPhoneCheck = enums.IsNot
	orgParam.CreatedAt = helpers.GetDatetime()
	orgParam.Code = md5_encrypt.MD5(uuid.New().String() + orgParam.CreatedAt)
	if crd.masterModel.InsertItem(orgParam) {
		crd.masterModel.Rollback()
		return -1, errors.New("保存企业失败")
	}
	fmt.Println("===========================")
	fmt.Println(orgParam.Id)

	//新增企业扩展
	if CreateOrgExtendFactory().AddDefault(orgParam.Id) {
		crd.masterModel.Rollback()
		return -2, errors.New("保存企业配置失败")
	}

	//新增企业用户
	userParam.OrgId = orgParam.Id
	userParam.CreatedAt = orgParam.CreatedAt
	if user_service.CreateCurdFactory().Register(userParam) {
		crd.masterModel.Rollback()
		return -3, errors.New("注册企业用户失败")
	}

	//新增企业日志
	logFact := CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgParam, enum_organization.LogActionPlatformAdd, netIp)
	orgLogParam.CreatedAt = orgParam.CreatedAt
	orgLogParam.OperateSource = orgParam.Source
	orgLogParam.OperateId = orgParam.CreateId
	orgLogParam.OperateName = orgParam.CreateName
	if user_service.CreateCurdFactory().Register(userParam) {
		crd.masterModel.Rollback()
		return -3, errors.New("保存企业日志失败")
	}
	crd.masterModel.Commit()
	return errCode, err
}

// GetItemByOrg 通过企业编号获取
func (crd *OrgCurd) GetItemByOrg(orgId int64) (model *model_organization.OrganizationModel) {
	model, _ = crd.masterModel.GetItemById(orgId)
	return model
}

// GetItemByName 通过企业名称获取
func (crd *OrgCurd) GetItemByName(name string) (model *model_organization.OrganizationModel) {
	model, _ = crd.masterModel.GetItemByName(name)
	return model
}

// Page 分页
func (crd *OrgCurd) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_organization.OrganizationModel) {
	where := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}
