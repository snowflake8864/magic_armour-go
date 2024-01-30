package database

import (
	"github.com/spf13/cobra"
	"goskeleton/app/global/variable"
	modelAgent "goskeleton/app/model/agent"
	modelAlarm "goskeleton/app/model/alarm"
	modelPlatform "goskeleton/app/model/platform"
	modelSyslog "goskeleton/app/model/syslog"
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
		_ = modelPlatform.CreateMenuFactory().InitData()
	case "dictionary":
		_ = modelPlatform.CreateDictionaryFactory().InitData()
	case "append":
		variable.ZapLog.Sugar().Debug("action:" + action)
	}

}

// 迁移数据
var migrateMap map[string]string

// 不在迁移表中
var notExistMigrateList []modelPlatform.MigrateModel

// 获取迁移数据
func getMigrationMap() {
	migrateFact := modelPlatform.CreateMigrateFactory()
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
	migrateFact := modelPlatform.CreateMigrateFactory()
	_ = migrateFact.InsertBatch(&notExistMigrateList)
}

func existMigration(migration string) bool {
	_, ok := migrateMap[migration]
	if ok {
		return true
	}
	var st = modelPlatform.MigrateModel{
		Migration: migration,
		CreatedAt: time.Now().Format(variable.DateFormat),
	}
	notExistMigrateList = append(notExistMigrateList, st)
	return false
}

func initTableAndData() {
	//初始化迁移表
	migrateFact := modelPlatform.CreateMigrateFactory()
	migrateFact.CreateTable()

	//获取迁移表数据
	getMigrationMap()

	//下面开始创建表 (platform库)
	if !existMigration("db_platform_init_tables") {
		adminFact := modelPlatform.CreateAdminFactory()
		_ = adminFact.CreateTable()
		_ = adminFact.InitData()

		_ = modelPlatform.CreateAgentTypeFactory().CreateTable()
		_ = modelPlatform.CreateCustomTableFactory().CreateTable()
		_ = modelPlatform.CreateDepartmentFactory().CreateTable()
		_ = modelPlatform.CreateDictionaryFactory().CreateTable()
		_ = modelPlatform.CreateGrantFactory().CreateTable()
		_ = modelPlatform.CreateGrantModuleFactory().CreateTable()
		_ = modelPlatform.CreateGrantUsedFactory().CreateTable()
		_ = modelPlatform.CreateMenuFactory().CreateTable()
		_ = modelPlatform.CreateOperateLogFactory().CreateTable()

		roleFact := modelPlatform.CreateRoleFactory()
		_ = roleFact.CreateTable()
		_ = roleFact.InitData()

		_ = modelPlatform.CreateServerMonitorFactory().CreateTable()
		_ = modelPlatform.CreateSmsLogFactory().CreateTable()
		_ = modelPlatform.CreateSoftwareFactory().CreateTable()
		_ = modelPlatform.CreateSoftwareDownloadFactory().CreateTable()
		_ = modelPlatform.CreateUploadFactory().CreateTable()
	}

	//下面开始创建表 (agent库)
	if !existMigration("db_agent_init_tables") {
		_ = modelAgent.CreateAccessLogFactory().CreateTable()
		_ = modelAgent.CreateAgentFactory().CreateTable()
		_ = modelAgent.CreateBlackmailGroupFactory().CreateTable()
		_ = modelAgent.CreateBlackmailPolicyFactory().CreateTable()
		_ = modelAgent.CreateDllFactory().CreateTable()
		_ = modelAgent.CreateDllRunningFactory().CreateTable()
		_ = modelAgent.CreateErrorFactory().CreateTable()
		_ = modelAgent.CreateMonitorFactory().CreateTable()

		outreachGlobalFact := modelAgent.CreateOutreachGlobalFactory()
		_ = outreachGlobalFact.CreateTable()
		_ = outreachGlobalFact.InitData()

		_ = modelAgent.CreatePlugDynamicFactory().CreateTable()
		_ = modelAgent.CreatePlugIpGlobalFactory().CreateTable()
		_ = modelAgent.CreatePlugIpPolicyFactory().CreateTable()
		_ = modelAgent.CreatePortVirtualGlobalFactory().CreateTable()
		_ = modelAgent.CreatePortVirtualPolicyFactory().CreateTable()
		_ = modelAgent.CreatePortWorkingFactory().CreateTable()
		_ = modelAgent.CreateProcGlobalFactory().CreateTable()
		_ = modelAgent.CreateProcSampleFactory().CreateTable()
		_ = modelAgent.CreateProcessRunningFactory().CreateTable()
		_ = modelAgent.CreateProcessFactory().CreateTable()
		_ = modelAgent.CreateTamperGlobalFactory().CreateTable()
		_ = modelAgent.CreateTamperPolicyFactory().CreateTable()
		_ = modelAgent.CreateTaskFactory().CreateTable()
		_ = modelAgent.CreateThreatenCacheFactory().CreateTable()
		_ = modelAgent.CreateThreatenPushFactory().CreateTable()
		_ = modelAgent.CreateTrustGlobalFactory().CreateTable()
		_ = modelAgent.CreateTrustPolicyFactory().CreateTable()
		_ = modelAgent.CreateUsbFactory().CreateTable()
		_ = modelAgent.CreateUsbGlobalFactory().CreateTable()
		_ = modelAgent.CreateUsbLogFactory().CreateTable()
	}

	//下面开始创建表 (alarm库)
	if !existMigration("db_alarm_init_tables") {
		_ = modelAlarm.CreateAlarmProcessFactory().CreateTable()
	}

	//下面开始创建表 (syslog库)
	if !existMigration("db_syslog_init_tables") {
		_ = modelSyslog.CreateSyslogDnsFactory().CreateTable()
	}

	//保存迁移数据
	saveMigrations()
}
