package database

import (
	"github.com/spf13/cobra"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/model/model_alarm"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/model/model_syslog"
	"time"
)

var (
	Action string
)

//数据迁移

var DbMigrate = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"migrate"}, // 定义别名
	Short:   "数据库迁移",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。 
			2.执行 go  run  cmd/cli/main.go  migrate -h  //可以查看使用指南
			3.执行 go  run  cmd/cli/main.go  migrate  -A init  // 通过 Action 动作执行相应的命令
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		start(Action)
	},
}

// 注册命令、初始化参数
func init() {
	DbMigrate.Flags().StringVarP(&Action, "Action", "A", "init", "-A 指定参数,例如：-A init")
}

func TestStart(action string) {
	start(action)
}

// 开始执行业务
func start(action string) {
	variable.ZapLog.Sugar().Debug("action:" + action)
	switch action {
	case "init":
		initTableAndData()
	case "menu":
		_ = model_platform.CreateMenuFactory().InitData()
	case "dictionary":
		_ = model_platform.CreateDictionaryFactory().InitData()
	case "append":
		variable.ZapLog.Sugar().Debug("action:" + action)
	}

}

// 迁移数据
var migrateMap map[string]string

// 不在迁移表中
var notExistMigrateList []model_platform.MigrateModel

// 获取迁移数据
func getMigrationMap() {
	migrateFact := model_platform.CreateMigrateFactory()
	migrateMap, _ = migrateFact.GetAllMaps()
	if migrateMap == nil {
		variable.ZapLog.Sugar().Debug("migrate GetAllMaps fail")
	}
}

// 保存本次迁移数据
func saveMigrations() {
	if len(notExistMigrateList) == 0 {
		variable.ZapLog.Sugar().Debug("本次迁移数据为空")
		return
	}
	migrateFact := model_platform.CreateMigrateFactory()
	_ = migrateFact.InsertBatch(&notExistMigrateList)
}

func existMigration(migration string) bool {
	_, ok := migrateMap[migration]
	if ok {
		return true
	}
	var st = model_platform.MigrateModel{
		Migration: migration,
		CreatedAt: time.Now().Format(variable.DateFormat),
	}
	notExistMigrateList = append(notExistMigrateList, st)
	return false
}

func initTableAndData() {
	//初始化迁移表
	migrateFact := model_platform.CreateMigrateFactory()
	migrateFact.CreateTable()

	//获取迁移表数据
	getMigrationMap()

	//下面开始创建表 (platform库)
	if !existMigration("db_platform_init_tables") {
		adminFact := model_platform.CreateAdminFactory()
		_ = adminFact.CreateTable()
		_ = adminFact.InitData()

		_ = model_platform.CreateAgentTypeFactory().CreateTable()
		_ = model_platform.CreateCustomTableFactory().CreateTable()
		_ = model_platform.CreateDepartmentFactory().CreateTable()
		_ = model_platform.CreateDictionaryFactory().CreateTable()
		_ = model_platform.CreateGrantFactory().CreateTable()
		_ = model_platform.CreateGrantModuleFactory().CreateTable()
		_ = model_platform.CreateGrantUsedFactory().CreateTable()
		_ = model_platform.CreateMenuFactory().CreateTable()
		_ = model_platform.CreateOperateLogFactory().CreateTable()

		roleFact := model_platform.CreateRoleFactory()
		_ = roleFact.CreateTable()
		_ = roleFact.InitData()

		_ = model_platform.CreateServerMonitorFactory().CreateTable()
		_ = model_platform.CreateSmsLogFactory().CreateTable()
		_ = model_platform.CreateSoftwareFactory().CreateTable()
		_ = model_platform.CreateSoftwareDownloadFactory().CreateTable()
		_ = model_platform.CreateUploadFactory().CreateTable()
	}

	//下面开始创建表 (agent库)
	if !existMigration("db_agent_init_tables") {
		_ = model_agent.CreateAccessLogFactory().CreateTable()
		_ = model_agent.CreateAgentFactory().CreateTable()
		_ = model_agent.CreateBlackmailGroupFactory().CreateTable()
		_ = model_agent.CreateBlackmailPolicyFactory().CreateTable()
		_ = model_agent.CreateDllFactory().CreateTable()
		_ = model_agent.CreateDllRunningFactory().CreateTable()
		_ = model_agent.CreateErrorFactory().CreateTable()
		_ = model_agent.CreateMonitorFactory().CreateTable()

		outreachGlobalFact := model_agent.CreateOutreachGlobalFactory()
		_ = outreachGlobalFact.CreateTable()
		_ = outreachGlobalFact.InitData()

		_ = model_agent.CreatePlugDynamicFactory().CreateTable()
		_ = model_agent.CreatePlugIpGlobalFactory().CreateTable()
		_ = model_agent.CreatePlugIpPolicyFactory().CreateTable()
		_ = model_agent.CreatePortVirtualGlobalFactory().CreateTable()
		_ = model_agent.CreatePortVirtualPolicyFactory().CreateTable()
		_ = model_agent.CreatePortWorkingFactory().CreateTable()
		_ = model_agent.CreateProcGlobalFactory().CreateTable()
		_ = model_agent.CreateProcSampleFactory().CreateTable()
		_ = model_agent.CreateProcessRunningFactory().CreateTable()
		_ = model_agent.CreateProcessFactory().CreateTable()
		_ = model_agent.CreateTamperGlobalFactory().CreateTable()
		_ = model_agent.CreateTamperPolicyFactory().CreateTable()
		_ = model_agent.CreateTaskFactory().CreateTable()
		_ = model_agent.CreateThreatenCacheFactory().CreateTable()
		_ = model_agent.CreateThreatenPushFactory().CreateTable()
		_ = model_agent.CreateTrustGlobalFactory().CreateTable()
		_ = model_agent.CreateTrustPolicyFactory().CreateTable()
		_ = model_agent.CreateUsbFactory().CreateTable()
		_ = model_agent.CreateUsbGlobalFactory().CreateTable()
		_ = model_agent.CreateUsbLogFactory().CreateTable()
	}

	//下面开始创建表 (alarm库)
	if !existMigration("db_alarm_init_tables") {
		_ = model_alarm.CreateAlarmProcessFactory().CreateTable()
	}

	//下面开始创建表 (syslog库)
	if !existMigration("db_syslog_init_tables") {
		_ = model_syslog.CreateSyslogDnsFactory().CreateTable()
	}

	//保存迁移数据
	saveMigrations()
}
