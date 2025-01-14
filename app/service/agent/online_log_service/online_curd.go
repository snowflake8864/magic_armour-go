package online_log_service

import (
	"errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_online_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *OnlineLogCurd {
	return &OnlineLogCurd{agent_online_models.CreateOnlineLogFactory()}
}

type OnlineLogCurd struct {
	masterModel *agent_online_models.OnlineLogModel
}

// Add 新增
func (crd *OnlineLogCurd) Add(param *agent_online_models.OnlineLogModel) error {
	param.Id = variable.SnowFlake.GetId()
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}
	return nil
}

// BatchAdd 新增
func (crd *OnlineLogCurd) BatchAdd(list *[]agent_online_models.OnlineLogModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []agent_online_models.OnlineLogModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	return nil
}
