package usb_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
	"strings"
)

func CreateGlobalCurdFactory() *GlobalCurd {
	return &GlobalCurd{masterModel: usb_models.CreateGlobalFactory()}
}

type GlobalCurd struct {
	masterModel *usb_models.GlobalModel
}

// 新增任务
func (crd *GlobalCurd) addTask(orgId int64, remark string) {
	isOnline := enums.IsYes
	conditions := agent_models.SearchConditions{
		OrgId:    orgId,
		IsOnline: &isOnline,
	}
	if ids := agent_models.CreateAgentFactory().GetIdsByByConditions(conditions); ids != nil || len(ids) > 0 {
		tasks := []int{
			enum_task.TypeUsbWhiteDownload,
			enum_task.TypeUsbBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgents(orgId, ids, tasks, remark)
	}
}

// Add 新增
func (crd *GlobalCurd) Add(param *usb_models.GlobalModel) error {
	param.Id = strings.ToLower(param.Id)
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	//验证数据
	conditions := usb_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("已存在相同策略，不能重复添加")
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	//新增任务
	crd.addTask(param.OrgId, "全局:新增")

	return nil
}

// Edit 修改
func (crd *GlobalCurd) Edit(param *usb_models.GlobalModel) error {
	param.UpdatedAt = date_helper.CurDatetime()

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(usb_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	}

	//修改
	updates := make(map[string]interface{})
	updates["alias"] = param.Alias
	updates["remark"] = param.Remark
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(usb_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	return nil
}

// SetStatus 设置状态
func (crd *GlobalCurd) SetStatus(param *usb_models.GlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(usb_models.SearchConditions{
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
	updates["status"] = param.Status
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(usb_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//新增任务
	crd.addTask(hashModel.OrgId, "全局:修改状态")

	return nil
}

// DeleteById 删除
func (crd *GlobalCurd) DeleteById(param *usb_models.GlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(usb_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New("只能删除自己添加的策略(管理员除外)")
	}

	//修改
	if !crd.masterModel.DeleteByConditions(usb_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}) {
		return errors.New("删除失败")
	}

	//新增任务
	crd.addTask(hashModel.OrgId, "全局:删除")

	return nil
}
