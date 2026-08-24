package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/dir_models"
	"goskeleton/app/service/agent/dir_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// FileDirController 文件目录相关
type FileDirController struct {
	ControllerBase
}

// GetTypes 类型
func (u *FileDirController) GetTypes(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("search_type", enum_tamper.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("search_type", enum_tamper.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, true))
}

// Page 分页
func (u *FileDirController) Page(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)
	dir := dir_service.CreateDirCurdFactory().DealSearchDir(ctx.GetString(consts.ValidatorPrefixOrg + "dir"))
	//查询条件
	condition := dir_models.SearchConditions{
		OrgId:        loginInfo.OrgId,
		AgentId:      agentId,
		SearchDir:    dir,
		FileName:     ctx.GetString(consts.ValidatorPrefixOrg + "fileName"),
		FileType:     int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "fileType")),
		StartTimeGte: ctx.GetString(consts.ValidatorPrefixOrg + "startTime"),
		StartTimeLte: ctx.GetString(consts.ValidatorPrefixOrg + "endTime"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := dir_models.CreateDirFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.FileDirTrans{Lang: cur_userinfo.GetLocales(ctx)}).TranslPage(total, models))
}

// HasUpload 判断是否上传
func (u *FileDirController) HasUpload(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	exist := enums.IsNot
	if val := dir_service.CreateDirCurdFactory().GetBySearch(loginInfo.OrgId, agentId); string_helper.IsEmpty(val) {
		exist = enums.IsYes
	}

	response.SuccessReturnData(ctx, gin.H{
		"exist": exist,
	})
}

// Refresh 刷新
func (u *FileDirController) Refresh(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)
	dir := ctx.GetString(consts.ValidatorPrefixOrg + "dir")
	force := false
	if val := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "force")); val == 1 {
		force = true
	}

	if err := dir_service.CreateDirCurdFactory().SaveBySearch(loginInfo.OrgId, agentId, dir, force); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
