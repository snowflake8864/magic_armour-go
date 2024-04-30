package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/model_agent"
)

type Threaten struct {
	transformer.Common
}

func (t *Threaten) pageCacheGeneral(model *model_agent.ThreatenCacheModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["level"] = model.Level
	data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Level, enum_threaten.EnumMap)
	data["source"] = model.Source
	data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("source", model.Source, enum_threaten.EnumMap)
	data["dir"] = model.Dir
	data["parent_dir"] = model.ParentDir
	data["os"] = model.Os
	data["data_from"] = model.DataFrom
	data["check_at"] = model.CheckAt
	data["created_at"] = model.CreatedAt
	return data
}

// TransCachePage 转换成管理页需要的数据
func (t *Threaten) TransCachePage(total int64, models []*model_agent.ThreatenCacheModel) interface{} {
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
