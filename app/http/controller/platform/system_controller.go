package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/monitor_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
	"os/exec"
)

// SystemController 系统管理
type SystemController struct {
	ControllerBase
}

// MonitorRecords 监控记录
func (u *SystemController) MonitorRecords(ctx *gin.Context) {
	date := ctx.GetString(consts.ValidatorPrefixPlatform + "date")

	//页面上传
	conditions := monitor_service.SearchConditions{}
	conditions.MonitorDay = date

	models := monitor_service.CreateCurdFactory().GetList(conditions)

	response.SuccessReturnData(ctx, (&platform_trans.Monitor{}).TransList(models))
}

// TimeGet  服务器器时间
func (u *SystemController) TimeGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})

}

// TimeSet  服务器器时间
func (u *SystemController) TimeSet(ctx *gin.Context) {
	date := ctx.GetString(consts.ValidatorPrefixPlatform + "date")

	command := " sudo /bin/timedatectl set-timezone Asia/Shanghai"
	command += " && sudo /sbin/hwclock --set --date '" + date + "'"
	command += " && sudo /sbin/hwclock --hctosys"
	command += " && sudo /sbin/clock -w"

	cmd := exec.Command(command)
	_, err := cmd.CombinedOutput() // 执行命令并获取输出结果
	if err != nil {
		response.FailData(ctx, "exec command failed", err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SshGet  SSH服务
func (u *SystemController) SshGet(ctx *gin.Context) {
	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.PlatformSshSwitch)

	response.SuccessReturnData(ctx, gin.H{
		"switch": helpers.String2Int(val),
	})
}

// SshSet  SSH服务
func (u *SystemController) SshSet(ctx *gin.Context) {
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "switch"))

	//删除以前的规则
	commandDelDrop := "sudo iptables -D INPUT -p Tcp --dport 22 -j DROP"
	commandDelACCEPT := "sudo iptables -D INPUT -p Tcp --dport 22 -j ACCEPT"
	for i := 0; i <= 2; i++ {
		_ = exec.Command(commandDelDrop).Run()
		_ = exec.Command(commandDelACCEPT).Run()
	}

	//按照开关添加规则
	commandAdd := ""
	if Switch == enums.SwitchOn {
		commandAdd += "sudo iptables -I INPUT -p Tcp --dport 22 -j ACCEPT"
	} else {
		commandAdd += "sudo iptables -I INPUT -p Tcp --dport 22 -j DROP"
	}
	cmd := exec.Command(commandAdd)
	_, err := cmd.CombinedOutput() // 执行命令并获取输出结果
	if err != nil {
		response.FailData(ctx, "exec command failed", err.Error())
		return
	}

	//修改数据字典
	dictionary_service.CreateCurdFactory().UpdateValById(enum_dictionary.PlatformSshSwitch, Switch)

	response.SuccessSimplify(ctx)
}

// NetworkGet  网络服务
func (u *SystemController) NetworkGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// NetworkSet  网络服务
func (u *SystemController) NetworkSet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// AgentLogConfigGet  终端日志配置
func (u *SystemController) AgentLogConfigGet(ctx *gin.Context) {
	dictIds := []string{
		enum_dictionary.AgentAlarmExpire,
		enum_dictionary.AgentSyslogExpire,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	response.SuccessReturnData(ctx, vals)
}

// AgentLogConfigSet  终端日志配置
func (u *SystemController) AgentLogConfigSet(ctx *gin.Context) {
	idAndVal := map[string]interface{}{
		enum_dictionary.AgentAlarmExpire:  ctx.GetFloat64(consts.ValidatorPrefixPlatform + "agent_alarm_expire"),
		enum_dictionary.AgentSyslogExpire: ctx.GetFloat64(consts.ValidatorPrefixPlatform + "agent_syslog_expire"),
	}

	if !dictionary_service.CreateCurdFactory().BatchUpdate(idAndVal) {
		response.FailSimplify(ctx, "保持失败")
		return
	}

	response.SuccessSimplify(ctx)
}

// HuanshiGet  幻势配置
func (u *SystemController) HuanshiGet(ctx *gin.Context) {
	dictIds := []string{
		enum_dictionary.HuanshiCallbackUrl,
		enum_dictionary.HuanshiCallIp,
		enum_dictionary.HuanshiSwitch,
		enum_dictionary.HuanshiUrl,
		enum_dictionary.HuanshiAlarmSwitch,
		enum_dictionary.HuanshiAlarmIp,
		enum_dictionary.HuanshiAlarmPort,
		enum_dictionary.HuanshiLogSwitch,
		enum_dictionary.HuanshiLogIp,
		enum_dictionary.HuanshiLogPort,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	response.SuccessReturnData(ctx, vals)
}

// HuanshiOnline  幻势配置
func (u *SystemController) HuanshiOnline(ctx *gin.Context) {
	config := huanshi_platform.RequestConfig{
		HuanshiCallbackUrl: ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_callback_url"),
		HuanshiCallIp:      "192.168.17.250",
		HuanshiUrl:         ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_url"),
	}
	threatenPlatform := huanshi_platform.CreateRequestFactory(config)
	if err := threatenPlatform.Online(); err != nil {
		response.FailData(ctx, "网络连接失败", err.Error())
		return
	}
	if err := threatenPlatform.GenerateToken(); err != nil {
		response.FailData(ctx, "配置不正确", err.Error())
		return
	}
	param := map[string]interface{}{
		"deviceMark":     "f963864bebc14dd2b4fe9a00047314c4",
		"agentNum":       4,
		"onlineAgentNum": 2,
	}
	if err := threatenPlatform.SyncDeviceByTcp(param); err != nil {
		response.FailData(ctx, "同步设备", err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// HuanshiSet  幻势配置
func (u *SystemController) HuanshiSet(ctx *gin.Context) {
	idAndVal := map[string]interface{}{
		enum_dictionary.HuanshiCallbackUrl: ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_callback_url"),
		enum_dictionary.HuanshiCallIp:      ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_callback_ip"),
		enum_dictionary.HuanshiSwitch:      ctx.GetFloat64(consts.ValidatorPrefixPlatform + "huanshi_switch"),
		enum_dictionary.HuanshiUrl:         ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_url"),
		enum_dictionary.HuanshiAlarmSwitch: ctx.GetFloat64(consts.ValidatorPrefixPlatform + "huanshi_alarm_switch"),
		enum_dictionary.HuanshiAlarmIp:     ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_alarm_ip"),
		enum_dictionary.HuanshiAlarmPort:   ctx.GetFloat64(consts.ValidatorPrefixPlatform + "huanshi_alarm_port"),
		enum_dictionary.HuanshiLogSwitch:   ctx.GetFloat64(consts.ValidatorPrefixPlatform + "huanshi_log_switch"),
		enum_dictionary.HuanshiLogIp:       ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_log_ip"),
		enum_dictionary.HuanshiLogPort:     ctx.GetFloat64(consts.ValidatorPrefixPlatform + "huanshi_log_port"),
	}

	if !dictionary_service.CreateCurdFactory().BatchUpdate(idAndVal) {
		response.FailSimplify(ctx, "保持失败")
		return
	}

	response.SuccessSimplify(ctx)
}
