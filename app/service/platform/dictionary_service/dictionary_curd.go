package dictionary_service

import (
	"goskeleton/app/model/model_platform"
)

func CreateDictionaryCurdFactory() *DictionaryCurd {
	return &DictionaryCurd{model_platform.CreateDictionaryFactory()}
}

type DictionaryCurd struct {
	masterModel *model_platform.DictionaryModel
}

func (crd *DictionaryCurd) GetMapByModule(module string) (result map[string]string) {
	items := crd.masterModel.GetItemsByModule(module)
	for _, obj := range items {
		result[obj.Id] = obj.Val
	}
	return result
}
