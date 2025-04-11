package tamper

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/http/validator/common/data_type"
	"goskeleton/app/http/validator/common/upload_files"
	"goskeleton/app/utils/response"
)

type PolicyImport struct {
	data_type.AgentSingle
	upload_files.UploadBase
}

func (n PolicyImport) CheckParams(context *gin.Context) {
	errCode, errMsg, data := n.CommCheck(context)
	if errCode != 0 {
		response.Fail(context, errCode, errMsg, data)
		return
	}

	(&organization.TamperController{}).ImportForPolicy(context)

}
