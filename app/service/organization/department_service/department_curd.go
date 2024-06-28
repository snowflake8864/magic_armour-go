package department_service

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
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

// GetIdsSelfWithParentById 获取本身以及否父级的id集合
func (crd *DepartmentCurd) GetIdsSelfWithParentById(orgId int64, id int64) []int64 {
	var result []int64
	model, err := crd.masterModel.GetItemById(orgId, id)
	if err != nil || model == nil {
		return result
	}

	for _, val := range strings.Split(model.FullCode, ",") {
		result = append(result, string_helper.ToInt64(val))
	}
	return result
}

// GetIdsSelfWithChildById 获取本身以及子集的id集合
func (crd *DepartmentCurd) GetIdsSelfWithChildById(orgId int64, id int64) []int64 {
	conditions := department_models.SearchConditions{
		OrgId:    orgId,
		FullCode: int_helper.Int64ToStr(id),
	}
	var result []int64
	models, err := crd.masterModel.GetListByConditions(conditions)
	if err != nil || models == nil {
		return result
	}

	for _, model := range models {
		result = append(result, model.Id)
	}
	return result
}

// GetIdsSelfWithChildByIds 获取本身以及子集的id集合
func (crd *DepartmentCurd) GetIdsSelfWithChildByIds(orgId int64, ids []int64) []int64 {
	conditions := department_models.SearchConditions{
		OrgId: orgId,
		Ids:   ids,
	}
	var result []int64
	models, err := crd.masterModel.GetListByConditions(conditions)
	if err != nil || models == nil {
		return result
	}

	for _, model := range models {
		result = append(result, crd.GetIdsSelfWithChildById(orgId, model.Id)...)
	}
	return result
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

// UpdateById 通过id修改
func (crd *DepartmentCurd) UpdateById(param *department_models.DepartmentModel) error {
	conditions := department_models.SearchConditions{
		OrgId:    param.OrgId,
		FullCode: int_helper.Int64ToStr(param.Id),
	}
	models, err := crd.masterModel.GetListByConditions(conditions)
	if err != nil || models == nil || len(models) == 0 {
		return errors.New("菜单不存在")
	}
	err = crd.masterModel.Transaction(func(tx *gorm.DB) error {
		for _, model := range models {
			fullCodeArr := strings.Split(model.FullCode, SEPARATOR)
			fullNameArr := json_helper.ToSliceString(model.FullName)
			if len(fullCodeArr) != len(fullNameArr) {
				continue
			}
			for index, code := range fullCodeArr {
				if string_helper.ToInt64(code) != param.Id {
					continue
				}
				fullNameArr[index] = param.Name
				conditions := department_models.SearchConditions{
					OrgId: param.OrgId,
					Id:    model.Id,
				}
				updates := map[string]interface{}{
					"full_name":   json_helper.Encode(fullNameArr),
					"create_id":   param.CreateId,
					"create_name": param.CreateName,
				}
				if model.Id == param.Id {
					updates["name"] = param.Name
					updates["ip"] = param.Ip

				}

				if !crd.masterModel.UpdateByConditions(conditions, updates) {
					return errors.New("保存失败")
				}

			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// DelSelfWithChildById 删除本身以及子集
func (crd *DepartmentCurd) DelSelfWithChildById(orgId, id int64) bool {
	conditions := department_models.SearchConditions{
		OrgId:    orgId,
		FullCode: int_helper.Int64ToStr(id),
	}

	return crd.masterModel.DeleteByConditions(conditions)

}

// DeleteByOrg 通过企业编号删除
func (crd *DepartmentCurd) DeleteByOrg(orgId int64) bool {
	conditions := department_models.SearchConditions{
		OrgId: orgId,
	}

	return crd.masterModel.DeleteByConditions(conditions)
}
