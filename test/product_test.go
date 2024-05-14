package test

import (
	"goskeleton/app/model/platform/prodcut_models"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestProductInit(t *testing.T) {
	productFact := prodcut_models.CreateProductFactory()
	_ = productFact.InitData()
}
