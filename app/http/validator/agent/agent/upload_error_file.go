package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/common/upload_files"
	"goskeleton/app/utils/response"
)

type UploadErrorFile struct {
	upload_files.UploadBase
}

func (n UploadErrorFile) CheckParams(context *gin.Context) {
	errCode, errMsg, data := n.CommCheck(context)
	if errCode != 0 {
		response.Fail(context, errCode, errMsg, data)
		return
	}

	(&agent.AgentController{}).UploadErrorFile(context)
}
