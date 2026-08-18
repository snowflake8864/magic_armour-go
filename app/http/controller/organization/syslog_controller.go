package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_syslog"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/syslog_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// SyslogController  通讯日志相关
type SyslogController struct {
	ControllerBase
}

// LoginStatus 登录状态
func (u *SyslogController) LoginStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("login_status", enum_syslog.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("login_status", enum_syslog.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// PageForSyslog 日志: 分页
func (u *SyslogController) PageForSyslog(ctx *gin.Context, module int) {
	loginInfo := u.getLogin(ctx)
	condition := syslog_models.SearchConditions{
		OrgId:       loginInfo.OrgId,
		CreateAtGte: ctx.GetString(consts.ValidatorPrefixOrg + "startTime"),
		CreateAtLte: ctx.GetString(consts.ValidatorPrefixOrg + "endTime"),
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
				response.FailSimplify(ctx, my_errors.CurdDataIsEmpty)
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

	var agentMap = make(map[string]*agent_models.AgentModel)
	var agentIds []string
	switch module {
	case enum_syslog.LogTypeDns:
		condition.ProcHash = ctx.GetString(consts.ValidatorPrefixOrg + "procHash")
		condition.DomainName = ctx.GetString(consts.ValidatorPrefixOrg + "domainName")

		fact := syslog_models.CreateDnsFactory()
		total, models := fact.Page(condition, offset, pageSize, order)
		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.SyslogTrans{}).TransDnsPage(total, models, agentMap, nil))
		break
	case enum_syslog.LogTypeInner:
		condition.ProcHash = ctx.GetString(consts.ValidatorPrefixOrg + "procHash")
		condition.SourceIp = ctx.GetString(consts.ValidatorPrefixOrg + "sourceIp")
		condition.SourcePort = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "sourcePort"))
		condition.TargetIp = ctx.GetString(consts.ValidatorPrefixOrg + "targetIp")
		condition.TargetPort = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "targetPort"))

		fact := syslog_models.CreateInnerFactory()
		total, models := fact.Page(condition, offset, pageSize, order)
		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.SyslogTrans{}).TransInnerPage(total, models, agentMap, nil))
		break
	case enum_syslog.LogTypeOuter:
		condition.ProcHash = ctx.GetString(consts.ValidatorPrefixOrg + "procHash")
		condition.SourceIp = ctx.GetString(consts.ValidatorPrefixOrg + "sourceIp")
		condition.SourcePort = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "sourcePort"))
		condition.TargetIp = ctx.GetString(consts.ValidatorPrefixOrg + "targetIp")
		condition.TargetPort = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "targetPort"))

		fact := syslog_models.CreateOuterFactory()
		total, models := fact.Page(condition, offset, pageSize, order)
		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.SyslogTrans{}).TransOuterPage(total, models, agentMap, nil))
		break
	case enum_syslog.LogTypeProcess:
		condition.ProcHash = ctx.GetString(consts.ValidatorPrefixOrg + "procHash")
		condition.ProcDir = ctx.GetString(consts.ValidatorPrefixOrg + "procDir")
		condition.ProcId = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "procId"))
		condition.ProcParam = ctx.GetString(consts.ValidatorPrefixOrg + "procParam")
		condition.ParentHash = ctx.GetString(consts.ValidatorPrefixOrg + "parentHash")
		condition.ParentDir = ctx.GetString(consts.ValidatorPrefixOrg + "parentDir")
		condition.ParentId = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "parentId"))
		//condition.ParentParam = ctx.GetString(consts.ValidatorPrefixOrg + "parentParam")

		fact := syslog_models.CreateProcessFactory()
		total, models := fact.Page(condition, offset, pageSize, order)
		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.SyslogTrans{}).TransProcessPage(total, models, agentMap, nil))
		break
	case enum_syslog.LogTypeLogin:
		condition.SourceIp = ctx.GetString(consts.ValidatorPrefixOrg + "sourceIp")
		loginStatus := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "loginStatus"))
		condition.LoginStatus = &loginStatus

		fact := syslog_models.CreateLoginFactory()
		total, models := fact.Page(condition, offset, pageSize, order)
		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.SyslogTrans{}).TransLoginPage(total, models, agentMap, nil))
		break
	default:
		response.FailSimplify(ctx, fmt.Sprintf("module(%d) not exists", module))
	}
}

// BatchDeleteForSyslog 日志: 批量删除
func (u *SyslogController) BatchDeleteForSyslog(ctx *gin.Context, module int) {
	loginInfo := u.getLogin(ctx)
	startTime := ctx.GetString(consts.ValidatorPrefixOrg + "startTime")
	endTime := ctx.GetString(consts.ValidatorPrefixOrg + "endTime")

	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "ids")
	idList := json_helper.ToSliceString(jsonStr)
	ids := []int64{}
	if len(idList) == 0 {
		response.FailSimplify(ctx, "参数错误")
		return
	} else {
		for _, id := range idList {
			ids = append(ids, string_helper.ToInt64(id))
		}
	}

	var err error

	switch module {
	case enum_syslog.LogTypeDns:
		err = syslog_service.CreateDnsCurdFactory().DeleteByIds(loginInfo.OrgId, ids, startTime, endTime)
		break
	case enum_syslog.LogTypeInner:
		err = syslog_service.CreateInnerCurdFactory().DeleteByIds(loginInfo.OrgId, ids, startTime, endTime)
		break
	case enum_syslog.LogTypeOuter:
		err = syslog_service.CreateOuterCurdFactory().DeleteByIds(loginInfo.OrgId, ids, startTime, endTime)
		break
	case enum_syslog.LogTypeProcess:
		err = syslog_service.CreateProcessCurdFactory().DeleteByIds(loginInfo.OrgId, ids, startTime, endTime)
		break
	case enum_syslog.LogTypeLogin:
		err = syslog_service.CreateLoginCurdFactory().DeleteByIds(loginInfo.OrgId, ids, startTime, endTime)
		break
	}

	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ClearForSyslog 告警: 清空
func (u *SyslogController) ClearForSyslog(ctx *gin.Context, module int) {
	loginInfo := u.getLogin(ctx)
	//用户管理的终端
	err, userAgentIds := u.getCurrUserAgentIds(loginInfo)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	switch module {
	case enum_syslog.LogTypeDns:
		err = syslog_service.CreateDnsCurdFactory().ClearByAgents(loginInfo.OrgId, userAgentIds)
		break
	case enum_syslog.LogTypeInner:
		err = syslog_service.CreateInnerCurdFactory().ClearByAgents(loginInfo.OrgId, userAgentIds)
		break
	case enum_syslog.LogTypeOuter:
		err = syslog_service.CreateOuterCurdFactory().ClearByAgents(loginInfo.OrgId, userAgentIds)
		break
	case enum_syslog.LogTypeProcess:
		err = syslog_service.CreateProcessCurdFactory().ClearByAgents(loginInfo.OrgId, userAgentIds)
		break
	case enum_syslog.LogTypeLogin:
		err = syslog_service.CreateLoginCurdFactory().ClearByAgents(loginInfo.OrgId, userAgentIds)
		break
	}

	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
