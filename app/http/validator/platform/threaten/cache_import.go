package threaten

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/http/validator/common/upload_files"
	"goskeleton/app/utils/response"
)

type CacheImport struct {
	upload_files.UploadBase
}

func (n CacheImport) CheckParams(context *gin.Context) {
	errCode, errMsg, data := n.CommCheck(context)
	if errCode != 0 {
		response.Fail(context, errCode, errMsg, data)
		return
	}

	(&platform.ThreatenController{}).ImportForCache(context)

}
