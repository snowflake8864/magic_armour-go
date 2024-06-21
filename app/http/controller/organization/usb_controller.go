package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/enums/enum_usb"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/agent/usb_service"
	"goskeleton/app/utils/response"
)

// UsbController 防篡改相关
type UsbController struct {
	ControllerBase
}

// EditStatus 编辑状态
func (u *UsbController) EditStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("edit_status", enum_usb.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("status", enum_usb.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PageForManage 分页
func (u *UsbController) PageForManage(ctx *gin.Context) {

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
		PolicyFact := usb_models.CreateUsbFactory()
		calcPolicy := PolicyFact.CalcByConditionsGroupAgentReturnMap(usb_models.SearchConditions{
			OrgId:    u.getLogin(ctx).OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc())

		for _, agentModel := range agentModels {
			gg := map[string]interface{}{
				"total_number": 0,
				"on_number":    0,
			}
			if item, isOk := calcPolicy[agentModel.AgentId]; isOk {
				gg["total_number"] = item.TotalNumber
				gg["white_number"] = item.WhiteNumber
				gg["black_number"] = item.BlackNumber
				gg["unset_number"] = item.UnsetNumber

			}
			extendInfo[agentModel.AgentId] = gg
		}
	}
	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// SetSwitchForManage 管理页: 设置开关
func (u *UsbController) SetSwitchForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		UsbSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "外设管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetModeForManage 管理页: 设置模式
func (u *UsbController) SetModeForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	agentIds := []string{agentId}
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))
	var param = agent_service.ManualColumn{
		UsbMode: Mode,
	}

	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "外设管理-模式")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *UsbController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))

	var param = agent_service.ManualColumn{
		UsbSwitch: &Switch,
		UsbMode:   Mode,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "外设管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForLog 外设插拔日志:分页
func (u *UsbController) PageForLog(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := usb_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
		UsbHash: ctx.GetString(consts.ValidatorPrefixOrg + "id"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := usb_models.CreateLogFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Usb{}).TransLogPage(total, models))
}

// HasUploadForPolicy 终端策略: 是否上传
func (u *UsbController) HasUploadForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	if !task_service.CreateCacheFactory().ExistUnfinished(loginInfo.OrgId, agentId, enum_task.TypeUsbUpload) {
		response.FailSimplify(ctx, "请继续等待")
		return
	}
	response.SuccessSimplify(ctx)
}

// RefreshForPolicy 终端策略: 刷新
func (u *UsbController) RefreshForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	task_service.CreateCurdFactory().AddTaskForAgent(loginInfo.OrgId, agentId, enum_task.TypeUsbUpload, "刷新")
	response.SuccessSimplify(ctx)
}

// PageForPolicy 终端策略: 分页
func (u *UsbController) PageForPolicy(ctx *gin.Context, isOnline int) {
	loginInfo := u.getLogin(ctx)

	condition := usb_models.SearchConditions{
		OrgId:    loginInfo.OrgId,
		AgentId:  u.getAgentIdFromGin(ctx),
		IsOnline: &isOnline,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := usb_models.CreateUsbFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Usb{}).TransPolicyPage(total, models))
}

// SetStatusForPolicy 终端策略: 设置状态
func (u *UsbController) SetStatusForPolicy(ctx *gin.Context, status int) {
	loginInfo := u.getLogin(ctx)

	param := usb_models.UsbModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		Status:     status,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := usb_service.CreateUsbCurdFactory().SetStatus(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetAliasForPolicy 终端策略: 别名
func (u *UsbController) SetAliasForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := usb_models.UsbModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		Alias:      ctx.GetString(consts.ValidatorPrefixOrg + "alias"),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := usb_service.CreateUsbCurdFactory().Edit(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGlobal 全局策略: 分页
func (u *UsbController) PageForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := usb_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := usb_models.CreateGlobalFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Usb{}).TransGlobalPage(total, models))
}

// AddForGlobal 全局策略: 添加
func (u *UsbController) AddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := usb_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		Status:     int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status")),
		Name:       ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		Type:       ctx.GetString(consts.ValidatorPrefixOrg + "type"),
		Intro:      ctx.GetString(consts.ValidatorPrefixOrg + "intro"),
		Alias:      ctx.GetString(consts.ValidatorPrefixOrg + "alias"),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := usb_service.CreateGlobalCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetStatusForGlobal 终端策略: 设置状态
func (u *UsbController) SetStatusForGlobal(ctx *gin.Context, status int) {
	loginInfo := u.getLogin(ctx)

	param := usb_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		Status:     status,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := usb_service.CreateGlobalCurdFactory().SetStatus(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetAliasForGlobal 终端策略: 别名
func (u *UsbController) SetAliasForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := usb_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		Alias:      ctx.GetString(consts.ValidatorPrefixOrg + "alias"),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := usb_service.CreateGlobalCurdFactory().Edit(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局策略: 删除
func (u *UsbController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := usb_models.GlobalModel{
		OrgId:      loginInfo.OrgId,
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := usb_service.CreateGlobalCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
