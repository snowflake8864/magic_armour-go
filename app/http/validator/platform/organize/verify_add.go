package organize

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type VerifyAdd struct {
	OrgType   int    `form:"orgType" json:"orgType" binding:"required,oneof=1 2"`
	OraName   string `form:"oraName" json:"oraName" binding:"required,min=4,max=128"`
	UserName  string `form:"userName" json:"userName" binding:"required,alphanum,min=4,max=128"`
	UserPw    string `form:"userPw"   json:"userPw" binding:"required"`
	UserPhone string `form:"userPhone" json:"userPhone" binding:"required,number,len=11"`
	UserEmail string `form:"userEmail" json:"userEmail" binding:"required,email,max=64"`
	Remark    string `form:"remark" json:"remark" binding:"omitempty,max=200"`
}

func (n VerifyAdd) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixPlatform, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[platform organize VerifyAdd]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&platform.OrganizeController{}).AddForVerify(extraAddBindDataContext)
	}

}
