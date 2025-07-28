package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/service/agent/safety_manage"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// SafetyManageController  安全管理相关
type SafetyManageController struct {
	ControllerBase
}

// GetForUser 用户配置: 获取
func (u *SafetyManageController) GetForUser(ctx *gin.Context) {
	vals := dictionary_service.CreateCurdFactory().GetMapByModule(enum_dictionary.ModuleOrganization)

	response.SuccessReturnData(ctx, gin.H{
		"login_expire":           string_helper.ToInt(vals[enum_dictionary.OrganizationLoginEXPIRE]),
		"login_error_max":        string_helper.ToInt(vals[enum_dictionary.OrganizationLoginErrorMax]),
		"login_lock_duration":    string_helper.ToInt(vals[enum_dictionary.OrganizationLoginLockDuration]),
		"pw_length_min":          string_helper.ToInt(vals[enum_dictionary.OrganizationPwLengthMin]),
		"pw_length_max":          string_helper.ToInt(vals[enum_dictionary.OrganizationPwLengthMax]),
		"pw_expire":              string_helper.ToInt(vals[enum_dictionary.OrganizationPwExpire]),
		"pw_cycle_modify_switch": string_helper.ToInt(vals[enum_dictionary.OrganizationPwCycleModifySwitch]),
		"pw_first_modify_switch": string_helper.ToInt(vals[enum_dictionary.OrganizationPwlFirstModifySwitch]),
	})
}

// SaveForUser 用户配置: 保存
func (u *SafetyManageController) SaveForUser(ctx *gin.Context) {
	idAndVal := make(map[string]interface{})
	idAndVal[enum_dictionary.OrganizationLoginEXPIRE] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "login_expire")
	idAndVal[enum_dictionary.OrganizationLoginErrorMax] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "login_error_max")
	idAndVal[enum_dictionary.OrganizationLoginLockDuration] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "login_lock_duration")
	idAndVal[enum_dictionary.OrganizationPwLengthMin] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "pw_length_min")
	idAndVal[enum_dictionary.OrganizationPwLengthMax] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "pw_length_max")
	idAndVal[enum_dictionary.OrganizationPwExpire] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "pw_expire")
	idAndVal[enum_dictionary.OrganizationPwCycleModifySwitch] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "pw_cycle_modify_switch")
	idAndVal[enum_dictionary.OrganizationPwlFirstModifySwitch] = ctx.GetFloat64(consts.ValidatorPrefixOrg + "pw_first_modify_switch")

	if !dictionary_service.CreateCurdFactory().BatchUpdate(idAndVal) {
		response.FailSimplify(ctx, "保存失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// ExportForBackup 数据备份: 导出
func (u *SafetyManageController) ExportForBackup(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	if fullPath, err := safety_manage.CreateBackupExcelFactory(loginInfo.OrgId).Export(); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else {
		response.FileDownload(ctx, fullPath)
	}
}

// ImportForBackup 数据备份: 导入
func (u *SafetyManageController) ImportForBackup(ctx *gin.Context) {

}
