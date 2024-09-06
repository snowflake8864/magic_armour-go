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
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/gomail_factory"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/helpers/string_helper"
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
	agentModel := agentFact.GetItemByOrgAndAgent(orgId, agentId, agentFact.ColumnConfig(true, false))
	if agentModel == nil || agentModel.LoginSolveSwitch == enums.SwitchOff {
		return nil
	}

	//登录口令破解配置
	loginConfig := agent_service.CreateExtendFactory().GetLoginSolve(orgId, agentId)
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

	//推送到幻势
	crd.PushHs(&model)

	//发送告警邮件
	crd.sendEmail(&model, agentModel)

	//动态封堵
	if err := plug_calc_service.CreatePlugCalcFactory().CalcByLogin(model.OrgId, model.AgentId, attackIp); err != nil {
		variable.ZapLog.Sugar().Debug(err.Error())
	}

	return nil
}

func (crd *AlarmLoginCalc) PushHs(model *alarm_models.LoginModel) {
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
}

func (crd *AlarmLoginCalc) sendEmail(loginModel *alarm_models.LoginModel, agentModel *agent_models.AgentModel) {
	//企业扩展
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	extendModel := OrgExtendFact.GetItemByOrg(agentModel.OrgId)
	if extendModel == nil || extendModel.Id == 0 {
		variable.ZapLog.Sugar().Debug("企业扩展不存在, orgId=", agentModel.OrgId)
		return
	}
	noticeConfig := OrgExtendFact.FormatAlarmNoticeConfig(extendModel.AlarmNoticeConfig)
	if noticeConfig.EmailSwitch == enums.SwitchOff {
		variable.ZapLog.Sugar().Debug("EmailSwitch关闭, orgId=", agentModel.OrgId)
		return
	}
	if string_helper.IsEmpty(noticeConfig.Email) {
		variable.ZapLog.Sugar().Debug("邮箱为空, orgId=", agentModel.OrgId)
		return
	}
	if noticeConfig.TimelySwitch == enums.SwitchOff {
		variable.ZapLog.Sugar().Debug("TimelySwitch关闭, orgId=", agentModel.OrgId)
		return
	}
	if len(noticeConfig.TimelyMod) == 0 {
		variable.ZapLog.Sugar().Debug("没有选中告警模块, orgId=", agentModel.OrgId)
		return
	}

	mailParam := gomail_factory.AlarmTimely{
		Date:          loginModel.CreatedAt,
		AgentId:       agentModel.AgentId,
		AgentIp:       agentModel.AgentIp,
		AgentHostname: agentModel.HostName,
		AlarmLogin:    map[string]int{loginModel.SourceIp: loginModel.SolveTotal},
	}

	//发送邮件
	subject, context := gomail_factory.CreateTemplateFactory().GetAlarmTimely(&mailParam)
	if err := gomail_factory.CreateMailerFactory().SendHtmlMsg(noticeConfig.Email, subject, context); err != nil {
		variable.ZapLog.Sugar().Debug("发送邮件错误, orgId=", agentModel, "【error】", err.Error())
		return
	}
}
