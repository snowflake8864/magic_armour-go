package test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

var (
	dbConn *sql.DB
	err    error
)

func TestGoMysql(t *testing.T) {
	dbConn, err = sql.Open("mysql", "root:veda2020@tcp(127.0.0.1:3306)/magic_armour_saas?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	// 设置数据库最大连接数
	dbConn.SetConnMaxLifetime(100)
	// 设置上数据库最大闲置连接数
	dbConn.SetMaxIdleConns(10)

	// 验证连接
	if err := dbConn.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}

	sql := "INSERT  INTO migrations (Migration,  batch) VALUES (?, ?)"
	// 防止sql注入
	stmIns, err := dbConn.Prepare(sql)
	if err != nil {
		fmt.Println("防止sql注入 fail")
		fmt.Println(err)
		return
	}
	_, err = stmIns.Exec("20240108000001_init_paltform", 50)
	if err != nil {
		fmt.Println("stmIns Exec fail")
	}
}

func TestGoDoris(t *testing.T) {
	dbConn, err = sql.Open("mysql", "root:veda2020@tcp(192.168.17.251:9030)/magic_armour_platform?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	// 设置数据库最大连接数
	dbConn.SetConnMaxLifetime(100)
	// 设置上数据库最大闲置连接数
	dbConn.SetMaxIdleConns(10)

	// 验证连接
	if err := dbConn.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}

	sql1 := "INSERT  INTO p_migrations (Migration,  created_at) VALUES (\"20240108000001_init_paltform\", \"2024-01-08 00:00:02\")"
	_, err = dbConn.Exec(sql1)
	if err != nil {
		fmt.Println("dbConn Exec fail")
		fmt.Println(err)
	}

	sql := "INSERT  INTO p_migrations (Migration,  created_at) VALUES (?, ?)"
	// 防止sql注入
	stmIns, err := dbConn.Prepare(sql)
	if err != nil {
		fmt.Println("防止sql注入 fail")
		fmt.Println(err)
		return
	}
	_, err = stmIns.Exec("20240108000001_init_paltform", "2024-01-08 00:00:01")
	if err != nil {
		fmt.Println("stmIns Exec fail")
	}
}
