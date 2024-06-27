package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/response"
	"strings"
)

// SoftwareController  软甲相关
type SoftwareController struct {
	ControllerBase
}

// PlatformTypes  平台类型
func (u *SoftwareController) PlatformTypes(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfInt("platform_type", enum_software.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("platform_type", enum_software.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// PublishedList 已发布列表
func (u *SoftwareController) PublishedList(ctx *gin.Context) {
	IsPublish := enums.IsYes
	condition := software_models.SearchConditions{
		IsPublish: &IsPublish,
	}

	models := software_models.CreateSoftwareFactory().GetItemsByConditions(condition, 0, 0, "")
	response.SuccessReturnData(ctx, (&org_trans.Software{}).TransList(models))
}

func (u *SoftwareController) OfflineDownload(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "softId"))
	model := software_service.CreateCurdFactory().GetItemById(softId)
	if model == nil || model.Id == 0 {
		response.FailSimplify(ctx, "软件不存在")
		return
	}

	dictIds := []string{
		enum_dictionary.AgentHttpUrl,
		enum_dictionary.AgentHttpPort,
	}
	mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	addr := strings.Trim(mapDict[enum_dictionary.AgentHttpUrl], "/") + ":" + mapDict[enum_dictionary.AgentHttpPort]

	verId := ""
	url := ""
	if model.PlatformType == enum_software.PlatformTypeLinux {
		url += "sudo wget --no-check-certificate -SO /opt/veda_armour.zip %s/v5/api/software/download?verId=%s"
	} else {
		url += "%s/v5/api/software/download?verId=%s"
	}

	response.SuccessReturnData(ctx, gin.H{
		"url": fmt.Sprintf(url, addr, verId),
	})
}

func (u *SoftwareController) OnlineDownload(ctx *gin.Context) {

}
