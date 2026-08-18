package jump

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

type PwJumpPut struct {
	Status int    `form:"status" json:"status"  binding:"required,oneof=1 2"`
	Reason string `form:"reason" json:"reason"`
	User   string `form:"user" json:"user"`
	Pw     string `form:"pw" json:"pw"`
}

func (n PwJumpPut) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	if n.Status == 1 {
		if string_helper.IsEmpty(n.User) {
			response.ValidatorError(context, fmt.Errorf("user not empyt"))
			return
		}
		if string_helper.IsEmpty(n.Pw) {
			response.ValidatorError(context, fmt.Errorf("pw not empyt"))
			return
		}
	}
	if n.Status == 2 {
		if string_helper.IsEmpty(n.Reason) {
			response.ValidatorError(context, fmt.Errorf("reason not empyt"))
			return
		}
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixAgent, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[agent pw_jump PwJumpPut]Form validator JSON conversion failed")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&agent.JumpController{}).PutPwJump(extraAddBindDataContext)
	}

}
