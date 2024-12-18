package system

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type HuanshiSet struct {
	HuanshiCallbackUrl string `form:"huanshi_callback_url" json:"huanshi_callback_url" binding:"omitempty,max=128"`
	HuanshiCallIp      string `form:"huanshi_callback_ip" json:"huanshi_callback_ip" binding:"omitempty,max=128"`
	HuanshiSwitch      *int   `form:"huanshi_switch" json:"huanshi_switch" binding:"required,oneof=0 1"`
	HuanshiUrl         string `form:"huanshi_url" json:"huanshi_url" binding:"omitempty,max=128"`
	HuanshiAlarmSwitch *int   `form:"huanshi_alarm_switch" json:"huanshi_alarm_switch" binding:"required,oneof=0 1"`
	HuanshiAlarmIp     string `form:"huanshi_alarm_ip" json:"huanshi_alarm_ip" binding:"omitempty,max=128"`
	HuanshiAlarmPort   int    `form:"huanshi_alarm_port" json:"huanshi_alarm_port" binding:"omitempty,max=65535"`
	HuanshiLogSwitch   *int   `form:"huanshi_log_switch" json:"huanshi_log_switch" binding:"required,oneof=0 1"`
	HuanshiLogIp       string `form:"huanshi_log_ip" json:"huanshi_log_ip" binding:"omitempty,max=128"`
	HuanshiLogPort     int    `form:"huanshi_log_port" json:"huanshi_log_port" binding:"omitempty,max=65535"`
	HuanshiThirdSwitch *int   `form:"huanshi_third_switch" json:"huanshi_third_switch" binding:"required,oneof=0 1"`
	HuanshiThirdIp     string `form:"huanshi_third_ip" json:"huanshi_third_ip" binding:"omitempty,max=128"`
	HuanshiThirdPort   int    `form:"huanshi_third_port" json:"huanshi_third_port" binding:"omitempty,max=65535"`
}

func (n HuanshiSet) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixPlatform, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[platform system HuanshiSet]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&platform.SystemController{}).HuanshiSet(extraAddBindDataContext)
	}

}
