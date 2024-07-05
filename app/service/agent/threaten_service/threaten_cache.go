package threaten_service

import (
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateCacheFactory() *ThreatenCache {
	return &ThreatenCache{masterModel: threaten_models.CreateCacheFactory()}
}

type ThreatenCache struct {
	masterModel *threaten_models.CacheModel
}

// Add 新增
func (crd *ThreatenCache) Add(param *threaten_models.CacheModel, checkExist bool) bool {
	if checkExist {
		if isExist, _ := crd.masterModel.GetLevelById(param.Id); isExist {
			return true
		}
	}

	param.Dir = path_helper.ToStandard(param.Dir)
	param.ParentDir = path_helper.ToStandard(param.ParentDir)
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	param.Level = enum_threaten.LevelDefault
	return crd.masterModel.InsertItem(param)
}

// DeleteById 通过ID删除
func (crd *ThreatenCache) DeleteById(id string) bool {
	return crd.masterModel.DeleteById(id)
}

// BatchDelete 通过集合ID删除
func (crd *ThreatenCache) BatchDelete(ids []string) bool {
	return crd.masterModel.DeleteByConditions(threaten_models.SearchConditions{
		Ids: ids,
	})
}

// BatchUnknown 批量设置成未知
func (crd *ThreatenCache) BatchUnknown(ids []string) bool {
	level := enum_threaten.LevelDefault
	wheres := threaten_models.SearchConditions{
		Ids:   ids,
		Level: &level,
	}
	updates := map[string]interface{}{
		"level":        enum_threaten.LevelUnknown,
		"check_source": enum_threaten.CheckSourceVeda,
		"check_at":     date_helper.CurDatetime(),
	}
	isOk, _ := crd.masterModel.UpdateByConditions(wheres, updates)
	return isOk
}

// Rewrite 根据第三方返回重写本地信息
func (crd *ThreatenCache) Rewrite(param *threaten_models.CacheModel) bool {
	checkAt := param.CheckAt
	if string_helper.IsEmpty(checkAt) {
		checkAt = date_helper.CurDatetime()
	}

	checkSource := enum_threaten.CheckSourceVeda
	if strings.Contains(param.Attribute, "malicious") && strings.Contains(param.Attribute, "detections") {
		checkSource = enum_threaten.CheckSourceInternet
	}

	updates := map[string]interface{}{
		"level":        param.Level,
		"attribute":    param.Attribute,
		"check_at":     checkAt,
		"check_source": checkSource,
	}
	if crd.masterModel.UpdateById(param.Id, updates) {
		return false
	}

	crd.rewriteAfter(param.Id, param.DataFrom, param.Level)

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
	wheres := threaten_models.SearchConditions{
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
	}

	updates := map[string]interface{}{
		"level": enum_threaten.LevelDefault,
	}
	isOk, _ := crd.masterModel.UpdateByConditions(wheres, updates)
	return isOk
}
