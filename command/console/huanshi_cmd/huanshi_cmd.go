package huanshi_cmd

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"time"
)

func CreateCmdFactory() *HuanshiCmd {
	return &HuanshiCmd{}
}

type HuanshiCmd struct {
}

// getHsCli  获取幻势客户端
func (crd *HuanshiCmd) getHsCli() (err error, cli *huanshi_platform.RequestHuanshi) {
	dictIds := []string{
		enum_dictionary.HuanshiCallbackUrl,
		enum_dictionary.HuanshiCallIp,
		enum_dictionary.HuanshiSwitch,
		enum_dictionary.HuanshiUrl,
		enum_dictionary.HuanshiAlarmSwitch,
		enum_dictionary.HuanshiAlarmIp,
		enum_dictionary.HuanshiAlarmPort,
		enum_dictionary.HuanshiLogSwitch,
		enum_dictionary.HuanshiLogIp,
		enum_dictionary.HuanshiLogPort,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if string_helper.ToInt(vals[enum_dictionary.HuanshiSwitch]) == enums.SwitchOff {
		return errors.New("幻势功能已关闭"), nil
	}

	//客户端
	cli = huanshi_platform.CreateRequestFactory(huanshi_platform.RequestConfig{
		HuanshiCallbackUrl: vals[enum_dictionary.HuanshiCallbackUrl],
		HuanshiCallIp:      vals[enum_dictionary.HuanshiCallIp],
		HuanshiUrl:         vals[enum_dictionary.HuanshiUrl],
		HuanshiAlarmIp:     vals[enum_dictionary.HuanshiAlarmIp],
		HuanshiAlarmPort:   vals[enum_dictionary.HuanshiAlarmPort],
		HuanshiLogIp:       vals[enum_dictionary.HuanshiLogIp],
		HuanshiLogPort:     vals[enum_dictionary.HuanshiLogIp],
	})

	//验证在线
	if err = cli.Online(); err != nil {
		return errors.New("网络连接失败： " + err.Error()), nil
	}

	//生成token
	if err = cli.GenerateToken(); err != nil {
		return errors.New("生成口令： " + err.Error()), nil
	}

	return nil, cli
}

func (crd *HuanshiCmd) getOrgModels() []*organization_models.OrganizationModel {
	huanshiSwitch := enums.SwitchOn
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList:    []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
		HuanshiSwitch: &huanshiSwitch,
	})
}

func (crd *HuanshiCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnBasic(true, false))
}

// Push  推送
func (crd *HuanshiCmd) Push() {
	variable.ZapLog.Sugar().Debug("push start")

	//获取客户端
	err, hsCli := crd.getHsCli()
	if err != nil {
		variable.ZapLog.Sugar().Debug("push ", err.Error())
		return
	}

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Sugar().Debug("push orgModels empty")
		return
	}

	//先同步设备数据
	for _, orgModel := range orgModels {
		if err1 := hsCli.SyncDeviceByTcp(crd.initDevice(orgModel)); err1 != err1 {
			variable.ZapLog.Sugar().Debug("push SyncDeviceByTcp：", err1.Error())
			return
		}
	}
	//休眠一下, 等待幻势处理数据
	time.Sleep(time.Second * 5)
	//终端类型
	jsonStr := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentTypes)
	typeMap := enum_agent.GetTypeMap(jsonStr)

	//同步终端数据
	for _, orgModel := range orgModels {
		//获取全部终端
		agentModels := crd.getAgentModels(orgModel.Id)
		if agentModels == nil || len(agentModels) == 0 {
			continue
		}

		for _, agentModel := range agentModels {
			agentType := ""
			if val, isOk := typeMap[agentModel.AgentType]; isOk {
				agentType = val
			}

			param := huanshi_platform.AgentData{
				DeviceMark:   orgModel.Code,
				OrgId:        orgModel.Id,
				AgentId:      agentModel.AgentId,
				Os:           agentModel.OsName,
				Ip:           agentModel.AgentIp,
				NetIp:        agentModel.NetIp,
				HostName:     agentModel.HostName,
				Status:       agentModel.IsOnline,
				OnlineTime:   agentModel.OnlineAt,
				Version:      agentModel.Ver,
				DeptCode:     "",
				DeptName:     "",
				AgentType:    agentType,
				AliasName:    agentModel.Alias,
				ContactName:  agentModel.ContactName,
				ContactPhone: agentModel.ContactTel,
			}
			variable.ZapLog.Sugar().Debug(json_helper.Encode(param))
		}
	}

	variable.ZapLog.Sugar().Debug("push end")
}
func (crd *HuanshiCmd) initDevice(orgModel *organization_models.OrganizationModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["deviceMark"] = orgModel.Code

	agentFact := agent_models.CreateAgentFactory()
	calc := agentFact.CalcByConditions(agent_models.SearchConditions{OrgId: orgModel.Id},
		agentFact.ColumnCalcAgents(false, true, false, false, false))
	data["agentNum"] = calc.TotalNumber
	data["onlineAgentNum"] = calc.OnlineNumber

	return data
}
