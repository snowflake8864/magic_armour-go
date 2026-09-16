package third_user_auth_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/date_helper"
)

// 三方用户认证表

func CreateAuthFactory() *AuthModel {
	return &AuthModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type AuthModel struct {
	model.BaseModel
	Id        int64  `json:"id" gorm:"primaryKey"`
	Addr      string `json:"addr"`
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	IsInit    int    `json:"is_init"`

	OrgSeq  string  `json:"org_seq"`
	OrgAt   *string `json:"org_at" gorm:"default:NULL"`
	UserSeq string  `json:"user_seq" gorm:"default:NULL"`
	UserAt  *string `json:"user_at"`

	NoticeSync int     `json:"notice_sync"`
	NoticeAt   *string `json:"notice_at" gorm:"default:NULL"`

	IsDelete  int     `json:"is_delete"`
	DeletedAt *string `json:"deleted_at" gorm:"default:NULL"`

	CreatedAt string `json:"created_at"`
}

// 表名
func (m *AuthModel) TableName() string {
	return "p_third_user_auth"
}
func (m *AuthModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			addr varchar(128) NOT NULL DEFAULT '' COMMENT '服务地址',
			app_id varchar(64) NOT NULL DEFAULT '' COMMENT '应用相关',
			app_key varchar(128) NOT NULL DEFAULT '' COMMENT '应用相关',
			app_secret varchar(128) NOT NULL DEFAULT '' COMMENT '应用相关',
			is_init tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否初始化',
			org_seq varchar(32) NOT NULL DEFAULT '' COMMENT '组织用户已同步序列',
			org_at datetime NULL COMMENT '组织上次同步时间',
			user_seq varchar(32) NOT NULL DEFAULT '' COMMENT '用户已同步序列',
			user_at datetime NULL COMMENT '用户上次同步时间',
			notice_sync tinyint(1) NOT NULL DEFAULT '0' COMMENT '通知同步的数据',
			notice_at datetime NULL COMMENT '通知同步的时间',
			is_delete tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
			deleted_at datetime NULL COMMENT '删除时间',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY unique_index (addr)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='三方用户认证表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *AuthModel) InitData() bool {
	//curTime := time.Now().Format(variable.DateFormatNow)
	list := []AuthModel{
		{
			Id:        1,
			Addr:      "https://10.32.2.4",
			AppId:     "789546861770502144",
			AppKey:    "56de926f281e4120b932476492b4d451",
			AppSecret: "48642278b1bf4fe28b74fdd65ab799e5",
			CreatedAt: date_helper.CurDatetime(),
		},
	}

	result := m.Create(&list)
	if result.RowsAffected > 0 {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

func (m *AuthModel) GetItemById(id int64) (item *AuthModel) {
	result := m.Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return item
	}
}

func (m *AuthModel) GetNewestItem() (item *AuthModel) {
	result := m.Where("is_delete=0").Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return item
	}
}

func (m *AuthModel) UpdateById(id int64, updates map[string]interface{}) (bool, int64) {
	result := m.Table(m.TableName()).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}
