package threaten_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *ThreatenCmd {
	return &ThreatenCmd{}
}

type ThreatenCmd struct {
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
	variable.ZapLog.Sugar().Debug("DeleteForCache")
}

// SyncDefaultForCache  本地缓存：同步默认状态(未同步)数据到威胁情报中心
func (crd *ThreatenCmd) SyncDefaultForCache() {
	variable.ZapLog.Sugar().Debug("SyncDefaultForCache")
}

// SyncUnknownForCache  本地缓存：同步未知状态数据到威胁情报中心
func (crd *ThreatenCmd) SyncUnknownForCache() {
	variable.ZapLog.Sugar().Debug("SyncUnknownForCache")
}
