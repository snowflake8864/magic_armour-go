package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/utils/response"
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

}

func (u *SoftwareController) OnlineDownload(ctx *gin.Context) {

}
