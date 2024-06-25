package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/threaten_models"
)

type Threaten struct {
	transformer.Common
}

func (t *Threaten) pageCacheGeneral(model *threaten_models.CacheModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["level"] = model.Level
	data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_threaten.EnumMap)
	data["source"] = model.Source
	data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("source", model.Source, enum_threaten.EnumMap)
	data["dir"] = model.Dir
	data["parent_dir"] = model.ParentDir
	data["os"] = model.Os
	data["introduce"] = model.Introduce
	data["copyright"] = model.Copyright
	data["signature"] = model.Signature
	data["attribute"] = model.Attribute
	data["data_from"] = model.DataFrom
	data["check_at"] = model.CheckAt
	data["created_at"] = model.CreatedAt
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
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["level"] = model.Level
	data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_threaten.EnumMap)

	data["dir"] = model.Dir
	data["parent_dir"] = model.ParentDir
	data["os"] = model.Os
	data["data_from"] = model.DataFrom
	data["push_at"] = model.PushAt
	data["push_status"] = model.PushStatus
	data["introduce"] = model.Introduce
	data["copyright"] = model.Copyright
	data["signature"] = model.Signature
	data["attribute"] = model.Attribute
	data["push_status"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("push_status", model.PushStatus, enum_threaten.EnumMap)
	data["create_id"] = model.CreateId
	data["create_name"] = model.CreateName
	data["created_at"] = model.CreatedAt
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
