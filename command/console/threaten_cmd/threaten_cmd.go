package threaten_cmd

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/threaten_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCmdFactory() *ThreatenCmd {
	return &ThreatenCmd{}
}

type ThreatenCmd struct {
}

// getThreatenCli  获取威胁情报中心客户端
func (crd *ThreatenCmd) getThreatenCli(isOnline, isToken bool) (err error, cli *threaten_platform.RequestThreaten) {
	dictIds := []string{
		enum_dictionary.ThreatenSwitch,
		enum_dictionary.ThreatenUrl,
		enum_dictionary.ThreatenAppId,
		enum_dictionary.ThreatenAppKey,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if string_helper.ToInt(vals[enum_dictionary.ThreatenSwitch]) == enums.SwitchOff {
		return errors.New("威胁情报功能已关闭"), nil
	}

	//客户端
	cli = threaten_platform.CreateRequestFactory(threaten_platform.RequestConfig{
		Domain:    vals[enum_dictionary.ThreatenUrl],
		AppId:     vals[enum_dictionary.ThreatenAppId],
		AppSecret: vals[enum_dictionary.ThreatenAppKey],
	})

	//验证在线
	if isOnline {
		if err = cli.Online(); err != nil {
			return errors.New("网络连接失败： " + err.Error()), nil
		}
	}

	//生成token
	if isToken {
		if err = cli.GenerateToken(); err != nil {
			return errors.New("生成口令： " + err.Error()), nil
		}
	}

	return nil, cli
}

// SyncForPushes  推送列表： 同步数据到威胁情报中心
func (crd *ThreatenCmd) SyncForPushes() {
	variable.ZapLog.Sugar().Debug("SyncForPushes")
}

// AddFromProcessForCache  本地缓存：从进程添加
func (crd *ThreatenCmd) AddFromProcessForCache() {
	variable.ZapLog.Sugar().Debug("AddFromProcessForCache")
}

// AddFromDllForCache  本地缓存：从模块添加
func (crd *ThreatenCmd) AddFromDllForCache() {
	variable.ZapLog.Sugar().Debug("AddFromDllForCache")
}

// DeleteForCache  本地缓存：删除
func (crd *ThreatenCmd) DeleteForCache() {
	variable.ZapLog.Sugar().Debug("DeleteForCache start")

	//获取客户端
	if err, _ := crd.getThreatenCli(false, false); err != nil {
		variable.ZapLog.Sugar().Debug("DeleteForCache ", err.Error())
		return
	}

	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.ThreatenCacheExpire)
	expire := string_helper.ToInt(val)
	variable.ZapLog.Sugar().Debug("DeleteForCache expire:", expire)
	_, _, dayEnd, _ := date_helper.CalcByDay(-expire)

	conditions := threaten_models.SearchConditions{
		CreateAtLte: dayEnd,
	}
	cacheFact := threaten_models.CreateCacheFactory()
	page := 0
	limit := enum_threaten.CenterSupportNumber
	for {
		page++
		offset := (page - 1) * limit
		models := cacheFact.GetItemsByConditions(conditions, offset, limit, "")
		if models == nil || len(models) == 0 {
			break
		}
		var ids []string
		for _, model := range models {
			ids = append(ids, model.Id)
		}

		threaten_service.CreateCacheFactory().BatchDelete(ids)
	}

	variable.ZapLog.Sugar().Debug("DeleteForCache end page:", page)
}

// SyncDefaultForCache  本地缓存：同步默认状态(未同步)数据到威胁情报中心
func (crd *ThreatenCmd) SyncDefaultForCache() {
	variable.ZapLog.Sugar().Debug("SyncDefaultForCache")
}

// SyncUnknownForCache  本地缓存：同步未知状态数据到威胁情报中心
func (crd *ThreatenCmd) SyncUnknownForCache() {
	variable.ZapLog.Sugar().Debug("SyncUnknownForCache")
}
