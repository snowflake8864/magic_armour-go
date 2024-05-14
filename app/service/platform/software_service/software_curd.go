package software_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/files"
	"goskeleton/app/utils/helpers"
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
func (crd *SoftwareCurd) GetMapsByPublished() (maps map[int]string) {
	models := crd.GetItemsByPublished()
	if models == nil {
		return nil
	}
	for _, model := range models {
		maps[model.OsType] = model.Ver
	}
	return maps
}

// GetItemsByPublished 最新发布数据,
func (crd *SoftwareCurd) GetItemsByPublished() (items []*software_models.SoftwareModel) {
	condition := map[string]interface{}{
		"is_publish": enums.IsYes,
	}
	items, _ = crd.masterModel.GetItemsByConditions(condition, 0, 0, "")
	return items
}

// Page 分页
func (crd *SoftwareCurd) Page(where map[string]interface{}, offset, pageSize int, orderField, orderSort string) (total int64, items []*software_models.SoftwareModel) {
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}

// Add 新增
func (crd *SoftwareCurd) Add(param *software_models.SoftwareModel) (errCode int, err error) {
	if helpers.StringIsEmpty(dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentHttpUrl)) {
		return -1, errors.New("上传软件前需要配置终端访问地址")
	}

	uploadModel := upload_models.CreateUploadFactory().GetItemById(param.UploadId)
	if uploadModel == nil || uploadModel.Id == 0 {
		return -1, errors.New("上传软件不存在")
	}
	if !files.ExistsFile(uploadModel.FullPath) {
		return -1, errors.New("文件不存在")
	}
	param.StoragePath = uploadModel.StoragePath
	param.FullPath = uploadModel.FullPath
	param.IsPublish = enums.IsNot
	param.CreatedAt = helpers.DatetimeCurr()

	if !crd.masterModel.InsertItem(param) {
		return -1, errors.New("保存失败")
	}

	return errCode, err
}

// DeleteById 删除
func (crd *SoftwareCurd) DeleteById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New("软件不存在")
	}

	if !crd.masterModel.DeleteById(id) {
		return -1, errors.New("删除失败")
	}

	return errCode, err
}

// PublishById 删除
func (crd *SoftwareCurd) PublishById(id int64) (errCode int, err error) {
	model := crd.GetItemById(id)
	if model == nil || model.Id == 0 {
		return -1, errors.New("软件不存在")
	}
	if model.IsPublish == enums.IsYes {
		return -1, errors.New("软件已经发布过了")
	}

	//关闭同系统类型的
	conditions := map[string]interface{}{
		"os_type":    model.OsType,
		"is_publish": enums.IsYes,
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
func (crd *SoftwareCurd) CancelById(id int64) (errCode int, err error) {
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
func (crd *SoftwareCurd) Clear() bool {
	CreateDownloadFactory().Clear()
	return crd.masterModel.TruncateTable()
}
