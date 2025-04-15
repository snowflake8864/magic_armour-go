package threaten_service

import (
	"goskeleton/app/model/agent/threaten_models"
)

func CreateImportFactory() *ThreatenImport {
	return &ThreatenImport{
		masterModel: threaten_models.CreateCacheFactory(),
	}
}

type ThreatenImport struct {
	masterModel *threaten_models.CacheModel
}

// ParseExcel 解析Excel
func (crd *ThreatenImport) ParseExcel(filename string) error {

	return nil
}
