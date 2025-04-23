package proc_sample

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/common/upload_files"
	"goskeleton/app/utils/response"
)

type Upload struct {
	upload_files.UploadBase
}

func (n Upload) CheckParams(context *gin.Context) {
	errCode, errMsg, data := n.CommCheck(context)
	if errCode != 0 {
		response.Fail(context, errCode, errMsg, data)
		return
	}

	(&agent.ProcSampleController{}).Upload(context)

}
