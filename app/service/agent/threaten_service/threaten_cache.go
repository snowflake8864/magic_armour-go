package threaten_service

import (
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/proc_service"
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

// BatchAdd 新增
func (crd *ThreatenCache) BatchAdd(list *[]threaten_models.CacheModel, checkExist bool) bool {
	if len(*list) == 0 {
		return true
	}
	hasMap := make(map[string]int)

	//校验存在的
	if checkExist {
		var ids []string
		for _, model := range *list {
			ids = append(ids, model.Id)
		}
		hasMap = crd.masterModel.GetLevelsByIds(ids)
	}
	createdAt := date_helper.CurDatetime()
	var inserts []threaten_models.CacheModel
	for _, model := range *list {
		if _, isOK := hasMap[model.Id]; !isOK {
			model.Dir = path_helper.ToStandard(model.Dir)
			model.ParentDir = path_helper.ToStandard(model.ParentDir)
			if string_helper.IsEmpty(model.CreatedAt) {
				model.CreatedAt = createdAt
			}
			model.Level = enum_threaten.LevelDefault

			if !string_helper.IsEmpty(model.Id) {
				inserts = append(inserts, model)
			}
		}
	}

	return crd.masterModel.InsertBatch(&inserts)
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
		"level":    enum_threaten.LevelUnknown,
		"check_at": date_helper.CurDatetime(),
	}
	isOk, _ := crd.masterModel.UpdateByConditions(wheres, updates)
	return isOk
}

// FastBatchSetHasLevelOnly 快速批量设置仅有等级的
func (crd *ThreatenCache) FastBatchSetHasLevelOnly(ids []string, level int) bool {
	wheres := threaten_models.SearchConditions{
		Ids: ids,
	}
	updates := map[string]interface{}{
		"level":    level,
		"check_at": date_helper.CurDatetime(),
	}
	isOk, _ := crd.masterModel.UpdateByConditions(wheres, updates)
	return isOk
}

// FastBatchSetHasLevelAndInfluenceProc 快速批量设置有等级和影响进程的
func (crd *ThreatenCache) FastBatchSetHasLevelAndInfluenceProc(ids []string, level, dataForm int) bool {
	wheres := threaten_models.SearchConditions{
		Ids: ids,
	}
	updates := map[string]interface{}{
		"level":    level,
		"check_at": date_helper.CurDatetime(),
	}
	if isOk, _ := crd.masterModel.UpdateByConditions(wheres, updates); !isOk {
		return false
	}

	if enum_threaten.IsWhite(level) || enum_threaten.IsBlack(level) {
		if dataForm == enum_threaten.DataFromProcess { //来源进程
			proc_service.CreateProcessCurdFactory().ThreatenCallbackByHashes(ids, level)
		} else if dataForm == enum_threaten.DataFromDll { //来源模块
			proc_service.CreateDllCurdFactory().ThreatenCallbackByHashes(ids, level)
		} else {
			proc_service.CreateProcessCurdFactory().ThreatenCallbackByHashes(ids, level)
			proc_service.CreateDllCurdFactory().ThreatenCallbackByHashes(ids, level)
		}
	}

	return true
}

// Rewrite 根据第三方返回重写本地信息
func (crd *ThreatenCache) Rewrite(param *threaten_models.CacheModel) bool {
	checkAt := param.CheckAt
	if string_helper.IsEmpty(checkAt) {
		checkAt = date_helper.CurDatetime()
	}

	checkSource := enum_threaten.CheckSourceVeda
	if strings.Contains(param.CheckInfo, "malicious") && strings.Contains(param.CheckInfo, "detections") {
		checkSource = enum_threaten.CheckSourceInternet
	}

	updates := map[string]interface{}{
		"level":        param.Level,
		"check_info":   param.CheckInfo,
		"check_at":     checkAt,
		"check_source": checkSource,
	}
	if !crd.masterModel.UpdateById(param.Id, updates) {
		return false
	}

	if enum_threaten.IsWhite(param.Level) || enum_threaten.IsBlack(param.Level) {
		if param.DataFrom == enum_threaten.DataFromProcess { //来源进程
			proc_service.CreateProcessCurdFactory().ThreatenCallbackByHash(param.Id, param.Level)
		} else if param.DataFrom == enum_threaten.DataFromDll { //来源模块
			proc_service.CreateDllCurdFactory().ThreatenCallbackByHash(param.Id, param.Level)
		} else {
			proc_service.CreateProcessCurdFactory().ThreatenCallbackByHash(param.Id, param.Level)
			proc_service.CreateDllCurdFactory().ThreatenCallbackByHash(param.Id, param.Level)
		}
	}

	return true
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
