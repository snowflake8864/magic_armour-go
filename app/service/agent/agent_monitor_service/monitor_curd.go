package agent_monitor_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/gomail_factory"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreateMonitorFactory() *MonitorCurd {
	return &MonitorCurd{
		masterModel:  agent_monitor_models.CreateMonitorFactory(),
		redisClient:  redis_factory.GetOneRedisClient(),
		newestPrefix: "agent_newest_monitor_",
		newestExpire: 600,
	}
}

type MonitorCurd struct {
	masterModel  *agent_monitor_models.MonitorModel
	redisClient  *redis_factory.RedisClient
	newestPrefix string
	newestExpire int
}

func (crd *MonitorCurd) keyBynNewest(orgId int64, agentId string) string {
	return crd.newestPrefix + int_helper.Int64ToStr(orgId) + "_" + agentId
}

// cacheNewest 缓存最新一次监控数据
func (crd *MonitorCurd) cacheNewest(param *agent_monitor_models.MonitorModel) error {
	key := crd.keyBynNewest(param.OrgId, param.AgentId)
	data := json_helper.Encode(param)
	//新增缓存
	if isOK, err := crd.redisClient.SetEx(key, data, crd.newestExpire); !isOK {
		return errors.New("cache newest monitor fail: " + err.Error())
	}

	return nil
}

// GetNewest 获取最新一次监控数据
func (crd *MonitorCurd) GetNewest(orgId int64, agentId string) (err error, data *agent_monitor_models.MonitorModel) {
	jsonStr := crd.redisClient.Get(crd.keyBynNewest(orgId, agentId))
	if string_helper.IsEmpty(jsonStr) {
		return errors.New(my_errors.CurdCacheNotExist), data
	}

	err = json_helper.ToStruct(jsonStr, &data)
	if err != nil {
		return errors.New(my_errors.CurdCacheFailed), data
	}
	return nil, data
}

// Add 新增
func (crd *MonitorCurd) Add(param *agent_monitor_models.MonitorModel) error {

	//入库
	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//缓存数据
	if err := crd.cacheNewest(param); err != nil {
		variable.ZapLog.Warn(err.Error())
	}

	//发送告警邮件
	crd.sendEmail(param)

	return nil
}

func (crd *MonitorCurd) sendEmail(model *agent_monitor_models.MonitorModel) {
	//企业扩展
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	extendModel := OrgExtendFact.GetItemByOrg(model.OrgId)
	if extendModel == nil || extendModel.Id == 0 {
		variable.ZapLog.Sugar().Info("企业扩展不存在, orgId=", model.OrgId)
		return
	}
	noticeConfig := OrgExtendFact.FormatHardwareNoticeConfig(extendModel.HardwareNoticeConfig)
	if noticeConfig.EmailSwitch == enums.SwitchOff {
		//variable.ZapLog.Sugar().Info("EmailSwitch关闭, orgId=", model.OrgId)
		return
	}
	if noticeConfig.TimelySwitch == enums.SwitchOff {
		//variable.ZapLog.Sugar().Info("TimelySwitch关闭, orgId=", model.OrgId)
		return
	}
	if string_helper.IsEmpty(noticeConfig.Email) {
		variable.ZapLog.Sugar().Info("邮箱为空, orgId=", model.OrgId)
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
		variable.ZapLog.Sugar().Info("不满足条件, orgId=", model.OrgId)
		return
	}

	//终端信息
	agentFact := agent_models.CreateAgentFactory()
	if agentModel := agentFact.GetItemByOrgAndAgent(model.OrgId, model.AgentId, agentFact.ColumnBasic(true, false)); agentModel == nil || agentModel.OrgId == 0 {
		variable.ZapLog.Sugar().Info("终端不存在, orgId=", model.OrgId, ", agentId=", model.AgentId)
		return
	} else {
		mailParam.AgentIp = agentModel.AgentIp
		mailParam.AgentHostname = agentModel.HostName
	}

	//发送邮件
	subject, context := gomail_factory.CreateTemplateFactory().GetHardwareTimely(&mailParam)
	if err := gomail_factory.CreateMailerFactory().SendHtmlMsg(noticeConfig.Email, subject, context); err != nil {
		variable.ZapLog.Sugar().Info("发送邮件错误, orgId=", model.OrgId, "【error】", err.Error())
		return
	}

}
