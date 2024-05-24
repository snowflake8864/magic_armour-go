package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/agent/agent_post_service"
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
