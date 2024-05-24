package proc_sample

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/common/upload_files"
)

type Upload struct {
	upload_files.UploadBase
}

func (n Upload) CheckParams(context *gin.Context) {
	errCode, errMsg, data := n.CommCheck(context)
	if errCode != 0 {
		variable.ZapLog.Sugar().Debug(errCode, errMsg, data)
	}

	(&agent.ProcSampleController{}).Upload(context)

}
