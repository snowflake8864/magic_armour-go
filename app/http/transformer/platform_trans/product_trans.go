package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/platform/prodcut_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Product struct {
	transformer.Common
	Lang string
}

func (t *Product) pageGeneral(model *prodcut_models.ProductModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)

	productGrantConfig := enum_product.ProductGrantConfig{}
	_ = json_helper.ToStruct(model.GrantModule, &productGrantConfig)
	data["grant_module"] = productGrantConfig

	if t.Lang == consts.LanguageZh {
		data["grant_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.GrantType, enum_product.EnumIntMap)
		data["grant_module_txt"] = productGrantConfig.SelectedTxt()
	} else {
		data["grant_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.GrantType, enum_product.EnumIntMap)
		data["grant_module_txt"] = productGrantConfig.SelectedTxtEn()
	}

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
