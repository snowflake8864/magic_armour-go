package usb_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_usb"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateUsbCurdFactory() *UsbCurd {
	return &UsbCurd{
		masterModel: usb_models.CreateUsbFactory(),
		logModel:    usb_models.CreateLogFactory(),
	}
}

type UsbCurd struct {
	masterModel *usb_models.UsbModel
	logModel    *usb_models.LogModel
}

// 终端下线设备
func (crd *UsbCurd) getOnlineMapByAgent(orgId int64, agentId string) map[string]*usb_models.UsbModel {
	maps := make(map[string]*usb_models.UsbModel)
	condition := usb_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		IsOnline: enums.IsYes,
	}
	models := crd.masterModel.GetListByConditions(condition)
	if models == nil || len(models) == 0 {
		return maps
	}
	for _, model := range models {
		maps[model.UsbHash] = model
	}
	return maps
}

// 终端hash存在设备
func (crd *UsbCurd) getMapByHashes(orgId int64, agentId string, hashList []string) map[string]*usb_models.UsbModel {
	maps := make(map[string]*usb_models.UsbModel)
	condition := usb_models.SearchConditions{
		OrgId:       orgId,
		AgentId:     agentId,
		UsbHashList: hashList,
		IsOnline:    enums.AllDefaultInt,
	}
	models := crd.masterModel.GetListByConditions(condition)
	if models == nil || len(models) == 0 {
		return maps
	}
	for _, model := range models {
		maps[model.UsbHash] = model
	}
	return maps
}

// 设置终端全部设备离线
func (crd *UsbCurd) setOfflineByAgent(orgId int64, agentId string) {
	condition := usb_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		IsOnline: enums.AllDefaultInt,
	}
	updates := map[string]interface{}{
		"is_online": enums.IsNot,
	}
	crd.masterModel.UpdateByConditions(condition, updates)
}
func (crd *UsbCurd) initLog(orgId int64, agentId, at string, plugStatus int, usbModel usb_models.UsbModel) usb_models.LogModel {
	return usb_models.LogModel{
		OrgId:      orgId,
		AgentId:    agentId,
		UsbHash:    strings.ToLower(usbModel.UsbHash),
		Name:       usbModel.Name,
		Status:     usbModel.Status,
		PlugStatus: plugStatus,
		CreatedAt:  at,
	}
}

// 设置设备名称并上线
func (crd *UsbCurd) setNameByHash(orgId int64, agentId, hash string, param usb_models.UsbModel) {
	condition := usb_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		UsbHash:  hash,
		IsOnline: enums.AllDefaultInt,
	}
	updates := map[string]interface{}{
		"name":      param.Name,
		"intro":     param.Intro,
		"type":      param.Type,
		"is_online": enums.IsYes,
	}
	crd.masterModel.UpdateByConditions(condition, updates)
}

// 设置设备上线
func (crd *UsbCurd) setOnlineByHashes(orgId int64, agentId string, hashes []string) {
	condition := usb_models.SearchConditions{
		OrgId:       orgId,
		AgentId:     agentId,
		UsbHashList: hashes,
		IsOnline:    enums.AllDefaultInt,
	}
	updates := map[string]interface{}{
		"is_online": enums.IsYes,
	}
	crd.masterModel.UpdateByConditions(condition, updates)
}

// BatchAdd 新增
/**
1. 获取终端在线设备
2. 设置终端全部设备离线
3. 本次上传设备设置成上线
4. 确认那些设备离线, 并记录拔出日志
5. 遍历本次设备数据：
	(1). 上传设备存在, 修改设备信息
    (2). 上传设备不存在, 新增设备信息，并记录插入日志

*/
func (crd *UsbCurd) BatchAdd(orgId int64, agentId string, hashList []string, params []usb_models.UsbModel) error {
	//终端在线设备
	onlineMap := crd.getOnlineMapByAgent(orgId, agentId)

	//设置终端全部设备离线
	crd.setOfflineByAgent(orgId, agentId)

	//本次hash
	curHashMap := make(map[string]struct{})
	if len(hashList) > 0 {
		//设备上线
		crd.setOnlineByHashes(orgId, agentId, hashList)

		for _, hash := range hashList {
			curHashMap[hash] = struct{}{}
		}
	}

	now := date_helper.CurDatetime()
	var masters []usb_models.UsbModel
	var logs []usb_models.LogModel

	//确认那些设备离线, 并记录日志
	if len(onlineMap) > 0 {
		for hash, model := range onlineMap {
			if _, isOk := curHashMap[hash]; !isOk {
				logs = append(logs, crd.initLog(orgId, agentId, now, enum_usb.PlugStatusOut, *model))
			}
		}
	}

	//遍历本次上传数据
	if len(params) > 0 {
		//获取存在设备
		existsMap := crd.getMapByHashes(orgId, agentId, hashList)
		for _, curModel := range params {
			hash := strings.ToLower(curModel.UsbHash)
			//本次上传设备存在
			if hasModel, isOk := existsMap[hash]; isOk {
				//如果设备名称不为空，则更新设备信息 @todo bug7675
				if string_helper.IsEmpty(hasModel.Name) && !string_helper.IsEmpty(curModel.Name) {
					crd.setNameByHash(orgId, agentId, hash, curModel)
				}
			} else {
				//新增设备信息
				masters = append(masters, usb_models.UsbModel{
					Id:        variable.SnowFlake.GetId(),
					OrgId:     orgId,
					AgentId:   agentId,
					UsbHash:   hash,
					Name:      curModel.Name,
					Intro:     curModel.Intro,
					Type:      curModel.Type,
					Status:    enum_usb.StatusDefault,
					IsOnline:  enums.IsYes,
					CreatedAt: now,
					UpdatedAt: now,
				})
				logs = append(logs, crd.initLog(orgId, agentId, now, enum_usb.PlugStatusIn, curModel))
			}
		}
	}

	if len(masters) > 0 {
		crd.masterModel.InsertBatch(&masters)
	}
	if len(logs) > 0 {
		crd.logModel.InsertBatch(&logs)
	}

	return nil
}

// DelByOrg 删除：企业
func (crd *UsbCurd) DelByOrg(orgId int64) bool {
	search := usb_models.SearchConditions{
		OrgId:    orgId,
		IsOnline: enums.AllDefaultInt,
	}
	return crd.masterModel.DeleteByConditions(search)
}

// DelByOrgAndAgent 删除：企业和终端
func (crd *UsbCurd) DelByOrgAndAgent(orgId int64, agentId string) bool {
	search := usb_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		IsOnline: enums.AllDefaultInt,
	}
	return crd.masterModel.DeleteByConditions(search)
}
