package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/model/agent/agent_online_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/online_log_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/response"
	"sort"
)

// AgentController 终端相关
type AgentController struct {
	ControllerBase
}

// OnlineStatus 在线状态
func (u *AgentController) OnlineStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("is_online", enum_agent.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("is_online", enum_agent.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// SafeStatus 安全状态
func (u *AgentController) SafeStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("is_safe", enum_agent.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("is_safe", enum_agent.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// GrantStatus 授权状态
func (u *AgentController) GrantStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("grant_status", enum_agent.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("grant_status", enum_agent.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// CountOnline 统计在线和离线数量
func (u *AgentController) CountOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "CountOnline")
}

// TypeList 类型列表
func (u *AgentController) TypeList(ctx *gin.Context) {
	jsonStr := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentTypes)
	typeMap := enum_agent.GetTypeMap(jsonStr)

	//先取值
	var typeVals []int
	if typeMap != nil {
		for val, _ := range typeMap {
			typeVals = append(typeVals, val)
		}
	}
	sort.Ints(typeVals)

	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(typeVals, typeMap, false))
}

// PageForManage 管理: 分页
func (u *AgentController) PageForManage(ctx *gin.Context) {
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
	//软件最新发布
	softwareNewest := software_service.CreateCurdFactory().GetMapsByPublished()

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransManagePage(total, agentModels, departmentMaps, softwareNewest))
}

// DeleteForManage 管理: 删除
func (u *AgentController) DeleteForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	errCode, err := agent_service.CreateCurdFactory().BatchDelete(u.getLogin(ctx).OrgId, agentIds, false)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// UninstallForManage 管理: 卸载
func (u *AgentController) UninstallForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	errCode, err := agent_service.CreateCurdFactory().BatchDelete(u.getLogin(ctx).OrgId, agentIds, true)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// UpgradeForManage 管理: 升级
func (u *AgentController) UpgradeForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpgrade(u.getLogin(ctx).OrgId, agentIds)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchDeleteForManage 管理: 批量删除
func (u *AgentController) BatchDeleteForManage(ctx *gin.Context) {
	agentIds := u.getAgentIdsFromGin(ctx)
	errCode, err := agent_service.CreateCurdFactory().BatchDelete(u.getLogin(ctx).OrgId, agentIds, false)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchUninstallForManage 管理: 批量卸载
func (u *AgentController) BatchUninstallForManage(ctx *gin.Context) {
	agentIds := u.getAgentIdsFromGin(ctx)
	errCode, err := agent_service.CreateCurdFactory().BatchDelete(u.getLogin(ctx).OrgId, agentIds, true)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchUpgradeForManage 管理: 批量升级
func (u *AgentController) BatchUpgradeForManage(ctx *gin.Context) {
	agentIds := u.getAgentIdsFromGin(ctx)
	errCode, err := agent_service.CreateCurdFactory().BatchUpgrade(u.getLogin(ctx).OrgId, agentIds)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchModifyForManage 管理: 批量修改
func (u *AgentController) BatchModifyForManage(ctx *gin.Context) {
	agentIds := u.getAgentIdsFromGin(ctx)
	DeptId := int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId"))
	SyslogDnsSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogDnsSwitch"))
	SyslogOuterSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogOuterSwitch"))
	SyslogInnerSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogInnerSwitch"))
	SyslogProcessSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogProcessSwitch"))
	SyslogLoginSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogLoginSwitch"))
	SelfProtectSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "selfProtectSwitch"))
	InternetSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "internetSwitch"))
	AutoWhiteSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "autoWhiteSwitch"))

	var param = agent_service.ManualColumn{
		DeptId:              &DeptId,
		SyslogDnsSwitch:     &SyslogDnsSwitch,
		SyslogOuterSwitch:   &SyslogOuterSwitch,
		SyslogInnerSwitch:   &SyslogInnerSwitch,
		SyslogProcessSwitch: &SyslogProcessSwitch,
		SyslogLoginSwitch:   &SyslogLoginSwitch,
		SelfProtectSwitch:   &SelfProtectSwitch,
		InternetSwitch:      &InternetSwitch,
		AutoWhiteSwitch:     &AutoWhiteSwitch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "批量修改")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetForBasic 基本信息：获取
func (u *AgentController) GetForBasic(ctx *gin.Context) {
	agentId := u.getAgentIdFromGin(ctx)

	curdFact := agent_models.CreateAgentFactory()
	model := curdFact.GetItemByOrgAndAgent(u.getLogin(ctx).OrgId, agentId, curdFact.ColumnBasic(true, false))

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransSingle(model))
}

// SaveForBasic 基本信息：保存
func (u *AgentController) SaveForBasic(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	deptId := int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId"))
	agentType := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "agentType"))
	alias := ctx.GetString(consts.ValidatorPrefixOrg + "alias")
	contactName := ctx.GetString(consts.ValidatorPrefixOrg + "contactName")
	contactTel := ctx.GetString(consts.ValidatorPrefixOrg + "contactTel")
	position := ctx.GetString(consts.ValidatorPrefixOrg + "position")
	remark := ctx.GetString(consts.ValidatorPrefixOrg + "remark")

	var param = agent_service.ManualColumn{
		DeptId:      &deptId,
		AgentType:   &agentType,
		Alias:       &alias,
		ContactName: &contactName,
		ContactTel:  &contactTel,
		Position:    &position,
		Remark:      &remark,
	}
	errCode, err := agent_service.CreateCurdFactory().UpdateByAgent(loginInfo.OrgId, agentId, &param, "保存基本信息")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ResetUidForBasic 基本信息： 重置终端ID
func (u *AgentController) ResetUidForBasic(ctx *gin.Context) {
	agentId := u.getAgentIdFromGin(ctx)

	task_service.CreateCurdFactory().AddTaskForAgent(u.getLogin(ctx).OrgId, agentId, enum_task.TypeResetUUid, "")

	response.SuccessSimplify(ctx)
}

// GetForConfig 配置信息：获取
func (u *AgentController) GetForConfig(ctx *gin.Context) {
	agentId := u.getAgentIdFromGin(ctx)

	curdFact := agent_models.CreateAgentFactory()
	model := curdFact.GetItemByOrgAndAgent(u.getLogin(ctx).OrgId, agentId, curdFact.ColumnConfig(true, true))

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransSingle(model))
}

// SaveForConfig 配置信息：保存
func (u *AgentController) SaveForConfig(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	processSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "processSwitch"))
	processMode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "processMode"))
	dllSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "dllSwitch"))
	dllMode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "dllMode"))
	tamperSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "tamperSwitch"))
	tamperMode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "tamperMode"))
	blackmailSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "blackmailSwitch"))
	blackmailMode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "blackmailMode"))
	portSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "portSwitch"))
	plugSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "plugSwitch"))
	usbSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "usbSwitch"))
	usbMode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "usbMode"))
	loginSolveSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "loginSolveSwitch"))
	syslogDnsSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogDnsSwitch"))
	syslogOuterSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogOuterSwitch"))
	syslogInnerSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogInnerSwitch"))
	syslogProcessSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogProcessSwitch"))
	syslogLoginSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "syslogLoginSwitch"))
	outreachAlarmSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "outreachAlarmSwitch"))
	selfProtectSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "selfProtectSwitch"))
	internetSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "internetSwitch"))
	autoWhiteSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "autoWhiteSwitch"))
	commonTime := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "commonTime"))
	fastTime := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "fastTime"))

	if processSwitch == enums.SwitchOn && processMode == enum_agent.ModeProject {
		condition := proc_models.SearchConditions{
			OrgId:   orgId,
			AgentId: agentId,
			Status:  enum_process.ProcStatusUnset,
		}
		ProcessFact := proc_models.CreateProcessFactory()
		if ProcessFact.Count(condition) > 0 {
			response.FailSimplify(ctx, "开启进程保护模式, 不能存在未设置的进程")
			return
		}
	}

	var param = agent_service.ManualColumn{
		ProcessSwitch:       &processSwitch,
		ProcessMode:         processMode,
		DllSwitch:           &dllSwitch,
		DllMode:             dllMode,
		TamperSwitch:        &tamperSwitch,
		TamperMode:          tamperMode,
		BlackmailSwitch:     &blackmailSwitch,
		BlackmailMode:       blackmailMode,
		PortSwitch:          &portSwitch,
		PlugSwitch:          &plugSwitch,
		UsbSwitch:           &usbSwitch,
		UsbMode:             usbMode,
		OutreachAlarmSwitch: &outreachAlarmSwitch,
		LoginSolveSwitch:    &loginSolveSwitch,
		SyslogDnsSwitch:     &syslogDnsSwitch,
		SyslogOuterSwitch:   &syslogOuterSwitch,
		SyslogInnerSwitch:   &syslogInnerSwitch,
		SyslogProcessSwitch: &syslogProcessSwitch,
		SyslogLoginSwitch:   &syslogLoginSwitch,
		SelfProtectSwitch:   &selfProtectSwitch,
		InternetSwitch:      &internetSwitch,
		AutoWhiteSwitch:     &autoWhiteSwitch,
		CommonTime:          commonTime,
		FastTime:            fastTime,
	}
	errCode, err := agent_service.CreateCurdFactory().UpdateByAgent(orgId, agentId, &param, "保存配置信息")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ListForMonitor 监控：列表
func (u *AgentController) ListForMonitor(ctx *gin.Context) {
	//页面上传
	conditions := agent_monitor_models.SearchConditions{}
	conditions.OrgId = u.getLogin(ctx).OrgId
	conditions.AgentId = u.getAgentIdFromGin(ctx)
	conditions.MonitorDay = ctx.GetString(consts.ValidatorPrefixOrg + "date")

	models := agent_monitor_models.CreateMonitorFactory().GetListByConditions(conditions, "")

	response.SuccessReturnData(ctx, (&org_trans.Monitor{}).TransList(models))
}

// InfoForMonitor 监控：信息
func (u *AgentController) InfoForMonitor(ctx *gin.Context) {
	//页面上传
	conditions := agent_monitor_models.SearchConditions{}
	conditions.OrgId = u.getLogin(ctx).OrgId
	conditions.AgentId = u.getAgentIdFromGin(ctx)
	conditions.MonitorDay = ctx.GetString(consts.ValidatorPrefixOrg + "date")
	conditions.CreatedAt = conditions.MonitorDay + " " + ctx.GetString(consts.ValidatorPrefixOrg+"time")

	model := agent_monitor_models.CreateMonitorFactory().GetItemByConditions(conditions)

	response.SuccessReturnData(ctx, model)
}

// PageForOnlineLog 上线日志：分页
func (u *AgentController) PageForOnlineLog(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	//页面上传
	conditions := agent_online_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId
	conditions.AgentId = u.getAgentIdFromGin(ctx)
	conditions.AgentIp = ctx.GetString(consts.ValidatorPrefixOrg + "agentIp")
	conditions.NetIp = ctx.GetString(consts.ValidatorPrefixOrg + "netIp")
	isOnline := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isOnline"))
	conditions.IsOnline = &isOnline
	conditions.CreateAtGte = ctx.GetString(consts.ValidatorPrefixOrg + "startTime")
	conditions.CreateAtLte = ctx.GetString(consts.ValidatorPrefixOrg + "endTime")

	//分页参数和排序参数
	offset, pageSize, _ := u.getPageAndSort(ctx, "")

	//分页
	total, models := agent_online_models.CreateOnlineLogFactory().Page(conditions, offset, pageSize, "")

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

	response.SuccessReturnData(ctx, (&org_trans.AgentOnline{}).TransPage(total, models, departmentMaps))
}

// ClearForOnlineLog 上线日志：分页
func (u *AgentController) ClearForOnlineLog(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	if err := online_log_service.CreateCurdFactory().DeleteByOrg(loginInfo.OrgId); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
