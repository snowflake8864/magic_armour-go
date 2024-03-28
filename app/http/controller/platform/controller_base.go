package platform

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/platform/admin_service"
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
func (u *ControllerBase) getAgentSearchConditions(ctx *gin.Context) (agent_service.SearchConditions, error) {

	//页面上传
	searchConditions := agent_service.SearchConditions{}
	searchConditions.AgentId = ctx.GetString(consts.ValidatorPrefixPlatform + "agentId")
	searchConditions.HostName = ctx.GetString(consts.ValidatorPrefixPlatform + "hostName")
	searchConditions.AgentIp = ctx.GetString(consts.ValidatorPrefixPlatform + "agentIp")
	searchConditions.NetIp = ctx.GetString(consts.ValidatorPrefixPlatform + "netIp")
	searchConditions.Alias = ctx.GetString(consts.ValidatorPrefixPlatform + "alias")
	searchConditions.ContactName = ctx.GetString(consts.ValidatorPrefixPlatform + "contactName")
	searchConditions.IsOnline = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isOnline"))
	searchConditions.IsSafe = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isSafe"))

	var reqDeptIds []int
	if deptIds, isOk := ctx.Get(consts.ValidatorPrefixPlatform + "deptId"); isOk {
		var tmp []interface{}
		tmp, _ = deptIds.([]interface{})
		reqDeptIds = helpers.SliceAny2Int(tmp)
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
			intersection := helpers.SliceIntIntersect(reqDeptIds, departments)
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
func (u *ControllerBase) getCurrAdminDepartments(ctx *gin.Context) (bool, []int) {
	loginInfo := u.getLogin(ctx)
	if loginInfo.IsAgent == enums.IsYes {
		return true, nil
	}

	departments := helpers.Json2SliceInt(loginInfo.Departments)
	return false, departments
}

// getPageAndSort 获取分页参数和排序参数
func (u *ControllerBase) getPageAndSort(ctx *gin.Context, defaultField string) (offset, pageSize int, orderField, orderSort string) {
	page := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "page"))
	pageSize = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "pageSize"))
	orderField = ctx.GetString(consts.ValidatorPrefixPlatform + "field")
	orderSort = ctx.GetString(consts.ValidatorPrefixPlatform + "sort")
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
	return offset, pageSize, orderField, orderSort
}
