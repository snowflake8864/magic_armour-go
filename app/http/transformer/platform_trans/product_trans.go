package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/platform/prodcut_models"
	"goskeleton/app/utils/helpers/json_helper"
)

type Product struct {
	transformer.Common
	Lang string
}

func (t *Product) pageGeneral(model *prodcut_models.ProductModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["title"] = model.Title
	data["grant_type"] = model.GrantType
	data["grant_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.GrantType, enum_product.EnumIntMap)
	productGrantConfig := enum_product.ProductGrantConfig{}
	_ = json_helper.ToStruct(model.GrantModule, &productGrantConfig)
	data["grant_module"] = productGrantConfig
	data["grant_module_txt"] = productGrantConfig.SelectedTxt()
	data["price"] = model.Price
	data["is_publish"] = model.IsPublish
	data["remark"] = model.Remark
	data["create_id"] = model.CreateId
	data["create_name"] = model.CreateName
	data["created_at"] = model.CreatedAt
	return data
}

// TransPage 转换成管理页需要的数据
func (t *Product) TransPage(total int64, models []*prodcut_models.ProductModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransList 转换列表页需要的数据
func (t *Product) TransList(models []*prodcut_models.ProductModel) (items []map[string]interface{}) {
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}
	return items
}
