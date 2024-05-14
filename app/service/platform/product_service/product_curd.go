package product_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/organization/prodcut_models"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *ProductCurd {
	return &ProductCurd{masterModel: prodcut_models.CreateProductFactory()}
}

type ProductCurd struct {
	masterModel *prodcut_models.ProductModel
}

type SearchConditions struct {
	Id        int64
	GrantType int
	Title     string
	IsPublish int
}

// search2Where 查询条件转化
func (crd *ProductCurd) search2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.Id > 0 {
		wheres["id"] = search.Id
	}
	if search.GrantType > 0 && search.GrantType != enums.AllDefaultInt {
		wheres["grant_type"] = search.GrantType
	}
	if !helpers.StringIsEmpty(search.Title) {
		wheres["title like"] = "%" + search.Title + "%"
	}
	if search.IsPublish > 0 && search.IsPublish != enums.AllDefaultInt {
		wheres["is_publish"] = search.IsPublish
	}
	return wheres
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
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// GetPublished 获取已发布的
func (crd *ProductCurd) GetPublished() (items []*prodcut_models.ProductModel) {
	search := SearchConditions{
		IsPublish: enums.IsYes,
	}
	items, err := crd.masterModel.GetItemsByConditions(crd.search2Where(search), 0, 0, "")
	if err != nil || len(items) == 0 {
		return nil
	}
	return items
}

// Page 分页
func (crd *ProductCurd) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*prodcut_models.ProductModel) {
	where := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}
