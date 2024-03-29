package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
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
	response.SuccessReturnData(ctx, (&platform_trans.Common{}).FormatMapIntStr2List(list, true))
}

// SafeStatus 安全状态
func (u *AgentController) SafeStatus(ctx *gin.Context) {
	list := enums.Get1stIsStrReturnMapIntOfStr("is_safe", enum_agent.EnumIntMap)
	response.SuccessReturnData(ctx, (&platform_trans.Common{}).FormatMapIntStr2List(list, true))
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

// DeleteForManage 管理: 删除
func (u *AgentController) DeleteForManage(ctx *gin.Context) {
	agentIds := []string{
		ctx.GetString(consts.ValidatorPrefixPlatform + "agentId"),
	}
	errCode, err := agent_service.CreateCurdFactory().BatchDelete(agentIds)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// UninstallForManage 管理: 卸载
func (u *AgentController) UninstallForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UninstallForManage")
}

// UpgradeForManage 管理: 升级
func (u *AgentController) UpgradeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UpgradeForManage")
}

// BatchDeleteForManage 管理: 批量删除
func (u *AgentController) BatchDeleteForManage(ctx *gin.Context) {
	agentIds := ctx.GetStringSlice(consts.ValidatorPrefixPlatform + "agentIds")
	errCode, err := agent_service.CreateCurdFactory().BatchDelete(agentIds)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchUninstallForManage 管理: 批量卸载
func (u *AgentController) BatchUninstallForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchUninstallForManage")
}

// BatchUpgradeForManage 管理: 批量升级
func (u *AgentController) BatchUpgradeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchUpgradeForManage")
}

// BatchSetConfigForManage 管理: 批量设置配置
func (u *AgentController) BatchSetConfigForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchSetConfigForManage")
}

// GetForBasic 基本信息：获取
func (u *AgentController) GetForBasic(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetForBasic")
}

// SaveForBasic 基本信息：保存
func (u *AgentController) SaveForBasic(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SaveForBasic")
}

// GetForConfig 配置信息：获取
func (u *AgentController) GetForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetForConfig")
}

// SaveForConfig 配置信息：保存
func (u *AgentController) SaveForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SaveForConfig")
}

// ListForMonitor 监控：列表
func (u *AgentController) ListForMonitor(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ListForMonitor")
}

// InfoForMonitor 监控：信息
func (u *AgentController) InfoForMonitor(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "InfoForMonitor")
}

// PageForError 错误文件：分页
func (u *AgentController) PageForError(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForError")
}
