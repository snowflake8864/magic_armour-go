package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/response"
)

// LoginSolveController 口令破解相关
type LoginSolveController struct {
	ControllerBase
}

// PageForManage 分页
func (u *LoginSolveController) PageForManage(ctx *gin.Context) {

	//查询条件
	conditions, err := u.getAgentSearchConditions(ctx)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, enum_agent.AgentFieldDefault)

	//终端分页
	total, agentModels := agent_models.CreateAgentFactory().Page(conditions, offset, pageSize, order)
	//转化数据
	_, departmentMaps := u.parseAgentModels(agentModels, true)

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, nil))
}

// SetSwitchForManage 管理页: 设置开关
func (u *LoginSolveController) SetSwitchForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		LoginSolveSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "口令破解管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *LoginSolveController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))

	var param = agent_service.ManualColumn{
		LoginSolveSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "口令破解管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigGetForPolicy 终端策略
func (u *LoginSolveController) ConfigGetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := agent_service.CreateExtendFactory().GetLoginSolve(loginInfo.OrgId, agentId)

	response.SuccessReturnData(ctx, config)
}

// ConfigSetForPolicy 终端策略
func (u *LoginSolveController) ConfigSetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := enum_login_solve.LoginSolveConfig{
		Switch:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Gap:       int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "gap")),
		Threshold: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "threshold")),
	}
	if err := agent_service.CreateExtendFactory().SetLoginSolve(loginInfo.OrgId, agentId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigGetForGlobal 全局策略
func (u *LoginSolveController) ConfigGetForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	model := OrgExtendFact.GetItemByOrg(loginInfo.OrgId)
	response.SuccessReturnData(ctx, OrgExtendFact.FormatLoginSolveConfig(model.LoginSolveConfig))
}

// ConfigSetForGlobal 全局策略
func (u *LoginSolveController) ConfigSetForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	config := enum_login_solve.LoginSolveConfig{
		Switch:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Gap:       int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "gap")),
		Threshold: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "threshold")),
	}
	if err := organize_service.CreateOrgExtendFactory().SetLoginSolve(loginInfo.OrgId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//用户管理的终端
	if err, agentIds := u.getCurrUserAgentIds(loginInfo); err == nil {
		_ = agent_service.CreateExtendFactory().BatchSetLoginSolve(loginInfo.OrgId, agentIds, config)
	}

	response.SuccessSimplify(ctx)
}
