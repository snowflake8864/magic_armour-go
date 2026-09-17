package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_third_user_auth"
	"goskeleton/app/service/platform/third_user_auth_service"
	"goskeleton/app/utils/response"
)

// KoalController 格尔身份认证系统相关
type KoalController struct {
}

// SubNotice 订阅通知
func (u *KoalController) SubNotice(ctx *gin.Context) {
	resType := ctx.GetString(consts.ValidatorPrefixApi + "resType")

	var noticeSync int
	if resType == "user" {
		noticeSync = enum_third_user_auth.NoticeSyncUser
	} else {
		noticeSync = enum_third_user_auth.NoticeSyncOrg
	}

	if err := third_user_auth_service.CreateConfigFactory().ModifyNoticeById(noticeSync); err != nil {
		response.Returnkoal(ctx, 500, err.Error(), gin.H{
			"noticeSync": noticeSync,
		})
		return
	}

	response.Returnkoal(ctx, 200, "OK", gin.H{
		"noticeSync": noticeSync,
	})
}
