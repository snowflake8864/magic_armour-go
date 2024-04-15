package outreach_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_outreach"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
	"time"
)

func CreateGlobalFactory() *OutreachGlobalCurd {
	return &OutreachGlobalCurd{model_agent.CreateOutreachGlobalFactory()}
}

type OutreachGlobalCurd struct {
	masterModel *model_agent.OutreachGlobalModel
}

// AddDefault 新增默认数据
func (crd *OutreachGlobalCurd) AddDefault(orgId int64) bool {
	curTime := time.Now().Format(variable.DateFormat)

	list := []model_agent.OutreachGlobalModel{
		{
			Id:         "382d29ad257f55c55c75c9f3a906c99e",
			OrgId:      orgId,
			Ip:         "fe80::1/10",
			Type:       1,
			Title:      "系统默认",
			IsDefault:  enums.IsYes,
			OpenSwitch: enums.SwitchOn,
			CreateId:   0,
			CreateName: "",
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
			CreateId:   0,
			CreateName: "",
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
			CreateId:   0,
			CreateName: "",
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
			CreateId:   0,
			CreateName: "",
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
			CreateId:   0,
			CreateName: "",
			CreatedAt:  curTime,
			UpdatedAt:  curTime,
		},
	}

	isOk := crd.masterModel.InsertBatch(&list)

	return isOk
}

// Add 新增
func (crd *OutreachGlobalCurd) Add(param *model_agent.OutreachGlobalModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
