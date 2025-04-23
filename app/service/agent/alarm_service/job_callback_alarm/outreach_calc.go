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
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateOutreachCalcFactory() *AlarmOutreachCalc {
	return &AlarmOutreachCalc{
		masterModel: alarm_models.CreateOutreachFactory(),
	}
}

type AlarmOutreachCalc struct {
	alarmBase
	masterModel *alarm_models.OutreachModel
}

func (crd *AlarmOutreachCalc) inOutreachWhile(orgId int64, targetIpStr string) (bool, string) {
	SwitchOn := enums.SwitchOn
	conditions := outreach_models.SearchConditions{
		OrgId:      orgId,
		Type:       enum_outreach.IpTypeWhite,
		OpenSwitch: &SwitchOn,
	}
	models := outreach_models.CreateGlobalFactory().GetListByConditions(conditions)
	if models == nil || len(models) == 0 {
		return true, ""
	}

	targetIpArr := strings.Split(targetIpStr, ";") //192.168.1.1;fe80::1;15.155.11.11
	for _, targetIp := range targetIpArr {
		isInWhile := false
		for _, model := range models {
			if ip_helper.IsIPInRange(targetIp, model.Ip) {
				isInWhile = true
				break
			}
		}
		if !isInWhile {
			return false, targetIp
		}
	}

	return true, ""
}

// Calc 计算
func (crd *AlarmOutreachCalc) Calc(dnsModel *syslog_models.DnsModel) error {
	//终端开关
	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(dnsModel.OrgId, dnsModel.AgentId, agentFact.ColumnConfig(false, false))
	if agentModel == nil || agentModel.OutreachAlarmSwitch == enums.SwitchOff {
		return nil
	}
	//在外联白名单中
	isOk, targetIp := crd.inOutreachWhile(dnsModel.OrgId, dnsModel.TargetIp)
	if isOk {
		return nil
	}

	model := alarm_models.OutreachModel{
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
		TargetIp:     targetIp,
	}

	//保存
	if !crd.masterModel.InsertItem(&model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchOutreach); hsCli != nil {
		crd.PushHs(hsCli, &model)
	}

	return nil
}

// BatchAdd 新增
func (crd *AlarmOutreachCalc) BatchAdd(list *[]alarm_models.OutreachModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []alarm_models.OutreachModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()

		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	//推送到幻势
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchOutreach); hsCli != nil {
		for _, model := range inserts {
			crd.PushHs(hsCli, &model)
		}
	}

	return nil
}

func (crd *AlarmOutreachCalc) PushHs(hsCli *huanshi_platform.RequestHuanshi, model *alarm_models.OutreachModel) {
	//推送到幻势,
	if hsCli != nil {
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
			variable.ZapLog.Warn(err.Error())
		}
	}
}
