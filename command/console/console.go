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
	"goskeleton/command/console/sntp_cmd"
	"goskeleton/command/console/syslog_cmd"
	"goskeleton/command/console/threaten_cmd"
)

var (
	Action string
)

// 定时任务
var CrontabCmd = &cobra.Command{
	Use:     "crontab",
	Aliases: []string{"crontab"}, // 定义别名
	Short:   "定时任务",
	Long: `调用方法：
	1.进入项目根目录（Ginkeleton）。 
	2.执行 go  run  cmd/cli/main.go  crontab -h  //可以查看使用指南
	3.执行 go  run  cmd/cli/main.go  crontab -A  // 指定全部参数
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		start(Action)
	},
}

// 注册命令、初始化参数
func init() {
	CrontabCmd.Flags().StringVarP(&Action, "Action", "A", "", "-A 或者 --Action 值")
}

// 开始执行业务
func start(action string) {
	variable.ZapLog.Info(fmt.Sprintf("action=%s start", action))

	switch action {
	//以下建议每1分钟执行一次
	case "AgentOffline": //终端离线
		agent_cmd.CreateAgentFactory().Offline()
	case "ThreatenCacheDefaultSync": //威胁情报缓存同步默认状态（未同步）数据到威胁情报中心
		threaten_cmd.CreateCmdFactory().SyncDefaultForCache()
	case "ThreatenCache2Pushes": //本地缓存数据同步威胁情报推送列表 - 暂时不做
		threaten_cmd.CreateCmdFactory().Cache2Pushes()

	//以下建议每2分钟执行一次
	case "PlugExpired": //动态封堵已过期的自动解封
		plug_cmd.CreateCmdFactory().DynamicExpired()
	case "AgentSafe": //终端安全状态
		agent_cmd.CreateAgentFactory().HandleSafe()

	//以下建议每3分钟执行一次
	case "TaskAddFromProc": //从进程或模块添加任务
		proc_cmd.CreateCmdFactory().AddFromProcForTask()
	case "ServiceMonitor": //服务器状态监控
		service_cmd.CreateCmdFactory().Monitor()
	case "ThreatenPushesSync": //威胁情报推送列表同步数据到威胁情报中心
		threaten_cmd.CreateCmdFactory().SyncForPushes()

	//以下建议每5分钟执行一次
	case "HuanshiPush": //同步推送数据到幻势
		huanshi_cmd.CreateCmdFactory().Push()

	//以下建议每12分钟执行一次
	case "CollectOrgToday": //统计企业今天相关数据
		//统计告警今天数据
		alarm_cmd.CreateCmdFactory().TodayStatistics()
		//统计企业今天数据
		org_cmd.CreateCmdFactory().TodayStatistics()

	//以下建议每天时执行一次
	case "Daily": //每日定时清理数据
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
		//清理服务器监控
		service_cmd.CreateCmdFactory().Delete()
		//SNTP同步时间
		sntp_cmd.CreateCmdFactory().Sync()
		//清理终端文件目录
		agent_cmd.CreateDirFactory().Clear()

	case "OrderExpired": //未知支付订单自动过期
		order_cmd.CreateCmdFactory().OrderExpired()
	case "GrantExpired": //终端授权自动过期
		agent_cmd.CreateAgentFactory().GrantExpired()
	case "CollectOrgYesterday": //统计企业昨天相关数据
		//统计企业昨天数据
		org_cmd.CreateCmdFactory().YesterdayStatistics()
		//统计告警昨天数据
		alarm_cmd.CreateCmdFactory().YesterdayStatistics()
	case "AlarmEmailNotice": //告警邮件通知
		alarm_cmd.CreateCmdFactory().EmailNotice()
	case "ThreatenCacheDelete": //删除威胁情报缓存
		threaten_cmd.CreateCmdFactory().DeleteForCache()
	case "ThreatenCacheUnknownSync": //威胁情报缓存同步未知状态数据到威胁情报中心
		threaten_cmd.CreateCmdFactory().SyncUnknownForCache()
	case "ThreatenCacheAddFromProc": //从进程或模块添加到威胁情报缓存
		threaten_cmd.CreateCmdFactory().AddFromProcForCache()

	//以下建议特殊情况下执行
	case "ThreatenRepairProc": //威胁情报缓存修复进程模块状态
		threaten_cmd.CreateCmdFactory().RepairProcForCache()
	case "RepairProcByThreaten": //通过威胁情报修复进程状态
		proc_cmd.CreateCmdFactory().RepairStatusByThreaten()

	default:
		variable.ZapLog.Info(fmt.Sprintf("action=%s fail", action))
	}
	variable.ZapLog.Info(fmt.Sprintf("action=%s end", action))
}

// 开始执行业务
func TestStart(at string) {
	start(at)
}
