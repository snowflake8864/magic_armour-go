package test

import (
	"fmt"
	"goskeleton/app/global/enums/enum_pw_jump"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/helpers/json_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"strings"
	"testing"
)

func TestIsIPInRange(t *testing.T) {
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.6.6"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.6.7"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.6.7/24"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.28.131/16"))
	fmt.Println("ok:", ip_helper.IsIPInRange("192.168.6.6", "192.168.28.131/24"))
}

func TestPwJumpConfig(t *testing.T) {
	config := enum_pw_jump.PwJumpConfig{}

	config.Combine = 1
	config.Size = 12
	fmt.Println("pw1: ", config.GeneratePw())

	config.Combine = 2
	config.Size = 12
	fmt.Println("pw2: ", config.GeneratePw())

	config.Combine = 3
	config.Size = 14
	fmt.Println("pw3: ", config.GeneratePw())

	config.Combine = 4
	config.Size = 14
	fmt.Println("pw4: ", config.GeneratePw())

	config.Combine = 5
	config.Size = 15
	fmt.Println("pw5: ", config.GeneratePw())

	config.Combine = 6
	config.Size = 15
	fmt.Println("pw6: ", config.GeneratePw())

	config.Combine = 7
	config.Size = 20
	fmt.Println("pw7: ", config.GeneratePw())

	config.Combine = 8
	config.Size = 20
	fmt.Println("pw8: ", config.GeneratePw())

	config.Combine = 9
	config.Size = 20
	fmt.Println("pw9: ", config.GeneratePw())

	config.Combine = 10
	config.Size = 20
	fmt.Println("pw10: ", config.GeneratePw())

	config.Combine = 11
	config.Size = 20
	fmt.Println("pw11: ", config.GeneratePw())

	config.Combine = 12
	config.Size = 8
	fmt.Println("pw12: ", config.GeneratePw())

	config.Combine = 13
	config.Size = 8
	fmt.Println("pw13: ", config.GeneratePw())

	config.Combine = 14
	config.Size = 8
	fmt.Println("pw14: ", config.GeneratePw())

	config.Combine = 15
	config.Size = 8
	fmt.Println("pw15: ", config.GeneratePw())

}

func TestIPv4Range(t *testing.T) {
	testCases := []string{
		"192.168.1.5/24",
		"192.168.1.5/16",
		"10.0.0.1/255.255.255.0",
		"172.16.0.5 255.255.0.0",
		"203.0.113.17/31",
		"192.0.2.3/32",
	}
	for _, tc := range testCases {
		fmt.Println("Testing:", tc)
		iPv4Range := ip_helper.IPv4Range{IpWithMask: tc}
		err, networkIp, broadcastIp, firstUsableIp, lastUsableIp, firstUsableNumber, lastUsableNumber := iPv4Range.CalcRange()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Printf("NetworkIP: %-15s Broadcast: %-15s FirstIP: %-15s LastIP: %s FirstNumber: %-15d LastNumber: %-15d\n\n",
			networkIp, broadcastIp, firstUsableIp, lastUsableIp, firstUsableNumber, lastUsableNumber)
	}
}

func TestStringSpilt(t *testing.T) {
	fmt.Println(json_helper.Encode(strings.Split("1,24,4,", ",")))
	fmt.Println(json_helper.Encode(strings.Split("1", ",")))
	fmt.Println(json_helper.Encode(strings.Split("", ",")))
	fmt.Println(json_helper.Encode(strings.Split(",2,", ",")))
}
