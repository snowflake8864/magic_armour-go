package department_service

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"strings"
)

const SEPARATOR string = ","

func CreateCurdFactory() *DepartmentCurd {
	return &DepartmentCurd{masterModel: department_models.CreateDepartmentFactory()}
}

type DepartmentCurd struct {
	masterModel *department_models.DepartmentModel
}

// DeleteByOrg 通过企业编号删除
func (crd *DepartmentCurd) getItemById(orgId, id int64) *department_models.DepartmentModel {
	model, _ := crd.masterModel.GetItemById(orgId, id)
	return model
}

// GetMapsByIds 通过ids获取多条记录, 有映射关系
func (crd *DepartmentCurd) GetMapsByIds(orgId int64, ids []int64) map[int64]*department_models.DepartmentModel {
	conditions := department_models.SearchConditions{
		OrgId: orgId,
		Ids:   ids,
	}
	models, err := crd.masterModel.GetListByConditions(conditions)
	if err != nil {
		return nil
	}
	if models == nil {
		return nil
	}
	result := make(map[int64]*department_models.DepartmentModel)
	for _, model := range models {
		result[model.Id] = model
	}
	return result
}

// Add 新增
func (crd *DepartmentCurd) Add(param *department_models.DepartmentModel) error {
	//等级从1开始
	param.Level = 1
	param.CreatedAt = date_helper.CurDatetime()

	var fullName []string
	if param.ParentId > 0 {
		if parentModel := crd.getItemById(param.OrgId, param.ParentId); parentModel.Id > 0 {
			param.Level = parentModel.Level + 1
			fullName = json_helper.ToSliceString(parentModel.FullName)
			param.FullCode = parentModel.FullCode
		} else {
			return errors.New("父级菜单不存在")
		}
	}
	fullName = append(fullName, param.Name)
	param.FullName = json_helper.Encode(fullName)

	err := crd.masterModel.Transaction(func(tx *gorm.DB) error {
		if !crd.masterModel.InsertItem(param) {
			return errors.New("添加失败")
		}

		//修改
		conditions := department_models.SearchConditions{
			OrgId: param.OrgId,
			Id:    param.Id,
		}
		updates := map[string]interface{}{
			"full_code": strings.Trim(param.FullCode+SEPARATOR+int_helper.Int64ToStr(param.Id), SEPARATOR),
		}
		if !crd.masterModel.UpdateByConditions(conditions, updates) {
			return errors.New("保存失败")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteByOrg 通过企业编号删除
func (crd *DepartmentCurd) DeleteByOrg(orgId int64) bool {
	conditions := department_models.SearchConditions{
		OrgId: orgId,
	}
	return crd.masterModel.DeleteByConditions(conditions)
}
