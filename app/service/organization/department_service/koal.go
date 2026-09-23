package department_service

import (
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/organization/org_log_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"strings"
)

func CreateKoalFactory() *Koal {
	return &Koal{
		deptModel:    department_models.CreateDepartmentFactory(),
		orgModel:     organization_models.CreateOrganizationFactory(),
		thirdOrgMap:  make(map[string]*organization_models.OrganizationModel),
		thirdDeptMap: make(map[string]*department_models.DepartmentModel),
	}
}

type Koal struct {
	deptModel    *department_models.DepartmentModel
	orgModel     *organization_models.OrganizationModel
	thirdOrgMap  map[string]*organization_models.OrganizationModel //三方企业映射
	thirdDeptMap map[string]*department_models.DepartmentModel     //三方部门映射
}

func (crd *Koal) Sync(models []*department_models.DepartmentModel) error {
	if len(models) == 0 {
		return nil
	}

	for _, model := range models {
		//是否存在
		if hasModel, _ := crd.deptModel.GetItemByCode(model.Code); hasModel == nil || hasModel.Id == 0 {
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

func (crd *Koal) getOrgModel(thirdCode int64) *organization_models.OrganizationModel {
	code := int_helper.Int64ToStr(thirdCode)
	if model, ok := crd.thirdOrgMap[code]; !ok {
		if hasModel, _ := crd.orgModel.GetItemByCode(code); hasModel == nil || hasModel.Id == 0 {
			crd.thirdOrgMap[code] = nil
			return nil
		} else {
			crd.thirdOrgMap[code] = hasModel

			return hasModel
		}
	} else {
		return model
	}

}

func (crd *Koal) getDeptModel(thirdCode int64) *department_models.DepartmentModel {
	code := int_helper.Int64ToStr(thirdCode)
	if model, ok := crd.thirdDeptMap[code]; !ok {
		if hasModel, _ := crd.deptModel.GetItemByCode(code); hasModel == nil || hasModel.Id == 0 {
			crd.thirdDeptMap[code] = nil
			return nil
		} else {
			crd.thirdDeptMap[code] = hasModel
			return hasModel
		}
	} else {
		return model
	}
}

func (crd *Koal) add(deptParam *department_models.DepartmentModel) error {
	//获取企业
	orgModel := crd.getOrgModel(deptParam.OrgId)
	if orgModel == nil {
		return fmt.Errorf("org[%d] not exist", deptParam.OrgId)
	} else {
		//真正的企业ID
		deptParam.OrgId = orgModel.Id
	}

	//父级信息
	if deptParam.ParentId > 0 {
		if parentModel := crd.getDeptModel(deptParam.ParentId); parentModel == nil {
			return fmt.Errorf("org[%d] pid[%d] not exist", deptParam.OrgId, deptParam.ParentId)
		} else {
			//真正的父级ID
			deptParam.ParentId = parentModel.Id
			//继承父级
			deptParam.FullCode = parentModel.FullCode
		}
	}

	err := crd.orgModel.Transaction(func(tx *gorm.DB) error {
		//保存组织结构
		if err := tx.Create(deptParam).Error; err != nil {
			return err
		}

		//保存组织结构完整编号
		updates := map[string]interface{}{
			"full_code": strings.Trim(deptParam.FullCode+SEPARATOR+int_helper.Int64ToStr(deptParam.Id), SEPARATOR),
		}
		if err := tx.Model(&department_models.DepartmentModel{}).Where("id=?", deptParam.Id).Updates(updates).Error; err != nil {
			return err
		}

		//新增企业日志
		logFact := org_log_service.CreateOrgLogFactory()
		logModule := enum_organization.LogActionKoalSyncOrg
		orgLogParam := logFact.InitByOrg(orgModel, logModule, "")
		orgLogParam.CreatedAt = date_helper.CurDatetime()
		orgLogParam.Content = fmt.Sprintf("新增部门(%s)", deptParam.Code)
		orgLogParam.OperateSource = enums.OperateSourcekoal
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

func (crd *Koal) modify(deptParam *department_models.DepartmentModel) error {
	//获取企业
	orgModel := crd.getOrgModel(deptParam.OrgId)
	if orgModel == nil {
		return fmt.Errorf("org[%d] not exist", deptParam.OrgId)
	} else {
		//真正的企业ID
		deptParam.OrgId = orgModel.Id
	}

	//父级信息
	if deptParam.ParentId > 0 {
		if parentModel := crd.getDeptModel(deptParam.ParentId); parentModel == nil {
			return fmt.Errorf("org[%d] pid[%d] not exist", deptParam.OrgId, deptParam.ParentId)
		} else {
			//真正的父级ID
			deptParam.ParentId = parentModel.Id
			//继承父级
			deptParam.FullCode = parentModel.FullCode
		}
	}

	err := crd.orgModel.Transaction(func(tx *gorm.DB) error {

		//保存
		updates := map[string]interface{}{
			"org_id":    deptParam.OrgId,
			"parent_id": deptParam.ParentId,
			"name":      deptParam.Name,
			"full_code": strings.Trim(deptParam.FullCode+SEPARATOR+int_helper.Int64ToStr(deptParam.Id), SEPARATOR),
			"full_name": deptParam.FullName,
			"level":     deptParam.Level,
			"sort":      deptParam.Sort,
		}
		if err := tx.Model(&department_models.DepartmentModel{}).Where("id=?", deptParam.Id).Updates(updates).Error; err != nil {
			return err
		}

		//新增企业日志
		logFact := org_log_service.CreateOrgLogFactory()
		logModule := enum_organization.LogActionKoalSyncOrg
		orgLogParam := logFact.InitByOrg(orgModel, logModule, "")
		orgLogParam.CreatedAt = date_helper.CurDatetime()
		orgLogParam.Content = fmt.Sprintf("修改部门(%s)", deptParam.Code)
		orgLogParam.OperateSource = enums.OperateSourcekoal
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

func (crd *Koal) DeleteByCodes(codes []string) error {

	return nil
}
