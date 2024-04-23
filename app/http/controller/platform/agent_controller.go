package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/model/model_organization"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/platform/department_service"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
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

// StatByOrg 按企业统计终端
func (u *AgentController) StatByOrg(ctx *gin.Context) {
	searchConditions := agent_service.SearchConditions{
		OrgId:    int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId")),
		IsSafe:   enums.AllDefaultInt,
		IsOnline: enums.AllDefaultInt,
	}

	agentFact := agent_service.CreateCurdFactory()
	stat := agentFact.CalcAgents(searchConditions, agentFact.ColumnCalcAgents(false, true, false, false, false))
	response.SuccessReturnData(ctx, stat)
}

// TypeList 类型列表
func (u *AgentController) TypeList(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "TypeList")
}

// PageForOwn 终端本身: 分页
func (u *AgentController) PageForOwn(ctx *gin.Context) {
	//页面上传
	conditions := agent_service.SearchConditions{}
	conditions.OrgId = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "OrgId"))
	conditions.AgentId = ctx.GetString(consts.ValidatorPrefixPlatform + "agentId")
	conditions.HostName = ctx.GetString(consts.ValidatorPrefixPlatform + "hostName")
	conditions.AgentIp = ctx.GetString(consts.ValidatorPrefixPlatform + "agentIp")
	conditions.NetIp = ctx.GetString(consts.ValidatorPrefixPlatform + "netIp")
	conditions.Alias = ctx.GetString(consts.ValidatorPrefixPlatform + "alias")
	conditions.ContactName = ctx.GetString(consts.ValidatorPrefixPlatform + "contactName")
	conditions.IsOnline = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isOnline"))
	conditions.IsSafe = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isSafe"))

	var reqDeptIds []int64
	if deptIds, isOk := ctx.Get(consts.ValidatorPrefixPlatform + "deptId"); isOk {
		var tmp []interface{}
		tmp, _ = deptIds.([]interface{})
		reqDeptIds = helpers.SliceAny2Int64(tmp)
	}
	conditions.DeptIds = reqDeptIds

	//分页参数和排序参数
	offset, pageSize, orderField, orderSort := u.getPageAndSort(ctx, enum_agent.OrderFieldDefault)

	//终端分页
	total, agentModels := agent_service.CreateCurdFactory().Page(conditions, offset, pageSize, orderField, orderSort)
	//转化数据
	_, departmentMaps := u.parseAgentModels(agentModels, true)

	response.SuccessReturnData(ctx, (&platform_trans.Agent{}).TransAgentManagePage(total, agentModels, departmentMaps))
}

// parseAgentModels 解析终端集合并拿到相关数据
// 返回 agentIds 终端集合
// 返回 departments 组织架构集合
func (u *ControllerBase) parseAgentModels(agentModels []*model_agent.AgentModel, needDepartment bool) (agentIds []string, departments map[int64]*model_organization.DepartmentModel) {
	if agentModels == nil || len(agentModels) == 0 {
		return nil, nil
	}
	deptStruct := make(map[int64]struct{})
	var deptIds []int64
	for _, model := range agentModels {
		agentIds = append(agentIds, model.AgentId)

		if _, isOk := deptStruct[model.DeptId]; model.DeptId > 0 && !isOk {
			deptStruct[model.DeptId] = struct{}{}
			deptIds = append(deptIds, model.DeptId)
		}
	}

	if needDepartment && len(deptIds) > 0 {
		departments = department_service.CreateCurdFactory().GetMapsByIds(deptIds)
	}

	return agentIds, departments
}
