package organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/service/platform/department_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers"
)

// ControllerBase 基础
type ControllerBase struct {
}

// getLogin 获取登录数据
func (u *ControllerBase) getLogin(ctx *gin.Context) admin_service.AdminCacheLogin {
	return cur_userinfo.GetCurrentAdmin(ctx)
}

// getAgentSearchConditions 获取终端查询条件
func (u *ControllerBase) getAgentSearchConditions(ctx *gin.Context) (agent_models.SearchConditions, error) {

	//页面上传
	searchConditions := agent_models.SearchConditions{}
	searchConditions.AgentId = ctx.GetString(consts.ValidatorPrefixOrganization + "agentId")
	searchConditions.HostName = ctx.GetString(consts.ValidatorPrefixOrganization + "hostName")
	searchConditions.AgentIp = ctx.GetString(consts.ValidatorPrefixOrganization + "agentIp")
	searchConditions.NetIp = ctx.GetString(consts.ValidatorPrefixOrganization + "netIp")
	searchConditions.Alias = ctx.GetString(consts.ValidatorPrefixOrganization + "alias")
	searchConditions.ContactName = ctx.GetString(consts.ValidatorPrefixOrganization + "contactName")
	searchConditions.IsOnline = int(ctx.GetFloat64(consts.ValidatorPrefixOrganization + "isOnline"))
	searchConditions.IsSafe = int(ctx.GetFloat64(consts.ValidatorPrefixOrganization + "isSafe"))

	var reqDeptIds []int64
	if deptIds, isOk := ctx.Get(consts.ValidatorPrefixOrganization + "deptId"); isOk {
		var tmp []interface{}
		tmp, _ = deptIds.([]interface{})
		reqDeptIds = helpers.SliceAny2Int64(tmp)
	}

	//用户管理的组织架构
	if isAgent, departments := u.getCurrAdminDepartments(ctx); !isAgent {
		if departments == nil || len(departments) == 0 {
			return searchConditions, errors.New("没有需要负责的组织架构")
		}
		//没有查询组织架构
		if len(reqDeptIds) == 0 {
			searchConditions.DeptIds = departments
		} else {
			//用户管理组织架构与查询组织架构交集
			intersection := helpers.SliceInt64Intersect(reqDeptIds, departments)
			if intersection == nil || len(intersection) == 0 {
				return searchConditions, errors.New("暂无数据")
			} else {
				searchConditions.DeptIds = intersection
			}
		}
	} else {
		searchConditions.DeptIds = reqDeptIds
	}
	return searchConditions, nil
}

// getCurrAdminDepartments 获取当前用户的管理部门
func (u *ControllerBase) getCurrAdminDepartments(ctx *gin.Context) (bool, []int64) {

	return true, nil
	/*
		loginInfo := u.getLogin(ctx)
		if loginInfo.IsAgent == enums.IsYes {
			return true, nil
		}
		departments := helpers.Json2SliceInt64(loginInfo.Departments)
		return false, departments
	*/
}

// getPageAndSort 获取分页参数和排序参数
func (u *ControllerBase) getPageAndSort(ctx *gin.Context, defaultField string) (offset, pageSize int, order string) {
	page := int(ctx.GetFloat64(consts.ValidatorPrefixOrganization + "page"))
	pageSize = int(ctx.GetFloat64(consts.ValidatorPrefixOrganization + "pageSize"))
	orderField := ctx.GetString(consts.ValidatorPrefixOrganization + "field")
	orderSort := ctx.GetString(consts.ValidatorPrefixOrganization + "sort")
	if page <= 0 {
		page = 1
	}
	if pageSize <= 1 {
		page = 20
	}
	offset = (page - 1) * pageSize

	if helpers.StringIsEmpty(orderField) {
		if helpers.StringIsEmpty(defaultField) {
			orderField = "created_at"
		} else {
			orderField = defaultField
		}
	}
	if helpers.StringIsEmpty(orderSort) {
		orderSort = "DESC"
	}
	return offset, pageSize, orderField + " " + orderSort
}

// parseAgentModels 解析终端集合并拿到相关数据
// 返回 agentIds 终端集合
// 返回 departments 组织架构集合
func (u *ControllerBase) parseAgentModels(agentModels []*agent_models.AgentModel, needDepartment bool) (agentIds []string, departments map[int64]*department_models.DepartmentModel) {
	if agentModels == nil || len(agentModels) == 0 {
		return nil, nil
	}
	deptStruct := make(map[int64]struct{})
	var deptIds []int64
	for _, model := range agentModels {
		agentIds = append(agentIds, model.AgentId)

		if _, isOk := deptStruct[model.DeptId]; model.DeptId > 0 && !isOk {
			deptStruct[model.DeptId] = struct{}{}
			deptIds = append(deptIds, model.DeptId)
		}
	}

	if needDepartment && len(deptIds) > 0 {
		departments = department_service.CreateCurdFactory().GetMapsByIds(deptIds)
	}

	return agentIds, departments
}

// getAgentIdsFromGin 从gin上下文获取终端id集合
// 参数agentIds([]string) 在中间件经过 data_transfer.DataAddContext转义成([]interface{})
func (u *ControllerBase) getAgentIdsFromGin(ctx *gin.Context) (agentIds []string) {
	if val, ok := ctx.Get(consts.ValidatorPrefixOrganization + "agentIds"); ok && val != nil {
		var tmp []interface{}
		tmp, _ = val.([]interface{})
		agentIds = helpers.SliceAny2String(tmp)
	}
	return agentIds
}
