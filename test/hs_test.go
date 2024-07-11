package test

import (
	"fmt"
	"goskeleton/app/service/veda_other/huanshi_platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console/huanshi_cmd"
	"log"
	"net"
	"testing"
	"time"
)

func TestHsAccess(t *testing.T) {
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

func TestHsPush(t *testing.T) {
	huanshi_cmd.CreateCmdFactory().Push()
}

func TestHsAnalog(t *testing.T) {
	serverAddr := "192.168.17.119:514"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                // Syslog 服务器地址
	logMsg := "<190>{\"deviceMark\":\"6e2f364c4c6ca7ab00d47c05928a7dbc\",\"platform\":\"huanjia\",\"dataType\":\"zhongduan\",\"org_id\":1,\"agentId\":\"neokylin_test_001\",\"os\":\"NeoKylin Desktop 5.0_U4 (Loongson)\",\"ip\":\"192.168.17.154,fe80::392d:7a87:b6d9:fed6\",\"netIp\":\"127.0.0.1\",\"hostName\":\"NeoKylin\",\"status\":0,\"onlineTime\":\"2024-07-04 11:26:41\",\"version\":\"2.0.0.29\",\"deptCode\":\"\",\"deptName\":\"\",\"agentType\":\"\",\"aliasName\":\"\",\"contactName\":\"\",\"contactPhone\":\"\"}\n" // 要发送的日志信息

	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 设置写超时时间为5秒
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write([]byte(logMsg))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sent syslog message: %s", logMsg)
}
