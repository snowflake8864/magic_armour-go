package munu_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *MenuCurd {
	return &MenuCurd{masterModel: model_platform.CreateMenuFactory()}
}

type MenuCurd struct {
	masterModel *model_platform.MenuModel
}

// Add 新增
func (crd *MenuCurd) Add(param *model_platform.MenuModel) (errCode int, err error) {

	if param.ParentId > 0 {
		if hasModel, _ := crd.masterModel.GetItemById(param.ParentId); hasModel != nil {
			param.Level = hasModel.Level + 1
		} else {
			return -1, errors.New("父级菜单不存在")
		}
	} else {
		//等级从1开始
		param.Level = 1
	}
	param.CreatedAt = helpers.GetDatetime()

	isOk := crd.masterModel.InsertItem(param)
	if !isOk {
		return -1, errors.New("添加失败")
	}
	return errCode, err
}

func (crd *MenuCurd) GetAll(isRoot int) []model_platform.MenuModel {
	condition := make(map[string]interface{})
	if isRoot > 0 {
		condition["is_root"] = enums.IsYes
	}

	list, _ := crd.masterModel.GetItemsByConditions(condition)
	return list
}
