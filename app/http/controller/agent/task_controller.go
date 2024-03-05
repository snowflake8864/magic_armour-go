package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

// TaskController 任务相关
type TaskController struct {
}

// Get 获取任务
func (u *TaskController) Get(ctx *gin.Context) {
	currentAgent := cur_userinfo.GetCurrentAgent(ctx)

	//未完成的任务, 不支持的任务
	var unfinishedTasks, notSupportTasks []int
	//有获取配置任务, 有卸载任务, 有需要开启保护模式前下发的任务
	var hasConfig, hasUninstall, hasProjectModeBefore bool

	taskCacheFact := task_service.CreateCacheFactory()
	//缓存中的任务
	cacheTasks := taskCacheFact.GetUnfinished(currentAgent.AgentId)
	if len(cacheTasks) > 0 {
		for _, val := range cacheTasks {
			//是否有配置任务
			if val == enum_task.TypeAgentConfigDownload {
				hasConfig = true
				continue
			}
			//是否有卸载任务
			if val == enum_task.TypeAgentConfigDownload {
				hasUninstall = true
				break
			}
			//win xp 不支持的功能
			if currentAgent.OsType == enum_software.OsTypeWinXp && u.checkWinXpNotSupport(val) {
				notSupportTasks = append(notSupportTasks, val)
				continue
			}
			//linux不支持的功能
			if currentAgent.PlatformType == enum_software.PlatformTypeLinux && u.checkLinuxNotSupport(val) {
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

	//不支持的任务直接关闭
	if len(notSupportTasks) > 0 {
		
	}

	response.SuccessReturnData(ctx, gin.H{"tasklist": unfinishedTasks})
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

	response.SuccessReturnData(ctx, "关闭任务")
}
