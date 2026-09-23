package organize_service

import (
	"fmt"
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

func (crd *Koal) Sync(models []*organization_models.OrganizationModel) error {
	if len(models) == 0 {
		return nil
	}

	for _, model := range models {
		//是否存在
		if hasModel, _ := crd.orgModel.GetItemByCode(model.Code); hasModel == nil || hasModel.Id == 0 {
			if err := crd.add(model); err != nil {
				return fmt.Errorf("add [%s] fail: %v", model.Code, err)
			}

		} else {
			model.Id = hasModel.Id
			if err := crd.modify(model); err != nil {
				return fmt.Errorf("modify [%s] fail: %v", model.Code, err)
			}
		}
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
		orgLogParam.Content = "新增企业"
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

func (crd *Koal) modify(orgParam *organization_models.OrganizationModel) error {
	err := crd.orgModel.Transaction(func(tx *gorm.DB) error {
		//保存企业
		updates := map[string]interface{}{
			"name":   orgParam.Name,
			"status": orgParam.Status,
			"is_del": orgParam.IsDel,
		}
		if err := tx.Model(&organization_models.OrganizationModel{}).Where("id=?", orgParam.Id).Updates(updates).Error; err != nil {
			return err
		}

		//新增企业日志
		logFact := org_log_service.CreateOrgLogFactory()
		logModule := enum_organization.LogActionKoalSyncOrg
		orgLogParam := logFact.InitByOrg(orgParam, logModule, "")
		orgLogParam.CreatedAt = orgParam.CreatedAt
		orgLogParam.Content = "修改企业"
		orgLogParam.OperateSource = orgParam.Source
		if err := tx.Create(&orgLogParam).Error; err != nil {
			tx.Rollback()
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
