package modelPlatform

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
	Id         int    `gorm:"primaryKey" json:"id"`
	ParentId   int    `gorm:"column:ParentId" json:"parent_id"`
	Position   int    `json:"position"`
	Title      string `json:"title"`
	WebCode    string `gorm:"column:WebCode" json:"web_code"`
	Sort       int    `json:"sort"`
	IsShow     int    `gorm:"column:IsShow" json:"is_show"`
	IsRoot     int    `gorm:"column:IsRoot" json:"is_root"`
	Remark     string `json:"remark"`
	ServiceApi string `gorm:"column:ServiceApi" json:"service_api"`
	Level      int    `json:"level"`
	CreateId   int    `gorm:"column:CreateId" json:"create_id"`
	CreateName string `gorm:"column:CreateName"json:"create_name"`
	CreatedAt  string `gorm:"column:created_at" json:"created_at"`
}

// 表名

func (m *MenuModel) TableName() string {
	return "p_menus"
}

func (m *MenuModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int NOT NULL DEFAULT '0' COMMENT '自增id',
			parent_id int NOT NULL DEFAULT '0' COMMENT '父级id',
			position tinyint NOT NULL DEFAULT '1' COMMENT '方位(1左,2右,3顶部)',
			type tinyint NOT NULL DEFAULT '0' COMMENT '类型(1菜单,2按钮,3链接)',
			title varchar(128)  NOT NULL DEFAULT '' COMMENT '标题',
			web_code varchar(128)  NOT NULL DEFAULT '' COMMENT '前端页面:标识',
			sort int NOT NULL DEFAULT '1' COMMENT '排序',
			is_show tinyint NOT NULL DEFAULT '0' COMMENT '是否显示(0否,1是),前端判断用',
			is_root tinyint NOT NULL DEFAULT '0' COMMENT '仅超管可用(0否,1是),后端判断用',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			service_api text  NOT NULL COMMENT '服务端对应接口(,)隔开',
			level smallint NOT NULL DEFAULT '1' COMMENT '级别(从1开始)',
			INDEX index_parent (parent_id) USING INVERTED
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
