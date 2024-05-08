package alarm_service

import (
	"goskeleton/app/model/model_alarm"
	"goskeleton/app/utils/helpers"
)

func CreateStatDailyFactory() *AlarmStatDaily {
	return &AlarmStatDaily{masterModel: model_alarm.CreateAlarmStatDailyFactory()}
}

type AlarmStatDaily struct {
	AlarmBase
	masterModel *model_alarm.AlarmStatDailyModel
}

// Add 新增
func (crd *AlarmStatDaily) Add(param *model_alarm.AlarmStatDailyModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}

// SumByDay 计算： 按日期
func (crd *AlarmStatDaily) SumByDay(day, column string) (result *model_alarm.AlarmStatDailyModel) {
	search := SearchConditions{
		StatDay: day,
	}
	where := crd.search2Where(search)
	result, err := crd.masterModel.CalcByConditions(where, column)
	if err != nil {
		return result
	}
	return result
}
