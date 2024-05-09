package operate_log_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *OperateLogCurd {
	return &OperateLogCurd{masterModel: model_platform.CreateOperateLogFactory()}
}

type OperateLogCurd struct {
	masterModel *model_platform.OperateLogModel
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

// Page 分页
func (crd *OperateLogCurd) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_platform.OperateLogModel) {
	where := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}

// Add 新增
func (crd *OperateLogCurd) Add(param *model_platform.OperateLogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
