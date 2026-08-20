package blackmail_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateGroupCurdFactory() *GroupCurd {
	return &GroupCurd{masterModel: blackmail_models.CreateGroupFactory()}
}

type GroupCurd struct {
	masterModel *blackmail_models.GroupModel
}

// Add 新增
func (crd *GroupCurd) Add(param *blackmail_models.GroupModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// Edit 修改
func (crd *GroupCurd) Edit(param *blackmail_models.GroupModel, isRoot int) error {
	param.UpdatedAt = date_helper.CurDatetime()

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(blackmail_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New(my_errors.CurdOnlyModifyOneselfAdded)
	}

	//修改
	updates := make(map[string]interface{})
	updates["title"] = param.Title
	updates["suffix"] = param.Suffix
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(blackmail_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// DeleteById 删除
func (crd *GroupCurd) DeleteById(param *blackmail_models.GroupModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(blackmail_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New(my_errors.CurdOnlyDeleteOneselfAdded)
	}

	//删除
	if !crd.masterModel.DeleteByConditions(blackmail_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}) {
		return errors.New(my_errors.CurdDeleteFailed)
	}

	return nil
}

// DeleteByOrg 删除
func (crd *GroupCurd) DeleteByOrg(orgId int64) {
	crd.masterModel.DeleteByConditions(blackmail_models.SearchConditions{
		OrgId: orgId,
	})
}
