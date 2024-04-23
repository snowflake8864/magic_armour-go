package product

// PAdd 产品
type PAdd struct {
	GrantType int    `form:"grant_type" json:"grant_type" binding:"required,oneof=1 2"`
	Title     string `form:"title" json:"title" binding:"required,min=4,max=128"`
	//GrantModule GrantModule `form:"grant_module" json:"grant_module" binding:"required,dive"`
	Price  float64 `form:"price"   json:"price" price:"required,max=100000"`
	Remark string  `form:"remark" json:"remark" binding:"omitempty,max=200"`
	GrantModule
}

type GrantModule struct {
	Process   *int `json:"m_process" form:"m_process" binding:"required,oneof=0 1"`
	Tamper    *int `json:"m_tamper"  form:"m_tamper" binding:"required,oneof=0 1"`
	Blackmail *int `json:"m_blackmail" form:"m_blackmail" binding:"required,oneof=0 1"`
	Port      *int `json:"m_port" form:"m_port" binding:"required,oneof=0 1"`
	Plug      *int `json:"m_plug"  form:"m_plug" binding:"required,oneof=0 1"`
	Usb       *int `json:"m_usb"  form:"m_usb" binding:"required,oneof=0 1"`
	Syslog    *int `json:"m_syslog"  form:"m_syslog" binding:"required,oneof=0 1"`
	Outreach  *int `json:"m_outreach" form:"m_outreach" binding:"required,oneof=0 1"`
	Login     *int `json:"m_login" form:"m_login" binding:"required,oneof=0 1"`
}
