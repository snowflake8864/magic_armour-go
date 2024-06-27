package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_org_operate"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/org_operate_models"
)

type Operation struct {
	transformer.Common
}

func (t *Operation) pageGeneral(model *org_operate_models.LogModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["agent_id"] = model.AgentId
	data["path"] = model.Path
	data["module"] = model.Module
	data["module_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("module", model.Module, enum_org_operate.EnumMap)
	data["status"] = model.Status

	data["net_ip"] = model.NetIp
	data["created_at"] = model.CreatedAt

	return data
}

// TransPage 转换成管理页需要的数据
func (t *Operation) TransPage(total int64, models []*org_operate_models.LogModel) interface{} {
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
