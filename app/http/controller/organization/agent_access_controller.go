package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent_access"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_access_models"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

// AgentAccessController  操作日志相关
type AgentAccessController struct {
	ControllerBase
}

// Status  状态
func (u *AgentAccessController) Status(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("status", enum_agent_access.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("status", enum_agent_access.EnumMapInt)
	if cur_userinfo.GetLocales(ctx) == consts.LanguageEn {
		labels = enums.Get1stIsStrReturnMapIntOfStr("status_en", enum_agent_access.EnumMapInt)
	}
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, true))
}

// Urls  路径
func (u *AgentAccessController) Urls(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfStr("path", enum_agent_access.WebMapStr)
	labels := enums.Get1stIsStrReturnMapStrOfStr("path", enum_agent_access.EnumStrMap)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatStringList(keys, labels, true))
}

// PageForLog 日志：分页
func (u *AgentAccessController) PageForLog(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := agent_access_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
		Status:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status")),
		Path:    ctx.GetString(consts.ValidatorPrefixOrg + "path"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := agent_access_models.CreateAccessLogFactory().Page(condition, offset, pageSize, order)

	var agentIds []string
	var agentMap = make(map[string]*agent_models.AgentModel)
	if models != nil || len(models) > 0 {
		for _, model := range models {
			agentIds = append(agentIds, model.AgentId)
		}
		AgentFact := agent_models.CreateAgentFactory()
		agentMap = AgentFact.GetMapsByOrgAndAgents(loginInfo.OrgId, agentIds, AgentFact.ColumnBasic(true, false))
	}

	response.SuccessReturnData(ctx, (&org_trans.AgentAccess{}).TransPage(total, models, agentMap, nil))
}
