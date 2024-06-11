package organization

import (
	"errors"
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
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// TamperController 防篡改相关
type TamperController struct {
	ControllerBase
}

// FilePower 文件权限
func (u *TamperController) FilePower(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("power_file", enum_tamper.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("power_folder", enum_tamper.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// FolderPower 文件夹权限
func (u *TamperController) FolderPower(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("power_folder", enum_tamper.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("power_folder", enum_tamper.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
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
		PolicyFact := tamper_models.CreatePolicyFactory()
		extendInfo = PolicyFact.CalcByConditionsGroupAgentReturnMap(tamper_models.SearchConditions{
			OrgId:    u.getLogin(ctx).OrgId,
			AgentIds: agentIds,
		}, PolicyFact.ColumnCalc(true))
	}
	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// SetSwitchForManage 管理页: 设置开关
func (u *TamperController) SetSwitchForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		TamperSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "防篡改管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetModeForManage 管理页: 设置模式
func (u *TamperController) SetModeForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentId := u.getAgentIdFromGin(ctx)
	agentIds := []string{agentId}
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))
	var param = agent_service.ManualColumn{
		TamperMode: Mode,
	}

	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "防篡改管理-模式")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *TamperController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))

	var param = agent_service.ManualColumn{
		TamperSwitch: &Switch,
		TamperMode:   Mode,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "防篡改管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForPolicy 终端策略: 分页
func (u *TamperController) PageForPolicy(ctx *gin.Context) {
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

// SetSwitchForPolicy 终端策略: 设置开关
func (u *TamperController) SetSwitchForPolicy(ctx *gin.Context) {
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

func (u *TamperController) getProtectRw(ctx *gin.Context) (error, int) {
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "protectRw")
	protectRw := json_helper.ToSliceInt(jsonStr)
	if protectRw == nil || len(protectRw) == 0 {
		return errors.New("禁止目录权限不能为空"), 0
	}
	var result int
	for _, val := range protectRw {
		result += val
	}

	return nil, result
}

func (u *TamperController) getProcesses(ctx *gin.Context) (error, string) {
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "processes")
	if string_helper.IsEmpty(jsonStr) {
		return nil, ""
	}

	var processes []tamper_models.Process
	err := json_helper.ToAny(jsonStr, &processes)
	if err != nil {
		return errors.New("信任进程数据数据错误"), ""
	}
	return nil, json_helper.Encode(processes)
}

// AddFileForPolicy 终端策略: 添加
func (u *TamperController) AddFileForPolicy(ctx *gin.Context) {
	if err := u.addForPolicy(ctx, enum_tamper.TypeFile); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

func (u *TamperController) addForPolicy(ctx *gin.Context, Type int) error {
	loginInfo := u.getLogin(ctx)
	param := tamper_models.PolicyModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Type:       Type,
		Dir:        ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err, protectRw := u.getProtectRw(ctx); err != nil {
		return err
	} else {
		param.ProtectRw = protectRw
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		return err
	} else {
		param.Processes = processes
	}

	if Type == enum_tamper.TypeFolder {
		param.IncludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "includeSuffix")
		param.ExcludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "excludeSuffix")
		param.ExcludeFolder = ctx.GetString(consts.ValidatorPrefixOrg + "excludeFolder")
		param.IsExtend = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isExtend"))
	}

	if err := tamper_service.CreatePolicyCurdFactory().Add(&param); err != nil {
		return err
	}
	return nil
}

// AddFolderForPolicy 终端策略: 添加
func (u *TamperController) AddFolderForPolicy(ctx *gin.Context) {
	if err := u.addForPolicy(ctx, enum_tamper.TypeFolder); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditFileForPolicy 终端策略: 编辑
func (u *TamperController) EditFileForPolicy(ctx *gin.Context) {
	if err := u.editForPolicy(ctx, enum_tamper.TypeFile); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

func (u *TamperController) editForPolicy(ctx *gin.Context, Type int) error {
	loginInfo := u.getLogin(ctx)
	param := tamper_models.PolicyModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err, protectRw := u.getProtectRw(ctx); err != nil {
		return err
	} else {
		param.ProtectRw = protectRw
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		return err
	} else {
		param.Processes = processes
	}

	if Type == enum_tamper.TypeFolder {
		param.IncludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "includeSuffix")
		param.ExcludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "excludeSuffix")
		param.ExcludeFolder = ctx.GetString(consts.ValidatorPrefixOrg + "excludeFolder")
		param.IsExtend = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isExtend"))
	}

	if err := tamper_service.CreatePolicyCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		return err
	}
	return nil
}

// EditFolderForPolicy 终端策略: 编辑
func (u *TamperController) EditFolderForPolicy(ctx *gin.Context) {
	if err := u.editForPolicy(ctx, enum_tamper.TypeFolder); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForPolicy 终端策略: 删除
func (u *TamperController) DeleteForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.PolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := tamper_service.CreatePolicyCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGlobal 全局策略: 分页
func (u *TamperController) PageForGlobal(ctx *gin.Context) {
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

// SetSwitchForGlobal 全局策略: 设置开关
func (u *TamperController) SetSwitchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		OpenSwitch:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}

	if err := tamper_service.CreateGlobalCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddFileForGlobal 全局策略: 添加
func (u *TamperController) AddFileForGlobal(ctx *gin.Context) {
	if err := u.addForGlobal(ctx, enum_tamper.TypeFile); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

func (u *TamperController) addForGlobal(ctx *gin.Context, Type int) error {
	loginInfo := u.getLogin(ctx)
	param := tamper_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		Type:         Type,
		Dir:          ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Remark:       ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:     loginInfo.UserId,
		CreateName:   loginInfo.Username,
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}
	if err, protectRw := u.getProtectRw(ctx); err != nil {
		return err
	} else {
		param.ProtectRw = protectRw
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		return err
	} else {
		param.Processes = processes
	}

	if Type == enum_tamper.TypeFolder {
		param.IncludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "includeSuffix")
		param.ExcludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "excludeSuffix")
		param.ExcludeFolder = ctx.GetString(consts.ValidatorPrefixOrg + "excludeFolder")
		param.IsExtend = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isExtend"))
	}

	if err := tamper_service.CreateGlobalCurdFactory().Add(&param); err != nil {
		return err
	}
	return nil
}

// AddFolderForGlobal 全局策略: 添加
func (u *TamperController) AddFolderForGlobal(ctx *gin.Context) {
	if err := u.addForGlobal(ctx, enum_tamper.TypeFolder); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGlobal 全局策略: 编辑
func (u *TamperController) EditForGlobal(ctx *gin.Context) {
	if err := u.editForGlobal(ctx, enum_tamper.TypeFile); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

func (u *TamperController) editForGlobal(ctx *gin.Context, Type int) error {
	loginInfo := u.getLogin(ctx)
	param := tamper_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		Id:           int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		DeptId:       int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		Remark:       ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}
	if err, protectRw := u.getProtectRw(ctx); err != nil {
		return err
	} else {
		param.ProtectRw = protectRw
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		return err
	} else {
		param.Processes = processes
	}

	if Type == enum_tamper.TypeFolder {
		param.IncludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "includeSuffix")
		param.ExcludeSuffix = ctx.GetString(consts.ValidatorPrefixOrg + "excludeSuffix")
		param.ExcludeFolder = ctx.GetString(consts.ValidatorPrefixOrg + "excludeFolder")
		param.IsExtend = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isExtend"))
	}

	if err := tamper_service.CreateGlobalCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		return err
	}
	return nil
}

// EditFolderForGlobal 全局策略: 编辑
func (u *TamperController) EditFolderForGlobal(ctx *gin.Context) {
	if err := u.editForGlobal(ctx, enum_tamper.TypeFolder); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局策略: 删除
func (u *TamperController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := tamper_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
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
