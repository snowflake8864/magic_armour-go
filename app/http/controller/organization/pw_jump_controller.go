package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_pw_jump"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/response"
)

// PwJumpController 口令跳变相关
type PwJumpController struct {
	ControllerBase
}

// Combines  口令组合
func (u *PwJumpController) Combines(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfInt("combine", enum_pw_jump.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("combine", enum_pw_jump.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// Cycles  口令周期
func (u *PwJumpController) Cycles(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfInt("cycle", enum_pw_jump.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("cycle", enum_pw_jump.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PageForManage 分页
func (u *PwJumpController) PageForManage(ctx *gin.Context) {

	//查询条件
	conditions, err := u.getAgentSearchConditions(ctx)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//分页参数和排序参数
	offset, pageSize, _ := u.getPageAndSort(ctx, "")

	//终端分页
	total, agentModels := agent_models.CreateAgentFactory().Page(conditions, offset, pageSize, enum_agent.AgentFieldDefault)
	//转化数据
	_, departmentMaps := u.parseAgentModels(agentModels, true)

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, nil))
}

// SetSwitchForManage 管理页: 设置开关
func (u *PwJumpController) SetSwitchForManage(ctx *gin.Context) {
	agentId := u.getAgentIdFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		PwJumpSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().UpdateByAgent(u.getLogin(ctx).OrgId, agentId, &param, "口令跳变管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *PwJumpController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))

	var param = agent_service.ManualColumn{
		PwJumpSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "口令跳变管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ShowPwForManage 管理页: 显示密码
func (u *PwJumpController) ShowPwForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	response.SuccessReturnData(ctx, gin.H{
		"pw": agent_service.CreateExtendFactory().GetPwJumpNewest(loginInfo.OrgId, agentId),
	})
}

// GetForAgentConfig 终端策略: 获取
func (u *PwJumpController) GetForAgentConfig(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := agent_service.CreateExtendFactory().GetPwJumpConfig(loginInfo.OrgId, agentId)

	response.SuccessReturnData(ctx, config)
}

func (u *PwJumpController) getConfig(ctx *gin.Context) enum_pw_jump.PwJumpConfig {
	config := enum_pw_jump.PwJumpConfig{
		Switch:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Cycle:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "cycle")),
		Combine: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "combine")),
		Size:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "size")),
	}
	config.Format()
	return config
}

// SaveForAgentConfig 终端策略: 保存
func (u *PwJumpController) SaveForAgentConfig(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := u.getConfig(ctx)

	if err := agent_service.CreateExtendFactory().SetPwJumpConfig(loginInfo.OrgId, agentId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ExecuteForAgentConfig 终端策略: 执行
func (u *PwJumpController) ExecuteForAgentConfig(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := u.getConfig(ctx)

	if err := agent_service.CreateExtendFactory().SetPwJumpConfig(loginInfo.OrgId, agentId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	// 下发任务
	task_service.CreateCurdFactory().AddTaskForAgent(loginInfo.OrgId, agentId, enum_task.TypePwJump, "终端策略")

	response.SuccessSimplify(ctx)
}

// GetForGlobalConfig 全局策略: 获取
func (u *PwJumpController) GetForGlobalConfig(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	model := OrgExtendFact.GetItemByOrg(loginInfo.OrgId)
	response.SuccessReturnData(ctx, OrgExtendFact.FormatPwJumpConfig(model.PwJumpConfig))
}

// SaveForGlobalConfig 全局策略: 保存
func (u *PwJumpController) SaveForGlobalConfig(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	config := u.getConfig(ctx)

	if err := organize_service.CreateOrgExtendFactory().SetPwJumpConfig(loginInfo.OrgId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ExecuteForGlobalConfig 全局策略: 执行
func (u *PwJumpController) ExecuteForGlobalConfig(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	config := u.getConfig(ctx)

	if err := organize_service.CreateOrgExtendFactory().SetPwJumpConfig(loginInfo.OrgId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	agentFact := agent_models.CreateAgentFactory()
	pwJumpSwitch := enums.IsYes
	conditions := agent_models.SearchConditions{
		OrgId:        loginInfo.OrgId,
		PwJumpSwitch: &pwJumpSwitch,
	}

	// 下发任务
	if agentIds := agentFact.GetIdsByByConditions(conditions); len(agentIds) > 0 {
		task_service.CreateCurdFactory().AddTaskForAgents(loginInfo.OrgId, agentIds, enum_task.TypePwJump, "全局策略")
	}

	response.SuccessSimplify(ctx)
}
