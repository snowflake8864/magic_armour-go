package engine_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 软件下载表

func CreateEngineFactory() *EngineModel {
	return &EngineModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type EngineModel struct {
	model.BaseModel
	Id int64 `json:"id" gorm:"primaryKey"`

	EngineName string `json:"engine_name" gorm:"column:engine_name"`
	EngineCode string `json:"engine_code" gorm:"column:engine_code"`
	EngineType int    `json:"engine_type" gorm:"column:engine_type"`

	CurrentVersionID int    `json:"current_version_id" gorm:"column:current_version_id"`
	CurrentVersion   string `json:"current_version" gorm:"column:current_version"`
	LatestVersionID  int    `json:"latest_version_id" gorm:"column:latest_version_id"`
	LatestVersion    string `json:"latest_version" gorm:"column:latest_version"`

	Status        int `json:"status" gorm:"column:status"`
	PublishStatus int `json:"publish_status" gorm:"column:publish_status"`
	ISOfficial    int `json:"is_official" gorm:"column:is_official"`

	DefaultConfig  string `json:"default_config" gorm:"column:default_config"`
	ConfigTemplate string `json:"config_template" gorm:"column:config_template"`

	Description string `json:"description" gorm:"column:description"`

	CreateId   int64  `json:"create_id" gorm:"column:create_id"`
	CreateName string `json:"create_name" gorm:"column:create_name"`
	CreatedAt  string `json:"created_at" gorm:"column:created_at"`
	UpdateId   int64  `json:"update_id" gorm:"column:update_id"`
	UpdateName string `json:"update_name" gorm:"column:update_name"`
	UpdatedAt  string `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt  string `json:"deleted_at" json:"deleted_at"`
}

// 表名

func (m *EngineModel) TableName() string {
	return "p_engine"
}

func (m *EngineModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
			-- ========== 主键 ==========
			id BIGINT NOT NULL AUTO_INCREMENT COMMENT '引擎ID',
			
			-- ========== 基本信息 ==========
			engine_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '引擎名称',
			engine_code VARCHAR(64) NOT NULL DEFAULT '' COMMENT '引擎编码(唯一标识)',
			engine_type INT NOT NULL DEFAULT 0 COMMENT '引擎类型: 0-病毒查杀 1-扫描引擎 2-监控引擎 3-防御引擎 4-其他',
			
			-- ========== 版本信息 ==========
			current_version_id INT NOT NULL DEFAULT 0 COMMENT '当前使用版本ID',
			current_version VARCHAR(32) NOT NULL DEFAULT '' COMMENT '当前版本号',
			latest_version_id INT NOT NULL DEFAULT 0 COMMENT '最新版本ID',
			latest_version VARCHAR(32) NOT NULL DEFAULT '' COMMENT '最新版本号',
			
			-- ========== 状态信息 ==========
			status INT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用 1-启用 2-维护中',
			publish_status INT NOT NULL DEFAULT 0 COMMENT '发布状态: 0-草稿 1-已发布 2-已下架',
			is_official INT NOT NULL DEFAULT 1 COMMENT '是否官方引擎: 0-第三方 1-官方',

		 	-- ========== 配置信息 ==========
			config_template JSON DEFAULT NULL COMMENT '配置模板(JSON格式)',
			default_config JSON DEFAULT NULL COMMENT '默认配置',
			
			-- ========== 描述信息 ==========
			description VARCHAR(500) NOT NULL DEFAULT '' COMMENT '引擎描述',
			
			-- ========== 审计字段 ==========
			create_id BIGINT NOT NULL DEFAULT 0 COMMENT '创建人ID(admin.id)',
			create_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '创建人姓名(admin.username)',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			update_id BIGINT NOT NULL DEFAULT 0 COMMENT '更新人ID',
			update_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '更新人姓名',
			updated_at DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			deleted_at DATETIME DEFAULT NULL COMMENT '删除时间(逻辑删除)',
			
			-- ========== 索引定义 ==========
			PRIMARY KEY (id),
			UNIQUE KEY uk_engine_code (engine_code),
			KEY idx_engine_type (engine_type),
			KEY idx_status (status),
			KEY idx_publish_status (publish_status),
			KEY idx_is_official (is_official),
			KEY idx_created_at (created_at),
			KEY idx_deleted_at (deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='引擎主表';
`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *EngineModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *EngineModel) InsertItem(s *EngineModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *EngineModel) GetItemById(id int64) (item *EngineModel, err error) {
	result := m.Model(m).Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemByUnique 通过id获取一条记录
func (m *EngineModel) GetItemByUnique(orgId, softwareId, userId int64, serviceIp string) (item *EngineModel, err error) {
	result := m.Model(m).Where("org_id=? AND software_id=? AND create_id=? AND service_ip=?", orgId, softwareId, userId, serviceIp).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}
