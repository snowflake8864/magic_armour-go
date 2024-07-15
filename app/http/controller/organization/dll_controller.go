package organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/response"
)

// DllController 模块相关
type DllController struct {
	ControllerBase
}

// PageForManage 管理页: 分页
func (u *DllController) PageForManage(ctx *gin.Context) {
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

// AgentProcStatForManage 管理页: 终端进程统计
func (u *DllController) AgentProcStatForManage(ctx *gin.Context) {
	conditions := proc_models.SearchConditions{
		OrgId:    u.getLogin(ctx).OrgId,
		AgentIds: u.getAgentIdsFromGin(ctx),
	}
	Type := ctx.GetString(consts.ValidatorPrefixOrg + "type")
	DllFact := proc_models.CreateDllFactory()
	column := ""
	switch Type {
	case "total":
		column = DllFact.ColumnStatusCalc(true, false, false, false, false)
		break
	case "unset":
		column = DllFact.ColumnStatusCalc(false, true, false, false, false)
		conditions.Status = enum_process.ProcStatusUnset
		break
	case "unknown":
		column = DllFact.ColumnStatusCalc(false, false, true, false, false)
		conditions.Status = enum_process.ProcStatusUnknown
		break
	case "white":
		column = DllFact.ColumnStatusCalc(false, false, false, true, false)
		conditions.Status = enum_process.ProcStatusWhite
		break
	case "black":
		column = DllFact.ColumnStatusCalc(false, false, false, false, true)
		conditions.Status = enum_process.ProcStatusBlack
		break
	default:
		column = DllFact.ColumnStatusCalc(false, false, false, false, false)
		break
	}

	calc := DllFact.CalcByConditionsGroupAgent(conditions, column)

	response.SuccessReturnData(ctx, calc)
}

// SetSwitchForManage 管理页: 设置开关
func (u *DllController) SetSwitchForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		DllSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "模块管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetModeForManage 管理页: 设置模式
func (u *DllController) SetModeForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	agentIds := []string{agentId}
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))
	var param = agent_service.ManualColumn{
		DllMode: Mode,
	}

	if Mode == enum_agent.ModeProject {
		condition := proc_models.SearchConditions{
			OrgId:   orgId,
			AgentId: agentId,
			Status:  enum_process.ProcStatusUnset,
		}
		DllFact := proc_models.CreateDllFactory()
		if DllFact.Count(condition) > 0 {
			response.FailSimplify(ctx, "开启保护模式, 不能存在未设置的模块")
			return
		}
	}

	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "模块管理-模式")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *DllController) BatchSetConfigForManage(ctx *gin.Context) {
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
		DllFact := proc_models.CreateDllFactory()
		column := DllFact.ColumnStatusCalc(false, true, false, false, false)
		calcList := DllFact.CalcByConditionsGroupAgent(condition, column)
		if calcList == nil || len(calcList) == 0 {
			validIds = agentIds
		} else {
			if len(calcList) == len(agentIds) {
				response.FailSimplify(ctx, "开启保护模式, 不能存在未设置的模块")
				return
			} else {
				for _, calc := range calcList {
					validIds = append(validIds, calc.AgentId)
				}
			}
		}
	}

	var param = agent_service.ManualColumn{
		DllSwitch: &Switch,
		DllMode:   Mode,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, validIds, &param, "模块管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetWhiteForManage 管理页: 批量把终端未设置进程变成白名单
func (u *DllController) BatchSetWhiteForManage(ctx *gin.Context) {
	agentIds := u.getAgentIdsFromGin(ctx)
	err := proc_service.CreateDllCurdFactory().BatchUnsetToWhiteByAgents(u.getLogin(ctx).OrgId, agentIds)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchRefreshForManage 管理页: 批量刷新(运行中模块)
func (u *DllController) BatchRefreshForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	err, validIds := u.dllSwitchOnByAgents(orgId, agentIds)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgents(orgId, validIds, enum_task.TypeDllRefresh, "批量刷新")
	//删除运行数据
	//proc_service.CreateDllRunningFactory().DeleteByAgents(orgId, validIds)
	//修改刷新时间
	agent_service.CreateExtendFactory().BatchSetFlushAt(orgId, validIds, false)

	response.SuccessSimplify(ctx)
}

func (u *DllController) dllSwitchOnByAgents(orgId int64, agentIds []string) (error, []string) {
	var validIds []string
	SwitchOn := enums.SwitchOn
	condition := agent_models.SearchConditions{
		OrgId:     orgId,
		AgentIds:  agentIds,
		DllSwitch: &SwitchOn,
	}
	AgentFact := agent_models.CreateAgentFactory()
	agentModels := AgentFact.GetListByConditions(condition, AgentFact.ColumnSimplify())
	if agentModels == nil || len(agentModels) == 0 {
		return errors.New("需要终端模块开关全部开启"), nil
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

// BatchCollectForManage 管理页: 批量获取全部模块
func (u *DllController) BatchCollectForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	err, validIds := u.dllSwitchOnByAgents(orgId, agentIds)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgents(orgId, validIds, enum_task.TypeDllSuffixIsDll, "批量获取全部模块")
	//修改刷新时间
	agent_service.CreateExtendFactory().BatchSetFlushAt(orgId, validIds, false)

	response.SuccessSimplify(ctx)
}

// PageAllForAgent 单终端: 全部进程
func (u *DllController) PageAllForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		DllHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		Dir:     ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Status:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status")),
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateDllFactory().Page(condition, offset, pageSize, order)
	//威胁情报数据
	threatenMaps := make(map[string]*threaten_models.CacheModel)
	if models != nil {
		existHash := make(map[string]struct{})
		var hashArr []string
		for _, model := range models {
			if _, isOK := existHash[model.DllHash]; !isOK {
				existHash[model.DllHash] = struct{}{}
				hashArr = append(hashArr, model.DllHash)
			}
		}
		threatenMaps = threaten_models.CreateCacheFactory().GetMapsByIds(hashArr)
	}
	response.SuccessReturnData(ctx, (&org_trans.Dll{}).TransMasterPage(total, models, threatenMaps))
}

// PageRunningForAgent 单终端: 运行中进程
func (u *DllController) PageRunningForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		DllHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		Dir:     ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Status:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status")),
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateDllRunningFactory().Page(condition, offset, pageSize, order)

	//威胁情报数据
	threatenMaps := make(map[string]*threaten_models.CacheModel)
	if models != nil {
		existHash := make(map[string]struct{})
		var hashArr []string
		for _, model := range models {
			if _, isOK := existHash[model.DllHash]; !isOK {
				existHash[model.DllHash] = struct{}{}
				hashArr = append(hashArr, model.DllHash)
			}
		}
		threatenMaps = threaten_models.CreateCacheFactory().GetMapsByIds(hashArr)
	}
	response.SuccessReturnData(ctx, (&org_trans.Dll{}).TransRunningPage(total, models, threatenMaps))
}

// SetStatusForAgent 单终端: 设置状态
func (u *DllController) SetStatusForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	hash := ctx.GetString(consts.ValidatorPrefixOrg + "hash")
	status := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status"))

	err := proc_service.CreateDllCurdFactory().SetStatus(orgId, agentId, hash, status)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// CancelForAgent 单终端: 取消
func (u *DllController) CancelForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	hash := ctx.GetString(consts.ValidatorPrefixOrg + "hash")

	err := proc_service.CreateDllCurdFactory().Cancel(orgId, agentId, hash)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchWhiteForAgent 单终端: 批量白名单
func (u *DllController) BatchWhiteForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "hashes")
	hashes := json_helper.ToSliceString(jsonStr)
	if len(hashes) == 0 {
		response.FailSimplify(ctx, "数据不能为空")
		return
	}

	err := proc_service.CreateDllCurdFactory().BatchWhiteByHashes(orgId, agentId, hashes)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ClearAllForAgent 单终端: 清空全部
func (u *DllController) ClearAllForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	err := proc_service.CreateDllCurdFactory().DeleteByAgent(orgId, agentId, true)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ClearRunningForAgent 单终端: 清空运行中
func (u *DllController) ClearRunningForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	if !proc_service.CreateProcessRunningFactory().DeleteByAgent(orgId, agentId) {
		response.FailSimplify(ctx, "清空失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// RefreshForAgent 单终端: 刷新(运行中进程)
func (u *DllController) RefreshForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)

	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeDllRefresh, "刷新")
	//删除运行数据
	proc_service.CreateDllRunningFactory().DeleteByAgent(orgId, agentId)
	//修改刷新时间
	agent_service.CreateExtendFactory().SetFlushAt(orgId, agentId, false)

	response.SuccessSimplify(ctx)
}

// CollectForAgent 单终端: 获取全部进程
func (u *DllController) CollectForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	//添加任务
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeDllSuffixIsDll, "获取全部模块")
	//修改刷新时间
	agent_service.CreateExtendFactory().SetFlushAt(orgId, agentId, false)

	response.SuccessSimplify(ctx)
}

// RepeatStatForAgent 单终端: 重复模块统计
func (u *DllController) RepeatStatForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "hashes")
	hashes := json_helper.ToSliceString(jsonStr)
	if len(hashes) == 0 {
		response.FailSimplify(ctx, "数据不能为空")
		return
	}
	conditions := proc_models.SearchConditions{
		OrgId:     orgId,
		DllHashes: hashes,
	}
	DllFact := proc_models.CreateDllFactory()
	calc := DllFact.CalcByConditionsGroupHash(conditions)
	response.SuccessReturnData(ctx, calc)
}

// RepeatPageForAgent 单终端: 重复进程分页
func (u *DllController) RepeatPageForAgent(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		DllHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateDllFactory().Page(condition, offset, pageSize, order)

	var agentIds []string
	var agentMap = make(map[string]*agent_models.AgentModel)
	if models != nil || len(models) > 0 {
		for _, model := range models {
			agentIds = append(agentIds, model.AgentId)
		}
		AgentFact := agent_models.CreateAgentFactory()
		agentMap = AgentFact.GetMapsByOrgAndAgents(orgId, agentIds, AgentFact.ColumnBasic(true, false))
	}

	response.SuccessReturnData(ctx, (&org_trans.Dll{}).TransRepeatPage(total, models, agentMap))
}
