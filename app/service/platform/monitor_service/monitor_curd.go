package monitor_service

import (
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *MonitorCurd {
	return &MonitorCurd{masterModel: model_platform.CreateServerMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *model_platform.ServerMonitorModel
}

type SearchConditions struct {
	MonitorDay string `json:"monitor_day"`
	ServerCode string `json:"server_code"`
}

// search2Where 查询条件转化
func (crd *MonitorCurd) search2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if !helpers.StringIsEmpty(search.MonitorDay) {
		wheres["monitor_day"] = search.MonitorDay
	}
	if !helpers.StringIsEmpty(search.ServerCode) {
		wheres["ServerCode"] = search.ServerCode
	}
	return wheres
}

// GetList 获取列表
func (crd *MonitorCurd) GetList(search SearchConditions) (items []*model_platform.ServerMonitorModel) {
	where := crd.search2Where(search)
	order := " id ASC"
	items, err := crd.masterModel.GetAllByConditions(where, order)
	if err != nil {
		return nil
	}
	return items
}

// Add 新增
func (crd *MonitorCurd) Add(param *model_platform.ServerMonitorModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
