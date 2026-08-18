package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/service/agent/plug_service/plug_calc_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// PlugController 封堵相关
type PlugController struct {
	ControllerBase
}

// GetDynamic 下载动态封堵
func (u *PlugController) GetDynamic(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypePluggingDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	//win xp 不支持的功能
	if loginInfo.OsType == enum_software.OsTypeWinXp {
		response.SuccessSimplify(ctx)
		return
	}
	openSwitch := enums.SwitchOn
	var existsPolicy = make(map[string]struct{})
	result := []*PlugDynamicResp{} //需初始化

	//终端策略
	policyModels := plug_models.CreateDynamicFactory().GetListByConditions(plug_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    loginInfo.AgentId,
		OpenSwitch: &openSwitch,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			policy := u.formatDynamic(model.Source, model.Ip, model.Direction, model.Duration, model.StartAt, model.EndAt)
			if _, isOK := existsPolicy[policy.Ip]; !isOK {
				existsPolicy[policy.Ip] = struct{}{}
				result = append(result, policy)
			}
		}
	}

	response.SuccessReturnData(ctx, result)
}

type PlugDynamicResp struct {
	Source    string      `json:"type"` //必须是字符串
	TypeName  string      `json:"typeName"`
	Ip        string      `json:"ip"`
	Direction int         `json:"direction"`
	Duration  int         `json:"duration"`
	StartAt   string      `json:"starttime"`
	EndAt     interface{} `json:"endtime"`
}

func (u *PlugController) formatDynamic(source int, ip string, direction, duration int, startAt, endAt string) *PlugDynamicResp {
	data := &PlugDynamicResp{}
	data.Source = int_helper.IntToStr(source)
	data.Ip = ip
	data.Direction = direction
	data.Duration = duration
	data.StartAt = startAt
	if string_helper.IsEmpty(endAt) {
		data.EndAt = 0
	} else {
		data.EndAt = endAt
	}

	return data
}

// PutDynamic 上传动态封堵(调试用)
func (u *PlugController) PutDynamic(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	Ip := ctx.GetString(consts.ValidatorPrefixAgent + "ip")
	Source := int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "source"))
	PlugType := int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "plugType"))
	Duration := int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "duration"))

	if err := plug_calc_service.CreatePlugCalcFactory().Add(loginInfo.OrgId, loginInfo.AgentId, Ip, Source, PlugType, Duration); err != nil {
		response.FailRoutine(ctx, my_errors.AgentParamErrorCode, err.Error())
	}

	response.SuccessReturnData(ctx, "上传动态封堵(调试用)")
}

// GetIpWhite 下载封堵ip白名单(不在使用)
func (u *PlugController) GetIpWhite(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypePlugIpWhiteDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	response.SuccessReturnData(ctx, "下载封堵ip白名单(不在使用)")
}

type PlugIpResp struct {
	Ip        string `json:"ip"`
	Direction int    `json:"direction"`
}

// GetIpBlack 下载封堵ip黑名单
func (u *PlugController) GetIpBlack(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypePlugIpBlackDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	//win xp 不支持的功能
	if loginInfo.OsType == enum_software.OsTypeWinXp {
		response.SuccessSimplify(ctx)
		return
	}
	openSwitch := enums.SwitchOn
	var existsPolicy = make(map[string]struct{})
	result := []*PlugIpResp{} //需初始化

	//终端策略
	policyModels := plug_models.CreateIpPolicyFactory().GetListByConditions(plug_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    loginInfo.AgentId,
		IpType:     enum_plug.IpTypeBlack,
		OpenSwitch: &openSwitch,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			policy := u.formatIp(model.Ip, model.Direction)
			if _, isOK := existsPolicy[policy.Ip]; !isOK {
				existsPolicy[policy.Ip] = struct{}{}
				result = append(result, policy)
			}
		}
	}

	//全局策略
	globalModels := plug_models.CreateIpGlobalFactory().GetListByConditions(plug_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		DeptIds:    u.getDeptIdsByAgent(loginInfo.OrgId, loginInfo.AgentId),
		IpType:     enum_plug.IpTypeBlack,
		OpenSwitch: &openSwitch,
	})
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			policy := u.formatIp(model.Ip, model.Direction)
			if _, isOK := existsPolicy[policy.Ip]; !isOK {
				existsPolicy[policy.Ip] = struct{}{}
				result = append(result, policy)
			}
		}
	}

	response.SuccessReturnData(ctx, result)
}

func (u *PlugController) formatIp(ip string, direction int) *PlugIpResp {
	data := &PlugIpResp{}

	data.Ip = ip
	data.Direction = direction

	return data
}
