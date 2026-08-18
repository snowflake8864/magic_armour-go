package login_solve

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/http/validator/common/data_type"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/http/validator/organization/agent_common_data"
	"goskeleton/app/utils/response"
)

type PolicyConfigSet struct {
	data_type.AgentSingle
	AddConfig
	agent_common_data.CommonSwitch
}

func (n PolicyConfigSet) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixOrg, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[organization login_solve PolicyConfigSet]Form validator JSON conversion failed")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&organization.LoginSolveController{}).ConfigSetForPolicy(extraAddBindDataContext)
	}

}
