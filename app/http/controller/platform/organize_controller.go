package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/response"
)

// OrganizeController 终端相关
type OrganizeController struct {
	ControllerBase
}

// Types 企业类型
func (u *OrganizeController) Types(ctx *gin.Context) {
	list := enums.Get1stIsStrReturnMapIntOfStr("type", enum_organization.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatMapIntStr2List(list, true))
}

// PageForManage 管理: 分页
func (u *OrganizeController) PageForManage(ctx *gin.Context) {
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
