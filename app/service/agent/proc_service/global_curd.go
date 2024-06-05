package proc_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
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
