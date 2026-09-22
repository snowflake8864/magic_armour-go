package organize_service

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/organization/org_log_service"
	"goskeleton/app/service/organization/role_service"
)

func CreateKoalFactory() *Koal {
	return &Koal{orgModel: organization_models.CreateOrganizationFactory()}
}

type Koal struct {
	orgModel *organization_models.OrganizationModel
}

func (crd *Koal) Sync(orgParam *organization_models.OrganizationModel) error {
	//是否存在
	if hasModel, _ := crd.orgModel.GetItemByCode(orgParam.Code); hasModel == nil || hasModel.Id == 0 {
		if err := crd.add(orgParam); err != nil {
			return err
		}

	} else {

	}

	return nil
}

func (crd *Koal) add(orgParam *organization_models.OrganizationModel) error {
	err := crd.orgModel.Transaction(func(tx *gorm.DB) error {
		//保存企业
		if err := tx.Create(orgParam).Error; err != nil {
			return err
		}

		//保存扩展
		OrgExtendParam := CreateOrgExtendFactory().InitDefault(orgParam.Id)
		if err := tx.Create(&OrgExtendParam).Error; err != nil {
			return err
		}

		//新增企业日志
		logFact := org_log_service.CreateOrgLogFactory()
		logModule := enum_organization.LogActionKoalSyncOrg
		orgLogParam := logFact.InitByOrg(orgParam, logModule, "")
		orgLogParam.CreatedAt = orgParam.CreatedAt
		orgLogParam.Content = orgParam.Remark
		orgLogParam.OperateSource = orgParam.Source
		if err := tx.Create(&orgLogParam).Error; err != nil {
			tx.Rollback()
		}

		return nil
	})

	if err != nil {
		return err
	}

	//默认企业角色
	role_service.CreateCurdFactory().AddDefault(orgParam.Id)

	return nil
}
