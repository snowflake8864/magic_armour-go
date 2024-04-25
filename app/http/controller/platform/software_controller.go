package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/response"
)

// SoftwareController 终端相关
type SoftwareController struct {
	ControllerBase
}

// GetOsTypes  获取系统类型
func (u *SoftwareController) GetOsTypes(ctx *gin.Context) {
	platformType := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "platformType"))

	osTypes := enums.Get1stIsStrAnd2ndIsIntReturnSliceInt("map_platform_os", platformType, enum_software.EnumSlice)
	if osTypes == nil {
		response.FailSimplify(ctx, "platformType no exist")
		return
	}

	var result []map[string]interface{}
	for _, osType := range osTypes {

		//先获取系统类型保护的版本
		var children []map[string]interface{}
		typeVers := enums.Get1stIsStrAnd2ndIsIntReturnSliceInt("map_os_ver", osType, enum_software.EnumSlice)
		for _, typeVer := range typeVers {
			children = append(children, map[string]interface{}{
				"val":   typeVer,
				"label": enums.Get1stIsStrAnd2ndIsIntReturnStr("os_ver", typeVer, enum_software.EnumMap),
			})
		}

		//系统类型
		result = append(result, map[string]interface{}{
			"val":      osType,
			"label":    enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type", osType, enum_software.EnumMap),
			"children": children,
		})
	}

	response.SuccessReturnData(ctx, result)
}

// GetAgentConfig  获取终端配置
func (u *SoftwareController) GetAgentConfig(ctx *gin.Context) {
	dictIds := []string{
		enum_dictionary.AgentHttpUrl,
		enum_dictionary.AgentHttpPort,
		enum_dictionary.AgentCompatibleDefaultOrg,
		enum_dictionary.AgentCompatibleDefaultUser,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	response.SuccessReturnData(ctx, gin.H{
		"http_url":     vals[enum_dictionary.AgentHttpUrl],
		"http_port":    vals[enum_dictionary.AgentHttpPort],
		"default_org":  vals[enum_dictionary.AgentCompatibleDefaultOrg],
		"default_user": vals[enum_dictionary.AgentCompatibleDefaultUser],
	})
}

// SetAgentConfig  保存终端配置
func (u *SoftwareController) SetAgentConfig(ctx *gin.Context) {

}

// PageForManage 管理: 分页
func (u *SoftwareController) PageForManage(ctx *gin.Context) {

}

// AddForManage 管理:  添加
func (u *SoftwareController) AddForManage(ctx *gin.Context) {

}

// ClearForManage 管理: 清空
func (u *SoftwareController) ClearForManage(ctx *gin.Context) {

}

// DownloadForManage 管理: 下载
func (u *SoftwareController) DownloadForManage(ctx *gin.Context) {

}

// PublishForManage 管理: 发布
func (u *SoftwareController) PublishForManage(ctx *gin.Context) {

}

// CancelForManage 管理: 取消
func (u *SoftwareController) CancelForManage(ctx *gin.Context) {

}

// DeleteForManage 管理: 删除
func (u *SoftwareController) DeleteForManage(ctx *gin.Context) {

}
