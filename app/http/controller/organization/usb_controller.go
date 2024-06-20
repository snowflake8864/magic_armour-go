package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/tamper_service"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/response"
)

// UsbController 防篡改相关
type UsbController struct {
	ControllerBase
}

// FilePower 文件权限
func (u *UsbController) FilePower(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("power_file", enum_tamper.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("power_folder", enum_tamper.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// FolderPower 文件夹权限
func (u *UsbController) FolderPower(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("power_folder", enum_tamper.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("power_folder", enum_tamper.EnumMapInt)
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
		PolicyFact := tamper_models.CreatePolicyFactory()
		calcPolicy := PolicyFact.CalcByConditionsGroupAgentReturnMap(tamper_models.SearchConditions{
			OrgId:    u.getLogin(ctx).OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc(true))

		for _, agentModel := range agentModels {
			gg := map[string]interface{}{
				"total_number": 0,
				"on_number":    0,
			}
			if item, isOk := calcPolicy[agentModel.AgentId]; isOk {
				gg["total_number"] = item.TotalNumber
				gg["on_number"] = item.OnNumber

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

	condition := tamper_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := tamper_models.CreatePolicyFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.TamperTrans{}).TransPolicyPage(total, models))
}

// PageForPolicy 终端策略: 分页
func (u *UsbController) PageForPolicy(ctx *gin.Context, isOnline int) {
	loginInfo := u.getLogin(ctx)

	condition := tamper_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := tamper_models.CreatePolicyFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.TamperTrans{}).TransPolicyPage(total, models))
}

// SetStatusForPolicy 终端策略: 设置状态
func (u *UsbController) SetStatusForPolicy(ctx *gin.Context, status int) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.PolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := tamper_service.CreatePolicyCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetAliasForPolicy 终端策略: 别名
func (u *UsbController) SetAliasForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.PolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := tamper_service.CreatePolicyCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGlobal 全局策略: 分页
func (u *UsbController) PageForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := tamper_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := tamper_models.CreateGlobalFactory().Page(condition, offset, pageSize, order)

	//获取组织架构
	departmentMaps := map[int64]*department_models.DepartmentModel{}
	if models != nil && len(models) > 0 {
		var deptIds []int64
		for _, model := range models {
			if model.DeptId > 0 {
				deptIds = append(deptIds, model.DeptId)
			}
		}
		departmentMaps = u.getDepartmentMap(loginInfo.OrgId, slice_helper.Int64RemoveRepeatByMap(deptIds))
	}

	response.SuccessReturnData(ctx, (&org_trans.TamperTrans{}).TransGlobalPage(total, models, departmentMaps))
}

// AddForGlobal 全局策略: 添加
func (u *UsbController) AddForGlobal(ctx *gin.Context) error {
	loginInfo := u.getLogin(ctx)
	param := tamper_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		DeptId:       int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		Dir:          ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Remark:       ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:     loginInfo.UserId,
		CreateName:   loginInfo.Username,
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}

	if err := tamper_service.CreateGlobalCurdFactory().Add(&param); err != nil {
		return err
	}
	return nil
}

// / SetStatusForGlobal 终端策略: 设置状态
func (u *UsbController) SetStatusForGlobal(ctx *gin.Context, status int) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.PolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := tamper_service.CreatePolicyCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetAliasForGlobal 终端策略: 别名
func (u *UsbController) SetAliasForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.PolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := tamper_service.CreatePolicyCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局策略: 删除
func (u *UsbController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		Id:           int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}

	if err := tamper_service.CreateGlobalCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
