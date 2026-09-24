package user_service

import (
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/org_log_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
)

func CreateKoalFactory() *Koal {
	return &Koal{
		userModel:    user_models.CreateUserFactory(),
		deptModel:    department_models.CreateDepartmentFactory(),
		orgModel:     organization_models.CreateOrganizationFactory(),
		thirdDeptMap: make(map[string]*department_models.DepartmentModel),
	}
}

type Koal struct {
	userModel    *user_models.UserModel
	orgModel     *organization_models.OrganizationModel
	deptModel    *department_models.DepartmentModel
	thirdDeptMap map[string]*department_models.DepartmentModel //三方部门映射
}

func (crd *Koal) Sync(saveMap map[string][]*user_models.UserModel) error {
	if len(saveMap) == 0 {
		return nil
	}

	for thirdOrgCode, models := range saveMap {
		//获取企业
		orgModel := crd.getOrgModel(thirdOrgCode)
		if orgModel == nil || orgModel.Id == 0 {
			return fmt.Errorf("org [%s]  not exist", thirdOrgCode)
		}

		//处理用户
		for _, model := range models {
			//是否存在
			if hasModel, _ := crd.userModel.GetItemByCode(model.Code); hasModel == nil || hasModel.Id == 0 {
				if err := crd.add(model, orgModel); err != nil {
					return fmt.Errorf("add [%s] fail: %v", model.Code, err)
				}

			} else {
				model.Id = hasModel.Id
				if err := crd.modify(model, orgModel); err != nil {
					return fmt.Errorf("modify [%s] fail: %v", model.Code, err)
				}
			}
		} //end models

		//重新选举企业超级管理员
	}

	return nil
}

func (crd *Koal) getOrgModel(thirdCode string) *organization_models.OrganizationModel {
	if hasModel, _ := crd.orgModel.GetItemByCode(thirdCode); hasModel == nil || hasModel.Id == 0 {

		return nil
	} else {
		return hasModel
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

func (crd *Koal) add(userParam *user_models.UserModel, orgModel *organization_models.OrganizationModel) error {

	err := crd.userModel.Transaction(func(tx *gorm.DB) error {
		//保存组织结构
		if err := tx.Create(userParam).Error; err != nil {
			return err
		}

		//新增企业日志
		logFact := org_log_service.CreateOrgLogFactory()
		logModule := enum_organization.LogActionKoalSyncOrg
		orgLogParam := logFact.InitByOrg(orgModel, logModule, "")
		orgLogParam.CreatedAt = date_helper.CurDatetime()
		orgLogParam.Content = fmt.Sprintf("新增用户(%s)", userParam.Code)
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

func (crd *Koal) modify(userParam *user_models.UserModel, orgModel *organization_models.OrganizationModel) error {

	err := crd.orgModel.Transaction(func(tx *gorm.DB) error {

		/*//保存
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
		}*/

		//新增企业日志
		logFact := org_log_service.CreateOrgLogFactory()
		logModule := enum_organization.LogActionKoalSyncOrg
		orgLogParam := logFact.InitByOrg(orgModel, logModule, "")
		orgLogParam.CreatedAt = date_helper.CurDatetime()
		orgLogParam.Content = fmt.Sprintf("修改用户(%s)", userParam.Code)
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

func (crd *Koal) DeleteByCodes(codes []string) bool {
	if len(codes) == 0 {
		return true
	}

	conditions := user_models.SearchConditions{
		Codes: codes,
	}

	return crd.userModel.DeleteByConditions(conditions)
}
