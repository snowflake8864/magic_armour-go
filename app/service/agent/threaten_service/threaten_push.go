package threaten_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
	"strings"
)

func CreatePushFactory() *ThreatenPush {
	return &ThreatenPush{masterModel: model_agent.CreateThreatenPushFactory()}
}

type ThreatenPush struct {
	ThreatenBase
	masterModel *model_agent.ThreatenPushModel
}

type PushBo struct {
	Id        string `json:"id"`
	Level     int    `json:"level"`
	Dir       string `json:"dir"`
	ParentDir string `json:"parent_dir"`
	Os        string `json:"os"`
	Introduce string `json:"introduce"`
	Copyright string `json:"copyright"`
	Signature string `json:"signature"`
	Attribute string `json:"attribute"`
	DataFrom  int    `json:"data_from"`
}

// GetMapsByIds 通过ids获取, 有映射关系
func (crd *ThreatenPush) GetMapsByIds(ids []string) map[string]*model_agent.ThreatenPushModel {
	result := make(map[string]*model_agent.ThreatenPushModel)

	models, err := crd.masterModel.GetItemsByIds(ids)
	if err != nil || models == nil {
		return nil
	}
	for _, model := range models {
		result[model.Id] = model
	}
	return result
}

// Page 分页
func (crd *ThreatenPush) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_agent.ThreatenPushModel) {
	wheres := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(wheres, offset, pageSize, order)
	return total, items
}

// Clear 清空
func (crd *ThreatenPush) Clear() bool {
	return crd.masterModel.TruncateTable()
}

// BatchAdd 批量添加
func (crd *ThreatenPush) BatchAdd(list []model_agent.ThreatenPushModel, level int, operateId int64, operateName string) (errCode int, err error) {
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

	var inserts []model_agent.ThreatenPushModel
	for _, model := range list {
		tmpId := strings.ToLower(model.Id)
		if _, isOK := existPushes[tmpId]; !isOK {
			inserts = append(inserts, model_agent.ThreatenPushModel{
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
