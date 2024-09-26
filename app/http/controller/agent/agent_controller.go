package agent

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/jobs/agent_job"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/trust_dir_models"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/agent/agent_err_service"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/service/platform/upload_service"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
	"strings"
)

// AgentController 终端相关
type AgentController struct {
	ControllerBase
}

func (u *AgentController) getMode(mode int) int {
	if mode == enum_agent.ModeMonitor {
		return 0
	} else {
		return 1
	}
}

// GetConfig 获取配置信息
func (u *AgentController) GetConfig(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeAgentConfigDownload}, ctx.ClientIP())

	//终端配置
	agentCurdFact := agent_models.CreateAgentFactory()
	agentModel := agentCurdFact.GetItemByOrgAndAgent(loginInfo.OrgId, loginInfo.AgentId, agentCurdFact.ColumnConfig(true, false))

	//数据字典
	dictIds := []string{
		enum_dictionary.AgentHttpUrl,
		enum_dictionary.AgentHttpPort,
		enum_dictionary.AgentLogSendProto,
		enum_dictionary.AgentOutreachDetectTime,
	}
	mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	mapConf := make(map[string]interface{})
	//服务器地址
	mapConf["serveripport"] = mapDict[enum_dictionary.AgentHttpUrl] + ":" + mapDict[enum_dictionary.AgentHttpPort]
	//暂时无使用
	mapConf["logipport"] = ""
	mapConf["logproto"] = string_helper.ToInt(mapDict[enum_dictionary.AgentLogSendProto])
	mapConf["logsent"] = 0
	mapConf["api_port"] = string_helper.ToInt(mapDict[enum_dictionary.AgentHttpPort])
	mapConf["syslog_port"] = 0

	//进程
	mapConf["proc_switch"] = agentModel.ProcessSwitch
	mapConf["proc_protect"] = u.getMode(agentModel.ProcessMode)
	//模块
	mapConf["module_switch"] = agentModel.DllSwitch
	mapConf["module_protect"] = u.getMode(agentModel.DllMode)

	//防篡改
	mapConf["file_switch"] = agentModel.TamperSwitch
	mapConf["file_protect"] = u.getMode(agentModel.TamperMode)
	//勒索
	mapConf["extortion_switch"] = agentModel.BlackmailSwitch
	mapConf["extortion_protect"] = u.getMode(agentModel.BlackmailMode)
	//端口虚开
	mapConf["open_port_switch"] = agentModel.PortSwitch
	//封堵
	mapConf["dynamic_switch"] = agentModel.PlugSwitch
	//外设
	mapConf["usb_switch"] = agentModel.UsbSwitch
	mapConf["usb_protect"] = u.getMode(agentModel.UsbMode)
	//自保
	mapConf["self_protect_switch"] = agentModel.SelfProtectSwitch
	//内核升级
	mapConf["kernel_upgrade_switch"] = agentModel.KernelUpgradeSwitch

	//外联管控间隔时间(秒)
	mapConf["outreach_switch"] = agentModel.OutreachAlarmSwitch
	outreachTime := 300
	if val, isOK := mapDict[enum_dictionary.AgentOutreachDetectTime]; isOK {
		outreachTime = string_helper.ToInt(val)
	}
	mapConf["outreach_time"] = outreachTime

	//agent通讯日志
	mapConf["syslog_switch"] = enums.SwitchOn
	mapConf["syslog_dns_switch"] = agentModel.SyslogDnsSwitch
	mapConf["syslog_outer_switch"] = agentModel.SyslogOuterSwitch
	mapConf["syslog_inner_switch"] = agentModel.SyslogInnerSwitch
	mapConf["syslog_process_switch"] = agentModel.SyslogProcessSwitch
	mapConf["syslog_login_switch"] = agentModel.SyslogLoginSwitch
	//上传硬件信息间隔时间(秒)
	mapConf["hardware_switch"] = agentModel.HardwareSwitch
	mapConf["hardware_time"] = agentModel.HardwareTime

	//linux 不支持的功能
	if loginInfo.PlatformType == enum_software.PlatformTypeLinux {
		mapConf["module_switch"] = enums.SwitchOff
		mapConf["module_protect"] = u.getMode(agentModel.DllMode)
	}
	//win xp 不支持的功能
	if loginInfo.OsType == enum_software.OsTypeWinXp {
		mapConf["open_port_switch"] = enums.SwitchOff
		mapConf["dynamic_switch"] = enums.SwitchOff
	}

	//心跳间隔时间(秒)
	if heartFlag := task_service.CreateCacheFactory().GetHeartbeatFlag(loginInfo.OrgId, loginInfo.AgentId); heartFlag {
		mapConf["crontime"] = agentModel.CommonTime
	} else {
		mapConf["crontime"] = agentModel.FastTime
	}

	response.SuccessReturnData(ctx, gin.H{"conf": mapConf})
}

// GetInfo 获取终端信息(调试用)
func (u *AgentController) GetInfo(ctx *gin.Context) {
	response.SuccessReturnData(ctx, u.getLogin(ctx))
}

// Upgrade 升级
func (u *AgentController) Upgrade(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeAgentUpgrade}, ctx.ClientIP())

	//客户端最新版本
	IsPublish := enums.IsYes
	Conditions := software_models.SearchConditions{
		PlatformType: loginInfo.PlatformType,
		OsType:       loginInfo.OsType,
		IsPublish:    &IsPublish,
	}
	model := software_models.CreateSoftwareFactory().GetItemByConditions(Conditions)
	if model == nil || model.Id == 0 {
		response.FailRoutine(ctx, consts.AgentParamErrorCode, "对应版本已下架:"+model.Ver)
		return
	}
	if !enum_software.CanUpgradeDirect(model.Ver, loginInfo.Ver) {
		response.FailRoutine(ctx, consts.AgentParamErrorCode, "已经是最新版本,无需升级")
		return
	}

	downloadFact := software_service.CreateDownloadFactory()
	err, addr, _, _ := downloadFact.GetServiceAddr()
	if err != nil {
		response.FailRoutine(ctx, consts.AgentParamErrorCode, err.Error())
		return
	}

	param := software_models.DownloadModel{
		OrgId:      loginInfo.OrgId,
		SoftwareId: model.Id,
		Remark:     "升级",
	}
	if err1, downloadModel := software_service.CreateDownloadFactory().GetDownloadInfo(&param); err != nil {
		response.FailRoutine(ctx, consts.AgentParamErrorCode, err1.Error())
		return
	} else {
		response.SuccessReturnData(ctx, gin.H{
			"alias":    downloadModel.Alias,
			"size":     fmt.Sprintf("%d", downloadModel.FileSize),
			"hash":     downloadModel.FileHash,
			"download": addr + "/" + strings.TrimLeft(downloadModel.StoragePath, "/"),
		})
		return
	}
}

// Uninstall 卸载回调
func (u *AgentController) Uninstall(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	agent_job.Delete(loginInfo.OrgId, loginInfo.AgentId)

	response.SuccessReturnData(ctx, "Uninstall")
}

// GetTrustDir 获取信任目录
func (u *AgentController) GetTrustDir(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeTrustDir}, ctx.ClientIP())

	OpenSwitch := enums.SwitchOn
	conditions := trust_dir_models.SearchConditions{
		OrgId:        loginInfo.OrgId,
		PlatformType: loginInfo.PlatformType,
		DeptIds:      u.getDeptIdsByAgent(loginInfo.OrgId, loginInfo.AgentId),
		OpenSwitch:   &OpenSwitch,
	}

	var result []map[string]interface{}

	//全局目录
	globalModels := trust_dir_models.CreateGlobalFactory().GetListByConditions(conditions)
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			dir := model.Dir
			if loginInfo.PlatformType == enum_software.PlatformTypeWindows {
				dir = path_helper.ToWindows(dir)
			}
			result = append(result, map[string]interface{}{
				"type":      model.Type,
				"dir":       dir,
				"is_extend": model.IsExtend,
			})
		}
	}

	response.SuccessReturnData(ctx, result)
}

// UploadErrorFile 上传转储(错误) 文件
func (u *AgentController) UploadErrorFile(ctx *gin.Context) {
	remark, _ := ctx.GetPostForm("remark")
	savePath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath")
	if r, uploadParam := upload_file.Upload(ctx, savePath); r == true {
		uploadParam.Type = "agent_error"
		loginInfo := u.getLogin(ctx)
		uploadParam.OperateSource = enums.OperateSourceAgent
		uploadParam.OperateId = loginInfo.OrgId
		uploadParam.OperateName = loginInfo.AgentId
		if !upload_service.CreateCurdFactory().Add(&uploadParam) {
			response.FailSimplify(ctx, "保存失败")
			return
		}

		if err := agent_err_service.CreateErrorCurdFactory().Add(loginInfo, remark, &uploadParam); err != nil {
			response.FailSimplify(ctx, err.Error())
			return
		}
		response.SuccessSimplify(ctx)
	} else {
		response.Fail(ctx, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
	}
}

// MonitorHardware 监控硬件资源
func (u *AgentController) MonitorHardware(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "info")
	err := agent_post_service.CreateAgentFactory().MonitorToQuery(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
