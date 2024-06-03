package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/service/platform/dictionary_service"
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

// StatByOrg 按企业统计终端
func (u *AgentController) StatByOrg(ctx *gin.Context) {
	searchConditions := agent_models.SearchConditions{
		OrgId: int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId")),
	}

	agentFact := agent_models.CreateAgentFactory()
	stat := agentFact.CalcByConditions(searchConditions, agentFact.ColumnCalcAgents(false, true, false, false, false))
	response.SuccessReturnData(ctx, stat)
}

// PageForOwn 终端本身: 分页
func (u *AgentController) PageForOwn(ctx *gin.Context) {
	//页面上传
	conditions := agent_models.SearchConditions{}
	conditions.OrgId = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "OrgId"))
	conditions.AgentId = ctx.GetString(consts.ValidatorPrefixPlatform + "agentId")
	conditions.HostName = ctx.GetString(consts.ValidatorPrefixPlatform + "hostName")
	conditions.AgentIp = ctx.GetString(consts.ValidatorPrefixPlatform + "agentIp")
	conditions.NetIp = ctx.GetString(consts.ValidatorPrefixPlatform + "netIp")
	conditions.Alias = ctx.GetString(consts.ValidatorPrefixPlatform + "alias")
	conditions.ContactName = ctx.GetString(consts.ValidatorPrefixPlatform + "contactName")
	IsOnline := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isOnline"))
	conditions.IsOnline = &IsOnline
	IsSafe := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isSafe"))
	conditions.IsSafe = &IsSafe
	conditions.DeptIds = u.getWebSliceOfInt64(ctx, "deptId")

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, enum_agent.OrderFieldDefault)

	//终端分页
	total, agentModels := agent_models.CreateAgentFactory().Page(conditions, offset, pageSize, order)
	//转化数据
	_, departmentMaps := u.parseAgentModels(agentModels, true)

	response.SuccessReturnData(ctx, (&platform_trans.Agent{}).TransAgentManagePage(total, agentModels, departmentMaps))
}

// parseAgentModels 解析终端集合并拿到相关数据
// 返回 agentIds 终端集合
// 返回 departments 组织架构集合
func (u *ControllerBase) parseAgentModels(agentModels []*agent_models.AgentModel, needDepartment bool) (agentIds []string, departments map[int64]*department_models.DepartmentModel) {
	if agentModels == nil || len(agentModels) == 0 {
		return nil, nil
	}
	deptStruct := make(map[int64]struct{})
	var deptIds []int64
	var orgId int64
	for _, model := range agentModels {
		agentIds = append(agentIds, model.AgentId)
		orgId = model.OrgId

		if _, isOk := deptStruct[model.DeptId]; model.DeptId > 0 && !isOk {
			deptStruct[model.DeptId] = struct{}{}
			deptIds = append(deptIds, model.DeptId)
		}
	}

	if needDepartment && len(deptIds) > 0 {
		departments = department_service.CreateCurdFactory().GetMapsByIds(orgId, deptIds)
	}

	return agentIds, departments
}

// PageForError 错误文件：分页
func (u *AgentController) PageForError(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForError")
}

// DownloadForError 错误文件：下载
func (u *AgentController) DownloadForError(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DownloadForError")
}
