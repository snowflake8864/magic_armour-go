package organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// ProcessController 进程相关
type ProcessController struct {
	ControllerBase
}

// SearchStatus 查询状态
func (u *ProcessController) SearchStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("proc_status", enum_process.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("all_status", enum_process.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// ChangeStatus 可改变状态
func (u *ProcessController) ChangeStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("change_status", enum_process.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("all_status", enum_process.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// PageForManage 管理页: 分页
func (u *ProcessController) PageForManage(ctx *gin.Context) {
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
	_, departmentMaps := u.parseAgentModels(agentModels, true)

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, nil))
}

// AgentProcStatForManage 管理页: 终端进程统计
func (u *ProcessController) AgentProcStatForManage(ctx *gin.Context) {
	conditions := proc_models.SearchConditions{
		OrgId:    u.getLogin(ctx).OrgId,
		AgentIds: u.getAgentIdsFromGin(ctx),
	}
	Type := ctx.GetString(consts.ValidatorPrefixOrg + "type")
	ProcessFact := proc_models.CreateProcessFactory()
	column := ""
	switch Type {
	case "total":
		column = ProcessFact.ColumnStatusCalc(true, false, false, false, false)
		break
	case "unset":
		column = ProcessFact.ColumnStatusCalc(false, true, false, false, false)
		conditions.Status = enum_process.ProcStatusUnset
		break
	case "unknown":
		column = ProcessFact.ColumnStatusCalc(false, false, true, false, false)
		conditions.Status = enum_process.ProcStatusUnknown
		break
	case "white":
		column = ProcessFact.ColumnStatusCalc(false, false, false, true, false)
		conditions.Status = enum_process.ProcStatusWhite
		break
	case "black":
		column = ProcessFact.ColumnStatusCalc(false, false, false, false, true)
		conditions.Status = enum_process.ProcStatusBlack
		break
	default:
		column = ProcessFact.ColumnStatusCalc(false, false, false, false, false)
		break
	}

	calc := ProcessFact.CalcByConditionsGroupAgent(conditions, column)

	response.SuccessReturnData(ctx, calc)
}

// SetSwitchForManage 管理页: 设置开关
func (u *ProcessController) SetSwitchForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		ProcessSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "进程管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetModeForManage 管理页: 设置模式
func (u *ProcessController) SetModeForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	agentIds := []string{agentId}
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))
	var param = agent_service.ManualColumn{
		ProcessMode: Mode,
	}

	if Mode == enum_agent.ModeProject {
		condition := proc_models.SearchConditions{
			OrgId:   orgId,
			AgentId: agentId,
			Status:  enum_process.ProcStatusUnset,
		}
		ProcessFact := proc_models.CreateProcessFactory()
		if ProcessFact.Count(condition) > 0 {
			response.FailSimplify(ctx, "开启保护模式, 不能存在未设置的进程")
			return
		}
	}

	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "进程管理-模式")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *ProcessController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))
	var validIds []string
	if Switch == enums.SwitchOn && Mode == enum_agent.ModeProject {
		condition := proc_models.SearchConditions{
			OrgId:    orgId,
			AgentIds: agentIds,
			Status:   enum_process.ProcStatusUnset,
		}
		ProcessFact := proc_models.CreateProcessFactory()
		column := ProcessFact.ColumnStatusCalc(false, true, false, false, false)
		calcList := ProcessFact.CalcByConditionsGroupAgent(condition, column)
		if calcList == nil || len(calcList) == 0 {
			validIds = agentIds
		} else {
			if len(calcList) == len(agentIds) {
				response.FailSimplify(ctx, "开启保护模式, 不能存在未设置的进程")
				return
			} else {
				for _, calc := range calcList {
					validIds = append(validIds, calc.AgentId)
				}
			}
		}
	}

	var param = agent_service.ManualColumn{
		ProcessSwitch: &Switch,
		ProcessMode:   Mode,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, validIds, &param, "进程管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetWhiteForManage 管理页: 批量把终端未设置进程变成白名单
func (u *ProcessController) BatchSetWhiteForManage(ctx *gin.Context) {
	agentIds := u.getAgentIdsFromGin(ctx)
	err := proc_service.CreateProcessCurdFactory().BatchUnsetToWhiteByAgents(u.getLogin(ctx).OrgId, agentIds)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchRefreshForManage 管理页: 批量刷新(运行中进程)
func (u *ProcessController) BatchRefreshForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	err, validIds := u.processSwitchOnByAgents(orgId, agentIds)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgents(orgId, validIds, enum_task.TypeProcessRefresh, "批量刷新")
	//删除运行数据
	//proc_service.CreateProcessRunningFactory().DeleteByAgents(orgId, validIds)
	//修改刷新时间
	agent_service.CreateExtendFactory().BatchSetFlushAt(orgId, validIds, true)

	response.SuccessSimplify(ctx)
}

func (u *ProcessController) processSwitchOnByAgents(orgId int64, agentIds []string) (error, []string) {
	var validIds []string
	SwitchOn := enums.SwitchOn
	condition := agent_models.SearchConditions{
		OrgId:         orgId,
		AgentIds:      agentIds,
		ProcessSwitch: &SwitchOn,
	}
	AgentFact := agent_models.CreateAgentFactory()
	agentModels := AgentFact.GetListByConditions(condition, AgentFact.ColumnSimplify())
	if agentModels == nil || len(agentModels) == 0 {
		return errors.New("需要终端进程开关全部开启"), nil
	} else {
		if len(agentModels) == len(agentIds) {
			return nil, agentIds
		} else {
			for _, model := range agentModels {
				validIds = append(validIds, model.AgentId)
			}
			return nil, validIds
		}
	}

}

// BatchCollectForManage 管理页: 批量获取全部进程
func (u *ProcessController) BatchCollectForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	err, validIds := u.processSwitchOnByAgents(orgId, agentIds)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgents(orgId, validIds, enum_task.TypeProcessSuffixIsExe, "批量获取全部进程")
	//修改刷新时间
	agent_service.CreateExtendFactory().BatchSetFlushAt(orgId, validIds, true)

	response.SuccessSimplify(ctx)
}

// PageAllForAgent 单终端: 全部进程
func (u *ProcessController) PageAllForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		ProcHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		Dir:      ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Status:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status")),
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateProcessFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Process{}).TransMasterPage(total, models))
}

// PageRunningForAgent 单终端: 运行中进程
func (u *ProcessController) PageRunningForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		ProcHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		Dir:      ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Status:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status")),
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateProcessRunningFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Process{}).TransRunningPage(total, models))
}

// SearchWhileForAgent 单终端: 查询白名单
func (u *ProcessController) SearchWhileForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	Type := ctx.GetString(consts.ValidatorPrefixOrg + "type")
	Val := ctx.GetString(consts.ValidatorPrefixOrg + "val")
	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		Status:  enum_process.ProcStatusWhite,
	}

	switch Type {
	case "dir":
		condition.Dir = Val
		break
	case "hash":
		condition.ProcHash = Val
		break
	}
	var models []*proc_models.ProcessModel
	if string_helper.IsEmpty(agentId) {
		models = proc_models.CreateProcessFactory().GetListByConditionsDistinctHash(condition)
	} else {
		models = proc_models.CreateProcessFactory().GetListByConditions(condition)
	}

	var result []map[string]interface{}
	if models != nil || len(models) > 0 {
		for _, model := range models {
			result = append(result, map[string]interface{}{
				"val":   model.ProcHash,
				"label": model.Dir,
			})
		}

	}
	response.SuccessReturnData(ctx, result)
}

// SetStatusForAgent 单终端: 设置状态
func (u *ProcessController) SetStatusForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	hash := ctx.GetString(consts.ValidatorPrefixOrg + "hash")
	status := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status"))

	err := proc_service.CreateProcessCurdFactory().SetStatus(orgId, agentId, hash, status)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// CancelForAgent 单终端: 取消
func (u *ProcessController) CancelForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	hash := ctx.GetString(consts.ValidatorPrefixOrg + "hash")

	err := proc_service.CreateProcessCurdFactory().Cancel(orgId, agentId, hash)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchWhiteForAgent 单终端: 批量白名单
func (u *ProcessController) BatchWhiteForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "hashes")
	hashes := json_helper.ToSliceString(jsonStr)
	if len(hashes) == 0 {
		response.FailSimplify(ctx, "数据不能为空")
		return
	}

	err := proc_service.CreateProcessCurdFactory().BatchWhiteByHashes(orgId, agentId, hashes)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ClearAllForAgent 单终端: 清空全部
func (u *ProcessController) ClearAllForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	err := proc_service.CreateProcessCurdFactory().DeleteByAgent(orgId, agentId, true)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ClearRunningForAgent 单终端: 清空运行中
func (u *ProcessController) ClearRunningForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	if !proc_service.CreateProcessRunningFactory().DeleteByAgent(orgId, agentId) {
		response.FailSimplify(ctx, "清空失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// RefreshForAgent 单终端: 刷新(运行中进程)
func (u *ProcessController) RefreshForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeProcessRefresh, "刷新")
	//删除运行数据
	proc_service.CreateProcessRunningFactory().DeleteByAgent(orgId, agentId)
	//修改刷新时间
	agent_service.CreateExtendFactory().SetFlushAt(orgId, agentId, true)

	response.SuccessSimplify(ctx)
}

// CollectForAgent 单终端: 获取全部进程
func (u *ProcessController) CollectForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeProcessSuffixIsExe, "获取全部进程")
	//修改刷新时间
	agent_service.CreateExtendFactory().SetFlushAt(orgId, agentId, true)
}

// RepeatStatForAgent 单终端: 重复进程统计
func (u *ProcessController) RepeatStatForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "hashes")
	hashes := json_helper.ToSliceString(jsonStr)
	if len(hashes) == 0 {
		response.FailSimplify(ctx, "数据不能为空")
		return
	}
	conditions := proc_models.SearchConditions{
		OrgId:      orgId,
		ProcHashes: hashes,
	}
	ProcessFact := proc_models.CreateProcessFactory()
	calc := ProcessFact.CalcByConditionsGroupHash(conditions)
	response.SuccessReturnData(ctx, calc)
}

// RepeatPageForAgent 单终端: 重复进程分页
func (u *ProcessController) RepeatPageForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		ProcHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateProcessFactory().Page(condition, offset, pageSize, order)

	var agentIds []string
	var agentMap = make(map[string]*agent_models.AgentModel)
	if models != nil || len(models) > 0 {
		for _, model := range models {
			agentIds = append(agentIds, model.AgentId)
		}
		AgentFact := agent_models.CreateAgentFactory()
		agentMap = AgentFact.GetMapsByOrgAndAgents(orgId, agentIds, AgentFact.ColumnBasic(true, false))
	}

	response.SuccessReturnData(ctx, (&org_trans.Process{}).TransRepeatPage(total, models, agentMap))
}

// RunningDllExtractForAgent 单终端: 提取运行中模块
func (u *ProcessController) RunningDllExtractForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	hash := ctx.GetString(consts.ValidatorPrefixOrg + "hash")
	procId := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "procId"))

	AgentFact := agent_models.CreateAgentFactory()
	agentModel := AgentFact.GetItemByOrgAndAgent(orgId, agentId, AgentFact.ColumnProc(true))
	if agentModel == nil || agentModel.OrgId == 0 {
		response.FailSimplify(ctx, "终端不存在")
		return
	}
	if agentModel.PlatformType == enum_software.PlatformTypeLinux {
		response.FailSimplify(ctx, "Linux暂不支持功能")
		return
	}
	if agentModel.DllSwitch == enums.SwitchOff {
		response.FailSimplify(ctx, "模块功能没有开启")
		return
	}

	if err := proc_service.CreateProcessRunningFactory().ExtractDlls(orgId, agentId, hash, procId); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// RunningDllShowForAgent 单终端: 展示运行中模块
func (u *ProcessController) RunningDllShowForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	hash := ctx.GetString(consts.ValidatorPrefixOrg + "hash")
	procId := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "procId"))

	hasModel := proc_models.CreateProcessRunningFactory().GetItemByOrgAndAgentAndHashAndProcId(orgId, agentId, hash, procId)
	if hasModel == nil || hasModel.DllStatus == enum_process.CollectDllStatusDefault {
		response.SuccessSimplify(ctx)
		return
	}
	if hasModel.DllStatus != enum_process.CollectDllStatusFinish {
		response.FailRoutine(ctx, 100, "请耐心等待")
		return
	}

	var list []proc_models.ProcessRunningDll
	_ = json_helper.ToAny(hasModel.DllList, &list)

	response.SuccessReturnData(ctx, list)
}

// PageForSearch 搜索: 分页
func (u *ProcessController) PageForSearch(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForSearch")
}
