package org_cmd

import (
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/org_stat_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/org_stat_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateCmdFactory() *OrgCmd {
	return &OrgCmd{}
}

type OrgCmd struct {
}

func (crd *OrgCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

// TodayStatistics  今天统计
func (crd *OrgCmd) TodayStatistics() {
	variable.ZapLog.Sugar().Debug("TodayStatistics start")

	statDay, startAt, endAt, statMonth := date_helper.CalcByDay(0)
	crd.handleStat(statDay, startAt, endAt, statMonth)

	variable.ZapLog.Sugar().Debug("TodayStatistics end startAt:" + startAt + ", endAt:" + endAt)
}

// YesterdayStatistics  昨天统计
func (crd *OrgCmd) YesterdayStatistics() {
	variable.ZapLog.Sugar().Debug("DailyStatistics start")

	statDay, startAt, endAt, statMonth := date_helper.CalcByDay(-1)
	crd.handleStat(statDay, startAt, endAt, statMonth)

	variable.ZapLog.Sugar().Debug("DailyStatistics end startAt:" + startAt + ", endAt:" + endAt)
}

func (crd *OrgCmd) handleStat(statDay, startAt, endAt string, statMonth int) {
	orgModels := crd.getOrgModels()

	if len(orgModels) > 0 {
		for _, model := range orgModels {
			param := org_stat_models.DailyModel{
				OrgId:     model.Id,
				StatDay:   statDay,
				StatMonth: statMonth,
			}
			//计算终端
			crd.calcAgent(&param, startAt, endAt)

			//计算用户
			crd.calcUser(&param, startAt, endAt)

			org_stat_service.CreateStatDailyFactory().Add(&param)
		}
	}
}

func (crd *OrgCmd) calcAgent(param *org_stat_models.DailyModel, startAt, endAt string) {
	conditions := agent_models.SearchConditions{
		OrgId: param.OrgId,
	}
	agentFact := agent_models.CreateAgentFactory()

	//总数
	dataTotal := agentFact.CalcByConditions(conditions, agentFact.ColumnCalcAgents(false, true, false, true, true))
	agentFact.FillCalcData(dataTotal)
	param.AgentTotal = dataTotal.TotalNumber
	param.AgentOnline = dataTotal.OnlineNumber
	param.AgentLinuxTotal = dataTotal.LinuxNumber
	param.AgentLinuxOnline = dataTotal.LinuxOnlineNumber
	param.AgentWinTotal = dataTotal.WindowsNumber
	param.AgentWinOnline = dataTotal.WindowsOnlineNumber

	//新增数
	conditions.CreateAtGte = startAt
	conditions.CreateAtLte = endAt
	dataAdd := agentFact.CalcByConditions(conditions, agentFact.ColumnCalcAgents(false, true, false, true, true))
	agentFact.FillCalcData(dataAdd)
	param.AgentAdd = dataAdd.TotalNumber
	param.AgentLinuxAdd = dataAdd.LinuxNumber
	param.AgentWinAdd = dataAdd.LinuxNumber

}

func (crd *OrgCmd) calcUser(param *org_stat_models.DailyModel, startAt, endAt string) {

	userFact := user_models.CreateUserFactory()

	param.UserTotal = int(userFact.Count(user_models.SearchConditions{
		OrgId: param.OrgId,
	}))
	param.UserAdd = int(userFact.Count(user_models.SearchConditions{
		OrgId:       param.OrgId,
		CreateAtGte: startAt,
		CreateAtLte: endAt,
	}))

	param.UserOnline = int(userFact.Count(user_models.SearchConditions{
		OrgId:          param.OrgId,
		LastLoginAtGte: startAt,
		LastLoginAtLte: endAt,
	}))
}
