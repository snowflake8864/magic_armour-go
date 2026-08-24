package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_safety_report"
	"goskeleton/app/http/transformer"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/response"
)

// SafetyReportController  安全报告相关
type SafetyReportController struct {
	ControllerBase
}

// GetForAlarmNotice 告警通知: 获取
func (u *SafetyReportController) GetForAlarmNotice(ctx *gin.Context) {
	lang := cur_userinfo.GetLocales(ctx)
	isEn := lang == consts.LanguageEn

	//通知周期
	keys1 := enums.Get1stIsStrReturnSliceOfInt("notice_cycle", enum_safety_report.WebMapInt)
	labels1 := enums.Get1stIsStrReturnMapIntOfStr("notice_cycle", enum_safety_report.EnumIntMap)
	if isEn {
		labels1 = enums.Get1stIsStrReturnMapIntOfStr("notice_cycle_en", enum_safety_report.EnumIntMap)
	}
	noticeCycleList := (&transformer.Common{Lang: lang}).FormatIntList(keys1, labels1, false)

	//周期发送模块
	keys2 := enums.Get1stIsStrReturnSliceOfInt("alarm_cycle_mod", enum_safety_report.WebMapInt)
	labels2 := enums.Get1stIsStrReturnMapIntOfStr("module", enum_alarm.EnumMapInt)
	if isEn {
		labels2 = enums.Get1stIsStrReturnMapIntOfStr("module_en", enum_alarm.EnumMapInt)
	}
	cycleModList := (&transformer.Common{Lang: lang}).FormatIntList(keys2, labels2, false)

	//立即发送模块
	keys3 := enums.Get1stIsStrReturnSliceOfInt("alarm_timely_mod", enum_safety_report.WebMapInt)
	labels3 := enums.Get1stIsStrReturnMapIntOfStr("module", enum_alarm.EnumMapInt)
	if isEn {
		labels3 = enums.Get1stIsStrReturnMapIntOfStr("module_en", enum_alarm.EnumMapInt)
	}
	timelyModList := (&transformer.Common{Lang: lang}).FormatIntList(keys3, labels3, false)

	//配置
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	model := OrgExtendFact.GetItemByOrg(u.getLogin(ctx).OrgId)
	config := OrgExtendFact.FormatAlarmNoticeConfig(model.AlarmNoticeConfig)

	response.SuccessReturnData(ctx, gin.H{
		"config":            config,
		"notice_cycle_list": noticeCycleList,
		"cycle_mod_list":    cycleModList,
		"timely_mod_list":   timelyModList,
	})
}

// SaveForAlarmNotice 告警通知: 保存
func (u *SafetyReportController) SaveForAlarmNotice(ctx *gin.Context) {
	config := enum_safety_report.AlarmNoticeConfig{}
	config.EmailSwitch = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "emailSwitch"))
	config.Email = ctx.GetString(consts.ValidatorPrefixOrg + "email")
	config.CycleInterval = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "cycleInterval"))
	config.CycleMod = json_helper.ToSliceInt(ctx.GetString(consts.ValidatorPrefixOrg + "cycleMod"))
	config.TimelySwitch = enums.SwitchOn
	config.TimelyMod = json_helper.ToSliceInt(ctx.GetString(consts.ValidatorPrefixOrg + "timelyMod"))

	if err := organize_service.CreateOrgExtendFactory().SetAlarmNoticeConfig(u.getLogin(ctx).OrgId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetForHardwareNotice 硬件状态通知: 获取
func (u *SafetyReportController) GetForHardwareNotice(ctx *gin.Context) {
	//配置
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	model := OrgExtendFact.GetItemByOrg(u.getLogin(ctx).OrgId)
	config := OrgExtendFact.FormatHardwareNoticeConfig(model.HardwareNoticeConfig)
	response.SuccessReturnData(ctx, config)
}

// SaveForHardwareNotice  硬件状态通知: 保存
func (u *SafetyReportController) SaveForHardwareNotice(ctx *gin.Context) {
	config := enum_safety_report.HardwareNoticeConfig{}
	config.EmailSwitch = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "emailSwitch"))
	config.Email = ctx.GetString(consts.ValidatorPrefixOrg + "email")
	config.TimelySwitch = enums.SwitchOn
	config.CpuSwitch = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "cpuSwitch"))
	config.CpuThreshold = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "cpuThreshold"))
	config.MemSwitch = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "memSwitch"))
	config.MemThreshold = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "memThreshold"))
	config.HdSwitch = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "hdSwitch"))
	config.HdThreshold = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "hdThreshold"))

	if err := organize_service.CreateOrgExtendFactory().SetHardwareNoticeConfig(u.getLogin(ctx).OrgId, config); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
