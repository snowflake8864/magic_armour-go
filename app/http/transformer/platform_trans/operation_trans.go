package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_platform_operate"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/model_platform"
)

type Operation struct {
	transformer.Common
}

func (t *Operation) pageGeneral(model *model_platform.OperateLogModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["path"] = model.Path
	data["module"] = model.Module
	data["grant_type_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("module", model.Module, enum_platform_operate.EnumMap)
	data["net_ip"] = model.NetIp
	data["input"] = model.Input
	data["method"] = model.Method
	data["status"] = model.Status
	data["http_code"] = model.HttpCode
	data["http_code"] = model.HttpCode
	data["delay_time"] = model.DelayTime
	data["resp"] = model.Resp
	data["create_id"] = model.CreateId
	data["create_name"] = model.CreateName
	data["created_at"] = model.CreatedAt
	return data
}

// TransPage 转换成管理页需要的数据
func (t *Operation) TransPage(total int64, models []*model_platform.OperateLogModel) interface{} {
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
