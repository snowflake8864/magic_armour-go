package software_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_platform"
)

func CreateCurdFactory() *SoftwareCurd {
	return &SoftwareCurd{masterModel: model_platform.CreateSoftwareFactory()}
}

type SoftwareCurd struct {
	masterModel *model_platform.SoftwareModel
}

// GetItemById 通过id获取
func (crd *SoftwareCurd) GetItemById(id int64) (model *model_platform.SoftwareModel) {
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
func (crd *SoftwareCurd) GetItemsByPublished() (items []*model_platform.SoftwareModel) {
	condition := map[string]interface{}{
		"is_publish": enums.IsYes,
	}
	items, _ = crd.masterModel.GetItemsByConditions(condition, 0, 0, "")
	return items
}
