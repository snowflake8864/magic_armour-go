package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 软件表

func CreateSoftwareFactory() *SoftwareModel {
	return &SoftwareModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type SoftwareModel struct {
	model.BaseModel
	Id           string `gorm:"primaryKey" json:"id"`
	UploadId     int    `gorm:"column:UploadId" json:"upload_id"`
	Ver          string `json:"ver"`
	PlatformType int    `gorm:"column:PlatformType" json:"platform_type"`
	OsType       string `gorm:"column:OsType" json:"os_type"`
	IsPublish    int    `gorm:"column:IsPublish" json:"is_publish"`
	Remark       string `json:"remark"`
	CreateId     int    `gorm:"column:CreateId" json:"create_id"`
	CreateName   string `gorm:"column:CreateName"json:"create_name"`
	CreatedAt    string `gorm:"column:CreatedAt" json:"created_at"`
}

// 表名

func (m *SoftwareModel) TableName() string {
	return "p_software"
}

func (m *SoftwareModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			  id int NOT NULL DEFAULT '0' COMMENT '自增id',
			upload_id int NOT NULL DEFAULT '0' COMMENT '(upload.id)',
            ver varchar(64)  NOT NULL DEFAULT '' COMMENT '版本号(0.0.0.0)',
			platform_type tinyint NOT NULL DEFAULT '0' COMMENT '平台类型(win|linxu, 统一定义)',
			os_type smallint NOT NULL DEFAULT '0' COMMENT '系统类别(参考ubuntu|centos等,统一定义)',
			os_ver varchar(512)  NOT NULL DEFAULT '' COMMENT '系统版本(统一定义)',
			is_publish tinyint NOT NULL DEFAULT '0' COMMENT '是否发布',
            remark text  NOT NULL COMMENT '备注',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			INDEX index_ver (ver) USING INVERTED,
			INDEX index_platform_type (platform_type) USING INVERTED,
			INDEX index_os_type (os_type) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(id) 
		DISTRIBUTED BY HASH(id) BUCKETS 1 
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
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
