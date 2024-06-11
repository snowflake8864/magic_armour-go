package org_stat_service

import (
	"goskeleton/app/model/organization/org_stat_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateStatDailyFactory() *StatDaily {
	return &StatDaily{masterModel: org_stat_models.CreateDailyFactory()}
}

type StatDaily struct {
	masterModel *org_stat_models.DailyModel
}

// Add 新增
func (crd *StatDaily) Add(param *org_stat_models.DailyModel) bool {
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	return crd.masterModel.InsertItem(param)
}

// SumByDay 计算： 每天
func (crd *StatDaily) SumByDay(day, column string) (result *org_stat_models.DailyModel) {
	search := org_stat_models.SearchConditions{
		StatDay: day,
	}
	result = crd.masterModel.CalcByConditions(search, column)
	if result != nil {
		return result
	}
	return result
}

// SumByMonthGroupDay 计算： 每月并且按照天分组
func (crd *StatDaily) SumByMonthGroupDay(month int, column string) (items []*org_stat_models.DailyModel) {
	search := org_stat_models.SearchConditions{
		StatMonth: month,
	}
	items = crd.masterModel.CalcByConditionsGroupDay(search, column)
	if len(items) > 0 {
		return nil
	}
	return items
}
