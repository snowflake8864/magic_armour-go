package model_platform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 软件下载表

func CreateSoftwareDownloadFactory() *SoftwareDownloadModel {
	return &SoftwareDownloadModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type SoftwareDownloadModel struct {
	model.BaseModel
	Id          int64  `json:"id" gorm:"primaryKey"`
	SoftwareId  int64  `json:"software_id"`
	ServiceIp   string `json:"service_ip"`
	Alias       string `json:"alias"`
	FullPath    string `json:"full_path"`
	StoragePath string `json:"storage_path"`
	FileHash    string `json:"file_hash"`
	FileSize    int    `json:"file_size"`
	Remark      string `json:"remark"`
	CreateId    int64  `json:"create_id"`
	CreateName  string `json:"create_name"`
	CreatedAt   string `json:"created_at"`
}

// 表名

func (m *SoftwareDownloadModel) TableName() string {
	return "p_software_download"
}

func (m *SoftwareDownloadModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			software_id int(11) NOT NULL DEFAULT '0' COMMENT '软件编号(software.id)',
            service_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '服务器地址',
            alias varchar(128)  NOT NULL DEFAULT '' COMMENT '别名',
			full_path varchar(2048)  NOT NULL DEFAULT '' COMMENT '全路径',
			storage_path varchar(1024)  NOT NULL DEFAULT '' COMMENT '存储相对路径(项目下)',
            file_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '文件hash值',
			file_size int(11) NOT NULL DEFAULT '0' COMMENT '文件大小',
            remark varchar(512)  NOT NULL COMMENT '备注',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
		 	KEY index_software_create_service (software_id,create_id,service_ip)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='软件下载表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
