package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type UploadAntivirusEngine struct {
	EngineType   int    `form:"engine_type" json:"engine_type" binding:"required,max=10"`
	EngineVer    string `form:"engine_ver" json:"engine_ver" binding:"required,max=32"`
	EngineStatus int    `form:"engine_status" json:"engine_status" binding:"required,oneof=1 2 3 4"`
}

func (n UploadAntivirusEngine) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixAgent, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[agent agent UploadAntivirusEngine]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&agent.AgentController{}).UploadAntivirusEngine(extraAddBindDataContext)
	}

}
