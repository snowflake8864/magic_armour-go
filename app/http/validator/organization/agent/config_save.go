package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/http/validator/common/data_type"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type ConfigSave struct {
	data_type.AgentSingle
	ProcessSwitch       *int `form:"processSwitch" json:"processSwitch" binding:"required,oneof=0 1"`
	ProcessMode         int  `form:"processMode" json:"processMode" binding:"required,oneof=1 2"`
	DllSwitch           *int `form:"dllSwitch" json:"dllSwitch" binding:"required,oneof=0 1"`
	DllMode             int  `form:"dllMode" json:"dllMode" binding:"required,oneof=1 2"`
	TamperSwitch        *int `form:"tamperSwitch" json:"tamperSwitch" binding:"required,oneof=0 1"`
	TamperMode          int  `form:"tamperMode" json:"tamperMode" binding:"required,oneof=1 2"`
	BlackmailSwitch     *int `form:"blackmailSwitch" json:"blackmailSwitch" binding:"required,oneof=0 1"`
	BlackmailMode       int  `form:"blackmailMode" json:"blackmailMode" binding:"required,oneof=1 2"`
	PortSwitch          *int `form:"portSwitch" json:"portSwitch" binding:"required,oneof=0 1"`
	PlugSwitch          *int `form:"plugSwitch" json:"plugSwitch" binding:"required,oneof=0 1"`
	UsbSwitch           *int `form:"usbSwitch" json:"usbSwitch" binding:"required,oneof=0 1"`
	UsbMode             int  `form:"usbMode" json:"usbMode" binding:"required,oneof=1 2"`
	LoginSolveSwitch    *int `form:"loginSolveSwitch" json:"loginSolveSwitch" binding:"required,oneof=0 1"`
	SyslogDnsSwitch     *int `form:"syslogDnsSwitch" json:"syslogDnsSwitch" binding:"required,oneof=0 1"`
	SyslogOuterSwitch   *int `form:"syslogOuterSwitch" json:"syslogOuterSwitch" binding:"required,oneof=0 1"`
	SyslogInnerSwitch   *int `form:"syslogInnerSwitch" json:"syslogInnerSwitch" binding:"required,oneof=0 1"`
	SyslogProcessSwitch *int `form:"syslogProcessSwitch" json:"syslogProcessSwitch" binding:"required,oneof=0 1"`
	SyslogLoginSwitch   *int `form:"syslogLoginSwitch" json:"syslogLoginSwitch" binding:"required,oneof=0 1"`
	OutreachAlarmSwitch *int `form:"outreachAlarmSwitch" json:"outreachAlarmSwitch" binding:"required,oneof=0 1"`
	SelfProtectSwitch   *int `form:"selfProtectSwitch" json:"selfProtectSwitch" binding:"required,oneof=0 1"`
	AutoWhiteSwitch     *int `form:"autoWhiteSwitch" json:"autoWhiteSwitch" binding:"required,oneof=0 1"`
	InternetSwitch      *int `form:"internetSwitch" json:"internetSwitch" binding:"required,oneof=0 1"`
	CommonTime          int  `form:"commonTime" json:"commonTime" binding:"required,min=20,max=300"`
	FastTime            int  `form:"fastTime" json:"fastTime" binding:"required,min=4,max=60"`
}

func (n ConfigSave) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixOrg, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[organization agent ConfigSave]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&organization.AgentController{}).SaveForConfig(extraAddBindDataContext)
	}

}
