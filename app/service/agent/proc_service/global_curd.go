package proc_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateGlobalCurdFactory() *GlobalCurd {
	return &GlobalCurd{
		masterModel: proc_models.CreateProcGlobalFactory(),
	}
}

type GlobalCurd struct {
	masterModel *proc_models.ProcGlobalModel
}

// GetOne 获取一条
func (crd *GlobalCurd) GetOne(orgId int64, hash string) (model *proc_models.ProcGlobalModel) {
	model = crd.masterModel.GetItemByOrgAndId(orgId, hash)
	if model == nil || string_helper.IsEmpty(model.Id) {
		return nil
	}
	return model
}

// isMaxLimit 是否到达最大限制数
func (crd *GlobalCurd) isMaxLimit(orgId int64, dataFrom int, number int) error {
	limit := string_helper.ToInt64(dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentProcGlobalLimit))

	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		DataFrom: dataFrom,
	}
	total := crd.masterModel.Count(condition)

	if (limit - total - int64(number)) < 0 {
		return errors.New("全局名单库已达到最大限制条目")
	}
	return nil
}

// Add 添加
func (crd *GlobalCurd) Add(param *proc_models.ProcGlobalModel) error {
	param.Id = strings.ToLower(param.Id)
	param.Dir = path_helper.ToStandard(param.Dir)

	hasModel := crd.masterModel.GetItemByOrgAndId(param.OrgId, param.Id)
	if hasModel != nil && hasModel.OrgId != 0 {
		return errors.New("数据已存在")
	}
	if err := crd.isMaxLimit(param.OrgId, param.DataFrom, 1); err != nil {
		return err
	}

	param.OpenSwitch = enums.SwitchOn
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	//
	if param.DataFrom == enum_process.DataFromProcess {
		CreateProcessCurdFactory().GlobalCallbackAddByHash(param.OrgId, param.Type, param.Id)
	} else {
		CreateDllCurdFactory().GlobalCallbackAddByHash(param.OrgId, param.Type, param.Id)
	}

	return nil
}

// BatchAdd  批量添加
func (crd *GlobalCurd) BatchAdd(param *proc_models.ProcGlobalModel, list *[]proc_models.ProcGlobalModel) error {
	//过滤一遍重复数据
	var repeatMap = make(map[string]struct{})
	var hashList []string
	for _, model := range *list {
		id := strings.ToLower(model.Id)
		if _, isOk := repeatMap[id]; isOk {
			repeatMap[id] = struct{}{}
			hashList = append(hashList, id)
		}
	}

	//数据库已存在的
	var existMap = make(map[string]struct{})
	hasModels := crd.masterModel.GetItemsByOrgAndIds(param.OrgId, hashList)
	if hasModels != nil && len(hasModels) > 0 {
		for _, model := range hasModels {
			existMap[model.Id] = struct{}{}
		}
	}

	now := date_helper.CurDatetime()
	//要插入数据
	var inserts []proc_models.ProcGlobalModel
	var hashes []string
	for _, model := range *list {
		id := strings.ToLower(model.Id)
		if _, isOk := existMap[id]; !isOk {
			hashes = append(hashList, id)

			inserts = append(inserts, proc_models.ProcGlobalModel{
				Id:         id,
				OrgId:      param.OrgId,
				Type:       param.Type,
				Dir:        path_helper.ToStandard(model.Dir),
				DataFrom:   param.DataFrom,
				OpenSwitch: enums.SwitchOn,
				CreateId:   param.CreateId,
				CreateName: param.CreateName,
				UpdateId:   param.CreateId,
				UpdateName: param.CreateName,
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
	}
	number := len(inserts)
	if number == 0 {
		return nil
	}
	if err := crd.isMaxLimit(param.OrgId, param.DataFrom, number); err != nil {
		return err
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	if param.DataFrom == enum_process.DataFromProcess {
		CreateProcessCurdFactory().GlobalCallbackAddByHashes(param.OrgId, param.Type, hashes)
	} else {
		CreateDllCurdFactory().GlobalCallbackAddByHashes(param.OrgId, param.Type, hashes)
	}

	return nil
}

// SetSwitch 通过hash设置开关
func (crd *GlobalCurd) SetSwitch(orgId int64, id string, openSwitch int, dataFrom int) error {
	condition := proc_models.SearchConditions{
		OrgId: orgId,
		Id:    id,
	}
	updates := map[string]interface{}{
		"updated_at":  date_helper.CurDatetime(),
		"open_switch": openSwitch,
	}
	crd.masterModel.UpdateByConditions(condition, updates)

	if dataFrom == enum_process.DataFromProcess {
		if openSwitch == enums.SwitchOn {
			if hasModel := crd.masterModel.GetItemByOrgAndId(orgId, id); hasModel != nil && hasModel.OrgId > 0 {
				CreateProcessCurdFactory().GlobalCallbackAddByHash(orgId, hasModel.Type, id)
			}
		} else {
			CreateProcessCurdFactory().GlobalCallbackDeleteByHash(orgId, id)
		}
	} else {
		if openSwitch == enums.SwitchOn {
			if hasModel := crd.masterModel.GetItemByOrgAndId(orgId, id); hasModel != nil && hasModel.OrgId > 0 {
				CreateDllCurdFactory().GlobalCallbackAddByHash(orgId, hasModel.Type, id)
			}
		} else {
			CreateDllCurdFactory().GlobalCallbackDeleteByHash(orgId, id)
		}
	}

	return nil
}

// DeleteById 通过hash删除
func (crd *GlobalCurd) DeleteById(orgId int64, id string, dataFrom int) error {
	if !crd.masterModel.DeleteByOrgAndId(orgId, id) {
		return errors.New("删除失败")
	}

	if dataFrom == enum_process.DataFromProcess {
		CreateProcessCurdFactory().GlobalCallbackDeleteByHash(orgId, id)
	} else {
		CreateDllCurdFactory().GlobalCallbackDeleteByHash(orgId, id)
	}

	return nil
}

// DeleteByDataFrom 通过数据来源删除
func (crd *GlobalCurd) DeleteByDataFrom(orgId int64, dataFrom int) error {
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		DataFrom: dataFrom,
	}
	if !crd.masterModel.DeleteByConditions(condition) {
		return errors.New("删除失败")
	}
	if dataFrom == enum_process.DataFromProcess {
		CreateProcessCurdFactory().GlobalCallbackDeleteAll(orgId)
	} else {
		CreateDllCurdFactory().GlobalCallbackDeleteAll(orgId)
	}
	return nil
}

// DeleteByOrg 通过企业删除
func (crd *GlobalCurd) DeleteByOrg(orgId int64) {
	crd.masterModel.DeleteByOrg(orgId)
}
