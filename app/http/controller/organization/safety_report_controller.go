package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_safety_report"
	"goskeleton/app/http/transformer"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/response"
)

// SafetyReportController  安全报告相关
type SafetyReportController struct {
	ControllerBase
}

// GetForAlarmNotice 告警通知: 获取
func (u *SafetyReportController) GetForAlarmNotice(ctx *gin.Context) {
	//通知周期
	keys1 := enums.Get1stIsStrReturnSliceOfInt("notice_cycle", enum_safety_report.WebMapInt)
	labels1 := enums.Get1stIsStrReturnMapIntOfStr("notice_cycle", enum_safety_report.EnumIntMap)
	noticeCycleList := (&transformer.Common{}).FormatIntList(keys1, labels1, false)

	//周期发送模块
	keys2 := enums.Get1stIsStrReturnSliceOfInt("alarm_cycle_mod", enum_safety_report.WebMapInt)
	labels2 := enums.Get1stIsStrReturnMapIntOfStr("module", enum_alarm.EnumMapInt)
	cycleModList := (&transformer.Common{}).FormatIntList(keys2, labels2, false)

	//立即发送模块
	keys3 := enums.Get1stIsStrReturnSliceOfInt("alarm_timely_mod", enum_safety_report.WebMapInt)
	labels4 := enums.Get1stIsStrReturnMapIntOfStr("module", enum_alarm.EnumMapInt)
	timelyModList := (&transformer.Common{}).FormatIntList(keys3, labels4, false)

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

// GetForHardwareNotice 硬件状态通知: 获取
func (u *SafetyReportController) GetForHardwareNotice(ctx *gin.Context) {
	//配置
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	model := OrgExtendFact.GetItemByOrg(u.getLogin(ctx).OrgId)
	config := OrgExtendFact.FormatHardwareNoticeConfig(model.HardwareNoticeConfig)
	response.SuccessReturnData(ctx, config)
}
