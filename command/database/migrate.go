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

// 开始执行业务
func start(action string) {
	switch action {
	case "init":
		initTables()
	case "append":
		variable.ZapLog.Sugar().Debug("action:" + action)
	}

	saveMigrations()
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

func TestInit() {
	start("init")
}

func initTables() {
	//初始化迁移表
	migrateFact := modelPlatform.CreateMigrateFactory()
	migrateFact.CreateTable()

	//获取迁移表数据
	getMigrationMap()

	//下面开始创建表 (platform库)
	if !existMigration("db_platform_create-p_admin") {
		_ = modelPlatform.CreateAdminFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_agent_type") {
		_ = modelPlatform.CreateAgentTypeFactory().CreateTable()
	}
	if !existMigration("db_platform_create-p_custom_table") {
		_ = modelPlatform.CreateCustomTableFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_department") {
		_ = modelPlatform.CreateDepartmentFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_dictionary") {
		_ = modelPlatform.CreateDictionaryFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_menu") {
		_ = modelPlatform.CreateMenuFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_operate_log") {
		_ = modelPlatform.CreateOperateLogFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_role") {
		_ = modelPlatform.CreateRoleFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_server_monitor") {
		_ = modelPlatform.CreateServerMonitorFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_sms_log") {
		_ = modelPlatform.CreateSmsLogFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_software") {
		_ = modelPlatform.CreateSoftwareFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_software_download") {
		_ = modelPlatform.CreateSoftwareDownloadFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-p_upload") {
		_ = modelPlatform.CreateUploadFactory().CreateTable()
	}

	//下面开始创建表 (agent库)
	if !existMigration("db_platform_create_table-a_agent") {
		_ = modelAgent.CreateAgentFactory().CreateTable()
	}

	if !existMigration("db_platform_create_table-a_dll") {
		_ = modelAgent.CreateDllFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-a_threaten_cache") {
		_ = modelAgent.CreateThreatenCacheFactory().CreateTable()
	}
	if !existMigration("db_platform_create_table-a_threaten_push") {
		_ = modelAgent.CreateThreatenPushFactory().CreateTable()
	}

	//下面开始创建表 (alarm库)
	if !existMigration("db_alarm_create_table-a_process") {
		_ = modelAlarm.CreateAlarmProcessFactory().CreateTable()
	}

	//下面开始创建表 (syslog库)
	if !existMigration("db_syslog_create_table-a_dns") {
		_ = modelSyslog.CreateSyslogDnsFactory().CreateTable()
	}
}
