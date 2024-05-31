package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_outreach"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/ip_helper"
	"strings"
)

func CreateOutreachCalcFactory() *AlarmOutreachCalc {
	return &AlarmOutreachCalc{
		masterModel: alarm_models.CreateAlarmOutreachFactory(),
	}
}

type AlarmOutreachCalc struct {
	alarmBase
	masterModel *alarm_models.AlarmOutreachModel
}

func (crd *AlarmOutreachCalc) inOutreachWhile(orgId int64, targetIp string) bool {
	SwitchOn := enums.SwitchOn
	conditions := outreach_models.SearchConditions{
		OrgId:      orgId,
		Type:       enum_outreach.TypeWhite,
		OpenSwitch: &SwitchOn,
	}
	models := outreach_models.CreateOutreachGlobalFactory().GetListByConditions(conditions)
	if models == nil || len(models) == 0 {
		return true
	}

	ipArr := strings.Split(targetIp, ";")
	for _, ip := range ipArr {
		for _, model := range models {
			if ip_helper.IsIPInRange(ip, model.Ip) {
				return true
			}
		}
	}

	return false
}

// Calc 计算
func (crd *AlarmOutreachCalc) Calc(dnsModel *syslog_models.SyslogDnsModel) error {
	//终端开关
	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(dnsModel.OrgId, dnsModel.AgentId, agentFact.ColumnConfig(false, false))
	if agentModel == nil || agentModel.PlugSwitch == enums.SwitchOff {
		return nil
	}

	//在外联白名单中
	if crd.inOutreachWhile(dnsModel.OrgId, dnsModel.TargetIp) {
		return nil
	}

	model := alarm_models.AlarmOutreachModel{
		Id:           dnsModel.Id,
		OrgId:        dnsModel.OrgId,
		AgentId:      dnsModel.AgentId,
		Type:         enum_alarm.TypeOutreachNotInWhite,
		Level:        enum_alarm.LevelImportant,
		AgentAt:      dnsModel.AgentAt,
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    dnsModel.CreatedAt,
		ProcHash:     dnsModel.ProcHash,
		ProcDir:      dnsModel.ProcDir,
		ProcId:       dnsModel.ProcId,
		DomainName:   dnsModel.DomainName,
		TargetIp:     dnsModel.TargetIp,
	}

	//保存
	if !crd.masterModel.InsertItem(&model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchOutreach); hsCli != nil {
		hsData := huanshi_platform.AlarmOutreach{
			AlarmBase: huanshi_platform.AlarmBase{
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Type,
				NoticeTime: model.CreatedAt,
			},
			Phash:      model.ProcHash,
			Pid:        model.ProcId,
			Ppath:      model.ProcDir,
			Domain:     model.DomainName,
			AnalysisIp: model.TargetIp,
		}
		if err := hsCli.SyncAlarmOutreach(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	return nil
}
