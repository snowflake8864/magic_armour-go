package engine_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 软件下载表

func CreateEngineVersionFactory() *EngineVersionModel {
	return &EngineVersionModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type EngineVersionModel struct {
	model.BaseModel
	Id              int64  `json:"id" gorm:"primaryKey"`
	UploadID        int64  `json:"upload_id" gorm:"column:upload_id"`
	EngineID        int64  `json:"engine_id" gorm:"column:engine_id"`
	VersionNumber   string `json:"version_number" gorm:"column:version_number"`
	VersionMajor    int    `json:"version_major" gorm:"column:version_major"`
	VersionMinor    int    `json:"version_minor" gorm:"column:version_minor"`
	VersionPatch    int    `json:"version_patch" gorm:"column:version_patch"`
	VersionBuild    int    `json:"version_build" gorm:"column:version_build"`
	VersionOrder    int    `json:"version_order" gorm:"column:version_order"`
	EngineType      int    `json:"engine_type" gorm:"column:engine_type"`
	PlatformType    int    `json:"platform_type" gorm:"column:platform_type"`
	PlatformArch    int    `json:"platform_arch" gorm:"column:platform_arch"`
	PlatformVersion string `json:"platform_version" gorm:"column:platform_version"`
	PublishStatus   int    `json:"publish_status" gorm:"column:publish_status"`
	PublishTime     string `json:"publish_time" gorm:"column:publish_time"`
	PublisherID     int64  `json:"publisher_id" gorm:"column:publisher_id"`
	PublisherName   string `json:"publisher_name" gorm:"column:publisher_name"`
	FileSize        int64  `json:"file_size" gorm:"column:file_size"`
	FileMd5         string `json:"file_md5" gorm:"column:file_md5"`
	FileSha256      string `json:"file_sha256" gorm:"column:file_sha256"`
	StoragePath     string `json:"storage_path" gorm:"column:storage_path"`
	FullPath        string `json:"full_path" gorm:"column:full_path"`
	DownloadUrl     string `json:"download_url" gorm:"column:download_url"`
	Dependencies    string `json:"dependencies" gorm:"column:dependencies"`

	ReleaseNotes string `json:"release_notes" gorm:"column:release_notes"`
	ChangeLog    string `json:"change_log" gorm:"column:change_log"`
	Remark       string `json:"remark" gorm:"column:remark"`
	ISLatest     int    `json:"is_latest" gorm:"column:is_latest"`
	ISActive     int    `json:"is_active" gorm:"column:is_active"`
	ISDeleted    int    `json:"is_deleted" gorm:"column:is_deleted"`
	CreateId     int64  `json:"create_id" gorm:"column:create_id"`
	CreateName   string `json:"create_name" gorm:"column:create_name"`
	CreatedAt    string `json:"created_at" gorm:"column:created_at"`
	UpdateId     int64  `json:"update_id" gorm:"column:update_id"`
	UpdateName   string `json:"update_name" gorm:"column:update_name"`
	UpdatedAt    string `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    string `json:"deleted_at" json:"deleted_at"`
}

// 表名

func (m *EngineVersionModel) TableName() string {
	return "p_software_download"
}

func (m *EngineVersionModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
			-- ========== 主键 ==========
			id BIGINT NOT NULL AUTO_INCREMENT COMMENT '引擎版本ID',
			
			-- ========== 关联信息 ==========
			upload_id BIGINT NOT NULL DEFAULT 0 COMMENT '上传记录ID(关联upload表)',
			engine_id BIGINT NOT NULL DEFAULT 0 COMMENT '引擎ID(关联engines表)',
			
			-- ========== 版本信息 ==========
			version_number VARCHAR(32) NOT NULL DEFAULT '' COMMENT '版本号(格式: 0.0.0.0)',
			version_major INT NOT NULL DEFAULT 0 COMMENT '主版本号',
			version_minor INT NOT NULL DEFAULT 0 COMMENT '次版本号',
			version_patch INT NOT NULL DEFAULT 0 COMMENT '补丁版本号',
			version_build INT NOT NULL DEFAULT 0 COMMENT '构建版本号',
			version_order INT NOT NULL DEFAULT 0 COMMENT '版本排序值(数值越大版本越新)',
			
			-- ========== 引擎类型 ==========
			engine_type INT NOT NULL DEFAULT 0 COMMENT '引擎类型: 0-病毒查杀 1-扫描引擎 2-监控引擎 3-防御引擎 4-其他',
			
			-- ========== 平台信息 ==========
			platform_type INT NOT NULL DEFAULT 0 COMMENT '平台类型: 0-Windows 1-Linux 2-MacOS 3-Android 4-iOS 5-其他',
			platform_arch INT NOT NULL DEFAULT 0 COMMENT 'CPU架构: 0-ARM 1-AMD64 2-x86 3-ARM64 4-其他',
			platform_version VARCHAR(32) NOT NULL DEFAULT '' COMMENT '最低系统版本要求',
			
			-- ========== 发布信息 ==========
			publish_status INT NOT NULL DEFAULT 0 COMMENT '发布状态: 0-未发布 1-已发布 2-已下架 3-已撤回',
			publish_time DATETIME DEFAULT NULL COMMENT '发布时间',
			publisher_id BIGINT NOT NULL DEFAULT 0 COMMENT '发布人ID',
			publisher_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '发布人姓名',
			
			-- ========== 文件信息 ==========
			file_size BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
			file_md5 VARCHAR(32) NOT NULL DEFAULT '' COMMENT '文件MD5校验值',
			file_sha256 VARCHAR(64) NOT NULL DEFAULT '' COMMENT '文件SHA256校验值',
			storage_path VARCHAR(512) NOT NULL DEFAULT '' COMMENT '存储相对路径',
			full_path VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '存储完整绝对路径',
			download_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '下载地址',

			-- ========== 依赖配置 ==========
			dependencies JSON DEFAULT NULL COMMENT '依赖配置(JSON格式)',
			
			-- ========== 版本说明 ==========
			release_notes TEXT COMMENT '版本更新说明',
			change_log JSON DEFAULT NULL COMMENT '变更日志(结构化JSON)',
			remark TEXT COMMENT '备注信息',
			
			-- ========== 状态标识 ==========
			is_latest INT NOT NULL DEFAULT 0 COMMENT '是否最新版本: 0-否 1-是',
			is_active INT NOT NULL DEFAULT 1 COMMENT '是否启用: 0-禁用 1-启用',
			is_deleted INT NOT NULL DEFAULT 0 COMMENT '逻辑删除: 0-未删除 1-已删除',
			
			-- ========== 审计字段 ==========
			create_id BIGINT NOT NULL DEFAULT 0 COMMENT '创建人ID(admin.id)',
			create_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '创建人姓名(admin.username)',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			update_id BIGINT NOT NULL DEFAULT 0 COMMENT '更新人ID',
			update_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '更新人姓名',
			updated_at DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',
			
			-- ========== 索引定义 ==========
			PRIMARY KEY (id),
			UNIQUE KEY uk_engine_version (engine_id, version_number),
			KEY idx_upload_id (upload_id),
			KEY idx_engine_id (engine_id),
			KEY idx_version_number (version_number),
			KEY idx_engine_type (engine_type),
			KEY idx_platform_type (platform_type),
			KEY idx_platform_arch (platform_arch),
			KEY idx_publish_status (publish_status),
			KEY idx_publish_time (publish_time),
			KEY idx_is_latest (is_latest),
			KEY idx_is_active (is_active),
			KEY idx_is_deleted (is_deleted),
			KEY idx_created_at (created_at),
			KEY idx_engine_status (engine_id, is_deleted, is_active, publish_status),
			KEY idx_engine_latest (engine_id, is_deleted, is_latest, publish_status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='引擎版本表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *EngineVersionModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *EngineVersionModel) InsertItem(s *EngineVersionModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *EngineVersionModel) GetItemById(id int64) (item *EngineVersionModel, err error) {
	result := m.Model(m).Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemByUnique 通过id获取一条记录
func (m *EngineVersionModel) GetItemByUnique(orgId, softwareId, userId int64, serviceIp string) (item *EngineVersionModel, err error) {
	result := m.Model(m).Where("org_id=? AND software_id=? AND create_id=? AND service_ip=?", orgId, softwareId, userId, serviceIp).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}
