package upload_service

import (
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *UploadCurd {
	return &UploadCurd{masterModel: model_platform.CreateUploadFactory()}
}

type UploadCurd struct {
	masterModel *model_platform.UploadModel
}

// GetItemById 通过id获取
func (crd *UploadCurd) GetItemById(id int64) (model *model_platform.UploadModel) {
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// Add 新增
func (crd *UploadCurd) Add(param *model_platform.UploadModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
