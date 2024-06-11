package org_log_service

import (
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateOrgLogFactory() *OrgLog {
	return &OrgLog{masterModel: organization_models.CreateLogFactory()}
}

type OrgLog struct {
	masterModel *organization_models.LogModel
}

// InitByOrg 初始化企业日志结构体
func (crd *OrgLog) InitByOrg(orgParam *organization_models.OrganizationModel, module int, netIp string) organization_models.LogModel {
	param := organization_models.LogModel{}
	param.Module = module
	param.NetIp = netIp
	param.OrgId = orgParam.Id
	param.OrgName = orgParam.Name
	param.UserName = orgParam.UserName
	param.UserPhone = orgParam.UserPhone
	param.UserEmail = orgParam.UserEmail
	return param
}

// Add 新增企业日志
func (crd *OrgLog) Add(param *organization_models.LogModel) bool {
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	return crd.masterModel.InsertItem(param)
}

// NewestList 最新列表
func (crd *OrgLog) NewestList(limit int) (items []*organization_models.LogModel) {
	where := organization_models.SearchConditions{}

	items = crd.masterModel.GetItemsByConditions(where, 0, limit, "created_at DESC")

	return items
}
