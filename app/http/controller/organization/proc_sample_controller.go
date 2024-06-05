package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
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
	orgId := u.getLogin(ctx).OrgId

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId: orgId,
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateProcSampleFactory().Page(condition, offset, pageSize, order)

	var agentIds []string
	var agentMap = make(map[string]*agent_models.AgentModel)
	if models != nil || len(models) > 0 {
		for _, model := range models {
			agentIds = append(agentIds, model.AgentId)
		}
		AgentFact := agent_models.CreateAgentFactory()
		agentMap = AgentFact.GetMapsByOrgAndAgents(orgId, agentIds, AgentFact.ColumnBasic(true, false))
	}

	response.SuccessReturnData(ctx, (&org_trans.ProcSample{}).TransPage(total, models, agentMap))
}

// Delete 删除
func (u *ProcSampleController) Delete(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	Id := ctx.GetString(consts.ValidatorPrefixOrg + "id")
	agentId := u.getAgentIdFromGin(ctx)
	if !proc_service.CreateSampleCurdFactory().DeleteById(loginInfo.OrgId, agentId, Id) {
		response.FailSimplify(ctx, "删除失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// Clear 清空
func (u *ProcSampleController) Clear(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	if !proc_service.CreateSampleCurdFactory().DelByOrg(loginInfo.OrgId) {
		response.FailSimplify(ctx, "删除失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// Download 下载
func (u *ProcSampleController) Download(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	Id := ctx.GetString(consts.ValidatorPrefixOrg + "id")
	agentId := u.getAgentIdFromGin(ctx)

	model := proc_models.CreateProcSampleFactory().GetItemByOrgAndAgentAndHash(loginInfo.OrgId, agentId, Id)
	if model == nil || model.OrgId == 0 {
		response.FailSimplify(ctx, "数据不存在")
		return
	}
	if model.Status != enum_process.SampleStatusFinished {
		response.FailSimplify(ctx, "样本没有上传")
		return
	}

	response.FileDownload(ctx, model.FullPath)
}

// AddFromProcess 添加来源进程
func (u *ProcSampleController) AddFromProcess(ctx *gin.Context) {
	if err := u.add(ctx, enum_process.DataFromProcess); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)

}

// AddFromDll 添加来源模块
func (u *ProcSampleController) AddFromDll(ctx *gin.Context) {
	if err := u.add(ctx, enum_process.DataFromDll); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// add
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
		Attribute:  ctx.GetString(consts.ValidatorPrefixOrg + "attribute"),
		DataFrom:   dataFrom,
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	return proc_service.CreateSampleCurdFactory().Add(&param)
}
