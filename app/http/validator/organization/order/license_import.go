package order

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/http/validator/common/upload_files"
	"goskeleton/app/utils/response"
)

type LicenseImport struct {
	upload_files.UploadBase
}

func (n LicenseImport) CheckParams(context *gin.Context) {
	errCode, errMsg, data := n.CommCheck(context)
	if errCode != 0 {
		response.Fail(context, errCode, errMsg, data)
		return
	}

	(&organization.OrderController{}).ImportForLicense(context)

}
