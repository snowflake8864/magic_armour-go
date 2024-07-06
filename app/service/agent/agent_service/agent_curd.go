package agent_service

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
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
	condition := agent_models.SearchConditions{OnlineAtLte: floorAt}
	updates := map[string]interface{}{
		"is_online": enums.IsNot,
	}
	isOk, _ := crd.masterModel.UpdateByConditions(condition, updates)
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

// SetLivingByAgent 设置终端活动中
func (crd *AgentCurd) SetLivingByAgent(orgId int64, id string) bool {
	nowDate := date_helper.CurDatetime()

	updates := map[string]interface{}{
		"live_at":   nowDate,
		"online_at": gorm.Expr("IF(is_online=?, ?, online_at)", enums.IsNot, nowDate),
		"is_online": enums.IsYes,
	}
	isOk, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, id, updates)
	return isOk
}

// BatchDelete 批量删除
func (crd *AgentCurd) BatchDelete(orgId int64, ids []string, isUninstall bool) (errCode int, err error) {
	var validIds []string
	if agentModels := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnSimplify()); len(agentModels) == 0 {
		for _, model := range agentModels {
			validIds = append(validIds, model.AgentId)
		}
	}
	if len(validIds) == 0 {
		return -1, errors.New("终端不存在")
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
	if agentModels := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnProc(true)); len(agentModels) == 0 {
		for _, model := range agentModels {
			if isOk, _ := enum_software.CanUpgradeIndirect(softwareMaps, model.OsType, model.Ver); isOk {
				validIds = append(validIds, model.AgentId)
			}
		}
	}

	if len(validIds) > 0 {
		//添加任务
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, validIds, enum_task.TypeAgentUpgrade, "")
	}
	return errCode, err
}

// BatchUpdate 批量修改
func (crd *AgentCurd) BatchUpdate(orgId int64, ids []string, param *ManualColumn, remark string) (errCode int, err error) {
	models := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnConfig(true, true))
	if models == nil || len(models) == 0 {
		return -1, errors.New("终端不存在")
	}

	fmt.Println("参数================================")
	fmt.Println(json_helper.Encode(param))
	for _, model := range models {
		err, hasTask, openAutoWhite, updates := crd.getUpdates(model, param)
		fmt.Println("修改================================")
		fmt.Println("AgentId:", model.AgentId)
		fmt.Println(json_helper.Encode(updates))
		fmt.Println("hasTask:", hasTask)
		fmt.Println("openAutoWhite:", openAutoWhite)
		if err != nil {
			continue
		}
		if updates == nil || len(updates) == 0 {
			continue
		}
		//修改
		if isOk, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, model.AgentId, updates); !isOk {
			continue
		}
		//需要下发任务
		if hasTask {
			crd.addTask(model, false, remark)
		}
		//开启自动白名单
		if openAutoWhite {
			//丢入队列
			msg := &kafka_go.MessageStruct{
				Type:    kafka_go.TypeCommonAgentOpenAutoWhite,
				OrgId:   orgId,
				AgentId: model.AgentId,
			}
			producer_kafka.CreateProducerFactory().ProduceCommon(context.Background(), msg)

			return errCode, err
		}
	}

	return errCode, err
}
