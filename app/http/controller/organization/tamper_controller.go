package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/utils/response"
)

// TamperController 防篡改相关
type TamperController struct {
	ControllerBase
}

// FilePower 文件权限
func (u *TamperController) FilePower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}

// FolderPower 文件夹权限
func (u *TamperController) FolderPower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FolderPower")
}

// PageForManage 分页
func (u *TamperController) PageForManage(ctx *gin.Context) {

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
		PolicyFact := tamper_models.CreateTamperPolicyFactory()
		extendInfo = PolicyFact.CalcByConditionsGroupAgentReturnMap(tamper_models.SearchConditions{
			OrgId:    u.getLogin(ctx).OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc(true))
	}
	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// SetSwitchForManage 管理页: 设置开关
func (u *TamperController) SetSwitchForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSwitchForManage")
}

// SetModeForManage 管理页: 设置模式
func (u *TamperController) SetModeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetModeForManage")
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *TamperController) BatchSetConfigForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchModifyForManage")
}

// PageForPolicy 终端策略: 分页
func (u *TamperController) PageForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForPolicy")
}

// SetSwitchForPolicy 终端策略: 设置开关
func (u *TamperController) SetSwitchForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSwitchForPolicy")
}

// AddForPolicy 终端策略: 添加
func (u *TamperController) AddForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddForPolicy")
}

// EditForPolicy 终端策略: 编辑
func (u *TamperController) EditForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "EditForPolicy")
}

// DeleteForPolicy 终端策略: 删除
func (u *TamperController) DeleteForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForPolicy")
}

// PageForGlobal 全局策略: 分页
func (u *TamperController) PageForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForGlobal")
}

// SetSwitchForGlobal 全局策略: 设置开关
func (u *TamperController) SetSwitchForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSwitchForGlobal")
}

// AddForGlobal 全局策略: 添加
func (u *TamperController) AddForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddForGlobal")
}

// EditForGlobal 全局策略: 编辑
func (u *TamperController) EditForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "EditForGlobal")
}

// DeleteForGlobal 全局策略: 删除
func (u *TamperController) DeleteForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForGlobal")
}
