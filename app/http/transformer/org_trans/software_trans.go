package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/platform/software_models"
)

type Software struct {
	transformer.Common
}

func (t *Software) pageGeneral(model *software_models.SoftwareModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["ver"] = model.Ver
	data["version_number_label"] = enum_software.ShowVer(model.PlatformType, model.OsType, model.Ver)
	data["platform_type"] = model.PlatformType
	data["os_type"] = model.OsType
	data["os_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type", model.OsType, enum_software.EnumMap)
	data["os_type_alias"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type_alias", model.OsType, enum_software.EnumMap)
	data["os_ver"] = model.OsVer
	data["remark"] = model.Remark
	data["download_name"] = enum_software.DownloadName(model.PlatformType, model.OsType, model.Ver)

	data["created_at"] = model.CreatedAt
	return data
}

// TransPage 转换成管理页需要的数据
func (t *Software) TransPage(total int64, models []*software_models.SoftwareModel) interface{} {
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
func (t *Software) TransList(models []*software_models.SoftwareModel) (items []map[string]interface{}) {
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}
	return items
}
