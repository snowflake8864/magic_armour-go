package database

import (
	"fmt"
	"github.com/spf13/cobra"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_access_models"
	"goskeleton/app/model/agent/agent_error_models"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/model/agent/agent_online_models"
	"goskeleton/app/model/agent/antivirus_engine_task_models"
	"goskeleton/app/model/agent/backup_models"
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
	"goskeleton/app/model/platform/antivirus_engine_models"
	"goskeleton/app/model/platform/dictionary_models"
	"goskeleton/app/model/platform/migrate_models"
	"goskeleton/app/model/platform/operate_models"
	"goskeleton/app/model/platform/org_menu_models"
	"goskeleton/app/model/platform/prodcut_models"
	"goskeleton/app/model/platform/server_monitor_models"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/organization/role_service"
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
			3.执行 go  run  cmd/cli/main.go  migrate -A 行为  // 通过 Action 动作执行相应的命令
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		start(Action)
	},
}

// 注册命令、初始化参数
func init() {
	DbMigrate.Flags().StringVarP(&Action, "Action", "A", "", "-A 或者 --Action 选择行为")
}

func TestStart(action string) {
	start(action)
}

// 开始执行业务
func start(action string) {
	variable.ZapLog.Debug("action:" + action)
	switch action {
	case "init":
		initTableAndData()
	case "menu":
		_ = org_menu_models.CreateOrgMenuFactory().InitData()
	case "dictionary":
		_ = dictionary_models.CreateDictionaryFactory().InitData()
		_ = dictionary_service.CreateCacheFactory().Rebuild()
	case "append":
		appendTableAndData()
	case "admini":
		addDefaultOrg()
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
		variable.ZapLog.Debug("migrate GetAllMaps fail")
	}
}

// 保存本次迁移数据
func saveMigrations() {
	if len(notExistMigrateList) == 0 {
		variable.ZapLog.Debug("本次迁移数据为空")
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
		CreatedAt: time.Now().Format(variable.DateFormatNow),
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
		_ = operate_models.CreateLogFactory().CreateTable()
		_ = org_menu_models.CreateOrgMenuFactory().CreateTable()
		productFact := prodcut_models.CreateProductFactory()
		_ = productFact.CreateTable()
		_ = productFact.InitData()
		_ = server_monitor_models.CreateServerMonitorFactory().CreateTable()
		_ = sms_models.CreateSmsLogFactory().CreateTable()
		_ = software_models.CreateSoftwareFactory().CreateTable()
		_ = software_models.CreateDownloadFactory().CreateTable()
		_ = upload_models.CreateUploadFactory().CreateTable()
	}
	//下面开始创建表 (organization库)
	if !existMigration("db_organization_init_tables") {
		_ = agent_type_models.CreateAgentTypeFactory().CreateTable()
		_ = custom_table_models.CreateCustomTableFactory().CreateTable()
		_ = department_models.CreateDepartmentFactory().CreateTable()
		_ = org_operate_models.CreateLogFactory().CreateTable()
		_ = order_models.CreateLogFactory().CreateTable()
		_ = order_models.CreateOrderFactory().CreateTable()
		_ = organization_models.CreateExtendFactory().CreateTable()
		_ = organization_models.CreateLogFactory().CreateTable()
		_ = org_stat_models.CreateDailyFactory().CreateTable()
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
		_ = outreach_models.CreateDetectFactory().CreateTable()
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
		_ = alarm_models.CreateBlackmailFactory().CreateTable()
		_ = alarm_models.CreateDllFactory().CreateTable()
		_ = alarm_models.CreateLoginFactory().CreateTable()
		_ = alarm_models.CreateOutreachFactory().CreateTable()
		_ = alarm_models.CreatePlugFactory().CreateTable()
		_ = alarm_models.CreatePortFactory().CreateTable()
		_ = alarm_models.CreateProcessFactory().CreateTable()
		_ = alarm_models.CreateSelfFactory().CreateTable()
		_ = alarm_stat_models.CreateStatDailyFactory().CreateTable()
		_ = alarm_models.CreateTamperFactory().CreateTable()
		_ = alarm_models.CreateUsbFactory().CreateTable()
	}

	//下面开始创建表 (syslog库)
	if !existMigration("db_syslog_init_tables") {
		_ = syslog_models.CreateDnsFactory().CreateTable()
		_ = syslog_models.CreateInnerFactory().CreateTable()
		_ = syslog_models.CreateLoginFactory().CreateTable()
		_ = syslog_models.CreateOuterFactory().CreateTable()
		_ = syslog_models.CreateProcessFactory().CreateTable()
	}

	//保存迁移数据
	saveMigrations()
}

func appendTableAndData() {
	//获取迁移表数据
	getMigrationMap()
	if ver := "3.1.0_T4"; !existMigration(ver) {
		_ = agent_models.CreateAgentFactory().AlterTable(ver)
		_ = agent_online_models.CreateOnlineLogFactory().CreateTable()
		_ = org_operate_models.CreateLogFactory().AlterTable(ver)
	}
	if ver := "MAT_3.1.0_R1"; !existMigration(ver) {
		_ = agent_models.CreateAgentFactory().AlterTable(ver)
		_ = agent_models.CreateExtendFactory().AlterTable(ver)
		_ = organization_models.CreateExtendFactory().AlterTable(ver)
		_ = backup_models.CreateBackupFactory().CreateTable()
	}
	if ver := "MAT_3.1.0_R2"; !existMigration(ver) {
		_ = agent_models.CreateExtendFactory().AlterTable(ver)
	}
	if ver := "MA_3.1.0_R4_B5"; !existMigration(ver) {
		_ = agent_models.CreateAgentFactory().AlterTable(ver)
		_ = user_models.CreateUserFactory().AlterTable(ver)
	}

	if ver := "MA_3.1.0_R4_B8"; !existMigration(ver) {
		_ = agent_models.CreateAgentFactory().AlterTable(ver)
	}

	//长城卫士
	if ver := "cgw_3.1.0_R1"; !existMigration(ver) {
		_ = alarm_models.CreateVirusFactory().CreateTable()
		_ = alarm_stat_models.CreateStatDailyFactory().AlterTable(ver)
	}

	//长城卫士
	if ver := "cgw_3.1.0_R2_01"; !existMigration(ver) {
		_ = agent_models.CreateAgentFactory().AlterTable(ver)
		_ = antivirus_engine_models.CreateEngineFactory().CreateTable()
		_ = antivirus_engine_task_models.CreateUpgradeTaskFactory().CreateTable()
		_ = antivirus_engine_task_models.CreateUpgradeFailFactory().CreateTable()
	}

	//保存迁移数据
	saveMigrations()
}

// 添加默认企业
func addDefaultOrg() {

	//企业信息
	orgParam := organization_models.OrganizationModel{}
	orgParam.Name = ""
	orgParam.UserName = "admin"
	orgParam.UserPhone = "13000000000"
	orgParam.UserEmail = "dev@test.com"
	orgParam.Remark = "默认企业用户"
	orgParam.Type = enum_organization.TypePerson
	orgParam.Status = enum_organization.StatusNormal
	orgParam.Source = enums.OperateSourcePlatform

	//用户信息
	userParam := user_models.UserModel{}
	userParam.Username = orgParam.UserName
	userParam.Phone = orgParam.UserPhone
	userParam.Email = orgParam.UserEmail
	userParam.Password = "Root@2026"
	userParam.IsRoot = enums.IsYes

	if errCode, err := organize_service.CreateOrgCurdFactory().Add(&orgParam, &userParam, "0.0.0.0"); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("code=%d, msg=%s", errCode, err))
	} else {
		//默认企业角色
		role_service.CreateCurdFactory().AddDefault(1)
	}
}
