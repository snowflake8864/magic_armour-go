package platform

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/model/platform/server_monitor_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/files"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
	"goskeleton/app/utils/zip_factory"
	"os/exec"
	"path"
	"strings"
)

// SystemController 系统管理
type SystemController struct {
	ControllerBase
}

// MonitorRecords 监控记录
func (u *SystemController) MonitorRecords(ctx *gin.Context) {
	date := ctx.GetString(consts.ValidatorPrefixPlatform + "date")

	//页面上传
	conditions := server_monitor_models.SearchConditions{}
	conditions.MonitorDay = date

	models, _ := server_monitor_models.CreateServerMonitorFactory().GetAllByConditions(conditions, "id ASC")

	response.SuccessReturnData(ctx, (&platform_trans.Monitor{}).TransList(models))
}

// TimeGet  服务器器时间
func (u *SystemController) TimeGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": date_helper.CurDatetime(),
	})

}

// TimeSet  服务器器时间
func (u *SystemController) TimeSet(ctx *gin.Context) {
	date := ctx.GetString(consts.ValidatorPrefixPlatform + "date")

	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.PlatformTimeCmd)

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(val, date))
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
		"switch": string_helper.ToInt(val),
	})
}

// SshSet  SSH服务
func (u *SystemController) SshSet(ctx *gin.Context) {
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "switch"))

	//按照开关添加规则
	command := ""
	if Switch == enums.SwitchOn {
		command += dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.PlatformSshCmdStart)
	} else {
		command += dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.PlatformSshCmdStop)
	}

	if err := exec.Command("/bin/bash", "-c", command).Run(); err != nil {
		response.FailData(ctx, "exec command failed", err.Error())
		return
	}

	//修改数据字典
	dictionary_service.CreateCurdFactory().UpdateValById(enum_dictionary.PlatformSshSwitch, Switch)

	response.SuccessSimplify(ctx)
}

// NetworkGet  网络服务
func (u *SystemController) NetworkGet(ctx *gin.Context) {

	shellPath := cur_app.ShellNetworkGet()
	shellInit := cur_app.ShellInit()
	cmd := exec.Command(shellInit, shellPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	if err != nil {
		response.FailData(ctx, "exec command Run failed", err.Error())
		return
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if string_helper.IsEmpty(outStr) || !string_helper.IsEmpty(errStr) {
		response.FailData(ctx, "exec command output empty", fmt.Sprintf("out:\n%s\nerr:\n%s\n", outStr, errStr))
		return
	}

	var result []map[string]interface{}
	outList := strings.Split(outStr, "\n")
	for _, line := range outList {
		//过滤虚拟
		if strings.Contains(line, "veth") || strings.Contains(line, "br-") || strings.Contains(line, "docker") {
			continue
		}
		arr := strings.Split(line, " ")
		arrLen := len(arr)
		if arrLen < 4 || arrLen > 6 {
			continue
		}

		gateway := ""
		if arrLen > 4 {
			gateway = arr[4]
		}

		dns := ""
		if arrLen > 5 {
			dns = arr[5]
		}

		result = append(result, map[string]interface{}{
			"name":    arr[0],
			"proto":   arr[1],
			"ip":      arr[2],
			"mask":    arr[3],
			"gateway": gateway,
			"dns":     dns,
		})
	}
	response.SuccessReturnData(ctx, result)
}

// NetworkSet  网络服务
func (u *SystemController) NetworkSet(ctx *gin.Context) {
	name := ctx.GetString(consts.ValidatorPrefixPlatform + "name")
	proto := ctx.GetString(consts.ValidatorPrefixPlatform + "proto")
	ip := ctx.GetString(consts.ValidatorPrefixPlatform + "ip")
	mask := ctx.GetString(consts.ValidatorPrefixPlatform + "mask")
	gateway := ctx.GetString(consts.ValidatorPrefixPlatform + "gateway")
	dns := ctx.GetString(consts.ValidatorPrefixPlatform + "dns")

	shellPath := cur_app.ShellNetworkSet()

	cmd := exec.Command(cur_app.ShellInit(), shellPath, name, proto, ip, mask, gateway, dns)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	if err != nil {
		response.FailData(ctx, "exec command Run failed", err.Error())
		return
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if string_helper.IsEmpty(outStr) || !string_helper.IsEmpty(errStr) {
		response.FailData(ctx, "exec command output empty", fmt.Sprintf("out:\n%s\nerr:\n%s\n", outStr, errStr))
		return
	}

	response.SuccessSimplify(ctx)
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
		enum_dictionary.HuanshiThirdSwitch,
		enum_dictionary.HuanshiThirdIp,
		enum_dictionary.HuanshiThirdPort,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	response.SuccessReturnData(ctx, vals)
}

// HuanshiOnline  幻势配置
func (u *SystemController) HuanshiOnline(ctx *gin.Context) {
	config := huanshi_platform.RequestConfig{
		HuanshiCallbackUrl: ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_callback_url"),
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

// HuanshiSetThird  幻势配置
func (u *SystemController) HuanshiSetThird(ctx *gin.Context) {
	idAndVal := map[string]interface{}{
		enum_dictionary.HuanshiThirdSwitch: ctx.GetFloat64(consts.ValidatorPrefixPlatform + "huanshi_third_switch"),
		enum_dictionary.HuanshiThirdIp:     ctx.GetString(consts.ValidatorPrefixPlatform + "huanshi_third_ip"),
		enum_dictionary.HuanshiThirdPort:   ctx.GetFloat64(consts.ValidatorPrefixPlatform + "huanshi_third_port"),
	}

	if !dictionary_service.CreateCurdFactory().BatchUpdate(idAndVal) {
		response.FailSimplify(ctx, "保持失败")
		return
	}

	response.SuccessSimplify(ctx)
}

// UpgradeRun  系统升级
func (u *SystemController) UpgradeRun(ctx *gin.Context) {
	uploadId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "uploadId"))
	uploadModel := upload_models.CreateUploadFactory().GetItemById(uploadId)
	if uploadModel == nil || uploadModel.Id == 0 {
		response.FailSimplify(ctx, "上传升级包不存在")
		return
	}
	if !files.ExistsFile(uploadModel.FullPath) {
		response.FailSimplify(ctx, "文件不存在误")
		return
	}
	//解压
	fileNameArr := strings.Split(uploadModel.StorageName, ".")
	unzipDir := path.Dir(uploadModel.FullPath) + "/" + fileNameArr[0] + "_unzip"
	if err := zip_factory.UnzipFile(uploadModel.FullPath, unzipDir); err != nil {
		response.FailData(ctx, "解压失败", err.Error())
		return
	}

	//脚本目录
	shellDir := unzipDir + "/veda_armour_install"
	//脚本路径
	shellPath := shellDir + "/" + cur_app.ShellUpgrade()
	if !files.ExistsFile(shellPath) {
		response.FailData(ctx, "升级脚本不存在", shellPath)
		return
	}

	//执行命令
	cmd := exec.Command(cur_app.ShellInit(), shellPath, shellDir)
	if err := cmd.Run(); err != nil {
		response.FailData(ctx, "exec command failed", err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
