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
	Id           int64  `gorm:"primaryKey" json:"id"`
	UploadId     int64  `gorm:"column:UploadId" json:"upload_id"`
	Ver          string `json:"ver"`
	PlatformType int    `gorm:"column:PlatformType" json:"platform_type"`
	OsType       string `gorm:"column:OsType" json:"os_type"`
	IsPublish    int    `gorm:"column:IsPublish" json:"is_publish"`
	Remark       string `json:"remark"`
	CreateId     int64  `gorm:"column:CreateId" json:"create_id"`
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
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			upload_id int(11) NOT NULL DEFAULT '0' COMMENT '(upload.id)',
            ver varchar(64)  NOT NULL DEFAULT '' COMMENT '版本号(0.0.0.0)',
			platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(win|linxu, 统一定义)',
			os_type smallint(4) NOT NULL DEFAULT '0' COMMENT '系统类别(参考ubuntu|centos等,统一定义)',
			os_ver varchar(512)  NOT NULL DEFAULT '' COMMENT '系统版本(统一定义)',
			is_publish tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否发布',
            remark text  NOT NULL COMMENT '备注',
            create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			UNIQUE KEY(id),
			KEY index_ver (ver),
			KEY index_platform_type (platform_type),
			KEY index_os_type (os_type)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='软甲表';
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
