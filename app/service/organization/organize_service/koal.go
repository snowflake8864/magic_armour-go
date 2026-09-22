package organize_service

import (
	"goskeleton/app/model/organization/organization_models"
)

func CreateKoalFactory() *Koal {
	return &Koal{orgModel: organization_models.CreateOrganizationFactory()}
}

type Koal struct {
	orgModel *organization_models.OrganizationModel
}

func (crd *Koal) Sync(orgParam *organization_models.OrganizationModel) error {

	//是否存在
	if hasModel, _ := crd.orgModel.GetItemById(orgParam.Id); hasModel == nil || hasModel.Id == 0 {

	} else {

	}

	return nil
}
