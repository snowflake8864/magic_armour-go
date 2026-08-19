package test

import (
	"fmt"
	"goskeleton/app/utils/helpers/file_helper"
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

func TestFileListDir(t *testing.T) {
	names, err := file_helper.ListDir("D:/go_project/magic_armour-go/storage/app/uploaded/2025_04/7c1b9dd19807937de5f6894cbeb2470f_unzip")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(json_helper.Encode(names))

}
