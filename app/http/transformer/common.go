package transformer

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
	"strings"
)

type Common struct {
}

// FormatIntList 格式成标准前端数据格式：val是int
func (t *Common) FormatIntList(keys []int, labels map[int]string, appendAll bool) (result []map[string]interface{}) {
	if keys == nil {
		return result
	}
	if appendAll {
		result = append(result, map[string]interface{}{
			"val":   enums.AllDefaultInt,
			"label": "全部",
		})
	}
	for _, val := range keys {
		if label, isOk := labels[val]; isOk {
			result = append(result, map[string]interface{}{
				"val":   val,
				"label": label,
			})
		}

	}
	return result
}

// FormatStringList 格式成标准前端数据格式:val是string
func (t *Common) FormatStringList(keys []string, labels map[string]string, appendAll bool) (result []map[string]interface{}) {
	if keys == nil {
		return result
	}
	if appendAll {
		result = append(result, map[string]interface{}{
			"val":   enums.AllDefaultStr,
			"label": "全部",
		})
	}
	for _, val := range keys {
		if label, isOk := labels[val]; isOk {
			result = append(result, map[string]interface{}{
				"val":   val,
				"label": label,
			})
		}

	}
	return result
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

// TransThreatenInfo 转换威胁情报信息
// @result["type"]: 1空 2遍历换行, 3表格形式
// @result["attr"]: 根据type变化
func (t *Common) TransThreatenInfo(model *model_agent.ThreatenCacheModel) map[string]interface{} {
	result := make(map[string]interface{})
	if helpers.StringIsEmpty(model.CheckInfo) {
		result["type"] = 1
		result["attr"] = ""
		return result
	}

	attrJson := attributeJson
	err := helpers.Json2Struct(model.CheckInfo, &attrJson)
	if err != nil {
		result["type"] = 2
		result["attr"] = strings.Split(model.CheckInfo, "\n")
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
