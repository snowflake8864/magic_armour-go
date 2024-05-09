package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/monitor_service"
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
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// AgentLogConfigSet  终端日志配置
func (u *SystemController) AgentLogConfigSet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// HuanshiGet  幻势配置
func (u *SystemController) HuanshiGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// HuanshiOnline  幻势配置
func (u *SystemController) HuanshiOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// HuanshiSet  幻势配置
func (u *SystemController) HuanshiSet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}
