package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/plug_service/plug_calc_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"strings"
)

func CreateLoginCalcFactory() *AlarmLoginCalc {
	return &AlarmLoginCalc{
		masterModel: alarm_models.CreateLoginFactory(),
	}
}

type AlarmLoginCalc struct {
	alarmBase
	masterModel *alarm_models.LoginModel
}

func (crd *AlarmLoginCalc) inOutreachWhile(orgId int64, targetIp string) bool {
	SwitchOn := enums.SwitchOn
	conditions := outreach_models.SearchConditions{
		OrgId:      orgId,
		OpenSwitch: &SwitchOn,
	}
	models := outreach_models.CreateGlobalFactory().GetListByConditions(conditions)
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
func (crd *AlarmLoginCalc) Calc(orgId int64, agentId string, attackIp string, attackAt string) error {
	//终端开关
	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(orgId, agentId, agentFact.ColumnConfig(false, false))
	variable.ZapLog.Sugar().Debug("口令破解告警=============================================")
	variable.ZapLog.Sugar().Debug(json_helper.Encode(agentModel))
	if agentModel == nil || agentModel.LoginSolveSwitch == enums.SwitchOff {
		return nil
	}

	//登录口令破解配置
	loginConfig := agent_service.CreateExtendFactory().GetLoginSolve(orgId, agentId)
	variable.ZapLog.Sugar().Debug("口令破解配置=============================================")
	variable.ZapLog.Sugar().Debug(json_helper.Encode(loginConfig))
	if loginConfig.Switch == enums.SwitchOff {
		return nil
	}
	//从登录日志计算
	loginStatus := enums.IsNot
	search := syslog_models.SearchConditions{
		OrgId:       orgId,
		AgentId:     agentId,
		SourceIp:    attackIp,
		LoginStatus: &loginStatus,
		CreateAtGte: date_helper.CalcDatetimeByMinute(attackAt, 0-int64(loginConfig.Gap)),
	}
	fact := syslog_models.CreateLoginFactory()
	calcSyslog := fact.CalcByConditions(search, fact.ColumnCountLogin())
	variable.ZapLog.Sugar().Debug("登录日志=============================================")
	variable.ZapLog.Sugar().Debug(" calcSyslog.TotalNumber=", calcSyslog.TotalNumber)
	variable.ZapLog.Sugar().Debug(" loginConfig.Threshold=", loginConfig.Threshold)
	//不满足条件
	if calcSyslog.TotalNumber < loginConfig.Threshold {
		return nil
	}

	//告警数据
	model := alarm_models.LoginModel{
		Id:           variable.SnowFlake.GetId(),
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         enum_alarm.TypeLoginSolve,
		Level:        enum_alarm.LevelSerious,
		AgentAt:      attackAt,
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    attackAt,
		SourceIp:     attackIp,
		SolveTotal:   calcSyslog.TotalNumber,
	}
	if !crd.masterModel.InsertItem(&model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchLogin); hsCli != nil {
		hsData := huanshi_platform.AlarmLogin{
			AlarmBase: huanshi_platform.AlarmBase{
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Type,
				NoticeTime: model.CreatedAt,
			},
			SourceIp: model.SourceIp,
			Total:    model.SolveTotal,
		}
		if err := hsCli.SyncAlarmLogin(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	//
	//动态封堵
	if err := plug_calc_service.CreatePlugCalcFactory().CalcByLogin(model.OrgId, model.AgentId, attackIp); err != nil {
		variable.ZapLog.Sugar().Debug(err.Error())
	}

	return nil
}
