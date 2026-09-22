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
