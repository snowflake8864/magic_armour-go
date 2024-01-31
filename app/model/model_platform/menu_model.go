package model_platform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 菜单表

func CreateMenuFactory() *MenuModel {
	return &MenuModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type MenuModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	ParentId   int64  `json:"parent_id"`
	Position   int    `json:"position"`
	Title      string `json:"title"`
	WebCode    string `json:"web_code"`
	Sort       int    `json:"sort"`
	IsShow     int    `json:"is_show"`
	IsRoot     int    `json:"is_root"`
	Remark     string `json:"remark"`
	ServiceApi string `json:"service_api"`
	Level      int    `json:"level"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *MenuModel) TableName() string {
	return "p_menus"
}

func (m *MenuModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) NOT NULL DEFAULT '0' COMMENT '自增id',
			parent_id int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
			position tinyint(1) NOT NULL DEFAULT '1' COMMENT '方位(1左,2右,3顶部)',
			type tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型(1菜单,2按钮,3链接)',
			title varchar(128)  NOT NULL DEFAULT '' COMMENT '标题',
			web_code varchar(128)  NOT NULL DEFAULT '' COMMENT '前端页面:标识',
			sort int(11) NOT NULL DEFAULT '1' COMMENT '排序',
			is_show tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否显示(0否,1是),前端判断用',
			is_root tinyint(1) NOT NULL DEFAULT '0' COMMENT '仅超管可用(0否,1是),后端判断用',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			service_api text  NOT NULL COMMENT '服务端对应接口(,)隔开',
			level smallint(4) NOT NULL DEFAULT '1' COMMENT '级别(从1开始)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_parent (parent_id)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *MenuModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *MenuModel) InitData() bool {
	_ = m.TruncateTable()

	sql := `INSERT INTO ` + m.TableName() + `(id, parent_id, position, type, title, web_code, sort, is_show, is_root, remark, service_api, level, created_at) VALUES`
	sql += `
			(1, 1, 1, 1, 'title', 'web_code', 1, 0, 0, 'ff', '', 1, '2024-01-30 00:00:00');
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
