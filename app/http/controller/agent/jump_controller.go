package agent

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_ip_jump"
	"goskeleton/app/global/enums/enum_pw_jump"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/ip_jump_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/response"
)

// JumpController 封堵相关
type JumpController struct {
	ControllerBase
}

// GetPwJump 获取密码跳变
func (u *JumpController) GetPwJump(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypePwJump}, ctx.ClientIP())

	config := enum_pw_jump.PwJumpConfig{}
	scheme := 1
	//优先终端策略
	if config = agent_service.CreateExtendFactory().GetPwJumpConfig(loginInfo.OrgId, loginInfo.AgentId); config.Switch == enums.SwitchOff {
		OrgExtendFact := organize_service.CreateOrgExtendFactory()
		config = OrgExtendFact.FormatPwJumpConfig(OrgExtendFact.GetItemByOrg(loginInfo.OrgId).PwJumpConfig)
		if config.Switch == enums.SwitchOff {
			response.FailSimplify(ctx, "Switch Off")
			return
		}
		scheme = 2
	}

	pw := config.GeneratePw()

	response.SuccessReturnData(ctx, gin.H{
		"pw":     pw,
		"scheme": scheme,
		"config": config,
	})
}

// PutPwJump 上报密码跳变结果
func (u *JumpController) PutPwJump(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	newest := enum_pw_jump.PwJumpNewest{}
	newest.Status = int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "status"))
	newest.Reason = ctx.GetString(consts.ValidatorPrefixAgent + "reason")
	newest.User = ctx.GetString(consts.ValidatorPrefixAgent + "user")
	newest.Pw = ctx.GetString(consts.ValidatorPrefixAgent + "pw")

	//保存密码
	if err := agent_service.CreateExtendFactory().SetPwJumpNewest(loginInfo.OrgId, loginInfo.AgentId, newest); err != nil {
		response.FailSimplify(ctx, "save pw fail")
		return
	}

	response.SuccessSimplify(ctx)
}

// GetIpJump 获取IP跳变
func (u *JumpController) GetIpJump(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeIpJump}, ctx.ClientIP())

	//企业策略
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	config := OrgExtendFact.FormatIpJumpConfig(OrgExtendFact.GetItemByOrg(loginInfo.OrgId).IpJumpConfig)
	if config.Switch == enums.SwitchOff {
		response.FailSimplify(ctx, "Switch Off")
		return
	}
	sourceIp := ctx.ClientIP()
	//生成IP跳变池
	if err, targetIp := ip_jump_service.CreatCurdFactory(loginInfo.OrgId, &config).NewestIp(sourceIp); err != nil {
		response.FailSimplify(ctx, "Switch Off")
		return
	} else {
		response.SuccessReturnData(ctx, gin.H{
			"source_ip":  sourceIp,
			"target_ip":  fmt.Sprintf("%s/%d", targetIp, config.Mask),
			"gateway":    config.Gateway,
			"mode":       1,  //模式：1追加, 2修改
			"allow_size": 3,  //追加模式下允许保留数目
			"aging_time": 30, //追加模式下老化时间(分钟)
		})
		return
	}
}

// PutIpJump 上报地址跳变结果
func (u *JumpController) PutIpJump(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	newest := enum_ip_jump.IpJumpNewest{}
	newest.Status = int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "status"))
	newest.Reason = ctx.GetString(consts.ValidatorPrefixAgent + "reason")
	newest.SourceIp = ctx.GetString(consts.ValidatorPrefixAgent + "source_ip")
	newest.TargetIp = ctx.GetString(consts.ValidatorPrefixAgent + "target_ip")
	newest.Gateway = ctx.GetString(consts.ValidatorPrefixAgent + "gateway")
	newest.AgentIp = ctx.GetString(consts.ValidatorPrefixAgent + "agent_ip")

	//保存密码
	if err := agent_service.CreateExtendFactory().SetIpJumpNewest(loginInfo.OrgId, loginInfo.AgentId, newest); err != nil {
		response.FailSimplify(ctx, "save ip fail")
		return
	}

	//修改终端
	if newest.Status == enum_ip_jump.NewestStatusSuccess {
		agent_service.CreateCurdFactory().SetIpByAgent(loginInfo.OrgId, loginInfo.AgentId, newest.AgentIp, ctx.ClientIP())
	}

	response.SuccessSimplify(ctx)
}
