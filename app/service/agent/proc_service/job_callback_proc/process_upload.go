package job_callback_proc

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateProcessUploadFactory() *ProcessUpload {
	return &ProcessUpload{
		masterModel: proc_models.CreateProcessFactory(),
	}
}

type ProcessUpload struct {
	ProcessBase
	masterModel *proc_models.ProcessModel
}

// GetOne 获取一条
func (crd *ProcessUpload) GetOne(orgId int64, agentId, hash string) (model *proc_models.ProcessModel) {
	model = crd.masterModel.GetItemByOrgAndAgentAndHash(orgId, agentId, hash)
	if model == nil || string_helper.IsEmpty(model.ProcHash) {
		return nil
	}
	return model
}

// Add 添加
func (crd *ProcessUpload) Add(processModel *proc_models.ProcessModel) error {
	if string_helper.IsEmpty(processModel.CreatedAt) {
		processModel.CreatedAt = date_helper.CurDatetime()
	}
	processModel.Dir = path_helper.ToStandard(processModel.Dir)

	//如果存在， 忽略
	hasModel := crd.GetOne(processModel.OrgId, processModel.AgentId, processModel.ProcHash)
	if hasModel != nil && !string_helper.IsEmpty(hasModel.ProcHash) {
		return nil
	}

	//终端信息
	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(processModel.OrgId, processModel.AgentId, agentFact.ColumnProc(true))
	if agentModel == nil || string_helper.IsEmpty(agentModel.AgentId) {
		return errors.New("终端不存在")
	}

	//全局名单状态
	processModel.GlobalStatus = crd.getGlobalStatus(processModel.OrgId, processModel.ProcHash)
	//默认自身状态
	processModel.SelfStatus = enum_process.RelyStatusUnset

	//默认威胁情报开关
	processModel.GrantThreatenSwitch = enums.SwitchOn
	//冗余自动匹配白名单开关
	processModel.AutoWhiteSwitch = agentModel.AutoWhiteSwitch
	//威胁情报状态
	isExist := false
	if isExist, processModel.ThreatenStatus, processModel.ThreatenLevel = crd.getThreatenStatus(processModel.ProcHash); !isExist {
		//加入威胁情报缓存
		crd.pushThreaten(processModel, agentModel.OsName)
	}

	//最终呈现状态
	processModel.Status, processModel.SyncStatusAt = crd.getProcStatus(processModel.GlobalStatus, processModel.SelfStatus,
		processModel.ThreatenStatus, processModel.AutoWhiteSwitch)

	//保存
	if !crd.masterModel.InsertItem(processModel) {
		return errors.New("保存失败")
	}

	return nil
}

func (crd *ProcessUpload) pushThreaten(processModel *proc_models.ProcessModel, agentOs string) {
	param := threaten_models.CacheModel{
		Id:        processModel.ProcHash,
		Source:    enum_threaten.SourceAgentUpload,
		Dir:       processModel.Dir,
		ParentDir: "",
		Os:        agentOs,
		Introduce: processModel.Introduce,
		Copyright: processModel.Copyright,
		Signature: processModel.Signature,
		Attribute: "",
		DataFrom:  enum_process.DataFromProcess,
	}
	threaten_service.CreateCacheFactory().Add(&param, false)
}
