package test

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/model/platform/prodcut_models"
	"goskeleton/app/service/organization/order_service"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestProductInit(t *testing.T) {
	productFact := prodcut_models.CreateProductFactory()
	_ = productFact.InitData()
}

func TestProductGrantConfig(t *testing.T) {
	fact := enum_product.ProductGrantConfig{
		Process:   100,
		Tamper:    90,
		Blackmail: 80,
		Port:      70,
		Plug:      60,
		Usb:       50,
		Syslog:    100,
		Outreach:  100,
		Login:     100,
		Threaten:  0,
	}
	for true {
		minimum := fact.MinGreaterThanZero()
		if minimum == 0 {
			break
		}
		fmt.Println("min: ", minimum)
		fact.Sub(minimum)
		fmt.Println(fact)
	}

}

func TestImportLicense(t *testing.T) {

	param := order_models.OrderModel{}
	param.OrgId = 1
	param.Source = enums.OperateSourcePlatform

	var uploadId int64
	uploadId = 38
	uploadModel := upload_models.CreateUploadFactory().GetItemById(uploadId)
	if uploadModel == nil || uploadModel.Id == 0 {
		fmt.Println("上传证书不存在")
		return
	}
	param.PayOrderNo = uploadModel.FileHash

	if errCode, err := order_service.CreateCurdFactory().ImportLicense(&param, uploadModel.FullPath, 1, "root"); err != nil {
		fmt.Println(errCode)
		fmt.Println(err.Error())
		return
	}
}

func TestGetGrantTypeByOrderNo(t *testing.T) {
	fmt.Println(enum_order.GetGrantTypeByOrderNo("Lhj924121100061"))
	fmt.Println(enum_order.GetGrantTypeByOrderNo("Lhj924121100061"))
}
