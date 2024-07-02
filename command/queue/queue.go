package queue

import (
	"fmt"
	"github.com/spf13/cobra"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/kafka_go/consumer_kafka"
	"goskeleton/app/utils/kafka_go/topic_kakfa"
)

var (
	//1. 定义行为
	Action string
	//2. 主题名称
	Topic string
)

// 队列
var QueueInfo = &cobra.Command{
	Use:     "queue",
	Aliases: []string{"queue"}, // 定义别名
	Short:   "队列相关",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。 
			2.执行 go  run  cmd/cli/main.go  queue -h  //可以查看使用指南
			3.执行 go  run  cmd/cli/main.go  queue init // 通过 Action 动作执行相应的命令
			4.执行 go  run  cmd/cli/main.go  queue consume common // 指定参数运行queue
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		start(Action, Topic)
	},
}

// 注册命令、初始化参数
func init() {
	QueueInfo.Flags().StringVarP(&Action, "Action", "A", "init", "-A 或者 --Action 选择行为，例如:init consume")
	QueueInfo.Flags().StringVarP(&Topic, "Topic", "T", "common", "-T 或者 --Topic 指定主题名称 common process")
}

// 开始执行业务
func start(action, topic string) {
	variable.ZapLog.Sugar().Debug(fmt.Sprintf("action:%s, topic:%s", action, topic))
	switch action {
	case "init":
		topic_kakfa.CreateTopicFactory().Init()
		break
	case "consume":
		switch topic {
		case "common":
			consumer_kafka.CreateConsumerFactory().ConsumeCommon()
			break
		case "process":
			consumer_kafka.CreateConsumerFactory().ConsumeProcess()
			break
		case "threaten":
			consumer_kafka.CreateConsumerFactory().ConsumeThreaten()
			break
		case "alarm":
			consumer_kafka.CreateConsumerFactory().ConsumeAlarm()
			break
		case "syslog":
			consumer_kafka.CreateConsumerFactory().ConsumeSyslog()
			break
		}
		break
	}

}
