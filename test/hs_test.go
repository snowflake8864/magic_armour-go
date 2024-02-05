package test

import (
	"fmt"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_service"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
	"time"
)

func TestAccess(t *testing.T) {
	//获取备标识码
	dictCacheFact := dictionary_service.CreateCacheFactory()
	deviceId := dictCacheFact.GetValById(enum_dictionary.SystemDeviceId)
	fmt.Println("deviceId:", deviceId)

	currentTime := time.Now().Unix()
	fmt.Println("currentTime:", currentTime)

	hsAccessFact := huanshi_service.CreateAccessFactory(deviceId)
	secret := hsAccessFact.GenerateSecret(currentTime)
	fmt.Println("secret:", secret)

	if isOk, errMsg := hsAccessFact.CheckSecret(currentTime, secret); !isOk {
		fmt.Println("CheckSecret:", errMsg)
	}

	accessToken, err := hsAccessFact.GenerateToken(secret)
	fmt.Println(accessToken)
	if err != nil {
		fmt.Println("GenerateToken err:", err.Error())
	}
	fmt.Println("token:", accessToken)
}
