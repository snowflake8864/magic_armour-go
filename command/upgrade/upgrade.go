package upgrade

import (
	"fmt"
	"github.com/spf13/cobra"
	"goskeleton/app/global/variable"
)

var (
	//1. 定义行为
	Action string
)

// 队列
var UpgradeLog = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"upgrade"}, // 定义别名
	Short:   "系统升级相关",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。 
			2.执行 go  run  cmd/cli/main.go  upgrade -h  //可以查看使用指南
			3.执行 go  run  cmd/cli/main.go  queue -A 行为 // 指定参数运行upgrade
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		start(Action)
	},
}

// 注册命令、初始化参数
func init() {
	UpgradeLog.Flags().StringVarP(&Action, "Action", "A", "", "-A 或者 --Action 选择行为")
}

// 开始执行业务
func start(action string) {
	variable.ZapLog.Sugar().Debug(fmt.Sprintf("action:%s", action))
	switch action {
	case "log":
		fmt.Println("aa")
	}

}
