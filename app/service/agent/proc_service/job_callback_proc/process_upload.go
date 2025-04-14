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

// GetMaps 获取多条
func (crd *ProcessUpload) GetMaps(orgId int64, agentId string, hashes []string) map[string]int {
	result := make(map[string]int)

	models := crd.masterModel.GetListByConditions(proc_models.SearchConditions{
		OrgId:      orgId,
		AgentId:    agentId,
		ProcHashes: hashes,
	})
	if models == nil {
		return nil
	}
	for _, model := range models {
		result[model.ProcHash] = model.Status
	}
	return result
}

// BatchAdd 添加
func (crd *ProcessUpload) BatchAdd(orgId int64, agentId string, models *[]proc_models.ProcessModel) error {
	//终端信息
	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(orgId, agentId, agentFact.ColumnProc(true))
	if agentModel == nil || string_helper.IsEmpty(agentModel.AgentId) {
		return errors.New("终端不存在")
	}

	//整理相关数据
	var hashes []string
	insertMap := make(map[string]proc_models.ProcessModel)
	for _, model := range *models {
		if !string_helper.IsEmpty(model.ProcHash) {
			hashes = append(hashes, model.ProcHash)
			insertMap[model.ProcHash] = model
		}
	}
	if len(insertMap) == 0 {
		return nil
	}

	//去掉已存在数据
	if hasMap := crd.GetMaps(orgId, agentId, hashes); len(hasMap) > 0 {
		hashes = nil
		for hash := range hasMap {
			delete(insertMap, hash)
			hashes = append(hashes, hash)
		}
	}
	if len(insertMap) == 0 {
		return nil
	}

	//全局名单
	globalMap := crd.getGlobalStatusByHashes(orgId, hashes)
	//威胁情报
	threatenMap := crd.getThreatenStatusByHashes(hashes)

	//进程插入数据
	var insertProc []proc_models.ProcessModel
	//威胁情报插入数据
	var insertThreaten []threaten_models.CacheModel
	createdAt := date_helper.CurDatetime()
	for _, procModel := range insertMap {
		if string_helper.IsEmpty(procModel.CreatedAt) {
			procModel.CreatedAt = createdAt
		}
		procModel.Dir = path_helper.ToStandard(procModel.Dir)

		//全局名单状态
		procModel.GlobalStatus = enum_process.RelyStatusUnset
		if val, isOK := globalMap[procModel.ProcHash]; isOK {
			procModel.GlobalStatus = val
		}
		//默认自身状态
		procModel.SelfStatus = enum_process.RelyStatusUnset
		//默认威胁情报开关
		procModel.GrantThreatenSwitch = enums.SwitchOn
		//冗余自动匹配白名单开关
		procModel.AutoWhiteSwitch = agentModel.AutoWhiteSwitch
		//威胁情报状态
		procModel.ThreatenStatus = enum_process.RelyStatusUnset
		if val, isOK := threatenMap[procModel.ProcHash]; isOK {
			procModel.ThreatenStatus = val
		} else {
			//威胁情报插入数据
			insertThreaten = append(insertThreaten, crd.initThreaten(&procModel, agentModel.OsName))
		}
		//最终呈现状态
		procModel.Status, procModel.SyncStatusAt = crd.getProcStatus(procModel.GlobalStatus, procModel.SelfStatus, procModel.ThreatenStatus, procModel.AutoWhiteSwitch)

		insertProc = append(insertProc, procModel)
	}

	//保存
	if !crd.masterModel.InsertBatch(&insertProc) {
		return errors.New("保存失败")
	}

	//保存威胁情报
	if len(insertThreaten) > 0 {
		threaten_service.CreateCacheFactory().BatchAdd(&insertThreaten, false)
	}
	return nil
}

func (crd *ProcessUpload) initThreaten(processModel *proc_models.ProcessModel, agentOs string) threaten_models.CacheModel {
	data := threaten_models.CacheModel{
		Id:        processModel.ProcHash,
		Source:    enum_threaten.SourceAgentUpload,
		Dir:       processModel.Dir,
		ParentDir: "",
		Os:        agentOs,
		Introduce: processModel.Introduce,
		Copyright: processModel.Copyright,
		Signature: processModel.Signature,
		Attribute: "",
		DataFrom:  enum_threaten.DataFromProcess,
	}

	return data
}
