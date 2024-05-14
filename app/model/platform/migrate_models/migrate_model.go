package migrate_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 数据迁移表

func CreateMigrateFactory() *MigrateModel {
	return &MigrateModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type MigrateModel struct {
	model.BaseModel
	Id        int64  `json:"id" gorm:"primaryKey"`
	Migration string `gorm:"column:migration" json:"migration"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *MigrateModel) TableName() string {
	return "p_migrations"
}

// 创建表
func (m *MigrateModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			migration varchar(128) NOT NULL DEFAULT '',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY unique_index (migration)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}

}

func (m *MigrateModel) InsertBatch(list *[]MigrateModel) bool {
	result := m.Create(list)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *MigrateModel) GetAllMaps() (map[string]string, error) {
	var items []MigrateModel

	sql := "SELECT  * FROM  `p_migrations`"
	result := m.Raw(sql).Find(&items)
	if result.Error == nil {
		res := make(map[string]string)
		for _, obj := range items {
			res[obj.Migration] = obj.CreatedAt
		}
		return res, nil
	} else {
		return nil, result.Error
	}
}
