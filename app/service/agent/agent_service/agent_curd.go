package agent_service

import (
	"context"
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_online_models"
	"goskeleton/app/service/agent/online_log_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{masterModel: agent_models.CreateAgentFactory()}
}

type AgentCurd struct {
	AgentBase
	masterModel *agent_models.AgentModel
}

// SetOffline 设置离线
func (crd *AgentCurd) SetOffline(floorAt string) bool {
	isOnline := enums.IsYes
	condition := agent_models.SearchConditions{
		IsOnline:  &isOnline,
		LiveAtLte: floorAt,
	}

	var agentIds []string

	models := crd.masterModel.GetListByConditions(condition, crd.masterModel.ColumnBasic(true, false))
	if len(models) > 0 {
		nowDate := date_helper.CurDatetime()
		var onlineLogList []agent_online_models.OnlineLogModel
		for _, model := range models {
			agentIds = append(agentIds, model.AgentId)

			onlineLogList = append(onlineLogList, *crd.initOnlineLog(model, enums.IsNot, nowDate, "自动离线"))
		}
		_ = online_log_service.CreateCurdFactory().BatchAdd(&onlineLogList)
	}

	//修改
	if len(agentIds) > 0 {
		condition.AgentIds = agentIds
		updates := map[string]interface{}{
			"is_online": enums.IsNot,
		}
		crd.masterModel.UpdateByConditions(condition, updates)
	}

	return true
}

func (crd *AgentCurd) initOnlineLog(model *agent_models.AgentModel, isOnLine int, nowDate, remark string) *agent_online_models.OnlineLogModel {
	return &agent_online_models.OnlineLogModel{
		OrgId:       model.OrgId,
		AgentId:     model.AgentId,
		HostName:    model.HostName,
		MacId:       model.MacId,
		AgentIp:     model.AgentIp,
		NetIp:       model.NetIp,
		OsName:      model.OsName,
		Ver:         model.Ver,
		IsOnline:    isOnLine,
		DeptId:      model.DeptId,
		Alias:       model.Alias,
		ContactName: model.ContactName,
		ContactTel:  model.ContactTel,
		Remark:      remark,
		CreatedAt:   nowDate,
	}
}

// SetLivingByAgent 设置终端活动中
func (crd *AgentCurd) SetLivingByAgent(orgId int64, id string) bool {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, id, crd.masterModel.ColumnBasic(true, false))

	nowDate := date_helper.CurDatetime()

	//设置活动时间
	updates := map[string]interface{}{
		"live_at":   nowDate,
		"is_online": enums.IsYes,
	}
	isOk, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, id, updates)

	//记录上线日志
	if model.IsOnline == enums.IsNot {
		_ = online_log_service.CreateCurdFactory().Add(crd.initOnlineLog(model, enums.IsYes, nowDate, "心跳"))
	}

	return isOk
}

// SetSafeByAgent 设置终端安全状态
func (crd *AgentCurd) SetSafeByAgent(orgId int64, id string, isSafe int) bool {
	updates := map[string]interface{}{
		"is_safe": isSafe,
	}
	isOk, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, id, updates)
	return isOk
}

// BatchDelete 批量删除
func (crd *AgentCurd) BatchDelete(orgId int64, ids []string, isUninstall bool) (errCode int, err error) {
	var validIds []string
	if agentModels := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnSimplifyAppend(true)); agentModels != nil && len(agentModels) > 0 {
		for _, model := range agentModels {
			if isUninstall {
				//卸载必须在线
				if model.IsOnline == enums.IsYes {
					validIds = append(validIds, model.AgentId)
				}
			} else {
				//删除必须离线
				if model.IsOnline == enums.IsNot {
					validIds = append(validIds, model.AgentId)
				}
			}

		}
	}
	if len(validIds) == 0 {
		return -1, errors.New("没有符合条件的终端")
	}

	if isOK := crd.masterModel.DeleteByOrgAndAgents(orgId, validIds); !isOK {
		return -1, errors.New("删除失败")
	}

	//如是是卸载
	if isUninstall {
		//添加任务
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, validIds, enum_task.TypeAgentUninstall, "")
	}

	//丢入队列
	ctx := context.Background()
	for _, val := range validIds {
		msg := &kafka_go.MessageStruct{
			Type:    kafka_go.TypeCommonAgentDelete,
			OrgId:   orgId,
			AgentId: val,
		}
		producer_kafka.CreateProducerFactory().ProduceCommon(ctx, msg)
	}

	return errCode, err
}

// BatchUpgrade 批量升级
func (crd *AgentCurd) BatchUpgrade(orgId int64, ids []string) (errCode int, err error) {
	softwareMaps := software_service.CreateCurdFactory().GetMapsByPublished()
	if softwareMaps == nil || len(softwareMaps) == 0 {
		return errCode, err
	}
	var validIds []string
	if agentModels := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnBasic(true, true)); agentModels != nil && len(agentModels) > 0 {
		for _, model := range agentModels {
			if enum_agent.CanUpgrade(model.GrantStatus) {
				if isOk, _ := enum_software.CanUpgradeIndirect(softwareMaps, model.OsType, model.Ver); isOk {
					validIds = append(validIds, model.AgentId)
				}
			}
		}
	}

	if len(validIds) > 0 {
		//添加任务
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, validIds, enum_task.TypeAgentUpgrade, "")
	}
	return errCode, err
}

// UpdateByAgent 单个修改
func (crd *AgentCurd) UpdateByAgent(orgId int64, agentId string, param *ManualColumn, remark string) (errCode int, err error) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, crd.masterModel.ColumnConfig(true, true))
	if model == nil || model.OrgId == 0 {
		return -1, errors.New("终端不存在")
	}

	err1, hasTask, hasConfigTask, updates, openAutoWhite := crd.getUpdates(model, param)
	if err1 != nil {
		return -1, err1
	}
	if updates == nil || len(updates) == 0 {
		return 0, nil
	}

	//修改
	if isOk, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, model.AgentId, updates); !isOk {
		return -1, errors.New("修改失败")
	}

	//需要下发任务
	if hasTask {
		crd.addTask(model, param, hasConfigTask, false, remark)
	}
	//开启自动白名单
	if openAutoWhite {
		crd.openAutoWhite(orgId, agentId)
	}

	return 0, nil
}

// BatchUpdate 批量修改
func (crd *AgentCurd) BatchUpdate(orgId int64, ids []string, param *ManualColumn, remark string) (errCode int, err error) {
	models := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnConfig(true, true))
	if models == nil || len(models) == 0 {
		return -1, errors.New("终端不存在")
	}

	//修改的数据
	updateResult := make(map[string]map[string]interface{}) // 修改结果
	agentUpdate := make(map[string][]string)                //终端修改数据
	agentConfigTask := make(map[string]bool)                //终端修改数据
	var hasTaskList []*agent_models.AgentModel              //有任务的终端列表
	var openAutoWhiteList []string                          //开启自动白名单终端
	for _, model := range models {
		err1, hasTask, hasConfigTask, updates, openAutoWhite := crd.getUpdates(model, param)
		if err1 != nil || updates == nil || len(updates) == 0 {
			continue
		}

		// 修改结果
		hash := md5_encrypt.MD5(json_helper.Encode(updates))
		if _, isOk := updateResult[hash]; !isOk {
			updateResult[hash] = updates
		}

		//终端修改数据
		if _, isOk := agentUpdate[hash]; isOk {
			agentUpdate[hash] = append(agentUpdate[hash], model.AgentId)
		} else {
			agentUpdate[hash] = []string{model.AgentId}
		}

		//需要下发任务
		if hasTask {
			hasTaskList = append(hasTaskList, model)
			agentConfigTask[model.AgentId] = hasConfigTask
		}

		//开启自动白名单
		if openAutoWhite {
			openAutoWhiteList = append(openAutoWhiteList, model.AgentId)
		}
	}

	//批量修改终端数据
	if len(agentUpdate) > 0 {
		//按修改数据批量处理
		for hash, agentIds := range agentUpdate {
			if updates, isOk := updateResult[hash]; isOk && len(updates) > 0 {
				crd.masterModel.UpdateByOrgAndAgents(orgId, agentIds, updates)
			}
		}
	}

	//开启自动白名单
	if len(openAutoWhiteList) > 0 {
		for _, agentId := range openAutoWhiteList {
			crd.openAutoWhite(orgId, agentId)
		}
	}

	//任务
	if len(hasTaskList) > 0 {
		for _, model := range hasTaskList {
			hasConfigTask := false
			if val, isOK := agentConfigTask[model.AgentId]; isOK {
				hasConfigTask = val
			}

			crd.addTask(model, param, hasConfigTask, false, remark)
		}
	}

	return errCode, nil
}

func (crd *AgentCurd) openAutoWhite(orgId int64, agentId string) {
	//丢入队列
	msg := &kafka_go.MessageStruct{
		Type:    kafka_go.TypeCommonAgentOpenAutoWhite,
		OrgId:   orgId,
		AgentId: agentId,
	}
	producer_kafka.CreateProducerFactory().ProduceCommon(context.Background(), msg)
}
