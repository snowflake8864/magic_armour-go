package enum_order

import "goskeleton/app/global/enums"

const (
	PayMethodBankcard int = 10 //支付方式： 银行卡
	PayMethodWechat   int = 20 //支付方式： 微信

	OrderStatusCancel    int = 33 //订单状态:已取消
	OrderStatusUnpaid    int = 1  //订单状态:待支付
	OrderStatusPaid      int = 5  //订单状态:已支付
	OrderStatusCompleted int = 8  //订单状态:已完成

	ActionAdd           int = 1 //行为：添加
	ActionPayBankcard   int = 2 //行为：对公转账方式付款
	ActionCancel        int = 3 //行为：取消授权
	ActionDecline       int = 4 //行为：修改支付金额
	ActionConfirmPaid   int = 5 //行为：确认支付成功
	ActionConfirmUnpaid int = 6 //行为：确认未付款
	ActionAutoExpire    int = 7 //行为：自动过期
	ActionGrant2Agent   int = 8 //行为：授权给终端
)

var EnumMap = map[string]map[int]string{

	"order_status": {
		OrderStatusCancel:    "待支付",
		OrderStatusUnpaid:    "已支付",
		OrderStatusPaid:      "已完成",
		OrderStatusCompleted: "已取消",
	},
	"pay_method": {
		PayMethodBankcard: "对公转账",
	},
	"grant_duration": {
		1:  "1个月",
		2:  "2个月",
		3:  "3个月",
		4:  "4个月",
		5:  "5个月",
		6:  "6个月",
		12: "1年",
		24: "2年",
		36: "3年",
	},
	"action": {
		ActionAdd:           "购买授权",
		ActionPayBankcard:   "对公转账方式付款",
		ActionCancel:        "取消授权",
		ActionDecline:       "修改支付金额",
		ActionConfirmPaid:   "确认支付成功",
		ActionConfirmUnpaid: "确认未付款",
		ActionAutoExpire:    "自动过期",
		ActionGrant2Agent:   "授权给终端",
	},
	"expire_status": {
		enums.IsNot: "有效",
		enums.IsYes: "无效",
	},
}

var WebMapInt = map[string][]int{
	"order_status": {
		OrderStatusCancel,
		OrderStatusUnpaid,
		OrderStatusPaid,
		OrderStatusCompleted,
	},
	"grant_duration": {
		1,
		2,
		3,
		4,
		5,
		6,
		12,
		24,
		36,
	},
	"expire_status": {
		enums.IsNot,
		enums.IsYes,
	},
	"pay_method": {
		PayMethodBankcard,
	},
}
