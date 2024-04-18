package organize_service

import (
	"errors"
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
	if orgParam.Source == enum_organization.SourcePlatform {
		orgParam.TrusteeshipSwitch = enums.IsYes
		orgParam.HuanshiSwitch = enums.IsYes
	} else {
		orgParam.Status = enum_organization.StatusNormal
		orgParam.TrusteeshipSwitch = enums.IsNot
		orgParam.HuanshiSwitch = enums.IsNot
	}
	orgParam.Status = enum_organization.StatusVerifying
	orgParam.IsDel = enums.IsNot
	orgParam.IsClean = enums.IsNot
	orgParam.LoginPhoneCheck = enums.IsNot
	orgParam.CreatedAt = helpers.GetDatetime()
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
	userParam.OrgId = orgParam.Id
	userParam.LastLoginAt = orgParam.CreatedAt
	userParam.CreatedAt = orgParam.CreatedAt
	if err = tx.Create(&userParam).Error; err != nil {
		tx.Rollback()
		return -3, errors.New("注册企业用户失败")
	}

	//新增企业日志
	logFact := CreateOrgLogFactory()
	orgLogParam := logFact.InitByOrg(orgParam, enum_organization.LogActionPlatformAdd, netIp)
	orgLogParam.CreatedAt = orgParam.CreatedAt
	orgLogParam.OperateSource = orgParam.Source
	orgLogParam.OperateId = orgParam.CreateId
	orgLogParam.OperateName = orgParam.CreateName
	orgLogParam.Content = orgParam.Remark
	if err = tx.Create(&orgLogParam).Error; err != nil {
		crd.masterModel.Rollback()
		return -3, errors.New("保存企业日志失败")
	}

	tx.Commit()
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
