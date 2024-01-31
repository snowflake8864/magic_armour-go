package dictionary_service

import (
	modelPlatform "goskeleton/app/model/model_platform"
)

func CreateDictionaryCurdFactory() *DictionaryCurd {
	return &DictionaryCurd{modelPlatform.CreateDictionaryFactory()}
}

type DictionaryCurd struct {
	masterModel *modelPlatform.DictionaryModel
}

func (crd *DictionaryCurd) GetItemsById(id string) (result *modelPlatform.DictionaryModel) {
	result, _ = crd.masterModel.GetItemById(id)
	return result
}
