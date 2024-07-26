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

// DllController 模块相关
type DllController struct {
	ControllerBase
}

// PutAll 上传全部模块
func (u *DllController) PutAll(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	tasks := []int{enum_task.TypeDllSuffixIsDll}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "list")
	err := agent_post_service.CreateProcFactory().ToQueryByAll(enum_process.DataFromDll, agent_post_service.ProcTaskAll,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutRunning 上传运行中模块
func (u *DllController) PutRunning(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	tasks := []int{enum_task.TypeDllRefresh}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQueryByCommon(enum_process.DataFromDll, agent_post_service.ProcTaskRefresh,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutAuto 主动上传模块
func (u *DllController) PutAuto(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQueryByCommon(enum_process.DataFromDll, agent_post_service.ProcTaskAuto,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

type DllWhiteResp struct {
	//Id   int    `json:"id"`
	//User string `json:"user"`
	//Dir  string `json:"dir"`
	Hash string `json:"hash"`
}

// GetWhite 下载模块白名单
func (u *DllController) GetWhite(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeDllWhiteDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	var existsPolicy = make(map[string]struct{})
	result := []*DllWhiteResp{} //需初始化

	//终端策略
	policyModels := proc_models.CreateDllFactory().GetListByConditions(proc_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: loginInfo.AgentId,
		Status:  enum_process.ProcStatusWhite,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			if _, isOK := existsPolicy[model.DllHash]; !isOK {
				existsPolicy[model.DllHash] = struct{}{}
				result = append(result, &DllWhiteResp{
					Hash: model.DllHash,
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
			DataFrom:   enum_process.DataFromDll,
			OpenSwitch: &openSwitch,
		})
		if globalModels != nil && len(globalModels) > 0 {
			for _, model := range globalModels {
				if _, isOK := existsPolicy[model.Id]; !isOK {
					existsPolicy[model.Id] = struct{}{}
					result = append(result, &DllWhiteResp{
						Hash: model.Id,
					})
				}
			}
		}

	}

	response.SuccessReturnData(ctx, result)
}

// GetBlack 下载模块黑名单
func (u *DllController) GetBlack(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeDllBlackDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	var existsPolicy = make(map[string]struct{})
	result := []*DllWhiteResp{} //需初始化

	//终端策略
	policyModels := proc_models.CreateDllFactory().GetListByConditions(proc_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: loginInfo.AgentId,
		Status:  enum_process.ProcStatusBlack,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			if _, isOK := existsPolicy[model.DllHash]; !isOK {
				existsPolicy[model.DllHash] = struct{}{}
				result = append(result, &DllWhiteResp{
					Hash: model.DllHash,
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
			DataFrom:   enum_process.DataFromDll,
			OpenSwitch: &openSwitch,
		})
		if globalModels != nil && len(globalModels) > 0 {
			for _, model := range globalModels {
				if _, isOK := existsPolicy[model.Id]; !isOK {
					existsPolicy[model.Id] = struct{}{}
					result = append(result, &DllWhiteResp{
						Hash: model.Id,
					})
				}
			}
		}

	}

	response.SuccessReturnData(ctx, result)
}
