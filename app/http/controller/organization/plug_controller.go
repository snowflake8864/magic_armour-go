package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/plug_models"
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
	keys := enums.Get1stIsStrReturnSliceOfInt("plug_type", enum_plug.WebMapInt)
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
		loginInfo := u.getLogin(ctx)

		//规则统计数据
		PolicyFact := plug_models.CreateDynamicFactory()
		calcPolicy := PolicyFact.CalcByConditionsGroupAgentReturnMap(plug_models.SearchConditions{
			OrgId:    loginInfo.OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc(true))

		IpPolicyFact := plug_models.CreateIpPolicyFactory()
		calcWhite := IpPolicyFact.CalcByConditionsGroupAgentReturnMap(plug_models.SearchConditions{
			OrgId:    loginInfo.OrgId,
			AgentIds: agentIds,
			IpType:   enum_plug.IpTypeWhite,
		}, IpPolicyFact.ColumnCalc(false))

		calcBlack := IpPolicyFact.CalcByConditionsGroupAgentReturnMap(plug_models.SearchConditions{
			OrgId:    loginInfo.OrgId,
			AgentIds: agentIds,
			IpType:   enum_plug.IpTypeBlack,
		}, IpPolicyFact.ColumnCalc(false))

		for _, agentModel := range agentModels {
			gg := map[string]interface{}{
				"dynamic_total_number":  0,
				"dynamic_on_number":     0,
				"ip_white_total_number": 0,
				"ip_black_total_number": 0,
			}

			if item, isOk := calcPolicy[agentModel.AgentId]; isOk {
				gg["dynamic_total_number"] = item.TotalNumber
				gg["dynamic_on_number"] = item.OnNumber

			}
			if item, isOk := calcWhite[agentModel.AgentId]; isOk {
				gg["ip_white_total_number"] = item.TotalNumber
			}
			if item, isOk := calcBlack[agentModel.AgentId]; isOk {
				gg["ip_black_total_number"] = item.TotalNumber
			}
			extendInfo[agentModel.AgentId] = gg
		}

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
		PlugSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "动态封堵管理-开关")
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
		PlugSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "动态封堵管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigPortGetForPolicy 终端策略 - 端口虚开
func (u *PlugController) ConfigPortGetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := agent_service.CreateExtendFactory().GetPlugOfPort(loginInfo.OrgId, agentId)

	response.SuccessReturnData(ctx, config)
}

// ConfigPortSetForPolicy 终端策略 - 端口虚开
func (u *PlugController) ConfigPortSetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := enum_plug.PlugPortConfig{
		Switch:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Weight:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "weight")),
		Type:      int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		CountTime: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "countTime")),
		PlugTime:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "plugTime")),
	}
	if err := agent_service.CreateExtendFactory().SetPlugOfPort(loginInfo.OrgId, agentId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigLoginGetForPolicy 终端策略 - 口令破解
func (u *PlugController) ConfigLoginGetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := agent_service.CreateExtendFactory().GetPlugOfLogin(loginInfo.OrgId, agentId)

	response.SuccessReturnData(ctx, config)
}

// ConfigLoginSetForPolicy 终端策略 - 口令破解
func (u *PlugController) ConfigLoginSetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := enum_plug.PlugLoginConfig{
		Switch:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Type:     int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		PlugTime: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "plugTime")),
	}
	if err := agent_service.CreateExtendFactory().SetPlugOfLogin(loginInfo.OrgId, agentId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpPageForPolicy 终端策略 - ip名单库: 分页
func (u *PlugController) IpPageForPolicy(ctx *gin.Context, ipType int) {

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

// IpSetSwitchForPolicy 终端策略 - ip名单库: 设置开关
func (u *PlugController) IpSetSwitchForPolicy(ctx *gin.Context, ipType int) {
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

// IpAddForPolicy 终端策略 - ip名单库: 添加
func (u *PlugController) IpAddForPolicy(ctx *gin.Context, ipType int) {
	fmt.Println("===================ipType", ipType)

	loginInfo := u.getLogin(ctx)
	param := plug_models.IpPolicyModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Type:       ipType,
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpEditForPolicy 终端策略 - ip名单库: 编辑
func (u *PlugController) IpEditForPolicy(ctx *gin.Context, ipType int) {
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

// IpDeleteForPolicy 终端策略 - ip名单库: 删除
func (u *PlugController) IpDeleteForPolicy(ctx *gin.Context, ipType int) {
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

// IpClearForPolicy 终端策略 - ip名单库: 删除
func (u *PlugController) IpClearForPolicy(ctx *gin.Context, ipType int) {
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

// DynamicPageForPolicy 终端策略 - 动态封堵: 分页
func (u *PlugController) DynamicPageForPolicy(ctx *gin.Context, openSwitch int) {
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

// DynamicDisableForPolicy 终端策略 - 动态封堵: 解封
func (u *PlugController) DynamicDisableForPolicy(ctx *gin.Context) {
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

// DynamicDisableAllForPolicy 终端策略 - 动态封堵: 解封全部
func (u *PlugController) DynamicDisableAllForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
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

// DynamicDeleteForPolicy 终端策略 - 动态封堵: 删除
func (u *PlugController) DynamicDeleteForPolicy(ctx *gin.Context) {
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

// DynamicDeleteAllForPolicy 终端策略 - 动态封堵: 删除全部
func (u *PlugController) DynamicDeleteAllForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigPortGetForGlobal 全局策略 - 端口虚开
func (u *PlugController) ConfigPortGetForGlobal(ctx *gin.Context) {

}

// ConfigPortSetForGlobal 全局策略 - 端口虚开
func (u *PlugController) ConfigPortSetForGlobal(ctx *gin.Context) {

}

// ConfigLoginGetForGlobal 全局策略 - 口令破解
func (u *PlugController) ConfigLoginGetForGlobal(ctx *gin.Context) {

}

// ConfigLoginSetForGlobal 全局策略 - 口令破解
func (u *PlugController) ConfigLoginSetForGlobal(ctx *gin.Context) {

}

// IpPageForGlobal 全局策略 - ip名单库: 分页
func (u *PlugController) IpPageForGlobal(ctx *gin.Context, ipType int) {
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

// IpSetSwitchForGlobal 全局策略 - ip名单库: 设置开关
func (u *PlugController) IpSetSwitchForGlobal(ctx *gin.Context, ipType int) {
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

// IpAddForGlobal 全局策略 - ip名单库: 添加
func (u *PlugController) IpAddForGlobal(ctx *gin.Context, ipType int) {
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

// IpEditForGlobal 全局策略 - ip名单库: 编辑
func (u *PlugController) IpEditForGlobal(ctx *gin.Context, ipType int) {
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

// IpDeleteForGlobal 全局策略 - ip名单库: 删除
func (u *PlugController) IpDeleteForGlobal(ctx *gin.Context, ipType int) {
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

// IpClearForGlobal 全局策略 - ip名单库: 删除
func (u *PlugController) IpClearForGlobal(ctx *gin.Context, ipType int) {
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

// DynamicPageForGlobal 全局策略 - 动态封堵: 分页
func (u *PlugController) DynamicPageForGlobal(ctx *gin.Context, openSwitch int) {
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

// DynamicDisableForGlobal 全局策略 - 动态封堵: 解封
func (u *PlugController) DynamicDisableForGlobal(ctx *gin.Context) {
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

// DynamicDisableAllForGlobal 全局策略 - 动态封堵: 解封全部
func (u *PlugController) DynamicDisableAllForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
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

// DynamicDeleteForGlobal 全局策略 - 动态封堵: 删除
func (u *PlugController) DynamicDeleteForGlobal(ctx *gin.Context) {
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

// DynamicDeleteAllForGlobal 全局策略 - 动态封堵: 删除全部
func (u *PlugController) DynamicDeleteAllForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
