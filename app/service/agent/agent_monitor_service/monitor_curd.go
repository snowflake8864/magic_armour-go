package agent_monitor_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/gomail_factory"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateMonitorFactory() *MonitorCurd {
	return &MonitorCurd{agent_monitor_models.CreateMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *agent_monitor_models.MonitorModel
}

// Add 新增
func (crd *MonitorCurd) Add(param *agent_monitor_models.MonitorModel) error {

	//入库
	if crd.masterModel.InsertItem(param) {
		return nil
	}

	//发送告警邮件
	crd.sendEmail(param)

	return errors.New("保存失败")
}

func (crd *MonitorCurd) sendEmail(model *agent_monitor_models.MonitorModel) {
	//企业扩展
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	extendModel := OrgExtendFact.GetItemByOrg(model.OrgId)
	if extendModel == nil || extendModel.Id == 0 {
		variable.ZapLog.Sugar().Debug("企业扩展不存在, orgId=", model.OrgId)
		return
	}
	noticeConfig := OrgExtendFact.FormatHardwareNoticeConfig(extendModel.HardwareNoticeConfig)
	if noticeConfig.EmailSwitch == enums.SwitchOff {
		variable.ZapLog.Sugar().Debug("EmailSwitch关闭, orgId=", model.OrgId)
		return
	}
	if string_helper.IsEmpty(noticeConfig.Email) {
		variable.ZapLog.Sugar().Debug("邮箱为空, orgId=", model.OrgId)
		return
	}
	if noticeConfig.TimelySwitch == enums.SwitchOff {
		variable.ZapLog.Sugar().Debug("TimelySwitch关闭, orgId=", model.OrgId)
		return
	}

	mailParam := gomail_factory.HardwareTimely{
		Date:    model.CreatedAt,
		AgentId: model.AgentId,
	}

	if noticeConfig.CpuSwitch == enums.SwitchOn && float_helper.Float64(noticeConfig.CpuThreshold) <= model.CpuUsage {
		mailParam.UsageList = append(mailParam.UsageList, &gomail_factory.HardwareTimelyInfo{
			Name:      "CPU阈值",
			Threshold: noticeConfig.CpuThreshold,
			Usage:     model.CpuUsage,
		})
	}
	if noticeConfig.MemSwitch == enums.SwitchOn && float_helper.Float64(noticeConfig.MemThreshold) <= model.MemUsage {
		mailParam.UsageList = append(mailParam.UsageList, &gomail_factory.HardwareTimelyInfo{
			Name:      "内存阈值",
			Threshold: noticeConfig.MemThreshold,
			Usage:     model.MemUsage,
		})
	}
	if noticeConfig.HdSwitch == enums.SwitchOn && float_helper.Float64(noticeConfig.HdThreshold) <= model.HdUsage {
		mailParam.UsageList = append(mailParam.UsageList, &gomail_factory.HardwareTimelyInfo{
			Name:      "硬盘阈值",
			Threshold: noticeConfig.HdThreshold,
			Usage:     model.HdUsage,
		})
	}
	if mailParam.UsageList == nil || len(mailParam.UsageList) == 0 {
		variable.ZapLog.Sugar().Debug("不满足条件, orgId=", model.OrgId)
		return
	}

	//终端信息
	agentFact := agent_models.CreateAgentFactory()
	if agentModel := agentFact.GetItemByOrgAndAgent(model.OrgId, model.AgentId, agentFact.ColumnBasic(true, false)); agentModel == nil || agentModel.OrgId == 0 {
		variable.ZapLog.Sugar().Debug("终端不存在, orgId=", model.OrgId, ", agentId=", model.AgentId)
		return
	} else {
		mailParam.AgentIp = agentModel.AgentIp
		mailParam.AgentHostname = agentModel.HostName
	}

	//发送邮件
	subject, context := gomail_factory.CreateTemplateFactory().GetHardwareTimely(&mailParam)
	if err := gomail_factory.CreateMailerFactory().SendHtmlMsg(noticeConfig.Email, subject, context); err != nil {
		variable.ZapLog.Sugar().Debug("发送邮件错误, orgId=", model.OrgId, "【error】", err.Error())
		return
	}

}
