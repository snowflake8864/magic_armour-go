package proc_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/utils/helpers"
)

func CreateProcessFactory() *ProcessCurd {
	return &ProcessCurd{
		masterModel: proc_models.CreateProcessFactory(),
	}
}

type ProcessCurd struct {
	ProcessBase
	masterModel *proc_models.ProcessModel
}

// GetOne 获取一条
func (crd *ProcessCurd) GetOne(orgId int64, agentId, hash string) (model *proc_models.ProcessModel) {
	model = crd.masterModel.GetItemByOrgAndAgentAndHash(orgId, agentId, hash)
	if helpers.StringIsEmpty(model.ProcHash) {
		return nil
	}
	return model
}

// Add 添加
func (crd *ProcessCurd) Add(processModel *proc_models.ProcessModel) (errCode int, err error) {
	if helpers.StringIsEmpty(processModel.CreatedAt) {
		processModel.CreatedAt = helpers.DatetimeCurr()
	}
	//如果存在， 忽略
	hasModel := crd.GetOne(processModel.OrgId, processModel.AgentId, processModel.ProcHash)
	if hasModel == nil || helpers.StringIsEmpty(hasModel.ProcHash) {
		return errCode, err
	}

	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(processModel.OrgId, processModel.AgentId, agentFact.ColumnProc(true))
	if agentModel == nil || helpers.StringIsEmpty(agentModel.AgentId) {
		variable.ZapLog.Info("终端不存在")
		return
	}

	//威胁情报状态
	isExist := false
	if isExist, processModel.ThreatenStatus = crd.getThreatenStatus(processModel.ProcHash); !isExist {
		//加入威胁情报缓存
		crd.pushThreaten(processModel, agentModel.OsName)
	}
	//全局名单状态
	processModel.GlobalStatus = crd.getGlobalStatus(processModel.OrgId, processModel.ProcHash)

	//最终状态
	processModel.Status, processModel.SyncStatusAt = crd.getProcStatus(processModel.GlobalStatus, processModel.SelfStatus,
		processModel.ThreatenStatus, agentModel.AutoWhiteSwitch)

	if crd.masterModel.InsertItem(processModel) {
		return -1, errors.New("保存失败")
	}

	return errCode, err
}

func (crd *ProcessCurd) pushThreaten(processModel *proc_models.ProcessModel, agentOs string) {
	param := threaten_models.ThreatenCacheModel{
		Id:        processModel.ProcHash,
		Source:    enum_threaten.SourceAgentUpload,
		Dir:       processModel.Dir,
		ParentDir: "",
		Os:        agentOs,
		Introduce: processModel.Introduce,
		Copyright: processModel.Copyright,
		Signature: processModel.Signature,
		Attribute: "",
		DataFrom:  enum_process.SourceProcess,
	}
	go threaten_service.CreateCacheFactory().Add(&param, false)
}
