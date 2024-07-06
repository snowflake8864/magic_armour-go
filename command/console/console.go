package console

import (
	"fmt"
	"github.com/spf13/cobra"
	"goskeleton/app/global/variable"
	"goskeleton/command/console/agent_cmd"
	"goskeleton/command/console/alarm_cmd"
	"goskeleton/command/console/huanshi_cmd"
	"goskeleton/command/console/order_cmd"
	"goskeleton/command/console/org_cmd"
	"goskeleton/command/console/plug_cmd"
	"goskeleton/command/console/proc_cmd"
	"goskeleton/command/console/service_cmd"
	"goskeleton/command/console/syslog_cmd"
	"goskeleton/command/console/threaten_cmd"
)

const (
	UnitMinute string = "minute"
	UnitHour   string = "hour"
	UnitDay    string = "day"

	ModeEvery string = "every"
	ModeOnce  string = "once"

	AtZero    string = "zero"
	AtOne     string = "one"
	AtTwo     string = "two"
	AtThree   string = "three"
	AtFour    string = "four"
	AtFive    string = "five"
	AtSix     string = "six"
	AtSeven   string = "seven"
	AtEight   string = "eight"
	AtNine    string = "nine"
	AtTen     string = "ten"
	AtFifteen string = "fifteen"
	AtTwenty  string = "twenty"
	AtThirty  string = "thirty"
)

var (
	//1. 时间单位
	Unit string
	//2. 模式
	Mode string
	//3 间隔
	Interval string
	//4. 值
	At string
)

// 定时任务
var CrontabCmd = &cobra.Command{
	Use:     "crontab",
	Aliases: []string{"crontab"}, // 定义别名
	Short:   "定时任务",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。 
			2.执行 go  run  cmd/cli/main.go  crontab -h  //可以查看使用指南
			3.执行 go  run  cmd/cli/main.go  crontab -U 单位 -M 模式 -I 间隔 -A 具体值  // 指定全部参数
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		start(Unit, Mode, Interval, At)
	},
}

// 注册命令、初始化参数
func init() {
	CrontabCmd.Flags().StringVarP(&Unit, "Unit", "U", "", "-U 或者 --Unit 时间单位")
	CrontabCmd.Flags().StringVarP(&Mode, "Mode", "M", "every", "-M 或者 --Mode 模式")
	CrontabCmd.Flags().StringVarP(&Interval, "Interval", "I", "one", "-I 或者 --Interval 间隔")
	CrontabCmd.Flags().StringVarP(&At, "At", "A", "", "-A 或者 --At 值")
}

// 开始执行业务
func start(unit, mode, interval, at string) {
	variable.ZapLog.Sugar().Debug(fmt.Sprintf("unit:%s, mode:%s, interval:%s, at:%s", unit, mode, interval, at))

	if unit == UnitDay {
		handleDay(mode, interval, at)
	} else if unit == UnitHour {
		handleHour(mode, interval, at)
	} else if unit == UnitMinute {
		handleMinute(mode, interval, at)
	}

}

func handleDay(mode, interval, at string) {
	if mode == ModeOnce {
		switch at {
		case AtZero:
			//清理通讯日志
			syslog_cmd.CreateCmdFactory().Delete()
			//清理告警
			alarm_cmd.CreateCmdFactory().Delete()
			//清理任务
			agent_cmd.CreateTaskFactory().Delete()
			//清理终端监控
			agent_cmd.CreateMonitorFactory().Delete()
			//清理终端访问日志
			agent_cmd.CreateAccessFactory().Clear()
			//清理进程和模块的运行中数据
			proc_cmd.CreateCmdFactory().ClearRunning()
			break
		case AtOne:
			//未知支付订单自动过期
			order_cmd.CreateCmdFactory().UnpaidExpired()
			break
		case AtTwo:
			//终端授权自动过期
			agent_cmd.CreateAgentFactory().GrantExpired()
			break
		case AtThree:
			//统计企业昨天数据
			org_cmd.CreateCmdFactory().YesterdayStatistics()
			break
		case AtFive:
			//统计告警昨天数据
			alarm_cmd.CreateCmdFactory().YesterdayStatistics()
			break
		case AtSix:
			//删除威胁情报缓存
			threaten_cmd.CreateCmdFactory().DeleteForCache()
			break
		case AtSeven:
			//威胁情报缓存同步未知状态数据到威胁情报中心
			threaten_cmd.CreateCmdFactory().SyncUnknownForCache()
			break
		case AtEight:
			//告警通知 - 暂时不做
			alarm_cmd.CreateCmdFactory().EmailNotice()
			break
		}
	} else if mode == ModeEvery {
		switch at {
		case AtZero:
			//从进程添加到威胁情报缓存
			threaten_cmd.CreateCmdFactory().AddFromProcessForCache()
			break
		case AtOne:
			//从模块添加到威胁情报缓存
			threaten_cmd.CreateCmdFactory().AddFromDllForCache()
			break
		}
	}
}

func handleHour(mode, interval, at string) {

}

func handleMinute(mode, interval, at string) {
	if interval == AtOne {
		switch at {
		case AtZero:
			//威胁情报缓存同步默认状态（未同步）数据到威胁情报中心
			threaten_cmd.CreateCmdFactory().SyncDefaultForCache()
			break
		case AtOne:
			//终端离线
			agent_cmd.CreateAgentFactory().Offline()
			break
		case AtTwo:
			break
		}

	} else if interval == AtTwo {
		switch at {
		case AtZero:
			//动态封堵已过期的自动解封
			plug_cmd.CreateCmdFactory().DynamicExpired()
			break
		case AtOne:
			//终端安全状态
			agent_cmd.CreateAgentFactory().HandleSafe()
			break
		}

	} else if interval == AtThree {
		switch at {
		case AtZero:
			//从进程添加任务
			proc_cmd.CreateCmdFactory().AddFromProcessForTask()
			break
		case AtOne:
			//从模块添加任务
			proc_cmd.CreateCmdFactory().AddFromDllForTask()
			break
		case AtTwo:
			//服务器状态监控
			service_cmd.CreateCmdFactory().Monitor()
			break
		case AtThree:
			//威胁情报同步数据到威胁情报中心
			threaten_cmd.CreateCmdFactory().SyncForPushes()
			break
		}

	} else if interval == AtFive {
		switch at {
		case AtZero:
			//同步推送数据到幻势
			huanshi_cmd.CreateCmdFactory().Push()
			break
		}
	} else if interval == AtTen {
		switch at {
		case AtZero:
			//统计告警今天数据
			alarm_cmd.CreateCmdFactory().TodayStatistics()
			break
		}
	} else if interval == AtFifteen {
		switch at {
		case AtZero:
			//统计企业今天数据
			org_cmd.CreateCmdFactory().TodayStatistics()
			break
		}
	}
}
