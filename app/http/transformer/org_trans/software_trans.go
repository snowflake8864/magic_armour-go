package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type SoftwareTrans struct {
	transformer.Common
	Lang string
}

func (t *SoftwareTrans) pageGeneral(model *software_models.SoftwareModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	data["download_name"] = enum_software.DownloadName(model.PlatformType, model.OsType, model.Ver)
	data["version_number_label"] = enum_software.ShowVer(model.PlatformType, model.OsType, model.Ver)

	if t.Lang == consts.LanguageZh {
		data["os_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type", model.OsType, enum_software.EnumMap)
		data["os_type_alias"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type_alias", model.OsType, enum_software.EnumMap)
	} else {
		data["os_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type_en", model.OsType, enum_software.EnumMap)
		data["os_type_alias"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type_alias_en", model.OsType, enum_software.EnumMap)
	}

	return data
}

// TransPage 转换成管理页需要的数据
func (t *SoftwareTrans) TransPage(total int64, models []*software_models.SoftwareModel) interface{} {
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
func (t *SoftwareTrans) TransList(models []*software_models.SoftwareModel) (items []map[string]interface{}) {
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}
	return items
}
