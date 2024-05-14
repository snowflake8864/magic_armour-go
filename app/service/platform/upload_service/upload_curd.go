package upload_service

import (
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *UploadCurd {
	return &UploadCurd{masterModel: upload_models.CreateUploadFactory()}
}

type UploadCurd struct {
	masterModel *upload_models.UploadModel
}

// Add 新增
func (crd *UploadCurd) Add(param *upload_models.UploadModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
