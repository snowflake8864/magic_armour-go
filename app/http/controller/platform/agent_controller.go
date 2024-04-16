package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/response"
)

// AgentController 终端相关
type AgentController struct {
	ControllerBase
}

// OnlineStatus 在线状态
func (u *AgentController) OnlineStatus(ctx *gin.Context) {
	list := enums.Get1stIsStrReturnMapIntOfStr("is_online", enum_agent.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatMapIntStr2List(list, true))
}

// SafeStatus 安全状态
func (u *AgentController) SafeStatus(ctx *gin.Context) {
	list := enums.Get1stIsStrReturnMapIntOfStr("is_safe", enum_agent.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatMapIntStr2List(list, true))
}

// CountOnline 统计在线和离线数量
func (u *AgentController) CountOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "CountOnline")
}

// TypeList 类型列表
func (u *AgentController) TypeList(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "TypeList")
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
	offset, pageSize, orderField, orderSort := u.getPageAndSort(ctx, enum_agent.OrderFieldDefault)

	//终端分页
	total, agentModels := agent_service.CreateCurdFactory().Page(conditions, offset, pageSize, orderField, orderSort)
	//转化数据
	_, departmentMaps := u.parseAgentModels(agentModels, true)
	//软件最新发布
	softwareNewest := software_service.CreateCurdFactory().GetMapsByPublished()

	response.SuccessReturnData(ctx, (&platform_trans.Agent{}).TransAgentManagePage(total, agentModels, departmentMaps, softwareNewest))
}
