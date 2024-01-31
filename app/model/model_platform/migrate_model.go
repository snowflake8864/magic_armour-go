package model_platform

import (
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"time"
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

func (m *MigrateModel) InsertOne(s *MigrateModel) bool {
	result := m.Create(s)
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

func (m *MigrateModel) InsertItem(s *MigrateModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&MigrateModel{}).Create(s)
	})
	fmt.Println(sql)
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *MigrateModel) InsertItems(list *[]MigrateModel) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&MigrateModel{}).Create(list)
	})
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *MigrateModel) ShowOneItem() (*MigrateModel, error) {
	sql := "SELECT  * FROM  `p_migrations`  WHERE Migration='20240108000001_init_paltform111' LIMIT 1"
	result := m.Raw(sql).First(m)
	if result.Error == nil {
		return m, nil
	} else {
		return nil, result.Error
	}
}

func (m *MigrateModel) GetALLItems() (temp []MigrateModel, err error) {
	sql := "SELECT  * FROM  `p_migrations`"
	result := m.Raw(sql).Find(&temp)
	if result.Error == nil {
		return temp, nil
	} else {
		return nil, result.Error
	}
}

func (m *MigrateModel) GetAllMap() (map[string]MigrateModel, error) {
	var items []MigrateModel
	sql := "SELECT  * FROM  `p_migrations`"
	result := m.Raw(sql).Find(&items)
	if result.Error == nil {
		res := make(map[string]MigrateModel)
		for _, obj := range items {
			res[obj.Migration] = obj
		}
		return res, nil
	} else {
		return nil, result.Error
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

func (m *MigrateModel) TestToSQL() bool {
	fmt.Println("\n SELECT--------------------------------------------------")
	sql := ""
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&MigrateModel{}).Where("Migration = ?", "20240108000001_init_paltform111").Limit(10).Order("age desc").Find(&[]MigrateModel{})
	})
	fmt.Println(sql)
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		st := "SELECT * FROM p_migrations WHERE Migration = ?"
		return tx.Model(&MigrateModel{}).Raw(st, "20240108000001_init_paltform111")
	})
	fmt.Println(sql)
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		st := "SELECT * FROM p_migrations WHERE Migration = ?"
		return tx.Model(&MigrateModel{}).Raw(st, "20240108000001_init_paltform111")
	})
	fmt.Println(sql)
	fmt.Println("\n INSERT--------------------------------------------------")
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&MigrateModel{}).Exec("INSERT  INTO p_migrations (Migration,  created_at) VALUES (?, ?)",
			"20240108000001_init_paltforma' or '1'='1",
			"2024-01-08 00:00:01")
	})
	fmt.Println(sql)

	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var st = &MigrateModel{
			Migration: "20240108000001_init_paltforma' or '1'='1",
			CreatedAt: time.Now().Format(variable.DateFormat),
		}
		return tx.Model(&MigrateModel{}).Create(st)
	})
	fmt.Println(sql)

	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var st = &[]MigrateModel{
			{
				Migration: "20240108000001_init_paltforma' or '2'='2",
				CreatedAt: time.Now().Format(variable.DateFormat),
			},
			{
				Migration: "20240108000001_init_paltforma' or '3'='3",
				CreatedAt: time.Now().Format(variable.DateFormat),
			},
		}
		return tx.Model(&MigrateModel{}).Create(st)
	})
	fmt.Println(sql)

	fmt.Println("\n update--------------------------------------------------")
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		st := "update p_migrations set Migration=?,  created_at=?  WHERE Migration=? AND created_at=?"
		return tx.Model(&MigrateModel{}).Exec(st, "20240108000001_init_paltforma", variable.DateFormat, "20240108000001_init_paltforma", variable.DateFormat)
	})
	fmt.Println(sql)
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		// 定义更新字段的map， 键值对
		var relaValue = map[string]interface{}{
			"Migration":  "192.168.6.66",
			"created_at": time.Now().Format(variable.DateFormat),
		}
		return tx.Model(&MigrateModel{}).Select("Migration", "created_at").Where("Migration=?", 13).Updates(relaValue)
	})
	fmt.Println(sql)

	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var st = &MigrateModel{
			Migration: "20240108000001_init_paltforma' or '1'='1",
			CreatedAt: time.Now().Format(variable.DateFormat),
		}
		return tx.Model(&MigrateModel{}).Save(st)
	})
	fmt.Println(sql)

	fmt.Println("\n delete--------------------------------------------------")
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		st := "DELETE FROM  p_migrations WHERE Migration=? AND created_at=?"
		return tx.Model(&MigrateModel{}).Exec(st, "20240108000001_init_paltforma", variable.DateFormat)
	})
	fmt.Println(sql)
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var st = &MigrateModel{
			Migration: "20240108000001_init_paltforma' or '2'='2",
			CreatedAt: time.Now().Format(variable.DateFormat),
		}
		return tx.Model(&MigrateModel{}).Delete(st)
	})
	fmt.Println(sql)
	sql = m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&MigrateModel{}).Delete(&MigrateModel{}, 5)
	})
	fmt.Println(sql)
	return true
}
