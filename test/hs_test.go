package test

import (
	"fmt"
	"goskeleton/app/service/veda_other/huanshi_platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
	"time"
)

func TestAccess(t *testing.T) {
	//获取备标识码
	deviceId := ""
	fmt.Println("deviceId:", deviceId)

	currentTime := time.Now().Unix()
	fmt.Println("currentTime:", currentTime)

	hsAccessFact := huanshi_platform.CreateCallbackFactory(deviceId)
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
