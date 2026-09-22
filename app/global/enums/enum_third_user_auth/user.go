package enum_third_user_auth

const (
	StatusNormal  int = 2000  //状态:正常
	StatusLocking int = 3000  //状态:停用中（锁定中）
	StatusDeleted int = -1000 //状态:注销（已删除）

	SyncStatusNormal int = 2000 //状态:正常
	SyncStatusRevoke int = 800  //状态:撤回 (文档说明已没有使用权限)

	OrgTypeOrg  int = 1 //机构类型: 单位
	OrgTypeDept int = 2 //机构类型: 部门

	NoticeSyncOrg  int = 1 //通知同步: 组织
	NoticeSyncUser int = 2 //通知同步: 用户

)

/*
// Config 第三方配置
type Config struct {
	Addr      string `json:"addr"`
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

// SyncSequence 同步序列
type SyncSequence struct {
	Init       int    `json:"init"`        //是否完成初始化
	UserSeq    string `json:"user_seq"`    //用户已同步序列
	UserAt     string `json:"user_at"`     //用户上次同步时间
	OrgSeq     string `json:"org_seq"`     //组织已同步序列
	OrgAt      string `json:"org_at"`      //组织上次同步时间
	NoticeSync int    `json:"notice_sync"` //通知同步的数据
	NoticeAt   string `json:"notice_at"`   //通知同步的时间
}*/

func SumNoticeSync(val int) int {
	if max1 := NoticeSyncOrg + NoticeSyncUser; val > max1 {
		return max1
	} else {
		return val
	}
}
