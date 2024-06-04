package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/utils/response"
)

// ProcSampleController 进程或模块样本相关
type ProcSampleController struct {
	ControllerBase
}

// Page 分页
func (u *ProcSampleController) Page(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}

// Delete 删除
func (u *ProcSampleController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Delete")
}

// Clear 清空
func (u *ProcSampleController) Clear(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Clear")
}

// Download 下载
func (u *ProcSampleController) Download(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Download")
}

// AddFromProcess 添加来源进程
func (u *ProcSampleController) AddFromProcess(ctx *gin.Context) {
	if err := u.add(ctx, enum_process.SourceProcess); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)

}

// AddFromDll 添加来源模块
func (u *ProcSampleController) AddFromDll(ctx *gin.Context) {
	if err := u.add(ctx, enum_process.SourceDll); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddFromDll 添加来源模块
func (u *ProcSampleController) add(ctx *gin.Context, dataFrom int) error {
	loginInfo := u.getLogin(ctx)

	param := proc_models.ProcSampleModel{
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Dir:        ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Introduce:  ctx.GetString(consts.ValidatorPrefixOrg + "introduce"),
		Copyright:  ctx.GetString(consts.ValidatorPrefixOrg + "copyright"),
		Signature:  ctx.GetString(consts.ValidatorPrefixOrg + "signature"),
		Attribute:  ctx.GetString(consts.ValidatorPrefixOrg + "signature"),
		DataFrom:   dataFrom,
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	return proc_service.CreateSampleCurdFactory().Add(&param)
}
