package proc_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/agent/proc_models"
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

// Add 添加
func (crd *GlobalCurd) Add(param *proc_models.ProcGlobalModel) error {
	param.Id = strings.ToLower(param.Id)
	param.Dir = path_helper.ToStandard(param.Dir)

	hasModel := crd.masterModel.GetItemByOrgAndId(param.OrgId, param.Id)
	if hasModel != nil && hasModel.OrgId != 0 {
		fmt.Println(hasModel.OrgId)
		return errors.New("数据已存在1")
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
				CreatedAt:  now,
				UpdatedAt:  now,
			})
		}
	}
	if len(inserts) == 0 {
		return nil
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	if param.DataFrom == enum_process.DataFromProcess {
		CreateProcessCurdFactory().GlobalCallbackAddByHashes(param.OrgId, param.Type, hashes)
	} else {

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

	}
	return nil
}

// DeleteByOrg 通过企业删除
func (crd *GlobalCurd) DeleteByOrg(orgId int64) {
	crd.masterModel.DeleteByOrg(orgId)
}
