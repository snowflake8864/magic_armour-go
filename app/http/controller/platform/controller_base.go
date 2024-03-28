package platform

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
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
func (u *ControllerBase) getAgentSearchConditions(ctx *gin.Context) (map[string]interface{}, error) {
	conditions := make(map[string]interface{})

	//页面上传
	conditions["agentId"] = ctx.GetString(consts.ValidatorPrefixPlatform + "agentId")
	conditions["hostName"] = ctx.GetString(consts.ValidatorPrefixPlatform + "hostName")
	conditions["agentIp"] = ctx.GetString(consts.ValidatorPrefixPlatform + "agentIp")
	conditions["netIp"] = ctx.GetString(consts.ValidatorPrefixPlatform + "netIp")
	conditions["alias"] = ctx.GetString(consts.ValidatorPrefixPlatform + "alias")
	conditions["contactName"] = ctx.GetString(consts.ValidatorPrefixPlatform + "contactName")
	var reqDeptIds []int
	if deptIds, isOk := ctx.Get(consts.ValidatorPrefixPlatform + "deptId"); isOk {
		reqDeptIds, _ = deptIds.([]int)
	}
	conditions["isOnline"] = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isOnline"))
	conditions["isSafe"] = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isSafe"))

	//用户管理的组织架构
	if isAgent, departments := u.getCurrAdminDepartments(ctx); !isAgent {
		if departments == nil || len(departments) == 0 {
			return conditions, errors.New("没有需要负责的组织架构")
		}
		//没有查询组织架构
		if len(reqDeptIds) == 0 {
			conditions["deptIds"] = departments
		} else {
			//用户管理组织架构与查询组织架构交集
			intersection := helpers.SliceIntIntersect(reqDeptIds, departments)
			if intersection == nil || len(intersection) == 0 {
				return conditions, errors.New("暂无数据")
			} else {
				conditions["deptIds"] = intersection
			}
		}
	} else {
		conditions["deptIds"] = reqDeptIds
	}

	return conditions, nil
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
