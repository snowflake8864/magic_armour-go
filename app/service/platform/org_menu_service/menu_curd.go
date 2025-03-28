package org_menu_service

import (
	"errors"
	"goskeleton/app/model/platform/org_menu_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/tree"
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
	param.CreatedAt = date_helper.CurDatetime()

	isOk := crd.masterModel.InsertItem(param)
	if !isOk {
		return -1, errors.New("添加失败")
	}
	return errCode, err
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

func (crd *MenuCurd) ToTree(models []*org_menu_models.OrgMenuModel, selectedMap map[int64]struct{}) []*tree.Node {

	var nodes []*tree.Node
	for _, model := range models {
		//是否选中
		isSelect := 0
		if _, isOk := selectedMap[model.Id]; isOk {
			isSelect = 1
		}

		nodes = append(nodes, &tree.Node{
			ID:       model.Id,
			ParentId: model.ParentId,
			Data: map[string]interface{}{
				"name":        model.Name,
				"type":        model.Type,
				"code":        model.Code,
				"sort":        model.Sort,
				"is_show":     model.IsShow,
				"is_root":     model.IsRoot,
				"is_selected": isSelect,
				"usage_ver":   model.UsageVer,
				"remark":      model.Remark,
				"api":         model.Api,
				"level":       model.Level,
				"created_at":  model.CreatedAt,
			},
			Child: []*tree.Node{},
		})
	}

	return tree.RealTree(nodes)
}
