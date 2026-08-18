package backup

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type PutBackup struct {
	Id         string `form:"id" json:"id" binding:"required,max=64"`
	State      *int   `form:"state" json:"state" binding:"required,oneof=0 1"`
	StorageDir string `form:"storage_dir" json:"storage_dir" binding:"omitempty,max=512"`
	FileName   string `form:"file_name" json:"file_name" binding:"omitempty,max=128"`
	FileSize   string `form:"file_size" json:"file_size" binding:"omitempty,max=32"`
	FailReason string `form:"fail_reason" json:"fail_reason" binding:"omitempty,max=512"`
	Param      string `form:"param" json:"param"`
}

func (n PutBackup) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixAgent, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[agent backup PutBackup]Form validator JSON conversion failed")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&agent.BackupController{}).UploadBackup(extraAddBindDataContext)
	}

}
