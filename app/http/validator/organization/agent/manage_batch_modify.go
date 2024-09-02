package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/http/validator/common/data_type"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type ManageBatchModify struct {
	data_type.AgentBatch
	DeptId              *int `form:"deptId" json:"deptId" binding:"required"`
	SyslogDnsSwitch     *int `form:"syslogDnsSwitch" json:"syslogDnsSwitch" binding:"required,oneof=-100 0 1"`
	SyslogOuterSwitch   *int `form:"syslogOuterSwitch" json:"syslogOuterSwitch" binding:"required,oneof=-100 0 1"`
	SyslogInnerSwitch   *int `form:"syslogInnerSwitch" json:"syslogInnerSwitch" binding:"required,oneof=-100 0 1"`
	SyslogProcessSwitch *int `form:"syslogProcessSwitch" json:"syslogProcessSwitch" binding:"required,oneof=-100 0 1"`
	SyslogLoginSwitch   *int `form:"syslogLoginSwitch" json:"syslogLoginSwitch" binding:"required,oneof=-100 0 1"`
	//OutreachAlarmSwitch *int `form:"outreachAlarmSwitch" json:"outreachAlarmSwitch" binding:"required,oneof=-100 0 1"`
	//SelfProtectSwitch   *int `form:"selfProtectSwitch" json:"selfProtectSwitch" binding:"required,oneof=-100 0 1"`
	AutoWhiteSwitch *int `form:"autoWhiteSwitch" json:"autoWhiteSwitch" binding:"required,oneof=-100 0 1"`
}

func (n ManageBatchModify) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixOrg, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[organization agent ManageBatchModify]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&organization.AgentController{}).BatchModifyForManage(extraAddBindDataContext)
	}

}
