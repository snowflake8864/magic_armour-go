package queue

import (
	"github.com/spf13/cobra"
	"goskeleton/app/utils/kakfa_v2"
)

var (
	Action string
)

//数据迁移

var TopicInit = &cobra.Command{
	Use:     "topic",
	Aliases: []string{"topic"}, // 定义别名
	Short:   "队列主题",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。 
			2.执行 go  run  cmd/cli/main.go  topic -h  //可以查看使用指南
			3.执行 go  run  cmd/cli/main.go  topic init // 通过 Action 动作执行相应的命令
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		start(Action)
	},
}

// 注册命令、初始化参数
func init() {
	TopicInit.Flags().StringVarP(&Action, "Action", "A", "init", "-A 指定参数,例如：-A init")
}

// 开始执行业务
func start(action string) {
	switch action {
	case "init":
		kakfa_v2.CreateTopicFactory()
	}
}
