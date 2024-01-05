package test

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"os"
	"testing"
)

var (
	username    = "root"
	password    = "veda2020"
	hostname    = "192.168.17.251:9030"
	dbname      = "magic_armour_platform"
	errUpFailed = errors.New("migration up failed")
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", username, password, hostname, dbName)
}

func TestDbMigrate(t *testing.T) {
	// 新建migrate对象

	//m, err := migrate.New("file://D:\\go_project\\magic_armour-go\\database", "mysql://"+dsn(dbname))
	m, err := migrate.New("//D:\\go_project\\magic_armour-go\\database\\20230616164949_init_up.sql", "mysql://"+dsn(dbname))

	if err != nil {
		fmt.Printf("new migrate err: %+v \n", err)
		return
	}

	// 进行up操作
	err = m.Up()
	version, dirty, _ := m.Version()
	if err == nil || err == migrate.ErrNoChange {
		fmt.Printf("migrate up success, version: %+v, dirty: %+v \n", version, dirty)
		return
	}

	fmt.Printf("migrate up failed, version: %+v, dirty: %+v, err: %+v \n", version, dirty, err)

	// 只要up没有成功，后续都是失败
	// 如果up失败，尝试回滚一步
	version, dirty, _ = m.Version()
	err = m.Steps(-1)
	if err == nil || err == os.ErrNotExist {
		fmt.Printf("migrate down -1 success, version: %+v, dirty: %+v \n", version, dirty)
		return
	}

	fmt.Printf("migrate down -1 failed, version: %+v, dirty: %+v, err: %+v \n", version, dirty, err)
	// 如果回滚失败，判断是不是因为dirty
	er, ok := err.(migrate.ErrDirty)
	if !ok {
		// 不是dirty错误
		return
	}

	// 是dirty错误，那我们强制设置version，再利用这个版本进行回滚
	err = m.Force(er.Version)
	if err != nil {
		fmt.Printf("migrate force %+v err: %+v \n", er.Version, err)
		return
	}

	err = m.Steps(-1)
	if err == nil || err == os.ErrNotExist {
		fmt.Printf("migrate down -1 after force success, version: %+v, dirty: %+v \n", version, dirty)
		return
	}
	fmt.Printf("migrate down -1 after force failed, version: %+v, dirty: %+v, err: %+v \n", version, dirty, err)
	return

}
