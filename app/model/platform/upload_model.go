package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 上传记录表

func CreateUploadFactory() *UploadModel {
	return &UploadModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type UploadModel struct {
	model.BaseModel
	CreatedAt    string `gorm:"column:CreatedAt" json:"created_at"`
	Type         string `json:"type"`
	OriginalName string `gorm:"column:OriginalName" json:"original_name"`
	StoragePath  string `gorm:"column:StoragePath" json:"storage_path"`
	StorageName  string `gorm:"column:StorageName" json:"storage_name"`
	FullPath     string `gorm:"column:FullPath" json:"full_path"`
	FileHash     string `gorm:"column:FileHash" json:"file_hash"`
	Size         string `json:"size"`
	CreateId     int    `gorm:"column:CreateId" json:"create_id"`
	CreateName   string `gorm:"column:CreateName"json:"create_name"`
}

// 表名

func (m *UploadModel) TableName() string {
	return "p_uploads"
}

func (m *UploadModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			created_at datetime NULL,
			type varchar(32)  NOT NULL DEFAULT '' COMMENT '类型(统一定义)',
			original_name varchar(128)  NOT NULL DEFAULT '' COMMENT '原文件名',
			storage_path varchar(512)  NOT NULL DEFAULT '' COMMENT '存储相对路径(项目下)',
			storage_name varchar(128)  NOT NULL DEFAULT '' COMMENT '存储文件名',
			full_path varchar(1024)  NOT NULL DEFAULT '' COMMENT '全部路径',
			file_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '文件hash值',
			suffix varchar(8)  NOT NULL DEFAULT '' COMMENT '后缀名',
			size int NOT NULL DEFAULT '0' COMMENT '文件大小',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)'
		)
		DUPLICATE  KEY(created_at, type, original_name) 
		DISTRIBUTED BY HASH(type) BUCKETS 1 
		PROPERTIES (
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
