package organize_service

import (
	"goskeleton/app/model/model_organization"
	"goskeleton/app/utils/helpers"
)

func CreateOrgLogFactory() *OrgLog {
	return &OrgLog{masterModel: model_organization.CreateOrgLogFactory()}
}

type OrgLog struct {
	OrgBase
	masterModel *model_organization.OrgLogModel
}

// InitByOrg 初始化企业日志结构体
func (crd *OrgLog) InitByOrg(orgParam *model_organization.OrganizationModel, module int, netIp string) model_organization.OrgLogModel {
	param := model_organization.OrgLogModel{}
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
func (crd *OrgLog) Add(param *model_organization.OrgLogModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}

// Page 分页
func (crd *OrgLog) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_organization.OrgLogModel) {
	where := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}

// NewestList 最新列表
func (crd *OrgLog) NewestList(limit int) (items []*model_organization.OrgLogModel) {
	where := crd.search2Where(SearchConditions{})

	items, err := crd.masterModel.GetItemsByConditions(where, 0, limit, "created_at DESC")
	if err != nil {
		return nil
	}
	return items
}
