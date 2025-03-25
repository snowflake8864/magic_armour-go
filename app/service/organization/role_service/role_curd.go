package role_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/organization/role_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateRoleCurdFactory() *RoleCurd {
	return &RoleCurd{masterModel: role_models.CreateRoleFactory()}
}

type RoleCurd struct {
	masterModel *role_models.RoleModel
}

// AddDefault 新增默认
func (crd *RoleCurd) AddDefault(orgId int64) bool {
	now := date_helper.CurDatetime()
	remark := "默认"
	menus := "1,2,3,4,5" //最小权限

	params := []role_models.RoleModel{
		role_models.RoleModel{
			OrgId:      orgId,
			Name:       "管理员",
			Remark:     remark,
			Menus:      menus,
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
		role_models.RoleModel{
			OrgId:      orgId,
			Name:       "操作员",
			Remark:     remark,
			Menus:      menus,
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
		role_models.RoleModel{
			OrgId:      orgId,
			Name:       "审计员",
			Remark:     remark,
			Menus:      menus,
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
	}

	return crd.masterModel.InsertBatch(&params)
}

// Add 新增
func (crd *RoleCurd) Add(param *role_models.RoleModel) bool {
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	return crd.masterModel.InsertItem(param)
}

// DeleteById 删除
func (crd *RoleCurd) exists(id int64, orgId int64) *role_models.RoleModel {
	model := crd.masterModel.GetItemById(id)
	if model == nil || model.Id == 0 {
		return nil
	}
	if model.OrgId != orgId {
		return nil
	}

	return model
}

// UpdateById 删除
func (crd *RoleCurd) UpdateById(param *role_models.RoleModel) error {
	if existsModel := crd.exists(param.Id, param.OrgId); existsModel == nil {
		return errors.New("数据不存在")
	} else if existsModel.IsDefault == enums.IsYes {
		return errors.New("默认数据不让修改")
	}

	return nil
}

// DeleteById 删除
func (crd *RoleCurd) DeleteById(id int64, orgId int64) error {
	if existsModel := crd.exists(id, orgId); existsModel == nil {
		return errors.New("数据不存在")
	} else if existsModel.IsDefault == enums.IsYes {
		return errors.New("默认数据不让删除")
	}

	return nil
}

// SaveMenus 保存菜单
func (crd *RoleCurd) SaveMenus(param *role_models.RoleModel) error {

	return nil
}
