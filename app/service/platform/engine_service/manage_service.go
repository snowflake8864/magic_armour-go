package engine_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/model/platform/engine_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/files"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *ManageService {
	return &ManageService{masterModel: engine_models.CreateEngineFactory()}
}

type ManageService struct {
	masterModel *engine_models.EngineModel
}

// GetItemById 通过id获取
func (s *ManageService) GetItemById(id int64) (model *engine_models.EngineModel) {
	model, _ = s.masterModel.GetItemById(id)
	return model
}

// Add 新增
func (s *ManageService) Add(param *engine_models.EngineModel) (errCode int, err error) {

	param.CreatedAt = date_helper.CurDatetime()

	if !s.masterModel.InsertItem(param) {
		return -1, errors.New("保存失败")
	}

	return errCode, err
}

// DeleteById 删除
func (s *ManageService) DeleteById(id int64) (errCode int, err error) {
	model := s.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New("软件不存在")
	}
	// 获取当前引擎下所有的版本信息(删除引擎之前，对应的版本信息要清空之后)

	// 删除引擎
	if !s.masterModel.DeleteById(id) {
		return -1, errors.New("删除失败")
	}

	return errCode, err
}

// PublishById 发布
func (crd *ManageService) PublishById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New("软件不存在")
	}
	if model.PublishStatus == enums.IsYes {
		return -1, errors.New("软件已经发布过了")
	}

	//关闭同系统类型的
	IsPublish := enums.IsYes
	conditions := engine_models.SearchConditions{
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
		return -1, errors.New("保存失败")
	}
	return errCode, err
}

// CancelById  取消
func (crd *ManageService) CancelById(id int64) (errCode int, err error) {
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

// Clear 清空
func (crd *ManageService) Clear() bool {
	CreateDownloadFactory().Clear()
	return crd.masterModel.TruncateTable()
}
