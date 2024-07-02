package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// ProcessController 进程相关
type ProcessController struct {
	ControllerBase
}

// PutAll 上传全部进程
func (u *ProcessController) PutAll(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	tasks := []int{enum_task.TypeProcessSuffixIsExe}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "list")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.DataFromProcess, agent_post_service.ProcTaskAll,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutRunning 上传运行中进程
func (u *ProcessController) PutRunning(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	tasks := []int{enum_task.TypeProcessRefresh}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.DataFromProcess, agent_post_service.ProcTaskRefresh,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutAuto 主动上传进程
func (u *ProcessController) PutAuto(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.DataFromProcess, agent_post_service.ProcTaskAuto,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

type ProcessWhiteResp struct {
	Id   int    `json:"id"`
	User string `json:"user"`
	Dir  string `json:"dir"`
	Hash string `json:"hash"`
}
type ProcessAssignResp struct {
	Pid  int    `json:"pid"`
	Hash string `json:"hash"`
}

// GetWhite 下载进程白名单
func (u *ProcessController) GetWhite(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeProcessWhiteDownload}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	var existsPolicy = make(map[string]struct{})
	result := []*ProcessWhiteResp{}

	//终端策略
	policyModels := proc_models.CreateProcessFactory().GetListByConditions(proc_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: loginInfo.AgentId,
		Status:  enum_process.ProcStatusWhite,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			if _, isOK := existsPolicy[model.ProcHash]; !isOK {
				existsPolicy[model.ProcHash] = struct{}{}
				result = append(result, &ProcessWhiteResp{
					Hash: model.ProcHash,
				})
			}
		}
	}

	//是否开启同步全局名单库
	if val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentProcGlobalSwitch); string_helper.ToInt(val) == enums.SwitchOn {
		openSwitch := enums.SwitchOn
		//全局策略
		globalModels := proc_models.CreateProcGlobalFactory().GetListByConditions(proc_models.SearchConditions{
			OrgId:      loginInfo.OrgId,
			GlobalType: enum_process.GlobalTypeWhite,
			DataFrom:   enum_process.DataFromProcess,
			OpenSwitch: &openSwitch,
		})
		if globalModels != nil && len(globalModels) > 0 {
			for _, model := range globalModels {
				if _, isOK := existsPolicy[model.Id]; !isOK {
					existsPolicy[model.Id] = struct{}{}
					result = append(result, &ProcessWhiteResp{
						Hash: model.Id,
					})
				}
			}
		}

	}

	response.SuccessReturnData(ctx, gin.H{
		"proclist": result,
	})
}

// GetBlack 下载进程黑名单
func (u *ProcessController) GetBlack(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeProcessBlackDownload}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	var existsPolicy = make(map[string]struct{})
	result := []*ProcessWhiteResp{}

	//终端策略
	policyModels := proc_models.CreateProcessFactory().GetListByConditions(proc_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: loginInfo.AgentId,
		Status:  enum_process.ProcStatusBlack,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			if _, isOK := existsPolicy[model.ProcHash]; !isOK {
				existsPolicy[model.ProcHash] = struct{}{}
				result = append(result, &ProcessWhiteResp{
					Hash: model.ProcHash,
				})
			}
		}
	}

	//是否开启同步全局名单库
	if val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentProcGlobalSwitch); string_helper.ToInt(val) == enums.SwitchOn {
		openSwitch := enums.SwitchOn
		//全局策略
		globalModels := proc_models.CreateProcGlobalFactory().GetListByConditions(proc_models.SearchConditions{
			OrgId:      loginInfo.OrgId,
			GlobalType: enum_process.GlobalTypeBlack,
			DataFrom:   enum_process.DataFromProcess,
			OpenSwitch: &openSwitch,
		})
		if globalModels != nil && len(globalModels) > 0 {
			for _, model := range globalModels {
				if _, isOK := existsPolicy[model.Id]; !isOK {
					existsPolicy[model.Id] = struct{}{}
					result = append(result, &ProcessWhiteResp{
						Hash: model.Id,
					})
				}
			}
		}

	}

	response.SuccessReturnData(ctx, gin.H{
		"proclist": result,
	})
}

// GetAssign 指定进程:下载信息
func (u *ProcessController) GetAssign(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeProcessCollectDll}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	result := []*ProcessAssignResp{}

	//终端策略
	policyModels := proc_models.CreateProcessRunningFactory().GetItemsByOrgAndAgentAnDllStatus(loginInfo.OrgId, loginInfo.AgentId, enum_process.CollectDllStatusDefault)
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			result = append(result, &ProcessAssignResp{
				Pid:  model.ProcId,
				Hash: model.ProcHash,
			})
		}
	}

	response.SuccessReturnData(ctx, result)
}

// PutAssign 指定进程:上传对应模块
func (u *ProcessController) PutAssign(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	tasks := []int{enum_task.TypeProcessCollectDll}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().CollectDll(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
