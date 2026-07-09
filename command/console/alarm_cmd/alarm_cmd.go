package alarm_cmd

import (
	"context"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
)

func CreateCmdFactory() *AlarmCmd {
	return &AlarmCmd{
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type AlarmCmd struct {
	queryClient *producer_kafka.ProducerClient
}

func (crd *AlarmCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

func (crd *AlarmCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnBasic(true, false))
}

// TodayStatistics  今天统计
func (crd *AlarmCmd) TodayStatistics() {
	variable.ZapLog.Info("TodayStatistics start")

	statDay, startAt, endAt, statMonth := date_helper.CalcByDay(0)
	crd.handleStat("today", statDay, startAt, endAt, statMonth)

	variable.ZapLog.Info("TodayStatistics end startAt:" + startAt + ", endAt:" + endAt)
}

// YesterdayStatistics  昨天统计
func (crd *AlarmCmd) YesterdayStatistics() {
	variable.ZapLog.Info("DailyStatistics start")

	statDay, startAt, endAt, statMonth := date_helper.CalcByDay(-1)
	crd.handleStat("yesterday", statDay, startAt, endAt, statMonth)

	variable.ZapLog.Info("DailyStatistics end startAt:" + startAt + ", endAt:" + endAt)
}

func (crd *AlarmCmd) handleStat(sign string, statDay, startAt, endAt string, statMonth int) {
	orgModels := crd.getOrgModels()

	if len(orgModels) > 0 {
		ctx := context.Background()
		now := date_helper.CurDatetime()

		var ceilTs int64
		if sign == "today" {
			ceilTs = date_helper.AddTimestampByDays(0, 0, -2)
		} else {
			ceilTs = date_helper.AddTimestampByDays(0, -1, 0)
		}
		variable.ZapLog.Info(fmt.Sprintf("ceilTs: %s", date_helper.TimestampToDatetime(ceilTs)))

		for _, orgModel := range orgModels {
			if agentModels := crd.getAgentModels(orgModel.Id); len(agentModels) > 0 {
				for _, agentModel := range agentModels {

					//如果终端离线x天不在统计
					if agentModel.IsOnline == enums.IsNot && date_helper.DatetimeToTimestamp(agentModel.LiveAt, variable.DateFormatNow) <= ceilTs {
						continue
					}

					param := alarm_stat_models.StatDailyModel{
						OrgId:     orgModel.Id,
						AgentId:   agentModel.AgentId,
						StatDay:   statDay,
						StatMonth: statMonth,
						CreatedAt: now,
					}

					//进入队列
					crd.queryClient.ProduceCrontab(ctx, &kafka_go.MessageStruct{
						Type:    kafka_go.TypeCrontabAlarmStat,
						AgentId: startAt,
						AgentOS: endAt,
						Data:    json_helper.Encode(&param),
					})
				}
			}
		}
	}
}

// Delete  删除
func (crd *AlarmCmd) Delete() {
	variable.ZapLog.Info("Delete start")

	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentAlarmExpire)
	expire := string_helper.ToInt(val)
	variable.ZapLog.Info(fmt.Sprintf("Delete expire: %d", expire))
	if expire == 1 {
		alarm_models.CreateProcessFactory().TruncateTable()
		alarm_models.CreateDllFactory().TruncateTable()
		alarm_models.CreateTamperFactory().TruncateTable()
		alarm_models.CreateBlackmailFactory().TruncateTable()
		alarm_models.CreatePortFactory().TruncateTable()
		alarm_models.CreatePlugFactory().TruncateTable()
		alarm_models.CreateUsbFactory().TruncateTable()
		alarm_models.CreateOutreachFactory().TruncateTable()
		alarm_models.CreateLoginFactory().TruncateTable()
		alarm_models.CreateSelfFactory().TruncateTable()
		alarm_models.CreateVirusFactory().TruncateTable()
	} else {
		_, _, dayEnd, _ := date_helper.CalcByDay(-expire)
		alarm_models.CreateProcessFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateDllFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateTamperFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateBlackmailFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreatePortFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreatePlugFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateUsbFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateOutreachFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateLoginFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateSelfFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateVirusFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})

	}
	variable.ZapLog.Info("Delete end")
}

// EmailNotice  邮件通知
func (crd *AlarmCmd) EmailNotice() {
	variable.ZapLog.Info("EmailNotice start")

	orgModels := crd.getOrgModels()

	if len(orgModels) == 0 {
		return
	}
	ctx := context.Background()
	for _, orgModel := range orgModels {
		//进入队列
		crd.queryClient.ProduceCommon(ctx, &kafka_go.MessageStruct{
			Type:  kafka_go.TypeCommonEmailAlarmCycle,
			OrgId: orgModel.Id,
		})
	}

	variable.ZapLog.Info("EmailNotice end")
}
