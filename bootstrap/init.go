package bootstrap

import (
	"fmt"
	_ "goskeleton/app/core/destroy" // 监听程序退出信号，用于资源的释放
)

// 检查项目必须的非编译目录是否存在，避免编译后调用的时候缺失相关目录
func checkRequiredFolders() {
}

func init() {
	// 1. 初始化 项目根路径，参见 variable 常量包，相关路径：app\global\variable\variable.go
	fmt.Println("init")
	//2.检查配置文件以及日志目录等非编译性的必要条件
	checkRequiredFolders()
}
