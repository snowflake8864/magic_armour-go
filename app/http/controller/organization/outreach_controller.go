package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/outreach_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// OutreachController 外联管控相关
type OutreachController struct {
	ControllerBase
}

// PageForManage 分页
func (u *OutreachController) PageForManage(ctx *gin.Context) {

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

// SetSwitchForManage 管理页: 设置开关
func (u *OutreachController) SetSwitchForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		OutreachAlarmSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().UpdateByAgent(orgId, agentId, &param, "外联管控管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *OutreachController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))

	var param = agent_service.ManualColumn{
		OutreachAlarmSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "外联管控管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGlobal 全局策略: 分页
func (u *OutreachController) PageForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := outreach_models.SearchConditions{
		OrgId:  loginInfo.OrgId,
		IpLike: ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := outreach_models.CreateGlobalFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.OutreachTrans{}).TransGlobalPage(total, models))
}

// SetSwitchForGlobal 全局策略: 设置开关
func (u *OutreachController) SetSwitchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := outreach_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := outreach_service.CreateGlobalFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// FastAddForGlobal 全局策略: 快速添加
func (u *OutreachController) FastAddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	if err := outreach_service.CreateGlobalFactory().AddDefault(loginInfo.OrgId, loginInfo.UserId, loginInfo.Username); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddForGlobal 全局策略: 添加
func (u *OutreachController) AddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := outreach_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err := outreach_service.CreateGlobalFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGlobal 全局策略: 编辑
func (u *OutreachController) EditForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := outreach_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := outreach_service.CreateGlobalFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局策略: 删除
func (u *OutreachController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := outreach_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := outreach_service.CreateGlobalFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ListForDetect 主动探测: 列表
func (u *OutreachController) ListForDetect(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := outreach_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}

	models := outreach_models.CreateDetectFactory().GetListByConditions(condition)

	response.SuccessReturnData(ctx, (&org_trans.OutreachTrans{}).TransDetectList(models))
}

// SetSwitchForDetect 主动探测: 设置开关
func (u *OutreachController) SetSwitchForDetect(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := outreach_models.DetectModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := outreach_service.CreateDetectFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddForDetect 主动探测: 添加
func (u *OutreachController) AddForDetect(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := outreach_models.DetectModel{
		OrgId:      loginInfo.OrgId,
		Addr:       ctx.GetString(consts.ValidatorPrefixOrg + "addr"),
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err := outreach_service.CreateDetectFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForDetect 主动探测: 编辑
func (u *OutreachController) EditForDetect(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := outreach_models.DetectModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Addr:       ctx.GetString(consts.ValidatorPrefixOrg + "addr"),
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := outreach_service.CreateDetectFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForDetect 主动探测: 删除
func (u *OutreachController) DeleteForDetect(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := outreach_models.DetectModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := outreach_service.CreateDetectFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
