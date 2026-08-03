package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/platform/antivirus_engine_models"
	"goskeleton/app/utils/helpers/file_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type AntivirusEngine struct {
	transformer.Common
}

func (t *AntivirusEngine) pageGeneral(model *antivirus_engine_models.EngineModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	data["platform_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("platform_type", model.PlatformType, enum_software.EnumMap)

	data["file_size_label"] = file_helper.FormatFileSize(model.FileSize)

	if _, isOk := data["full_path"]; isOk {
		delete(data, "full_path")
	}

	return data
}

// TransPage 转换成管理页需要的数据
func (t *AntivirusEngine) TransPage(total int64, models []*antivirus_engine_models.EngineModel) interface{} {
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
func (t *AntivirusEngine) TransList(models []*antivirus_engine_models.EngineModel) (items []map[string]interface{}) {
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}
	return items
}
