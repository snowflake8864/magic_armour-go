package transformer

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
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
func (t *Common) TransThreatenInfo(model *threaten_models.CacheModel) map[string]interface{} {
	result := make(map[string]interface{})
	if string_helper.IsEmpty(model.CheckInfo) {
		result["type"] = 1
		result["attr"] = ""
		return result
	}

	attrJson := attributeJson
	err := json_helper.ToStruct(model.CheckInfo, &attrJson)
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
					"status":      0,
					"category":    item.Category,
					"engine_name": item.EngineName,
				})
			} else {
				undetectedList = append(undetectedList, map[string]interface{}{
					"status":      1,
					"category":    item.Category,
					"engine_name": item.EngineName,
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

func (t *Common) GetDeptFullname(id int64, departments map[int64]*department_models.DepartmentModel) interface{} {
	var result []interface{}
	if id == 0 {
		result = append(result, "全部")
		return result
	}
	if departments == nil && len(departments) == 0 {
		result = append(result, "组织架构为空")
		return result
	}
	if deptModel, ok := departments[id]; ok && !string_helper.IsEmpty(deptModel.FullName) {
		return json_helper.ToSliceString(deptModel.FullName)
	} else {
		result = append(result, "组织架构已删除")
		return result
	}
}

// AppendAgentInfo 追加终端信息
func (t *Common) AppendAgentInfo(agentId string, agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) map[string]interface{} {
	var data = make(map[string]interface{})
	//终端信息
	data["host_name"] = ""
	data["mac_id"] = ""
	data["agent_ip"] = ""
	data["net_ip"] = ""
	data["os_name"] = ""
	data["contact_name"] = ""
	data["contact_tel"] = ""
	data["alias"] = ""
	data["is_online"] = 0
	data["dept_id"] = 0
	data["dept_fullname"] = []interface{}{}
	if agentModel, isOK := agentModels[agentId]; isOK {
		data["host_name"] = agentModel.HostName
		data["mac_id"] = agentModel.MacId
		data["agent_ip"] = agentModel.AgentIp
		data["net_ip"] = agentModel.NetIp
		data["os_name"] = agentModel.OsName
		data["contact_name"] = agentModel.ContactName
		data["contact_tel"] = agentModel.ContactTel
		data["alias"] = agentModel.Alias
		data["is_online"] = agentModel.IsOnline

		if departments != nil && len(departments) > 0 {
			if deptModel, ok := departments[agentModel.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
				data["dept_id"] = deptModel.Id
				data["dept_fullname"] = json_helper.ToSliceString(deptModel.FullName)
			}
		}
	}
	return data
}

// AppendThreatenInfo 追加威胁情报信息
func (t *Common) AppendThreatenInfo(hash string, threatenMaps map[string]*threaten_models.CacheModel) map[string]interface{} {
	var data = make(map[string]interface{})
	data["threaten_level"] = enum_threaten.LevelDefault
	data["threaten_level_label"] = ""
	data["threaten_is_sync"] = enums.IsNot

	if model, isOK := threatenMaps[hash]; isOK {
		data["threaten_level"] = model.Level
		data["threaten_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_threaten.EnumMap)
		if model.CheckSource == enum_threaten.CheckSourceInternet {
			data["threaten_is_sync"] = enums.IsYes
		}
	}
	return data
}
