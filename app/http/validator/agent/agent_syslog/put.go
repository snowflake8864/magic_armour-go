package agent_syslog

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/agent/common_data"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type Put struct {
	common_data.SyslogUpload
}

func (n Put) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixAgent, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[agent Syslog Put]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&agent.SyslogController{}).Put(extraAddBindDataContext)
	}

}

type VulnPut struct {
	StartAt   string `form:"start_at" json:"start_at" binding:"required,datetime=2006-01-02 15:04:05"`
	EndAt     string `form:"end_at" json:"end_at" binding:"required,datetime=2006-01-02 15:04:05"`
	VulnTotal int    `form:"vuln_total" json:"vuln_total"`
	VulnList  string `form:"vuln_list" json:"vuln_list"`
}

func (n VulnPut) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixAgent, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[agent Syslog VulnPut]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&agent.SyslogController{}).VulnPut(extraAddBindDataContext)
	}

}
