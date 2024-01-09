package model

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"strings"
)

// 创建调用工厂

func CreateMigrateFactory() *MigrateModel {
	return &MigrateModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type MigrateModel struct {
	model.BaseModel
	Migration string `gorm:"column:migration" json:"migration"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *MigrateModel) TableName() string {
	return "p_migrations"
}

//创建表

func (m *MigrateModel) CreateTable() bool {
	m.PrepareStmt = false
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			migration varchar(128) NOT NULL DEFAULT '',
			created_at datetime NULL
		)
		UNIQUE KEY(migration)
		DISTRIBUTED BY HASH(migration) BUCKETS 1
		PROPERTIES (
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	//fmt.Println(sql)
	sql = strings.Replace(sql, "	", " ", -1)
	sql = strings.Replace(sql, "\t", "", -1)
	sql = strings.Replace(sql, "\n", "", -1)
	fmt.Println(sql)
	result := m.Exec(sql)

	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}

}

func (m *MigrateModel) AddOne() bool {
	//sql := "INSERT  INTO p_migrations (Migration,  created_at) VALUES (?, ?)"
	//result := m.Exec(sql, "20240108000001_init_paltform", "2024-01-08 00:00:01")
	sql := "INSERT  INTO p_migrations (Migration,  created_at) VALUES (\"20240108000001_init_paltform\", \"2024-01-09 00:01:01\")"
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *MigrateModel) InsertOne() bool {
	//sql := "INSERT  INTO p_migrations (Migration,  created_at) VALUES (?, ?)"
	//result := m.Exec(sql, "20240108000001_init_paltform", "2024-01-08 00:00:01")
	sql := "INSERT  INTO p_migrations (Migration,  created_at) VALUES (\"20240108000001_init_paltform\", \"2024-01-09 00:01:01\")"
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
