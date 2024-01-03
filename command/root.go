package cmd

import (
	"fmt"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {

}

func init() {
	fmt.Println("1----------------------------root")
	// 如果子命令是存在于子目录，那么就需要在入口统一添加；
	// 如果和 root.go 同目录，则不需要像下一行一样添加
	//RootCmd.AddCommand(demo.Demo1)
	//RootCmd.AddCommand(demo_simple.DemoSimple)

}
