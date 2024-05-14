package organize_service

import (
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/utils/helpers"
)

func CreateOrgLogFactory() *OrgLog {
	return &OrgLog{masterModel: organization_models.CreateOrgLogFactory()}
}

type OrgLog struct {
	OrgBase
	masterModel *organization_models.OrgLogModel
}

// InitByOrg 初始化企业日志结构体
func (crd *OrgLog) InitByOrg(orgParam *organization_models.OrganizationModel, module int, netIp string) organization_models.OrgLogModel {
	param := organization_models.OrgLogModel{}
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
func (crd *OrgLog) Add(param *organization_models.OrgLogModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}

// Page 分页
func (crd *OrgLog) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*organization_models.OrgLogModel) {
	where := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}

// NewestList 最新列表
func (crd *OrgLog) NewestList(limit int) (items []*organization_models.OrgLogModel) {
	where := crd.search2Where(SearchConditions{})

	items, err := crd.masterModel.GetItemsByConditions(where, 0, limit, "created_at DESC")
	if err != nil {
		return nil
	}
	return items
}
