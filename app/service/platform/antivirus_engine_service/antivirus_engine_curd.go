package antivirus_engine_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/model/platform/antivirus_engine_models"
	"goskeleton/app/utils/files"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
)

func CreateCurdFactory() *EngineCurd {
	return &EngineCurd{masterModel: antivirus_engine_models.CreateEngineFactory()}
}

type EngineCurd struct {
	masterModel *antivirus_engine_models.EngineModel
}

// GetItemById 通过id获取
func (crd *EngineCurd) GetItemById(id int64) (model *antivirus_engine_models.EngineModel) {
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// GetItemsByPublished 最新发布数据
func (crd *EngineCurd) GetItemsByPublished() (items []*antivirus_engine_models.EngineModel) {
	IsPublish := enums.IsYes
	condition := antivirus_engine_models.SearchConditions{
		IsPublish: &IsPublish,
	}
	items = crd.masterModel.GetItemsByConditions(condition, 0, 0, "")
	return items
}

// Add 新增
func (crd *EngineCurd) Add(param *antivirus_engine_models.EngineModel) (errCode int, err error) {
	uploadModel := upload_models.CreateUploadFactory().GetItemById(param.UploadId)
	if uploadModel == nil || uploadModel.Id == 0 {
		return -1, errors.New("上传软件不存在")
	}
	if !files.ExistsFile(uploadModel.FullPath) {
		//return -1, errors.New("文件不存在")
	}

	//如是是最新的
	if param.IsNewest == enums.IsYes {
		//关闭同系统类型的
		conditions := antivirus_engine_models.SearchConditions{
			PlatformType: param.PlatformType,
		}
		updates := map[string]interface{}{
			"is_newest": enums.IsNot,
		}
		crd.masterModel.UpdateByConditions(conditions, updates)
	}

	param.FileHash = uploadModel.FileHash
	param.FileSize = uploadModel.Size
	param.StoragePath = path_helper.ToStandard(uploadModel.StoragePath)
	param.FullPath = path_helper.ToStandard(uploadModel.FullPath)
	param.IsPublish = enums.IsNot
	param.CreatedAt = date_helper.CurDatetime()

	if !crd.masterModel.InsertItem(param) {
		return -1, errors.New("保存失败")
	}

	return errCode, err
}

// DeleteById 删除
func (crd *EngineCurd) DeleteById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New("软件不存在")
	}

	if !crd.masterModel.DeleteById(id) {
		return -1, errors.New("删除失败")
	}

	return errCode, err
}

// PublishById 发布
func (crd *EngineCurd) PublishById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New("软件不存在")
	}
	if model.IsPublish == enums.IsYes {
		return -1, errors.New("软件已经发布过了")
	}

	//关闭同系统类型的
	IsPublish := enums.IsYes
	conditions := antivirus_engine_models.SearchConditions{
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
		return -1, errors.New("保存失败")
	}
	return errCode, err
}

// CancelById  取消
func (crd *EngineCurd) CancelById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New("软件不存在")
	}

	//修改本身
	updates := map[string]interface{}{
		"is_publish": enums.IsNot,
	}
	if !crd.masterModel.UpdateById(id, updates) {
		return -1, errors.New("保存失败")
	}
	return errCode, err
}
