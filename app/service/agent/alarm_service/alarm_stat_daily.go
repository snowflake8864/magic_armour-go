package alarm_service

import (
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/utils/helpers"
)

func CreateStatDailyFactory() *AlarmStatDaily {
	return &AlarmStatDaily{masterModel: alarm_stat_models.CreateAlarmStatDailyFactory()}
}

type AlarmStatDaily struct {
	masterModel *alarm_stat_models.AlarmStatDailyModel
}

// Add 新增
func (crd *AlarmStatDaily) Add(param *alarm_stat_models.AlarmStatDailyModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}

// SumByDay 计算： 每天
func (crd *AlarmStatDaily) SumByDay(day, column string) (result *alarm_stat_models.AlarmStatDailyModel) {
	where := alarm_stat_models.SearchConditions{
		StatDay: day,
	}
	result = crd.masterModel.CalcByConditions(where, column)
	if result != nil {
		return nil
	}
	return result
}

// SumByMonthGroupDay 计算： 每月并且按照天分组
func (crd *AlarmStatDaily) SumByMonthGroupDay(month int, column string) (items []*alarm_stat_models.AlarmStatDailyModel) {
	where := alarm_stat_models.SearchConditions{
		StatMonth: month,
	}
	items = crd.masterModel.CalcByConditionsGroupDay(where, column)
	if items != nil {
		return nil
	}
	return items
}
