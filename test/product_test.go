package test

import (
	"goskeleton/app/model/model_platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestProductInit(t *testing.T) {
	productFact := model_platform.CreateProductFactory()
	_ = productFact.InitData()
}
