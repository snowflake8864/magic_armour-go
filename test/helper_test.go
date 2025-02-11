package test

import (
	"fmt"
	"goskeleton/app/utils/helpers/ip_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestIsIPInRange(t *testing.T) {
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.6.6"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.6.7"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.6.7/24"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.28.131/16"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.28.131/24"))
}
