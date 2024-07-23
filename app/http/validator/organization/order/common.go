package order

// SearchOrder 订单查询条件
type SearchOrder struct {
	ProductName string `form:"productName" json:"productName" binding:"omitempty,min=2,max=128"`
	IsExpire    *int   `form:"isExpire" json:"isExpire" binding:"omitempty,oneof=-100 0 1"`
	OrderStatus *int   `form:"orderStatus" json:"orderStatus" binding:"omitempty,oneof=-100 1 5 8"`
	StartTime   string `form:"startTime" json:"startTime" binding:"omitempty,datetime=2006-01-02 15:04:05"`
	EndTime     string `form:"endTime" json:"endTime" binding:"omitempty,datetime=2006-01-02 15:04:05"`
}

// SingleOrder 订单单个
type SingleOrder struct {
	OrderNo string `form:"orderNo" json:"orderNo" binding:"required,max=20"`
}
