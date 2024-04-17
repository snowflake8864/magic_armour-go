package organize_service

import (
	"goskeleton/app/model/model_organization"
)

func CreateOrgCurdFactory() *OrgCurd {
	return &OrgCurd{masterModel: model_organization.CreateOrganizationFactory()}
}

type OrgCurd struct {
	OrgBase
	masterModel *model_organization.OrganizationModel
}

// GetItemByOrg 通过企业编号获取
func (crd *OrgCurd) GetItemByOrg(orgId int64) (model *model_organization.OrganizationModel) {
	model, _ = crd.masterModel.GetItemById(orgId)
	return model
}

// GetItemByName 通过企业名称获取
func (crd *OrgCurd) GetItemByName(name string) (model *model_organization.OrganizationModel) {
	model, _ = crd.masterModel.GetItemByName(name)
	return model
}

// Page 分页
func (crd *OrgCurd) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_organization.OrganizationModel) {
	where := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}
