package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_user"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type User struct {
	transformer.Common
	Lang string
}

func (t *User) pageGeneral(model *user_models.UserModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	if t.Lang == consts.LanguageZh {
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_user.EnumMap)
	} else {
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_user.EnumMap)
	}
	return data
}

// TransPage 转换成管理页需要的数据
func (t *User) TransPage(total int64, models []*user_models.UserModel) interface{} {
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
