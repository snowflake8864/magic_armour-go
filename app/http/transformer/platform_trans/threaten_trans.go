package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
	"strings"
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

var attributeJson struct {
	Hash       string `json:"hash"`
	FileName   string `json:"fileName"`
	Undetected int    `json:"undetected"` //正常数量
	Malicious  int    `json:"malicious"`  //异常数量
	Detections []struct {
		Category   string `json:"category"`
		EngineName string `json:"engineName"`
	} `json:"detections"`
}

// TransCacheAttr 转换缓存属性
// @result["type"]: 1空 2遍历换行, 3表格形式
// @result["attr"]: 根据type变化
func (t *Threaten) TransCacheAttr(model *model_agent.ThreatenCacheModel) map[string]interface{} {
	result := make(map[string]interface{})
	if helpers.StringIsEmpty(model.Attribute) {
		result["type"] = 1
		result["attr"] = ""
		return result
	}

	attrJson := attributeJson
	err := helpers.Json2Struct(model.Attribute, &attrJson)
	if err != nil {
		result["type"] = 2
		result["attr"] = strings.Split(model.Attribute, "\n")
		return result
	} else {
		result["type"] = 3
		attr := make(map[string]interface{})

		var undetectedList []map[string]interface{}
		var maliciousList []map[string]interface{}
		for _, item := range attrJson.Detections {
			if strings.ToLower(item.Category) == "malicious" {
				maliciousList = append(maliciousList, map[string]interface{}{
					"status":     0,
					"category":   item.Category,
					"engineName": item.EngineName,
				})
			} else {
				undetectedList = append(undetectedList, map[string]interface{}{
					"status":     1,
					"category":   item.Category,
					"engineName": item.EngineName,
				})
			}
		}

		attr["undetected"] = len(undetectedList)
		attr["malicious"] = len(maliciousList)
		attr["detections"] = append(maliciousList, undetectedList...)
		result["attr"] = attr
		return result
	}
}
