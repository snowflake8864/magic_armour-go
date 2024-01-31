package dictionary_service

import (
	"fmt"
	"goskeleton/app/model/model_platform"
)

func CreateDictionaryCurdFactory() *DictionaryCurd {
	return &DictionaryCurd{model_platform.CreateDictionaryFactory()}
}

type DictionaryCurd struct {
	masterModel *model_platform.DictionaryModel
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
func (crd *DictionaryCurd) GetAll() (items []model_platform.DictionaryModel) {
	items = crd.masterModel.GetAll()
	return items
}

// UpdateValById 通过Id修改值
func (crd *DictionaryCurd) UpdateValById(id, val string) bool {
	err := crd.masterModel.UpdateValById(id, val)
	fmt.Println(err)
	if err != nil {
		return false
	}
	return true
}
