package product_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/platform/prodcut_models"
)

func CreateCurdFactory() *ProductCurd {
	return &ProductCurd{masterModel: prodcut_models.CreateProductFactory()}
}

type ProductCurd struct {
	masterModel *prodcut_models.ProductModel
}

// Add 新增
func (crd *ProductCurd) Add(param *prodcut_models.ProductModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}

// UpdateById 新增
func (crd *ProductCurd) UpdateById(param *prodcut_models.ProductModel) (errCode int, err error) {
	updates := make(map[string]interface{})
	updates["title"] = param.Title
	updates["grant_type"] = param.GrantType
	updates["grant_module"] = param.GrantModule
	updates["price"] = param.Price
	updates["remark"] = param.Remark
	updates["create_id"] = param.CreateId
	updates["create_name"] = param.CreateName
	isOk := crd.masterModel.UpdateById(param.Id, updates)
	if !isOk {
		return -1, errors.New("修改失败")
	}
	return errCode, err
}

// DeleteById 删除
func (crd *ProductCurd) DeleteById(id int64) bool {
	return crd.masterModel.DeleteById(id)
}

// GetItemById 通过编号称获取
func (crd *ProductCurd) GetItemById(id int64) (model *prodcut_models.ProductModel) {
	model = crd.masterModel.GetItemById(id)
	return model
}

// GetPublished 获取已发布的
func (crd *ProductCurd) GetPublished() (items []*prodcut_models.ProductModel) {
	IsPublish := enums.IsYes
	search := prodcut_models.SearchConditions{
		IsPublish: &IsPublish,
	}
	items = crd.masterModel.GetItemsByConditions(search, 0, 0, "")
	if items == nil || len(items) == 0 {
		return nil
	}
	return items
}
