package dictionary_service

import (
	"goskeleton/app/model/platform/dictionary_models"
)

func CreateCurdFactory() *DictionaryCurd {
	return &DictionaryCurd{dictionary_models.CreateDictionaryFactory()}
}

type DictionaryCurd struct {
	masterModel *dictionary_models.DictionaryModel
}

// GetMapByModule 通过模块获取map[id][val]
func (crd *DictionaryCurd) GetMapByModule(module string) map[string]string {
	result := make(map[string]string)

	items := crd.masterModel.GetItemsByModule(module)
	for _, obj := range items {
		result[obj.Id] = obj.Val
	}
	return result
}

// GetAll 获取全部
func (crd *DictionaryCurd) GetAll() (items []*dictionary_models.DictionaryModel) {
	items = crd.masterModel.GetAll()
	return items
}

// UpdateValById 通过Id修改值
func (crd *DictionaryCurd) UpdateValById(id string, val interface{}) bool {
	err := crd.masterModel.UpdateValById(id, val)
	if err != nil {
		return false
	}
	_ = CreateCacheFactory().Rebuild()
	return true
}

// BatchUpdate 批量修改
func (crd *DictionaryCurd) BatchUpdate(idAndVal map[string]interface{}) bool {
	for id, val := range idAndVal {
		_ = crd.masterModel.UpdateValById(id, val)
	}

	_ = CreateCacheFactory().Rebuild()

	return true
}
