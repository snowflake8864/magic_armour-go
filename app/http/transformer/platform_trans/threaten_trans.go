package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Threaten struct {
	transformer.Common
	Lang string
}

func (t *Threaten) pageCacheGeneral(model *threaten_models.CacheModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	if t.Lang == consts.LanguageZh {
		data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_threaten.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("source", model.Source, enum_threaten.EnumMap)
	} else {
		data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_threaten.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("source_en", model.Source, enum_threaten.EnumMap)
	}
	return data
}

// TransCachePage 转换成管理页需要的数据
func (t *Threaten) TransCachePage(total int64, models []*threaten_models.CacheModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageCacheGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

func (t *Threaten) pagePushGeneral(model *threaten_models.PushModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	if t.Lang == consts.LanguageZh {
		data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_threaten.EnumMap)
		data["push_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("push_status", model.PushStatus, enum_threaten.EnumMap)
	} else {
		data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_threaten.EnumMap)
		data["push_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("push_status_en", model.PushStatus, enum_threaten.EnumMap)
	}
	return data
}

// TransPushPage 转换成推送页需要的数据
func (t *Threaten) TransPushPage(total int64, models []*threaten_models.PushModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pagePushGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
