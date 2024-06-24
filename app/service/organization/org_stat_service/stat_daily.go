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
