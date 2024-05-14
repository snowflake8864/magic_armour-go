package org_menu_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/organization/org_menu_models"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *MenuCurd {
	return &MenuCurd{masterModel: org_menu_models.CreateOrgMenuFactory()}
}

type MenuCurd struct {
	masterModel *org_menu_models.OrgMenuModel
}

// Add 新增
func (crd *MenuCurd) Add(param *org_menu_models.OrgMenuModel) (errCode int, err error) {

	if param.ParentId > 0 {
		if hasModel, _ := crd.masterModel.GetItemById(param.ParentId); hasModel != nil && hasModel.Id > 0 {
			param.Level = hasModel.Level + 1
		} else {
			return -1, errors.New("父级菜单不存在")
		}
	} else {
		//等级从1开始
		param.Level = 1
	}
	param.CreatedAt = helpers.DatetimeCurr()

	isOk := crd.masterModel.InsertItem(param)
	if !isOk {
		return -1, errors.New("添加失败")
	}
	return errCode, err
}

// GetAll 获取全部
func (crd *MenuCurd) GetAll(isRoot int) []*org_menu_models.OrgMenuModel {
	condition := make(map[string]interface{})
	if isRoot == 0 {
		condition["is_root"] = enums.IsYes
	}

	list, _ := crd.masterModel.GetItemsByConditions(condition)
	return list
}

// GetItemsByIds 通过id获取
func (crd *MenuCurd) GetItemsByIds(ids []int64, isRoot int) []*org_menu_models.OrgMenuModel {
	condition := make(map[string]interface{})
	condition["id"] = ids
	if isRoot == 0 {
		condition["is_root"] = enums.IsYes
	}

	list, _ := crd.masterModel.GetItemsByConditions(condition)
	return list
}

// UpdateById 修改
func (crd *MenuCurd) UpdateById(param *org_menu_models.OrgMenuModel) (errCode int, err error) {
	updates := make(map[string]interface{})
	updates["type"] = param.Type
	updates["name"] = param.Name
	updates["code"] = param.Code
	updates["sort"] = param.Sort
	updates["is_show"] = param.IsShow
	updates["is_root"] = param.IsRoot
	updates["usage_ver"] = param.UsageVer
	updates["remark"] = param.Remark

	isOk := crd.masterModel.UpdateById(param.Id, updates)
	if !isOk {
		return -1, errors.New("修改失败")
	}
	return errCode, err
}

// DeleteById 删除
func (crd *MenuCurd) DeleteById(id int64) (errCode int, err error) {
	if hasModel, _ := crd.masterModel.GetItemByParent(id); hasModel != nil && hasModel.Id > 0 {
		return -1, errors.New("有子菜单或按钮，不能被删除")
	}
	isOk := crd.masterModel.DeleteById(id)
	if !isOk {
		return -1, errors.New("删除失败")
	}
	return errCode, err
}
