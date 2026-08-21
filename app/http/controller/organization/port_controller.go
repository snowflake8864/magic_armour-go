package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_port"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/port_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// PortController 端口相关
type PortController struct {
	ControllerBase
}

// Protocols 协议
func (u *PortController) Protocols(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfStr("protocol", enum_port.WebMapStr)
	labels := enums.Get1stIsStrReturnMapStrOfStr("protocol", enum_port.EnumMapStr)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatStringList(keys, labels, false))
}

// IpTypes IP类型
func (u *PortController) IpTypes(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfStr("ip_type", enum_port.WebMapStr)
	labels := enums.Get1stIsStrReturnMapStrOfStr("ip_type", enum_port.EnumMapStr)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatStringList(keys, labels, false))
}

// AlarmLevels 告警权重
func (u *PortController) AlarmLevels(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("alarm_level", enum_port.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("alarm_level", enum_port.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, false))
}

// ForwardTypes 转发类型
func (u *PortController) ForwardTypes(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("forward_type", enum_port.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("forward_type", enum_port.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, false))
}

// PageForManage 分页
func (u *PortController) PageForManage(ctx *gin.Context) {

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
	agentIds, departmentMaps := u.parseAgentModels(agentModels, true)

	var extendInfo = make(map[string]interface{})
	if len(agentIds) > 0 {
		//规则统计数据
		PolicyFact := port_models.CreateVirtualPolicyFactory()
		calcPolicy := PolicyFact.CalcByConditionsGroupAgentReturnMap(port_models.SearchConditions{
			OrgId:    u.getLogin(ctx).OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc(true))

		for _, agentModel := range agentModels {
			gg := map[string]interface{}{
				"total_number": 0,
				"on_number":    0,
			}
			if item, isOk := calcPolicy[agentModel.AgentId]; isOk {
				gg["total_number"] = item.TotalNumber
				gg["on_number"] = item.OnNumber

			}
			extendInfo[agentModel.AgentId] = gg
		}
	}
	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// SetSwitchForManage 管理页: 设置开关
func (u *PortController) SetSwitchForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		PortSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().UpdateByAgent(orgId, agentId, &param, "Manage Page Switch")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *PortController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))

	var param = agent_service.ManualColumn{
		PortSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "Manage Page Batch Modify")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForWorking 正常端口: 分页
func (u *PortController) PageForWorking(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := port_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
		Port:    int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "port"))),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := port_models.CreateWorkingFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.PortTrans{}).TransWorkingPage(total, models))
}

// RefreshForWorking 正常端口: 刷新
func (u *PortController) RefreshForWorking(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	if !port_service.CreateWorkingFactory().DelByOrgAndAgent(loginInfo.OrgId, agentId) {
		response.FailSimplify(ctx, my_errors.CurdRefreshFailed)
		return
	}

	task_service.CreateCurdFactory().AddTaskForAgent(loginInfo.OrgId, agentId, enum_task.TypePortWorkingUpload, "Refresh")

	response.SuccessSimplify(ctx)
}

// PageForPolicy 终端策略: 分页
func (u *PortController) PageForPolicy(ctx *gin.Context) {
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
func (u *PortController) SetSwitchForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory(cur_userinfo.GetLocales(ctx)).SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddForPolicy 终端策略: 添加
func (u *PortController) AddForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualPolicyModel{
		OrgId:       loginInfo.OrgId,
		AgentId:     u.getAgentIdFromGin(ctx),
		Protocol:    ctx.GetString(consts.ValidatorPrefixOrg + "protocol"),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		ForwardType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardType")),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		CreateId:    loginInfo.UserId,
		CreateName:  loginInfo.Username,
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory(cur_userinfo.GetLocales(ctx)).Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForPolicy 终端策略: 编辑
func (u *PortController) EditForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualPolicyModel{
		OrgId:       loginInfo.OrgId,
		AgentId:     u.getAgentIdFromGin(ctx),
		Id:          string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Protocol:    ctx.GetString(consts.ValidatorPrefixOrg + "protocol"),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		ForwardType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardType")),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory(cur_userinfo.GetLocales(ctx)).Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForPolicy 终端策略: 删除
func (u *PortController) DeleteForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreatePolicyCurdFactory(cur_userinfo.GetLocales(ctx)).DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGlobal 全局策略: 分页
func (u *PortController) PageForGlobal(ctx *gin.Context) {
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
func (u *PortController) SetSwitchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualGlobalModel{
		OrgId:      loginInfo.OrgId,
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreateGlobalCurdFactory(cur_userinfo.GetLocales(ctx)).SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddForGlobal 全局策略: 添加
func (u *PortController) AddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualGlobalModel{
		OrgId:       loginInfo.OrgId,
		Protocol:    ctx.GetString(consts.ValidatorPrefixOrg + "protocol"),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		ForwardType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardType")),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		DeptId:      int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		CreateId:    loginInfo.UserId,
		CreateName:  loginInfo.Username,
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}
	if err := port_service.CreateGlobalCurdFactory(cur_userinfo.GetLocales(ctx)).Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGlobal 全局策略: 编辑
func (u *PortController) EditForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := port_models.VirtualGlobalModel{
		OrgId:       loginInfo.OrgId,
		Id:          string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Protocol:    ctx.GetString(consts.ValidatorPrefixOrg + "protocol"),
		Type:        ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		SourcePort:  ctx.GetString(consts.ValidatorPrefixOrg + "sourcePort"),
		ForwardIp:   ctx.GetString(consts.ValidatorPrefixOrg + "forwardIp"),
		ForwardPort: int_helper.IntToStr(int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardPort"))),
		ForwardType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "forwardType")),
		AlarmLevel:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmLevel")),
		DeptId:      int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}

	if err := port_service.CreateGlobalCurdFactory(cur_userinfo.GetLocales(ctx)).Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局策略: 删除
func (u *PortController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := port_models.VirtualGlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := port_service.CreateGlobalCurdFactory(cur_userinfo.GetLocales(ctx)).DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
