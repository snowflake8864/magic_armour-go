package outreach_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_outreach"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/outreach_models"
	"time"
)

func CreateGlobalFactory() *OutreachGlobalCurd {
	return &OutreachGlobalCurd{outreach_models.CreateOutreachGlobalFactory()}
}

type OutreachGlobalCurd struct {
	masterModel *outreach_models.OutreachGlobalModel
}

// AddDefault 新增默认数据
func (crd *OutreachGlobalCurd) AddDefault(orgId int64, createId int64, createName string) bool {
	curTime := time.Now().Format(variable.DateFormat)

	list := []outreach_models.OutreachGlobalModel{
		{
			Id:         "382d29ad257f55c55c75c9f3a906c99e",
			OrgId:      orgId,
			Ip:         "fe80::1/10",
			Type:       enum_outreach.TypeWhite,
			Title:      "系统默认",
			IsDefault:  enums.IsYes,
			OpenSwitch: enums.SwitchOn,
			CreateId:   createId,
			CreateName: createName,
			CreatedAt:  curTime,
			UpdatedAt:  curTime,
		},
		{
			Id:         "0f39d743e689c10ca2fe6e8f57d3d2cd",
			OrgId:      orgId,
			Ip:         "127.0.0.0/8",
			Type:       enum_outreach.TypeWhite,
			Title:      "系统默认",
			IsDefault:  enums.IsYes,
			OpenSwitch: enums.SwitchOn,
			CreateId:   createId,
			CreateName: createName,
			CreatedAt:  curTime,
			UpdatedAt:  curTime,
		},
		{
			Id:         "738a5b3b5ed9abce50d289ee49772ddd",
			OrgId:      orgId,
			Ip:         "192.168.0.0/16",
			Type:       enum_outreach.TypeWhite,
			Title:      "系统默认",
			IsDefault:  enums.IsYes,
			OpenSwitch: enums.SwitchOn,
			CreateId:   createId,
			CreateName: createName,
			CreatedAt:  curTime,
			UpdatedAt:  curTime,
		},
		{
			Id:         "bb6febec1d21cbfd72a76d25a1ef6622",
			OrgId:      orgId,
			Ip:         "172.16.0.0/12",
			Type:       enum_outreach.TypeWhite,
			Title:      "系统默认",
			IsDefault:  enums.IsYes,
			OpenSwitch: enums.SwitchOn,
			CreateId:   createId,
			CreateName: createName,
			CreatedAt:  curTime,
			UpdatedAt:  curTime,
		},
		{
			Id:         "920069446bf4f6b26db744f88f1106bb",
			OrgId:      orgId,
			Ip:         "10.0.0.0/8",
			Type:       enum_outreach.TypeWhite,
			Title:      "系统默认",
			IsDefault:  enums.IsYes,
			OpenSwitch: enums.SwitchOn,
			CreateId:   createId,
			CreateName: createName,
			CreatedAt:  curTime,
			UpdatedAt:  curTime,
		},
	}

	isOk := crd.masterModel.InsertBatch(&list)

	return isOk
}

// Add 新增
func (crd *OutreachGlobalCurd) Add(param *outreach_models.OutreachGlobalModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
