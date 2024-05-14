package monitor_service

import (
	"goskeleton/app/model/platform/server_monitor_models"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *MonitorCurd {
	return &MonitorCurd{masterModel: server_monitor_models.CreateServerMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *server_monitor_models.ServerMonitorModel
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
func (crd *MonitorCurd) GetList(search SearchConditions) (items []*server_monitor_models.ServerMonitorModel) {
	where := crd.search2Where(search)
	order := " id ASC"
	items, err := crd.masterModel.GetAllByConditions(where, order)
	if err != nil {
		return nil
	}
	return items
}

// Add 新增
func (crd *MonitorCurd) Add(param *server_monitor_models.ServerMonitorModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
