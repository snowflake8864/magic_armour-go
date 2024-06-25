package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_syslog"
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
	agentId := u.getAgentIdFromGin(ctx)
	if !string_helper.IsEmpty(agentId) && len(userAgentIds) > 0 {
		if slice_helper.StringIn(agentId, userAgentIds) {
			condition.AgentId = agentId
		} else {
			response.FailSimplify(ctx, "数据为空")
			return
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

		break
	case enum_syslog.LogTypeOuter:

		break
	case enum_syslog.LogTypeProcess:

		break
	case enum_syslog.LogTypeLogin:

		break
	default:
		response.FailSimplify(ctx, fmt.Sprintf("module(%d) not exists", module))
	}
}

// BatchDeleteForSyslog 日志: 批量删除
func (u *SyslogController) BatchDeleteForSyslog(ctx *gin.Context, module int) {
	loginInfo := u.getLogin(ctx)

	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "ids")
	ids := json_helper.ToSliceInt64(jsonStr)
	if len(ids) == 0 {
		response.FailSimplify(ctx, "参数错误")
		return
	}

	var err error

	switch module {
	case enum_syslog.LogTypeDns:
		err = syslog_service.CreateDnsCurdFactory().DeleteByIds(loginInfo.OrgId, ids)
		break
	case enum_syslog.LogTypeInner:

		break
	case enum_syslog.LogTypeOuter:

		break
	case enum_syslog.LogTypeProcess:

		break
	case enum_syslog.LogTypeLogin:

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

		break
	case enum_syslog.LogTypeOuter:

		break
	case enum_syslog.LogTypeProcess:

		break
	case enum_syslog.LogTypeLogin:

		break
	}

	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
