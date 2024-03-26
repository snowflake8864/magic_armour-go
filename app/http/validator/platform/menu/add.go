package menu

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type Add struct {
	ParentId *int   `form:"parentId" json:"parentId" binding:"required,min=0,max=10000"`
	Type     int    `form:"type" json:"type" binding:"required,max=10000"`
	Title    string `form:"name" json:"name" binding:"required,max=100"`
	WebCode  string `form:"code" json:"code" binding:"required,max=100"`
	Sort     *int   `form:"sort" json:"sort" binding:"required,min=0,max=9999"`
	IsShow   *int   `form:"isShow" json:"isShow" binding:"required,min=0,max=1"`
	IsRoot   *int   `form:"isRoot" json:"isRoot" binding:"required,min=0,max=1"`
	Remark   string `form:"remark" json:"remark" binding:"omitempty,max=100"`
}

func (n Add) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixPlatform, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[platform menu Add]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&platform.MenuController{}).Add(extraAddBindDataContext)
	}

}
