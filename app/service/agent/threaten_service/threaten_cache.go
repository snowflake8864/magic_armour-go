package threaten_service

import (
	"fmt"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func CreateCacheFactory() *ThreatenCache {
	return &ThreatenCache{masterModel: model_agent.CreateThreatenCacheFactory()}
}

type ThreatenCache struct {
	ThreatenBase
	masterModel *model_agent.ThreatenCacheModel
}

// GetLevelById 通过id获取
func (crd *ThreatenCache) GetLevelById(id string) (isExist bool, level int) {
	model, _ := crd.masterModel.GetItemById(id, "id, level")
	if helpers.StringIsEmpty(model.Id) {
		return false, enum_threaten.LevelUnknown
	}
	return true, model.Level
}

// GetItemById 通过id获取
func (crd *ThreatenCache) GetItemById(id string) (model *model_agent.ThreatenCacheModel) {
	model, _ = crd.masterModel.GetItemById(id, "*")
	return model
}

// GetMapsByIds 通过ids获取, 有映射关系
func (crd *ThreatenCache) GetMapsByIds(ids []string) map[string]*model_agent.ThreatenCacheModel {
	result := make(map[string]*model_agent.ThreatenCacheModel)

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
func (crd *ThreatenCache) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_agent.ThreatenCacheModel) {
	wheres := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(wheres, offset, pageSize, order)
	return total, items
}

// Add 新增
func (crd *ThreatenCache) Add(param *model_agent.ThreatenCacheModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}

// DeleteById 通过订单删除
func (crd *ThreatenCache) DeleteById(ID string) bool {
	return crd.masterModel.DeleteById(ID)
}

// BatchUnknown 批量设置成未知
func (crd *ThreatenCache) BatchUnknown(ids []string) bool {
	wheres := crd.search2Where(SearchConditions{
		Ids:   ids,
		Level: enum_threaten.LevelNotSync,
	})
	updates := map[string]interface{}{
		"level":    enum_threaten.LevelUnknown,
		"check_at": helpers.DatetimeCurr(),
	}
	return crd.masterModel.UpdateByConditions(wheres, updates)
}

// Rewrite 根据第三方返回重写本地信息
func (crd *ThreatenCache) Rewrite(param *model_agent.ThreatenCacheModel) bool {
	updates := map[string]interface{}{
		"level":     param.Level,
		"attribute": param.Attribute,
		"check_at":  helpers.DatetimeCurr(),
	}
	if crd.masterModel.UpdateById(param.Id, updates) {
		return false
	}

	go crd.rewriteAfter(param.Id, param.DataFrom, param.Level)

	return true
}

func (crd *ThreatenCache) rewriteAfter(id string, dataFrom, level int) {
	if enum_threaten.IsWhite(level) {
		if dataFrom == enum_threaten.DataFromProcess {
			//进程白名单
		} else if dataFrom == enum_threaten.DataFromDll {
			//模块白名单
		} else if dataFrom == enum_threaten.DataFromUnknown {
			//进程白名单
			//模块白名单
		}
	} else if enum_threaten.IsBlack(level) {
		if dataFrom == enum_threaten.DataFromProcess {
			//进程黑名单
		} else if dataFrom == enum_threaten.DataFromDll {
			//模块黑名单
		} else if dataFrom == enum_threaten.DataFromUnknown {
			//进程黑名单
			//模块黑名单
		}
	}
}

// Clear 清空
func (crd *ThreatenCache) Clear() bool {
	return crd.masterModel.TruncateTable()
}

// Refresh 刷新
func (crd *ThreatenCache) Refresh() bool {
	wheres := crd.search2Where(SearchConditions{
		LevelList: []int{
			enum_threaten.LevelUnknown,
			enum_threaten.LevelWhite,
			enum_threaten.LevelDefaultOne,
			enum_threaten.LevelDefaultTwo,
			enum_threaten.LevelDefaultThree,
			enum_threaten.LevelDefaultFour,
			enum_threaten.LevelDefaultFive,
			enum_threaten.LevelDefaultSix,
			enum_threaten.LevelDefaultSeven,
			enum_threaten.LevelDefaultEight,
			enum_threaten.LevelDefaultNine,
			enum_threaten.LevelBlack,
		},
	})
	fmt.Println(wheres)
	updates := map[string]interface{}{
		"level": enum_threaten.LevelNotSync,
	}
	return crd.masterModel.UpdateByConditions(wheres, updates)
}
