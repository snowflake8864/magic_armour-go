package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/plug_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// PlugController 端口相关
type PlugController struct {
	ControllerBase
}

// PlugTypes 封堵类型
func (u *PlugController) PlugTypes(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("plug_type", enum_plug.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("plug_type", enum_plug.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PageForManage 分页
func (u *PlugController) PageForManage(ctx *gin.Context) {

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
	agentIds, departmentMaps := u.parseAgentModels(agentModels, true)

	var extendInfo = make(map[string]interface{})
	if len(agentIds) > 0 {
		loginInfo := u.getLogin(ctx)

		//规则统计数据
		PolicyFact := plug_models.CreateDynamicFactory()
		calcPolicy := PolicyFact.CalcByConditionsGroupAgentReturnMap(plug_models.SearchConditions{
			OrgId:    loginInfo.OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc(true))

		IpPolicyFact := plug_models.CreateIpPolicyFactory()
		calcWhite := IpPolicyFact.CalcByConditionsGroupAgentReturnMap(plug_models.SearchConditions{
			OrgId:    loginInfo.OrgId,
			AgentIds: agentIds,
			IpType:   enum_plug.IpTypeWhite,
		}, IpPolicyFact.ColumnCalc(false))

		calcBlack := IpPolicyFact.CalcByConditionsGroupAgentReturnMap(plug_models.SearchConditions{
			OrgId:    loginInfo.OrgId,
			AgentIds: agentIds,
			IpType:   enum_plug.IpTypeBlack,
		}, IpPolicyFact.ColumnCalc(false))

		for _, agentModel := range agentModels {
			gg := map[string]interface{}{
				"dynamic_total_number":  0,
				"dynamic_on_number":     0,
				"ip_white_total_number": 0,
				"ip_black_total_number": 0,
			}

			if item, isOk := calcPolicy[agentModel.AgentId]; isOk {
				gg["dynamic_total_number"] = item.TotalNumber
				gg["dynamic_on_number"] = item.OnNumber

			}
			if item, isOk := calcWhite[agentModel.AgentId]; isOk {
				gg["ip_white_total_number"] = item.TotalNumber
			}
			if item, isOk := calcBlack[agentModel.AgentId]; isOk {
				gg["ip_black_total_number"] = item.TotalNumber
			}
			extendInfo[agentModel.AgentId] = gg
		}

	}

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// SetSwitchForManage 管理页: 设置开关
func (u *PlugController) SetSwitchForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		PlugSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "动态封堵管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *PlugController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))

	var param = agent_service.ManualColumn{
		PlugSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "动态封堵管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigPortGetForPolicy 终端策略 - 端口虚开
func (u *PlugController) ConfigPortGetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := agent_service.CreateExtendFactory().GetPlugOfPort(loginInfo.OrgId, agentId)

	response.SuccessReturnData(ctx, config)
}

// ConfigPortSetForPolicy 终端策略 - 端口虚开
func (u *PlugController) ConfigPortSetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := enum_plug.PlugPortConfig{
		Switch:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Weight:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "weight")),
		Type:      int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		CountTime: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "countTime")),
		PlugTime:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "plugTime")),
	}
	if err := agent_service.CreateExtendFactory().SetPlugOfPort(loginInfo.OrgId, agentId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigLoginGetForPolicy 终端策略 - 口令破解
func (u *PlugController) ConfigLoginGetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := agent_service.CreateExtendFactory().GetPlugOfLogin(loginInfo.OrgId, agentId)

	response.SuccessReturnData(ctx, config)
}

// ConfigLoginSetForPolicy 终端策略 - 口令破解
func (u *PlugController) ConfigLoginSetForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	config := enum_plug.PlugLoginConfig{
		Switch:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Type:     int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		PlugTime: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "plugTime")),
	}
	if err := agent_service.CreateExtendFactory().SetPlugOfLogin(loginInfo.OrgId, agentId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpPageForPolicy 终端策略 - ip名单库: 分页
func (u *PlugController) IpPageForPolicy(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)
	condition := plug_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
		IpType:  ipType,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := plug_models.CreateIpPolicyFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.PlugTrans{}).TranIpPolicyPage(total, models))
}

// IpSetSwitchForPolicy 终端策略 - ip名单库: 设置开关
func (u *PlugController) IpSetSwitchForPolicy(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.IpPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err := plug_service.CreatePolicyIpCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpAddForPolicy 终端策略 - ip名单库: 添加
func (u *PlugController) IpAddForPolicy(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)
	param := plug_models.IpPolicyModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Type:       ipType,
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyIpCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpEditForPolicy 终端策略 - ip名单库: 编辑
func (u *PlugController) IpEditForPolicy(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)
	param := plug_models.IpPolicyModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Type:       ipType,
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyIpCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpDeleteForPolicy 终端策略 - ip名单库: 删除
func (u *PlugController) IpDeleteForPolicy(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.IpPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Type:       ipType,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyIpCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpClearForPolicy 终端策略 - ip名单库: 删除
func (u *PlugController) IpClearForPolicy(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.IpPolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Type:       ipType,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyIpCurdFactory().DeleteByAgentAndType(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicPageForPolicy 终端策略 - 动态封堵: 分页
func (u *PlugController) DynamicPageForPolicy(ctx *gin.Context, openSwitch int) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := plug_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		OpenSwitch: &openSwitch,
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "plugIp"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := plug_models.CreateDynamicFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.PlugTrans{}).TranDynamicPolicyPage(total, models))
}

// DynamicDisableForPolicy 终端策略 - 动态封堵: 解封
func (u *PlugController) DynamicDisableForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.DynamicModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DisableById(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicDisableBatchForPolicy 终端策略 - 动态封堵: 批量解封
func (u *PlugController) DynamicDisableBatchForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.DynamicModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DisableByAgentAndIds(&param, u.getIdsOfInt64(ctx)); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicDeleteForPolicy 终端策略 - 动态封堵: 删除
func (u *PlugController) DynamicDeleteForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.DynamicModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DeleteById(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicDeleteBatchForPolicy 终端策略 - 动态封堵: 批量删除
func (u *PlugController) DynamicDeleteBatchForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.DynamicModel{
		OrgId:      loginInfo.OrgId,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DeleteByAgentsAndIds(&param, u.getIdsOfInt64(ctx), u.getAgentIdsFromGin(ctx)); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ConfigPortGetForGlobal 全局策略 - 端口虚开
func (u *PlugController) ConfigPortGetForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	model := OrgExtendFact.GetItemByOrg(loginInfo.OrgId)
	response.SuccessReturnData(ctx, OrgExtendFact.FormatPlugPortConfig(model.PlugPortConfig))
}

// ConfigPortSetForGlobal 全局策略 - 端口虚开
func (u *PlugController) ConfigPortSetForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	config := enum_plug.PlugPortConfig{
		Switch:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Weight:    int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "weight")),
		Type:      int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		CountTime: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "countTime")),
		PlugTime:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "plugTime")),
	}
	if err := organize_service.CreateOrgExtendFactory().SetPlugOfPort(loginInfo.OrgId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//用户管理的终端
	if err, agentIds := u.getCurrUserAgentIds(loginInfo); err == nil {
		_ = agent_service.CreateExtendFactory().BatchSetPlugOfPort(loginInfo.OrgId, agentIds, config)
	}

	response.SuccessSimplify(ctx)
}

// ConfigLoginGetForGlobal 全局策略 - 口令破解
func (u *PlugController) ConfigLoginGetForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	model := OrgExtendFact.GetItemByOrg(loginInfo.OrgId)
	response.SuccessReturnData(ctx, OrgExtendFact.FormatPlugLoginConfig(model.PlugLoginConfig))
}

// ConfigLoginSetForGlobal 全局策略 - 口令破解
func (u *PlugController) ConfigLoginSetForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	config := enum_plug.PlugLoginConfig{
		Switch:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		Type:     int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		PlugTime: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "plugTime")),
	}
	if err := organize_service.CreateOrgExtendFactory().SetPlugOfLogin(loginInfo.OrgId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//用户管理的终端
	if err, agentIds := u.getCurrUserAgentIds(loginInfo); err == nil {
		_ = agent_service.CreateExtendFactory().BatchSetPlugOfLogin(loginInfo.OrgId, agentIds, config)
	}

	response.SuccessSimplify(ctx)
}

// IpPageForGlobal 全局策略 - ip名单库: 分页
func (u *PlugController) IpPageForGlobal(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)
	condition := plug_models.SearchConditions{
		OrgId:  loginInfo.OrgId,
		IpType: ipType,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := plug_models.CreateIpGlobalFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.PlugTrans{}).TranIpGlobalPage(total, models))
}

// IpSetSwitchForGlobal 全局策略 - ip名单库: 设置开关
func (u *PlugController) IpSetSwitchForGlobal(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.IpGlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err := plug_service.CreateGlobalIpCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpAddForGlobal 全局策略 - ip名单库: 添加
func (u *PlugController) IpAddForGlobal(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)
	param := plug_models.IpGlobalModel{
		OrgId:      loginInfo.OrgId,
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Type:       ipType,
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreateGlobalIpCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpEditForGlobal 全局策略 - ip名单库: 编辑
func (u *PlugController) IpEditForGlobal(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)
	param := plug_models.IpGlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		Type:       ipType,
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreateGlobalIpCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpDeleteForGlobal 全局策略 - ip名单库: 删除
func (u *PlugController) IpDeleteForGlobal(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.IpGlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreateGlobalIpCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// IpClearForGlobal 全局策略 - ip名单库: 删除
func (u *PlugController) IpClearForGlobal(ctx *gin.Context, ipType int) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.IpPolicyModel{
		OrgId:      loginInfo.OrgId,
		Type:       ipType,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreateGlobalIpCurdFactory().DeleteByAgentAndType(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicPageForGlobal 全局策略 - 动态封堵: 分页
func (u *PlugController) DynamicPageForGlobal(ctx *gin.Context, openSwitch int) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := plug_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		OpenSwitch: &openSwitch,
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "plugIp"),
	}
	//用户管理的终端
	err, userAgentIds := u.getCurrUserAgentIds(loginInfo)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	//查询终端比较用户管理的终端
	if agentId := u.getAgentIdFromGin(ctx); !string_helper.IsEmpty(agentId) { //存在查询终端
		if len(userAgentIds) > 0 { //不是管理员
			if slice_helper.StringIn(agentId, userAgentIds) { // 查询终端在管理终端列表下
				condition.AgentId = agentId
			} else {
				response.FailSimplify(ctx, "数据为空")
				return
			}
		} else { //是管理员
			condition.AgentId = agentId
		}
	} else {
		condition.AgentIds = userAgentIds
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := plug_models.CreateDynamicFactory().Page(condition, offset, pageSize, order)

	var agentMap = make(map[string]*agent_models.AgentModel)
	var agentIds []string
	if models != nil || len(models) > 0 {
		for _, model := range models {
			agentIds = append(agentIds, model.AgentId)
		}
		AgentFact := agent_models.CreateAgentFactory()
		agentMap = AgentFact.GetMapsByOrgAndAgents(loginInfo.OrgId, agentIds, AgentFact.ColumnBasic(true, false))
	}

	response.SuccessReturnData(ctx, (&org_trans.PlugTrans{}).TranDynamicGlobalPage(total, models, agentMap, nil))
}

// DynamicDisableForGlobal 全局策略 - 动态封堵: 解封
func (u *PlugController) DynamicDisableForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.DynamicModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DisableById(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicDisableBatchForGlobal 全局策略 - 动态封堵: 批量解封
func (u *PlugController) DynamicDisableBatchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentIds := u.getAgentIdsFromGin(ctx)
	ids := u.getIdsOfInt64(ctx)
	param := plug_models.DynamicModel{
		OrgId:      loginInfo.OrgId,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DisableByAgentsAndIds(&param, ids, agentIds); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicDeleteForGlobal 全局策略 - 动态封堵: 删除
func (u *PlugController) DynamicDeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := plug_models.DynamicModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DeleteById(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DynamicDeleteBatchForGlobal 全局策略 - 动态封堵: 批量删除
func (u *PlugController) DynamicDeleteBatchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentIds := u.getAgentIdsFromGin(ctx)
	ids := u.getIdsOfInt64(ctx)
	param := plug_models.DynamicModel{
		OrgId:      loginInfo.OrgId,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := plug_service.CreatePolicyDynamicCurdFactory().DeleteByAgentsAndIds(&param, ids, agentIds); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
