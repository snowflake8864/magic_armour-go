package usb_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	UsbHash     string
	UsbHashList []string
	Status      int
	IsOnline    int
}

type baseUsb struct {
}

// conditions2Where 查询条件转化
func (m *baseUsb) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !string_helper.IsEmpty(search.UsbHash) {
		wheres["usb_hash"] = search.UsbHash
	}
	if len(search.UsbHashList) > 0 {
		wheres["usb_hash in"] = search.UsbHashList
	}
	if search.Status > 0 {
		wheres["status"] = search.Status
	}
	if search.IsOnline != enums.AllDefaultInt {
		wheres["is_online"] = search.IsOnline
	}

	return wheres
}
