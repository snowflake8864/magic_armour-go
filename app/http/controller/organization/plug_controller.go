package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/port_service"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/response"
)

// PlugController 端口相关
type PlugController struct {
	ControllerBase
}

// PlugTypes 封堵类型
func (u *PlugController) PlugTypes(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("alarm_level", enum_plug.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("plug_type", enum_plug.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PageForManage 分页
func (u *PlugController) PageForManage(ctx *gin.Context) {

	//查询条件
	conditions, err := u.getAgentSearchConditions(ctx)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, enum_agent.OrderFieldDefault)

	//终端分页
	total, agentModels := agent_models.CreateAgentFactory().Page(conditions, offset, pageSize, order)
	//转化数据
	agentIds, departmentMaps := u.parseAgentModels(agentModels, true)

	var extendInfo = make(map[string]interface{})
	if len(agentIds) > 0 {
		//规则统计数据
		PolicyFact := port_models.CreateVirtualPolicyFactory()
		extendInfo = PolicyFact.CalcByConditionsGroupAgentReturnMap(port_models.SearchConditions{
			OrgId:    u.getLogin(ctx).OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc(true))
	}
	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// SetSwitchForManage 管理页: 设置开关
func (u *PlugController) SetSwitchForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		PortSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "端口虚开管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *PlugController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))

	var param = agent_service.ManualColumn{
		PortSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "端口虚开管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForPolicy 终端策略: 分页
func (u *PlugController) PageForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := port_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := port_models.CreateVirtualPolicyFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.PortTrans{}).TransPolicyPage(total, models))
}

// SetSwitchForPolicy 终端策略: 设置开关
func (u *PlugController) SetSwitchForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddForPolicy 终端策略: 添加
func (u *PlugController) AddForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualPolicyModel{
		OrgId:       loginInfo.OrgId,
		AgentId:     u.getAgentIdFromGin(ctx),
		Protocol:    ctx.GetString(consts.ValidatorPrefixOrg + "protocol"),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		CreateId:    loginInfo.UserId,
		CreateName:  loginInfo.Username,
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForPolicy 终端策略: 编辑
func (u *PlugController) EditForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualPolicyModel{
		OrgId:       loginInfo.OrgId,
		AgentId:     u.getAgentIdFromGin(ctx),
		Id:          int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		Protocol:    ctx.GetString(consts.ValidatorPrefixOrg + "protocol"),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForPolicy 终端策略: 删除
func (u *PlugController) DeleteForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGlobal 全局策略: 分页
func (u *PlugController) PageForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := port_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := port_models.CreateVirtualGlobalFactory().Page(condition, offset, pageSize, order)

	//获取组织架构
	departmentMaps := map[int64]*department_models.DepartmentModel{}
	if models != nil && len(models) > 0 {
		var deptIds []int64
		for _, model := range models {
			if model.DeptId > 0 {
				deptIds = append(deptIds, model.DeptId)
			}
		}
		departmentMaps = u.getDepartmentMap(loginInfo.OrgId, slice_helper.Int64RemoveRepeatByMap(deptIds))
	}

	response.SuccessReturnData(ctx, (&org_trans.PortTrans{}).TransGlobalPage(total, models, departmentMaps))
}

// SetSwitchForGlobal 全局策略: 设置开关
func (u *PlugController) SetSwitchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualGlobalModel{
		OrgId:      loginInfo.OrgId,
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreateGlobalCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddForGlobal 全局策略: 添加
func (u *PlugController) AddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualGlobalModel{
		OrgId:       loginInfo.OrgId,
		Protocol:    ctx.GetString(consts.ValidatorPrefixOrg + "protocol"),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		DeptId:      int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		CreateId:    loginInfo.UserId,
		CreateName:  loginInfo.Username,
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}
	if err := port_service.CreateGlobalCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGlobal 全局策略: 编辑
func (u *PlugController) EditForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualGlobalModel{
		OrgId:       loginInfo.OrgId,
		Id:          int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		DeptId:      int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}

	if err := port_service.CreateGlobalCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局策略: 删除
func (u *PlugController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualGlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreateGlobalCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
