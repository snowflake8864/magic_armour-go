package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_user"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/model_organization"
)

type User struct {
	transformer.Common
}

func (t *User) pageGeneral(model *model_organization.UserModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["username"] = model.Username
	data["email"] = model.Email
	data["phone"] = model.Phone
	data["realname"] = model.Realname
	data["position"] = model.Position
	data["intro"] = model.Intro
	data["last_login_at"] = model.LastLoginAt
	data["status"] = model.Status
	data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_user.EnumMap)
	data["is_root"] = model.IsRoot
	data["created_at"] = model.CreatedAt
	return data
}

// TransPage 转换成管理页需要的数据
func (t *User) TransPage(total int64, models []*model_organization.UserModel) interface{} {
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
