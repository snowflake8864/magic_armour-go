package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Agent struct {
	transformer.Common
	Lang string
}

func (t *Agent) pageGeneral(model *agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, agentExtend map[string]interface{}) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	delete(data, "id")
	data["agent_id"] = model.AgentId
	data["ver"] = enum_software.ShowVer(model.PlatformType, model.OsType, model.Ver)

	data["dept_fullname"] = []interface{}{}
	if departments != nil && len(departments) > 0 {
		if deptModel, ok := departments[model.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
			data["dept_fullname"] = json_helper.ToSliceString(deptModel.FullName)
		}
	}

	productGrantConfig := enum_product.ProductGrantConfig{}
	_ = json_helper.ToStruct(model.GrantConfig, &productGrantConfig)
	data["grant_config"] = productGrantConfig

	if t.Lang == consts.LanguageZh {
		data["grant_config_txt"] = productGrantConfig.SelectedTxt()
		data["process_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.ProcessMode, enum_agent.EnumIntMap)
		data["dll_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.DllMode, enum_agent.EnumIntMap)
		data["tamper_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.TamperMode, enum_agent.EnumIntMap)
		data["blackmail_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.BlackmailMode, enum_agent.EnumIntMap)
		data["usb_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.UsbMode, enum_agent.EnumIntMap)
	} else {
		data["grant_config_txt"] = productGrantConfig.SelectedTxtEn()
		data["process_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode_en", model.ProcessMode, enum_agent.EnumIntMap)
		data["dll_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode_en", model.DllMode, enum_agent.EnumIntMap)
		data["tamper_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode_en", model.TamperMode, enum_agent.EnumIntMap)
		data["blackmail_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode_en", model.BlackmailMode, enum_agent.EnumIntMap)
		data["usb_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode_en", model.UsbMode, enum_agent.EnumIntMap)
	}

	//扩展信息
	var tmp interface{}
	data["extend_info"] = tmp
	if agentExtend != nil && len(agentExtend) > 0 {
		if extend, ok := agentExtend[model.AgentId]; ok {
			data["extend_info"] = extend
		}
	}

	return data
}

// TransManagePage 转换成管理分页需要的数据
func (t *Agent) TransManagePage(total int64, models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, softwareNewest map[int]string) interface{} {
	var items []map[string]interface{}
	if models != nil {
		//终端类型
		typeJsonStr := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentTypes)
		agentTypes := enum_agent.GetTypeMap(typeJsonStr)

		for _, model := range models {
			tmp := t.pageGeneral(model, departments, nil)

			canUpgrade := 0
			if enum_agent.CanUpgrade(model.GrantStatus) {
				if isOk, _ := enum_software.CanUpgradeIndirect(softwareNewest, model.OsType, model.Ver); isOk {
					canUpgrade = 1
				}
			}
			tmp["can_upgrade"] = canUpgrade
			tmp["agent_type_label"] = enum_agent.GetTypeLabel(agentTypes, model.AgentType)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransOtherPage 转换成其他分页需要的数据
func (t *Agent) TransOtherPage(total int64, models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, extendInfo map[string]interface{}) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments, extendInfo)

			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransList 转换列表需要的数据l
func (t *Agent) TransList(models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, extendInfo map[string]interface{}) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments, extendInfo)

			items = append(items, tmp)
		}
	}

	return items
}

// TransSingle 转换单个需要的数据
func (t *Agent) TransSingle(model *agent_models.AgentModel) interface{} {
	return t.pageGeneral(model, nil, nil)
}
