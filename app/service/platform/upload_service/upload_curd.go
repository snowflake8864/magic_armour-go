package upload_service

import (
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *UploadCurd {
	return &UploadCurd{masterModel: upload_models.CreateUploadFactory()}
}

type UploadCurd struct {
	masterModel *upload_models.UploadModel
}

// Add 新增
func (crd *UploadCurd) Add(param *upload_models.UploadModel) bool {
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
