package software_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/files"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *SoftwareCurd {
	return &SoftwareCurd{masterModel: software_models.CreateSoftwareFactory()}
}

type SoftwareCurd struct {
	masterModel *software_models.SoftwareModel
}

// GetItemById 通过id获取
func (crd *SoftwareCurd) GetItemById(id int64) (model *software_models.SoftwareModel) {
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// GetMapsByPublished 最新发布数据, 有映射关系 os_type:ver
func (crd *SoftwareCurd) GetMapsByPublished() map[int]string {
	models := crd.GetItemsByPublished()
	if models == nil {
		return nil
	}
	var maps = make(map[int]string)
	for _, model := range models {
		maps[model.OsType] = model.Ver
	}
	return maps
}

// GetItemsByPublished 最新发布数据,
func (crd *SoftwareCurd) GetItemsByPublished() (items []*software_models.SoftwareModel) {
	IsPublish := enums.IsYes
	condition := software_models.SearchConditions{
		IsPublish: &IsPublish,
	}
	items = crd.masterModel.GetItemsByConditions(condition, 0, 0, "")
	return items
}

// Add 新增
func (crd *SoftwareCurd) Add(param *software_models.SoftwareModel) (errCode int, err error) {
	if string_helper.IsEmpty(dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentHttpUrl)) {
		return -1, errors.New(my_errors.CurdSoftwareUploadNeedConfig)
	}

	uploadModel := upload_models.CreateUploadFactory().GetItemById(param.UploadId)
	if uploadModel == nil || uploadModel.Id == 0 {
		return -1, errors.New(my_errors.CurdSoftwareNotExist)
	}
	if !files.ExistsFile(uploadModel.FullPath) {
		return -1, errors.New(my_errors.FilesUploadFailMsg)
	}
	param.StoragePath = path_helper.ToStandard(uploadModel.StoragePath)
	param.FullPath = path_helper.ToStandard(uploadModel.FullPath)
	param.IsPublish = enums.IsNot
	param.CreatedAt = date_helper.CurDatetime()

	if !crd.masterModel.InsertItem(param) {
		return -1, errors.New(my_errors.CurdSaveFailed)
	}

	return errCode, err
}

// DeleteById 删除
func (crd *SoftwareCurd) DeleteById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New(my_errors.CurdSoftwareNotExist)
	}

	if !crd.masterModel.DeleteById(id) {
		return -1, errors.New(my_errors.CurdDeleteFailed)
	}

	return errCode, err
}

// PublishById 发布
func (crd *SoftwareCurd) PublishById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New(my_errors.CurdSoftwareNotExist)
	}
	if model.IsPublish == enums.IsYes {
		return -1, errors.New(my_errors.CurdSoftwarePublished)
	}

	//关闭同系统类型的
	IsPublish := enums.IsYes
	conditions := software_models.SearchConditions{
		OsType:    model.OsType,
		IsPublish: &IsPublish,
	}
	updates := map[string]interface{}{
		"is_publish": enums.IsNot,
	}
	crd.masterModel.UpdateByConditions(conditions, updates)

	//修改本身
	updates = map[string]interface{}{
		"is_publish": enums.IsYes,
	}
	if !crd.masterModel.UpdateById(id, updates) {
		return -1, errors.New(my_errors.CurdSaveFailed)
	}
	return errCode, err
}

// CancelById  取消
func (crd *SoftwareCurd) CancelById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New(my_errors.CurdSoftwareNotExist)
	}

	//修改本身
	updates := map[string]interface{}{
		"is_publish": enums.IsNot,
	}
	if !crd.masterModel.UpdateById(id, updates) {
		return -1, errors.New(my_errors.CurdSaveFailed)
	}
	return errCode, err
}

// Clear 清空
func (crd *SoftwareCurd) Clear() bool {
	CreateDownloadFactory().Clear()
	return crd.masterModel.TruncateTable()
}
