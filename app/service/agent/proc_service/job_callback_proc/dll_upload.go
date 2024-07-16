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

func CreateDllUploadFactory() *DllUpload {
	return &DllUpload{
		masterModel: proc_models.CreateDllFactory(),
	}
}

type DllUpload struct {
	ProcessBase
	masterModel *proc_models.DllModel
}

// GetOne 获取一条
func (crd *DllUpload) GetOne(orgId int64, agentId, hash string) (model *proc_models.DllModel) {
	model = crd.masterModel.GetItemByOrgAndAgentAndHash(orgId, agentId, hash)
	if model == nil || string_helper.IsEmpty(model.DllHash) {
		return nil
	}
	return model
}

// Add 添加
func (crd *DllUpload) Add(dllModel *proc_models.DllModel) error {
	if string_helper.IsEmpty(dllModel.CreatedAt) {
		dllModel.CreatedAt = date_helper.CurDatetime()
	}
	dllModel.Dir = path_helper.ToStandard(dllModel.Dir)
	dllModel.ProcDir = path_helper.ToStandard(dllModel.ProcDir)
	dllModel.GrantThreatenSwitch = enums.SwitchOn

	//如果存在， 忽略
	hasModel := crd.GetOne(dllModel.OrgId, dllModel.AgentId, dllModel.DllHash)
	if hasModel != nil && !string_helper.IsEmpty(hasModel.DllHash) {
		return nil
	}

	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(dllModel.OrgId, dllModel.AgentId, agentFact.ColumnProc(true))
	if agentModel == nil || string_helper.IsEmpty(agentModel.AgentId) {
		return errors.New("终端不存在")
	}
	//冗余自动匹配白名单开关
	dllModel.AutoWhiteSwitch = agentModel.AutoWhiteSwitch

	//威胁情报状态
	isExist := false
	if isExist, dllModel.ThreatenStatus, dllModel.ThreatenLevel = crd.getThreatenStatus(dllModel.DllHash); !isExist {
		//加入威胁情报缓存
		crd.pushThreaten(dllModel, agentModel.OsName)
	}
	//全局名单状态
	dllModel.GlobalStatus = crd.getGlobalStatus(dllModel.OrgId, dllModel.DllHash)

	//默认威胁情报开关
	dllModel.GrantThreatenSwitch = enums.SwitchOn

	//最终状态
	dllModel.Status, dllModel.SyncStatusAt = crd.getProcStatus(dllModel.GlobalStatus, dllModel.SelfStatus,
		dllModel.ThreatenStatus, agentModel.AutoWhiteSwitch)

	//保存
	if !crd.masterModel.InsertItem(dllModel) {
		return errors.New("保存失败")
	}

	return nil
}

func (crd *DllUpload) pushThreaten(dllModel *proc_models.DllModel, agentOs string) {
	param := threaten_models.CacheModel{
		Id:        dllModel.DllHash,
		Source:    enum_threaten.SourceAgentUpload,
		Dir:       dllModel.Dir,
		ParentDir: dllModel.ProcDir,
		Os:        agentOs,
		Introduce: dllModel.Introduce,
		Copyright: dllModel.Copyright,
		Signature: dllModel.Signature,
		Attribute: dllModel.Attribute,
		DataFrom:  enum_process.DataFromDll,
	}
	threaten_service.CreateCacheFactory().Add(&param, false)
}
