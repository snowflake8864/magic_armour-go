package huanshi_cmd

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/huanshi_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
	"time"
)

func CreateCmdFactory() *HuanshiCmd {
	return &HuanshiCmd{}
}

type HuanshiCmd struct {
}

// getHsCli  获取幻势客户端
func (crd *HuanshiCmd) getHsCli() (err error, cli *huanshi_platform.RequestHuanshi) {
	cli = huanshi_helper.GetHuanshiCli("")

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

func (crd *HuanshiCmd) getDepartmentMap(orgId int64, agentModels []*agent_models.AgentModel) map[int64]*department_models.DepartmentModel {
	var deptIds []int64
	deptStruct := make(map[int64]struct{})
	for _, model := range agentModels {
		if _, isOk := deptStruct[model.DeptId]; model.DeptId > 0 && !isOk {
			deptStruct[model.DeptId] = struct{}{}
			deptIds = append(deptIds, model.DeptId)
		}
	}
	if len(deptIds) == 0 {
		return nil
	}
	return department_service.CreateCurdFactory().GetMapsByIds(orgId, deptIds)
}

// Push  推送
func (crd *HuanshiCmd) Push() {
	variable.ZapLog.Sugar().Debug("push start")
	//
	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Sugar().Debug("push orgModels empty")
		return
	}

	//获取客户端
	err, hsCli := crd.getHsCli()
	if err != nil {
		variable.ZapLog.Sugar().Debug("push ", err.Error())
		return
	}

	//先同步设备数据
	for _, orgModel := range orgModels {
		if err1 := hsCli.SyncDeviceByTcp(crd.initDevice(orgModel)); err1 != nil {
			variable.ZapLog.Sugar().Debug("push SyncDeviceByTcp：", err1.Error())
			return
		}
	}
	//休眠一下, 等待幻势处理设备数据
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
		//企业下组织架构
		departments := crd.getDepartmentMap(orgModel.Id, agentModels)

		for _, agentModel := range agentModels {
			agentType := ""
			if val, isOk := typeMap[agentModel.AgentType]; isOk {
				agentType = val
			}

			//终端对应组织架构
			var deptName, deptCode string
			if departments != nil && len(departments) > 0 {
				if deptModel, ok := departments[agentModel.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
					deptName = strings.Join(json_helper.ToSliceString(deptModel.FullName), "|")
					if string_helper.IsEmpty(deptModel.Code) {
						deptCode = int_helper.Int64ToStr(agentModel.DeptId)
					} else {
						deptCode = deptModel.Code
					}
				}
			}

			//缓存终端数据
			cacheAgent := huanshi_platform.CacheAgent{
				OrgCode:     orgModel.Code,
				AgentIp:     agentModel.AgentIp,
				NetIp:       agentModel.NetIp,
				DeptName:    deptName,
				ContactName: agentModel.ContactName,
				AgentType:   agentType,
			}
			if err1 := hsCli.SaveAgent(orgModel.Id, agentModel.AgentId, &cacheAgent); err1 != nil {
				variable.ZapLog.Sugar().Debug("push SaveAgent：", err1.Error())
				return
			}

			//幻势需要终端数据
			agentData := huanshi_platform.AgentData{
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
				DeptCode:     deptCode,
				DeptName:     deptName,
				AgentType:    agentType,
				AliasName:    agentModel.Alias,
				ContactName:  agentModel.ContactName,
				ContactPhone: agentModel.ContactTel,
			}
			//variable.ZapLog.Sugar().Debug(json_helper.Encode(agentData))
			if err1 := hsCli.SyncAgentBySyslog(&agentData); err1 != nil {
				variable.ZapLog.Sugar().Debug("push SyncAgentBySyslog：", err1.Error())
				return
			}
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
