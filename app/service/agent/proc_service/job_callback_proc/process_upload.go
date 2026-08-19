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

// getProcMaps 获取多条
func (crd *ProcessUpload) getProcMaps(orgId int64, agentId string, hashes []string) map[string]*proc_models.ProcessModel {
	result := make(map[string]*proc_models.ProcessModel)

	models := crd.masterModel.GetListByConditions(proc_models.SearchConditions{
		OrgId:      orgId,
		AgentId:    agentId,
		ProcHashes: hashes,
	})
	if models == nil {
		return nil
	}
	for _, model := range models {
		result[model.ProcHash] = model
	}
	return result
}

func (crd *ProcessUpload) isChanged(existModel, newestModel *proc_models.ProcessModel) bool {
	if existModel.Status != newestModel.Status {
		return true
	}
	if existModel.GlobalStatus != newestModel.GlobalStatus {
		return true
	}
	if existModel.SelfStatus != newestModel.SelfStatus {
		return true
	}
	if existModel.ThreatenStatus != newestModel.ThreatenStatus {
		return true
	}
	if existModel.AutoWhiteSwitch != newestModel.AutoWhiteSwitch {
		return true
	}

	return false
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

	//全局名单
	globalMap := crd.getGlobalStatusByHashes(orgId, hashes)
	//威胁情报
	threatenMap := crd.getThreatenStatusByHashes(hashes)
	//已存在的进程
	existsProcMap := crd.getProcMaps(orgId, agentId, hashes)

	//进程插入数据
	var insertProc []proc_models.ProcessModel
	//威胁情报插入数据
	var insertThreaten []threaten_models.CacheModel
	createdAt := date_helper.CurDatetime()
	for _, newModel := range insertMap {
		if string_helper.IsEmpty(newModel.CreatedAt) {
			newModel.CreatedAt = createdAt
		}
		//路径标准化
		newModel.Dir = path_helper.ToStandard(newModel.Dir)
		//默认威胁情报开关
		newModel.GrantThreatenSwitch = enums.SwitchOn
		//冗余自动匹配白名单开关
		newModel.AutoWhiteSwitch = agentModel.AutoWhiteSwitch

		//全局名单状态
		newModel.GlobalStatus = enum_process.RelyStatusUnset
		if val, isOK := globalMap[newModel.ProcHash]; isOK {
			newModel.GlobalStatus = val
		}

		//自身状态
		newModel.SelfStatus = enum_process.RelyStatusUnset
		if existsProc, isOK := existsProcMap[newModel.ProcHash]; isOK {
			newModel.SelfStatus = existsProc.SelfStatus
		}

		//威胁情报状态
		newModel.ThreatenStatus = enum_process.RelyStatusUnset
		if val, isOK := threatenMap[newModel.ProcHash]; isOK {
			newModel.ThreatenStatus = val
		} else {
			//威胁情报插入数据
			insertThreaten = append(insertThreaten, crd.initThreaten(&newModel, agentModel.OsName))
		}

		//最终呈现状态
		newModel.Status, newModel.SyncStatusAt = crd.getProcStatus(newModel.GlobalStatus, newModel.SelfStatus, newModel.ThreatenStatus, newModel.AutoWhiteSwitch)

		//已存在
		if existsProc, isOK := existsProcMap[newModel.ProcHash]; isOK {
			//发生改变,利用doris主键模型特性
			if crd.isChanged(existsProc, &newModel) {
				insertProc = append(insertProc, newModel)
			}
		} else {
			insertProc = append(insertProc, newModel)
		}
	}

	//保存进程
	if len(insertProc) > 0 {
		if !crd.masterModel.InsertBatch(&insertProc) {
			return errors.New(my_errors.CurdSaveFailed)
		}
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
