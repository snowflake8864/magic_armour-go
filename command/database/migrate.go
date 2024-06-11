package database

import (
	"github.com/spf13/cobra"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_access_models"
	"goskeleton/app/model/agent/agent_error_models"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/model/agent/dir_models"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/model/agent/trust_dir_models"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/model/common/sms_models"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/model/organization/agent_type_models"
	"goskeleton/app/model/organization/custom_table_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/model/organization/org_operate_models"
	"goskeleton/app/model/organization/org_stat_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/role_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/model/platform/admin_models"
	"goskeleton/app/model/platform/dictionary_models"
	"goskeleton/app/model/platform/migrate_models"
	"goskeleton/app/model/platform/operate_models"
	"goskeleton/app/model/platform/org_menu_models"
	"goskeleton/app/model/platform/prodcut_models"
	"goskeleton/app/model/platform/server_monitor_models"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/platform/dictionary_service"
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
		_ = org_menu_models.CreateOrgMenuFactory().InitData()
	case "dictionary":
		_ = dictionary_models.CreateDictionaryFactory().InitData()
		_ = dictionary_service.CreateCacheFactory().Rebuild()
	case "append":
		variable.ZapLog.Sugar().Debug("action:" + action)
	}

}

// 迁移数据
var migrateMap map[string]string

// 不在迁移表中
var notExistMigrateList []migrate_models.MigrateModel

// 获取迁移数据
func getMigrationMap() {
	migrateFact := migrate_models.CreateMigrateFactory()
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
	migrateFact := migrate_models.CreateMigrateFactory()
	_ = migrateFact.InsertBatch(&notExistMigrateList)
}

func existMigration(migration string) bool {
	_, ok := migrateMap[migration]
	if ok {
		return true
	}
	var st = migrate_models.MigrateModel{
		Migration: migration,
		CreatedAt: time.Now().Format(variable.DateFormat),
	}
	notExistMigrateList = append(notExistMigrateList, st)
	return false
}

func initTableAndData() {
	//初始化迁移表
	migrateFact := migrate_models.CreateMigrateFactory()
	migrateFact.CreateTable()

	//获取迁移表数据
	getMigrationMap()

	//下面开始创建表 (platform库)
	if !existMigration("db_platform_init_tables") {
		adminFact := admin_models.CreateAdminFactory()
		_ = adminFact.CreateTable()
		_ = adminFact.InitData()
		_ = dictionary_models.CreateDictionaryFactory().CreateTable()
		_ = operate_models.CreateOperateLogFactory().CreateTable()
		_ = org_menu_models.CreateOrgMenuFactory().CreateTable()
		productFact := prodcut_models.CreateProductFactory()
		_ = productFact.CreateTable()
		_ = productFact.InitData()
		_ = server_monitor_models.CreateServerMonitorFactory().CreateTable()
		_ = sms_models.CreateSmsLogFactory().CreateTable()
		_ = software_models.CreateSoftwareFactory().CreateTable()
		_ = software_models.CreateSoftwareDownloadFactory().CreateTable()
		_ = upload_models.CreateUploadFactory().CreateTable()
	}
	//下面开始创建表 (organization库)
	if !existMigration("db_organization_init_tables") {
		_ = agent_type_models.CreateAgentTypeFactory().CreateTable()
		_ = custom_table_models.CreateCustomTableFactory().CreateTable()
		_ = department_models.CreateDepartmentFactory().CreateTable()
		_ = org_operate_models.CreateOperateLogFactory().CreateTable()
		_ = order_models.CreateOrderLogFactory().CreateTable()
		_ = order_models.CreateOrderFactory().CreateTable()
		_ = organization_models.CreateOrgExtendFactory().CreateTable()
		_ = organization_models.CreateOrgLogFactory().CreateTable()
		_ = org_stat_models.CreateOrgStatDailyFactory().CreateTable()
		_ = organization_models.CreateOrganizationFactory().CreateTable()
		_ = role_models.CreateRoleFactory().CreateTable()
		_ = user_models.CreateUserFactory().CreateTable()
	}
	//下面开始创建表 (agent库)
	if !existMigration("db_agent_init_tables") {
		_ = agent_access_models.CreateAccessLogFactory().CreateTable()
		_ = agent_models.CreateAgentFactory().CreateTable()
		_ = agent_models.CreateExtendFactory().CreateTable()
		_ = blackmail_models.CreateGlobalFactory().CreateTable()
		_ = blackmail_models.CreateGroupFactory().CreateTable()
		_ = blackmail_models.CreatePolicyFactory().CreateTable()
		_ = dir_models.CreateDirFactory().CreateTable()
		_ = proc_models.CreateDllFactory().CreateTable()
		_ = proc_models.CreateDllRunningFactory().CreateTable()
		_ = agent_error_models.CreateErrorFactory().CreateTable()
		_ = agent_monitor_models.CreateMonitorFactory().CreateTable()
		_ = outreach_models.CreateGlobalFactory().CreateTable()
		_ = plug_models.CreateDynamicFactory().CreateTable()
		_ = plug_models.CreateIpGlobalFactory().CreateTable()
		_ = plug_models.CreateIpPolicyFactory().CreateTable()
		_ = port_models.CreateVirtualGlobalFactory().CreateTable()
		_ = port_models.CreateVirtualPolicyFactory().CreateTable()
		_ = port_models.CreateWorkingFactory().CreateTable()
		_ = proc_models.CreateProcGlobalFactory().CreateTable()
		_ = proc_models.CreateProcSampleFactory().CreateTable()
		_ = proc_models.CreateProcessRunningFactory().CreateTable()
		_ = proc_models.CreateProcessFactory().CreateTable()
		_ = tamper_models.CreateGlobalFactory().CreateTable()
		_ = tamper_models.CreatePolicyFactory().CreateTable()
		_ = task_models.CreateTaskFactory().CreateTable()
		_ = threaten_models.CreateCacheFactory().CreateTable()
		_ = threaten_models.CreatePushFactory().CreateTable()
		_ = trust_dir_models.CreateGlobalFactory().CreateTable()
		_ = trust_dir_models.CreatePolicyFactory().CreateTable()
		_ = usb_models.CreateUsbFactory().CreateTable()
		_ = usb_models.CreateGlobalFactory().CreateTable()
		_ = usb_models.CreateLogFactory().CreateTable()
	}

	//下面开始创建表 (alarm库)
	if !existMigration("db_alarm_init_tables") {
		_ = alarm_models.CreateAlarmBlackmailFactory().CreateTable()
		_ = alarm_models.CreateAlarmDllFactory().CreateTable()
		_ = alarm_models.CreateAlarmLoginFactory().CreateTable()
		_ = alarm_models.CreateAlarmOutreachFactory().CreateTable()
		_ = alarm_models.CreateAlarmPlugFactory().CreateTable()
		_ = alarm_models.CreateAlarmPortFactory().CreateTable()
		_ = alarm_models.CreateAlarmProcessFactory().CreateTable()
		_ = alarm_models.CreateAlarmSelfFactory().CreateTable()
		_ = alarm_stat_models.CreateAlarmStatDailyFactory().CreateTable()
		_ = alarm_models.CreateAlarmTamperFactory().CreateTable()
		_ = alarm_models.CreateAlarmUsbFactory().CreateTable()
	}

	//下面开始创建表 (syslog库)
	if !existMigration("db_syslog_init_tables") {
		_ = syslog_models.CreateSyslogDnsFactory().CreateTable()
		_ = syslog_models.CreateSyslogInnerFactory().CreateTable()
		_ = syslog_models.CreateSyslogLoginFactory().CreateTable()
		_ = syslog_models.CreateSyslogOuterFactory().CreateTable()
		_ = syslog_models.CreateSyslogProcessFactory().CreateTable()
	}

	//保存迁移数据
	saveMigrations()
}
