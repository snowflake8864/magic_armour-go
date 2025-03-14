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
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/outreach_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/string_helper"
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
	agentId := u.getAgentIdFromGin(ctx)

	response.SuccessReturnData(ctx, gin.H{
		"pw": agentId,
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

// AddForGlobal 全局策略: 添加
func (u *PwJumpController) AddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := outreach_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err := outreach_service.CreateGlobalFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGlobal 全局策略: 编辑
func (u *PwJumpController) EditForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := outreach_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := outreach_service.CreateGlobalFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
