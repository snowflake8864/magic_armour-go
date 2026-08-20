package outreach_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_outreach"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateDetectFactory() *DetectCurd {
	return &DetectCurd{outreach_models.CreateDetectFactory()}
}

type DetectCurd struct {
	masterModel *outreach_models.DetectModel
}

// 新增任务
func (crd *DetectCurd) addTask(orgId int64, remark string) {
	isOnline := enums.IsYes
	switchOn := enums.SwitchOn
	conditions := agent_models.SearchConditions{
		OrgId:               orgId,
		IsOnline:            &isOnline,
		OutreachAlarmSwitch: &switchOn,
	}
	if ids := agent_models.CreateAgentFactory().GetIdsByByConditions(conditions); ids != nil && len(ids) > 0 {
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, ids, enum_task.TypeOutreachDetectDownload, remark)
	}
}

// Add 新增
func (crd *DetectCurd) Add(param *outreach_models.DetectModel) error {

	//是否超过限制
	limit := string_helper.ToInt64(dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentOutreachDetectMax))
	if limit > 0 {
		if crd.masterModel.Count(outreach_models.SearchConditions{OrgId: param.OrgId}) > limit {
			return errors.New(my_errors.CurdDetectListMaxLimit)
		}
	}

	param.Id = variable.SnowFlake.GetId()
	param.Type = enum_outreach.DetectTypeDomain
	param.Method = enum_outreach.MethodGet
	param.OpenSwitch = enums.SwitchOn
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	//验证数据
	conditions := outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Addr:  param.Addr,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New(my_errors.CurdSamePolicyExistCannotAdd)
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//新增任务
	crd.addTask(param.OrgId, "Add")

	return nil
}

// Edit 修改
func (crd *DetectCurd) Edit(param *outreach_models.DetectModel, isRoot int) error {
	param.UpdatedAt = date_helper.CurDatetime()

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyModifyOneselfAdded)
	} else if hashModel.Addr != param.Addr {
		conditions := outreach_models.SearchConditions{
			OrgId: hashModel.OrgId,
			Addr:  param.Addr,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New(my_errors.CurdSamePolicyExistCannotAdd)
		}
	}

	//修改
	updates := make(map[string]interface{})
	updates["addr"] = param.Addr
	updates["title"] = param.Title
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(outreach_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//新增任务
	crd.addTask(param.OrgId, "Edit")

	return nil
}

// SetSwitch 开关
func (crd *DetectCurd) SetSwitch(param *outreach_models.DetectModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyModifyOneselfAdded)
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
		return errors.New(my_errors.CurdSaveFailed)
	}

	crd.addTask(hashModel.OrgId, "Switch")

	return nil
}

// DeleteById 删除
func (crd *DetectCurd) DeleteById(param *outreach_models.DetectModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(outreach_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyDeleteOneselfAdded)
	}

	//删除
	if !crd.masterModel.DeleteByConditions(outreach_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}) {
		return errors.New(my_errors.CurdDeleteFailed)
	}

	crd.addTask(hashModel.OrgId, "Switch")

	return nil
}
