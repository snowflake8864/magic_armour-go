package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// TaskController 任务相关
type TaskController struct {
	ControllerBase
}

// Get 获取任务
func (u *TaskController) Get(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	//未完成的任务, 不支持的任务
	var unfinishedTasks, notSupportTasks []int
	//有获取配置任务, 有卸载任务, 需要开启保护模式前下发的任务, 有重置终端标识任务, 有删除的任务
	var hasConfig, hasUninstall, hasProjectModeBefore, hasResetUuid, hasDelete bool

	taskCacheFact := task_service.CreateCacheFactory()
	//缓存中的任务
	cacheTasks := taskCacheFact.GetUnfinished(loginInfo.OrgId, loginInfo.AgentId)
	if len(cacheTasks) > 0 {
		for _, val := range cacheTasks {
			//是否有配置任务
			if val == enum_task.TypeAgentConfigDownload {
				hasConfig = true
				continue
			}
			//是否有卸载任务
			if val == enum_task.TypeAgentUninstall {
				hasUninstall = true
				break
			}
			//是否有重置终端标识任务
			if val == enum_task.TypeResetUUid {
				hasResetUuid = true
				break
			}
			//是否有删除的任务
			if val == enum_task.TypeAgentDelete {
				hasDelete = true
				break
			}
			//win xp 不支持的功能
			if loginInfo.OsType == enum_software.OsTypeWinXp && u.checkWinXpNotSupport(val) {
				notSupportTasks = append(notSupportTasks, val)
				continue
			}
			//linux不支持的功能
			if loginInfo.PlatformType == enum_software.PlatformTypeLinux && u.checkLinuxNotSupport(val) {
				notSupportTasks = append(notSupportTasks, val)
				continue
			}
			//有需要开启保护模式前下发的任务
			if u.hasProjectModeBefore(val) {
				hasProjectModeBefore = true
			}
			unfinishedTasks = append(unfinishedTasks, val)
		}
	}

	taskCurdFact := task_service.CreateCurdFactory()

	//有卸载任务
	if hasUninstall {
		//只下发卸载任务
		unfinishedTasks = []int{enum_task.TypeAgentUninstall}

		//任务加入不支持列表，等待关闭
		notSupportTasks = append(notSupportTasks, enum_task.TypeAgentUninstall)
	} else if hasResetUuid { //重置终端标识任务
		//只下重置终端标识任务
		unfinishedTasks = []int{enum_task.TypeResetUUid}
	} else if hasDelete { // 删除的任务
		//删除所有任务相关数据，保证下次正常登录后没有遗留任务
		task_service.CreateCurdFactory().DeleteForAgent(loginInfo.OrgId, loginInfo.AgentId)

		response.ErrorTokenFail(ctx, my_errors.AgentLoginExpiredCode, "agent delete")
		return
	} else {
		if taskCacheFact.LockAgentLive(loginInfo.OrgId, loginInfo.AgentId) {
			//设置存活状态
			if isExist, _ := agent_service.CreateCurdFactory().SetLivingByAgent(loginInfo.OrgId, loginInfo.AgentId); !isExist {
				response.ErrorTokenFail(ctx, my_errors.AgentLoginExpiredCode, "agent not exist")
				return
			}
		}

		//是否有未完成的任务
		var hasUnfinishedTask bool
		if len(unfinishedTasks) > 0 || hasConfig {
			hasUnfinishedTask = true
		}

		//设置心跳并判断是否需要追加修改配置的任务
		appendConfig := taskCacheFact.SetHeartbeat(loginInfo.OrgId, loginInfo.AgentId, hasUnfinishedTask)
		//需追加配置 && 这次没有配置任务
		if appendConfig && !hasConfig {
			hasConfig = true
		}

		//@todo 下载配置信息放在最后一位 http://192.168.19.9/zentao/story-view-636.html
		if hasConfig {
			//@todo 如果指令集中同时存在保护模式和相关内容，配置暂缓下发，等待获取相关内容完成
			if hasProjectModeBefore {
				//全是监控模式
				if !u.checkAgentHasProject(loginInfo.OrgId, loginInfo.AgentId) {
					unfinishedTasks = append(unfinishedTasks, enum_task.TypeAgentConfigDownload)
				}
			} else {
				unfinishedTasks = append(unfinishedTasks, enum_task.TypeAgentConfigDownload)
			}
		}
	}

	//不支持的任务直接关闭
	if len(notSupportTasks) > 0 {
		taskCurdFact.CloseTasks(loginInfo.OrgId, loginInfo.AgentId, notSupportTasks, ctx.ClientIP())
	}

	if unfinishedTasks == nil {
		unfinishedTasks = []int{}
	}
	response.SuccessReturnData(ctx, gin.H{"tasklist": unfinishedTasks})
}

// 验证终端是否有保护模式
func (u *TaskController) checkAgentHasProject(orgId int64, agentId string) bool {
	agentCurdFact := agent_models.CreateAgentFactory()
	agentModel := agentCurdFact.GetItemByOrgAndAgent(orgId, agentId, agentCurdFact.ColumnSwitchAndModeForAlarm(true))
	if agentModel.ProcessSwitch == enums.SwitchOn && agentModel.ProcessMode == enum_agent.ModeProject {
		return true
	}
	if agentModel.DllSwitch == enums.SwitchOn && agentModel.DllMode == enum_agent.ModeProject {
		return true
	}
	if agentModel.TamperSwitch == enums.SwitchOn && agentModel.TamperMode == enum_agent.ModeProject {
		return true
	}
	if agentModel.BlackmailSwitch == enums.SwitchOn && agentModel.BlackmailMode == enum_agent.ModeProject {
		return true
	}
	if agentModel.UsbSwitch == enums.SwitchOn && agentModel.UsbMode == enum_agent.ModeProject {
		return true
	}
	return false
}

// win xp 不支持的功能
func (u *TaskController) checkWinXpNotSupport(task int) bool {
	list := map[int]struct{}{
		enum_task.TypePortWorkingUpload:         {},
		enum_task.TypePortVirtualPolicyDownload: {},
		enum_task.TypePluggingDownload:          {},
		enum_task.TypePlugIpBlackDownload:       {},
	}

	_, support := list[task]
	return support
}

// linux不支持的功能
func (u *TaskController) checkLinuxNotSupport(task int) bool {
	list := map[int]struct{}{
		enum_task.TypeDllRefresh:        {},
		enum_task.TypeDllWhiteDownload:  {},
		enum_task.TypeDllBlackDownload:  {},
		enum_task.TypeDllSuffixIsDll:    {},
		enum_task.TypeProcessCollectDll: {},
	}
	_, support := list[task]
	return support
}

// 有需要开启保护模式前下发的任务
func (u *TaskController) hasProjectModeBefore(task int) bool {
	list := map[int]struct{}{
		enum_task.TypeProcessWhiteDownload:    {},
		enum_task.TypeProcessBlackDownload:    {},
		enum_task.TypeDllWhiteDownload:        {},
		enum_task.TypeDllBlackDownload:        {},
		enum_task.TypeTamperPolicyDownload:    {},
		enum_task.TypeBlackmailPolicyDownload: {},
		enum_task.TypeUsbWhiteDownload:        {},
		enum_task.TypeUsbBlackDownload:        {},
	}
	_, support := list[task]
	return support
}

// Close 关闭任务
func (u *TaskController) Close(ctx *gin.Context) {
	taskType := int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "tasklist"))
	types := []int{taskType}

	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, types, ctx.ClientIP())

	response.SuccessReturnData(ctx, my_errors.CurdCloseTask)

}
