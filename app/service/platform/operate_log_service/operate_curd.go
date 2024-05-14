package operate_log_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/platform/operate_models"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *OperateLogCurd {
	return &OperateLogCurd{masterModel: operate_models.CreateOperateLogFactory()}
}

type OperateLogCurd struct {
	masterModel *operate_models.OperateLogModel
}

type SearchConditions struct {
	Module     string
	NetIp      string
	CreateName string
}

// search2Where 查询条件转化
func (crd *OperateLogCurd) search2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if !helpers.StringIsEmpty(search.Module) && search.Module != enums.AllDefaultStr {
		wheres["module"] = search.Module
	}
	if !helpers.StringIsEmpty(search.NetIp) {
		wheres["net_ip"] = search.NetIp
	}
	if !helpers.StringIsEmpty(search.CreateName) {
		wheres["create_name"] = search.CreateName
	}
	return wheres
}

// Add 新增
func (crd *OperateLogCurd) Add(param *operate_models.OperateLogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
