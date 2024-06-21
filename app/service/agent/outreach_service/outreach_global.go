package outreach_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_outreach"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/utils/helpers/date_helper"
	"time"
)

func CreateGlobalFactory() *GlobalCurd {
	return &GlobalCurd{outreach_models.CreateGlobalFactory()}
}

type GlobalCurd struct {
	masterModel *outreach_models.GlobalModel
}

// AddDefault 新增默认数据
func (crd *GlobalCurd) AddDefault(orgId int64, createId int64, createName string) error {
	curTime := time.Now().Format(variable.DateFormat)
	ipList := []string{
		"fe80::1/10",
		"127.0.0.0/8",
		"192.168.0.0/16",
		"172.16.0.0/12",
		"10.0.0.0/8",
	}
	var inserts []outreach_models.GlobalModel
	for _, ip := range ipList {
		inserts = append(inserts, outreach_models.GlobalModel{
			Id:         variable.SnowFlake.GetId(),
			OrgId:      orgId,
			Ip:         ip,
			Type:       enum_outreach.TypeWhite,
			Title:      "快速添加",
			IsDefault:  enums.IsYes,
			OpenSwitch: enums.SwitchOn,
			CreateId:   createId,
			CreateName: createName,
			CreatedAt:  curTime,
			UpdateId:   createId,
			UpdateName: createName,
			UpdatedAt:  curTime,
		})
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	return nil
}

// Add 新增
func (crd *GlobalCurd) Add(param *outreach_models.GlobalModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.Type = enum_outreach.TypeWhite
	param.IsDefault = enums.IsNot
	param.OpenSwitch = enums.SwitchOn
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	//验证数据
	conditions := outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Ip:    param.Ip,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("已存在相同策略，不能重复添加")
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	return nil
}

// Edit 修改
func (crd *GlobalCurd) Edit(param *outreach_models.GlobalModel, isRoot int) error {
	param.UpdatedAt = date_helper.CurDatetime()

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	} else if hashModel.Ip != param.Ip {
		conditions := outreach_models.SearchConditions{
			OrgId: hashModel.OrgId,
			Ip:    param.Ip,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New("已存在相同策略，不能重复添加")
		}
	}

	//修改
	updates := make(map[string]interface{})
	updates["ip"] = param.Ip
	updates["title"] = param.Title
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(outreach_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	return nil
}

// SetSwitch 开关
func (crd *GlobalCurd) SetSwitch(param *outreach_models.GlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	}

	//修改
	updates := make(map[string]interface{})
	updates["open_switch"] = param.OpenSwitch
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(outreach_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	return nil
}

// DeleteById 删除
func (crd *GlobalCurd) DeleteById(param *outreach_models.GlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New("只能删除自己添加的策略(管理员除外)")
	}

	//删除
	if !crd.masterModel.DeleteByConditions(outreach_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}) {
		return errors.New("删除失败")
	}

	return nil
}
