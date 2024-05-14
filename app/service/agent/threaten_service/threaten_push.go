package threaten_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/utils/helpers"
	"strings"
)

func CreatePushFactory() *ThreatenPush {
	return &ThreatenPush{masterModel: threaten_models.CreateThreatenPushFactory()}
}

type ThreatenPush struct {
	masterModel *threaten_models.ThreatenPushModel
}

// GetMapsByIds 通过ids获取, 有映射关系
func (crd *ThreatenPush) GetMapsByIds(ids []string) map[string]*threaten_models.ThreatenPushModel {
	result := make(map[string]*threaten_models.ThreatenPushModel)

	models := crd.masterModel.GetItemsByIds(ids)
	if models == nil {
		return nil
	}
	for _, model := range models {
		result[model.Id] = model
	}
	return result
}

// Clear 清空
func (crd *ThreatenPush) Clear() bool {
	return crd.masterModel.TruncateTable()
}

// BatchAdd 批量添加
func (crd *ThreatenPush) BatchAdd(list []threaten_models.ThreatenPushModel, level int, operateId int64, operateName string) (errCode int, err error) {
	//取出id
	var ids []string
	for _, model := range list {
		ids = append(ids, strings.ToLower(model.Id))
	}
	//本地威胁情报缓存中是有等级， 忽略
	//existCaches := CreateCacheFactory().GetMapsByIds(ids)
	////本地威胁情报推送中已存在， 忽略
	existPushes := crd.GetMapsByIds(ids)
	currTime := helpers.DatetimeCurr()

	var inserts []threaten_models.ThreatenPushModel
	for _, model := range list {
		tmpId := strings.ToLower(model.Id)
		if _, isOK := existPushes[tmpId]; !isOK {
			inserts = append(inserts, threaten_models.ThreatenPushModel{
				Id:         tmpId,
				Level:      level,
				Dir:        model.Dir,
				ParentDir:  model.ParentDir,
				Os:         model.Os,
				Introduce:  model.Introduce,
				Copyright:  model.Copyright,
				Signature:  model.Signature,
				Attribute:  model.Attribute,
				DataFrom:   model.DataFrom,
				PushStatus: enum_threaten.PushStatusYes,
				CreateId:   operateId,
				CreateName: operateName,
				CreatedAt:  currTime,
			})
		}
	}
	if len(inserts) == 0 {
		return errCode, err
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return -3, errors.New("保数据失败")
	}

	return errCode, err
}
