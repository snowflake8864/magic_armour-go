package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/syslog_service/job_callback_syslog"
	"goskeleton/app/utils/response"
)

// SyslogController 通信日志相关
type SyslogController struct {
	ControllerBase
}

// Put 上传通讯日志
func (u *SyslogController) Put(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "list")
	err := agent_post_service.CreateSyslogFactory().ToQuery(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// Put 上传通讯日志
func (u *SyslogController) VulnPut(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	model := syslog_models.VulnModel{
		BaseModel:   model.BaseModel{},
		OrgId:       loginInfo.OrgId,
		AgentId:     loginInfo.AgentId,
		ScanStartAt: ctx.GetString(consts.ValidatorPrefixAgent + "start_at"),
		ScanEndAt:   ctx.GetString(consts.ValidatorPrefixAgent + "start_at"),
		VulnTotal:   int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "vuln_total")),
		VulnInfo:    ctx.GetString(consts.ValidatorPrefixAgent + "vuln_list"),
	}

	err := job_callback_syslog.CreateVulnUploadFactory().Add(&model)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
