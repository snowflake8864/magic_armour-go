package organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_blackmail"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/blackmail_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// BlackmailController 勒索防护相关相关
type BlackmailController struct {
	ControllerBase
}

// Types 类别
func (u *BlackmailController) Types(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type", enum_blackmail.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_blackmail.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, false))
}

// PageForGroup 文件类型组: 分页
func (u *BlackmailController) PageForGroup(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := blackmail_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := blackmail_models.CreateGroupFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.BlackmailTrans{}).TransGroupPage(total, models))
}

// GetAllForGroup 文件类型组: 全部
func (u *BlackmailController) GetAllForGroup(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := blackmail_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}

	models := blackmail_models.CreateGroupFactory().GetListByConditions(condition)

	response.SuccessReturnData(ctx, (&org_trans.BlackmailTrans{}).TransGroupList(models))
}

// AddForGroup 文件类型组: 添加
func (u *BlackmailController) AddForGroup(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := blackmail_models.GroupModel{
		OrgId:      loginInfo.OrgId,
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		Suffix:     ctx.GetString(consts.ValidatorPrefixOrg + "suffix"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := blackmail_service.CreateGroupCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGroup 文件类型组: 编辑
func (u *BlackmailController) EditForGroup(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := blackmail_models.GroupModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Title:      ctx.GetString(consts.ValidatorPrefixOrg + "title"),
		Suffix:     ctx.GetString(consts.ValidatorPrefixOrg + "suffix"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := blackmail_service.CreateGroupCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGroup 文件类型组: 删除
func (u *BlackmailController) DeleteForGroup(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := blackmail_models.GroupModel{
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := blackmail_service.CreateGroupCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForManage 分页
func (u *BlackmailController) PageForManage(ctx *gin.Context) {

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
		//规则统计数据
		PolicyFact := blackmail_models.CreatePolicyFactory()
		calcPolicy := PolicyFact.CalcByConditionsGroupAgentReturnMap(blackmail_models.SearchConditions{
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
func (u *BlackmailController) SetSwitchForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		BlackmailSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().UpdateByAgent(loginInfo.OrgId, agentId, &param, "Manage Page Switch")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetModeForManage 管理页: 设置模式
func (u *BlackmailController) SetModeForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))
	var param = agent_service.ManualColumn{
		BlackmailMode: Mode,
	}

	errCode, err := agent_service.CreateCurdFactory().UpdateByAgent(loginInfo.OrgId, agentId, &param, "Manage Page Mode")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *BlackmailController) BatchSetConfigForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	Mode := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "mode"))

	var param = agent_service.ManualColumn{
		BlackmailSwitch: &Switch,
		BlackmailMode:   Mode,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "Manage Page Batch Modify")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForPolicy 终端策略: 分页
func (u *BlackmailController) PageForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := blackmail_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := blackmail_models.CreatePolicyFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.BlackmailTrans{}).TransPolicyPage(total, models))
}

// SetSwitchForPolicy 终端策略: 设置开关
func (u *BlackmailController) SetSwitchForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := blackmail_models.PolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		OpenSwitch: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := blackmail_service.CreatePolicyCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

func (u *BlackmailController) getProcesses(ctx *gin.Context) (error, string) {
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "processes")

	var processes []tamper_models.Process
	err := json_helper.ToAny(jsonStr, &processes)
	if err != nil {
		return errors.New(my_errors.CurdTrustProcessDataError), ""
	}
	return nil, json_helper.Encode(processes)
}

// AddForPolicy 终端策略: 添加
func (u *BlackmailController) AddForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := blackmail_models.PolicyModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Type:       int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		TypeVal:    ctx.GetString(consts.ValidatorPrefixOrg + "typeVal"),
		GroupId:    string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "groupId")),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else {
		param.Processes = processes
	}

	if err := blackmail_service.CreatePolicyCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForPolicy 终端策略: 编辑
func (u *BlackmailController) EditForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := blackmail_models.PolicyModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Type:       int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		TypeVal:    ctx.GetString(consts.ValidatorPrefixOrg + "typeVal"),
		GroupId:    string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "groupId")),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else {
		param.Processes = processes
	}
	if err := blackmail_service.CreatePolicyCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForPolicy 终端策略: 删除
func (u *BlackmailController) DeleteForPolicy(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := blackmail_models.PolicyModel{
		AgentId:    u.getAgentIdFromGin(ctx),
		OrgId:      loginInfo.OrgId,
		Id:         string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	if err := blackmail_service.CreatePolicyCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGlobal 全局策略: 分页
func (u *BlackmailController) PageForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := blackmail_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := blackmail_models.CreateGlobalFactory().Page(condition, offset, pageSize, order)

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

	response.SuccessReturnData(ctx, (&org_trans.BlackmailTrans{}).TransGlobalPage(total, models, departmentMaps))
}

// SetSwitchForGlobal 全局策略: 设置开关
func (u *BlackmailController) SetSwitchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := blackmail_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		Id:           string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		OpenSwitch:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}

	if err := blackmail_service.CreateGlobalCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddForGlobal 全局策略: 添加
func (u *BlackmailController) AddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := blackmail_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		Type:         int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		TypeVal:      ctx.GetString(consts.ValidatorPrefixOrg + "typeVal"),
		GroupId:      string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "groupId")),
		Remark:       ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		DeptId:       int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		CreateId:     loginInfo.UserId,
		CreateName:   loginInfo.Username,
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else {
		param.Processes = processes
	}
	if err := blackmail_service.CreateGlobalCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGlobal 全局策略: 编辑
func (u *BlackmailController) EditForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := blackmail_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		Id:           string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		Type:         int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		TypeVal:      ctx.GetString(consts.ValidatorPrefixOrg + "typeVal"),
		GroupId:      string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "groupId")),
		Remark:       ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		DeptId:       int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}
	if err, processes := u.getProcesses(ctx); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else {
		param.Processes = processes
	}
	if err := blackmail_service.CreateGlobalCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局策略: 删除
func (u *BlackmailController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := blackmail_models.GlobalModel{
		OrgId:        loginInfo.OrgId,
		Id:           string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixOrg + "id")),
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		UpdateId:     loginInfo.UserId,
		UpdateName:   loginInfo.Username,
	}

	if err := blackmail_service.CreateGlobalCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
