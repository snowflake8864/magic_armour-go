package role_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/organization/role_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *RoleCurd {
	return &RoleCurd{masterModel: role_models.CreateRoleFactory()}
}

type RoleCurd struct {
	masterModel *role_models.RoleModel
}

func (crd *RoleCurd) defaultMenus() string {
	return "1,2,3,4,5" //最小权限
}

// AddDefault 新增默认
func (crd *RoleCurd) AddDefault(orgId int64) bool {
	now := date_helper.CurDatetime()
	remark := "默认"

	params := []role_models.RoleModel{
		{
			OrgId:      orgId,
			Name:       "管理员",
			Remark:     remark,
			Menus:      crd.defaultMenus(),
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
		{
			OrgId:      orgId,
			Name:       "操作员",
			Remark:     remark,
			Menus:      crd.defaultMenus(),
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
		{
			OrgId:      orgId,
			Name:       "审计员",
			Remark:     remark,
			Menus:      crd.defaultMenus(),
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
	}

	return crd.masterModel.InsertBatch(&params)
}

// Add 新增
func (crd *RoleCurd) Add(param *role_models.RoleModel) error {
	if crd.existsName(param.Name, param.OrgId, 0) {
		return errors.New("名称已存在")
	}

	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	if string_helper.IsEmpty(param.Menus) {
		param.Menus = crd.defaultMenus()
	}
	param.IsDefault = enums.IsNot

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}
	return nil
}

func (crd *RoleCurd) existsName(name string, orgId int64, id int64) bool {
	model := crd.masterModel.GetItemByConditions(role_models.SearchConditions{
		OrgId: orgId,
		Name:  name,
	})
	//数据库没有
	if model == nil || model.Id == 0 {
		return false
	}
	//不是当前id
	if model.Id != id {
		return true
	}
	return false
}

func (crd *RoleCurd) HasId(id int64, orgId int64) *role_models.RoleModel {
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
	if existsModel := crd.HasId(param.Id, param.OrgId); existsModel == nil {
		return errors.New("数据不存在")
	} else if existsModel.IsDefault == enums.IsYes {
		return errors.New("默认数据不让修改")
	}

	if crd.existsName(param.Name, param.OrgId, param.Id) {
		return errors.New("名称已存在")
	}

	//修改
	updates := map[string]interface{}{
		"name":   param.Name,
		"remark": param.Remark,
	}

	if isOk, _ := crd.masterModel.UpdateByConditions(role_models.SearchConditions{
		Id: param.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}
	return nil
}

// DeleteById 删除
func (crd *RoleCurd) DeleteById(id int64, orgId int64) error {
	if existsModel := crd.HasId(id, orgId); existsModel == nil {
		return errors.New("数据不存在")
	} else if existsModel.IsDefault == enums.IsYes {
		return errors.New("默认数据不让删除")
	}

	if !crd.masterModel.DeleteByConditions(role_models.SearchConditions{Id: id}) {
		return errors.New("删除")
	}

	return nil
}

// SaveMenus 保存菜单
func (crd *RoleCurd) SaveMenus(param *role_models.RoleModel) error {
	if existsModel := crd.HasId(param.Id, param.OrgId); existsModel == nil {
		return errors.New("数据不存在")
	}
	//修改
	updates := map[string]interface{}{
		"menus": param.Menus,
	}

	if isOk, _ := crd.masterModel.UpdateByConditions(role_models.SearchConditions{
		Id: param.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}
	return nil
}
