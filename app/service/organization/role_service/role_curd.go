package role_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/my_errors"
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

func (crd *RoleCurd) defaultMinMenus() string {
	return "1,2,3,4,5" //最小权限
}

func (crd *RoleCurd) defaultAuditorMenus() string {
	return "1,2,3,4,5,53,54"
}

func (crd *RoleCurd) defaultOperatorMenus() string {
	return "1,2,3,4,5,10,11,12,87,13,14,15,69,17,18,19,20,21,71,22,23,68,25,26,27,28,29,30,31,32,33,34,35,36,37,38,73,62,63,64,51,52"
}

func (crd *RoleCurd) defaultAdministratorMenus() string {
	return "1,2,3,4,5,10,11,12,13,14,15,17,18,19,20,21,22,23,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,67,68,71,73,75,81,82"
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
			Menus:      crd.defaultAdministratorMenus(),
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
		{
			OrgId:      orgId,
			Name:       "操作员",
			Remark:     remark,
			Menus:      crd.defaultOperatorMenus(),
			IsDefault:  enums.IsYes,
			CreateId:   0,
			CreateName: "",
			CreatedAt:  now,
		},
		{
			OrgId:      orgId,
			Name:       "审计员",
			Remark:     remark,
			Menus:      crd.defaultAuditorMenus() + ",53,54",
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
		param.Menus = crd.defaultMinMenus()
	}
	param.IsDefault = enums.IsNot

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
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
		return errors.New(my_errors.CurdSaveFailed)
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
		return errors.New(my_errors.CurdSaveFailed)
	}
	return nil
}
