package initMigrate

import (
	"github.com/spf13/cobra"
	"goskeleton/app/global/variable"
	modelPlatform "goskeleton/app/model/platform"
	"time"
)

var (
	DbName string
)

// DbMigrate 数据迁移
var DbMigrate = &cobra.Command{
	Use:     "db_migrate",
	Aliases: []string{"db_migrate"}, // 定义别名
	Short:   "数据库迁移实现",
	Long: `调用方法：
			1.进入项目根目录（Ginkeleton）。 
			2.执行 go  run  cmd/cli/main.go  db_migrate -h  //可以查看使用指南
			3.执行 go  run  cmd/cli/main.go  db_migrate  -A platform  // 通过 Action 动作执行相应的命令
		`,
	// Run 命令是 核心 命令，其余命令都是为该命令服务，可以删除，由您自由选择
	Run: func(cmd *cobra.Command, args []string) {
		Start(DbName)
	},
}

// 注册命令、初始化参数
func init() {
	DbMigrate.Flags().StringVarP(&DbName, "dbName", "D", "platform", "-D 指定参数,例如：-D platform")
}

// 开始执行业务
func Start(DbName string) {
	ok := CreateTableMigrationIfNotExist()
	if !ok {
		return
	}

	switch DbName {
	case "platform":
		InitDbPlatform()
	case "alarm":
		variable.ZapLog.Info(DbName)
	case "syslog":
		variable.ZapLog.Info(DbName)
	}
}

// 迁移数据
var migrateMap map[string]string

// 不在迁移表中
var notExistMigrateList []modelPlatform.MigrateModel

// 创建表：migration 如果不存在, 并获取迁移数据

func CreateTableMigrationIfNotExist() bool {
	migrateFact := modelPlatform.CreateMigrateFactory()
	migrateFact.CreateTable()

	migrateMap, _ := migrateFact.GetAllMaps()
	if migrateMap != nil {
		variable.ZapLog.Sugar().Debug("migrate GetAllMaps fail")
		return false
	}
	return true
}

func ExistMigration(migration string) bool {

	var st = modelPlatform.MigrateModel{
		Migration: migration,
		CreatedAt: time.Now().Format(variable.DateFormat),
	}

	if migrateMap == nil {
		notExistMigrateList = append(notExistMigrateList, st)
		return false
	}

	_, ok := migrateMap[migration]
	if ok {
		return true
	} else {
		notExistMigrateList = append(notExistMigrateList, st)
		return false
	}
}

func InitDbPlatform() {
	if !ExistMigration("create_table_agent") {
		_ = modelPlatform.CreateAgentFactory().CreateTable()
	}
	if !ExistMigration("create_table_dictionary") {
		_ = modelPlatform.CreateDictionaryFactory().CreateTable()
	}

}
