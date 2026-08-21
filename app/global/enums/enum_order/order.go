package enum_order

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

const (
	PayMethodDefault  int = 1  //支付方式： 后台人工
	PayMethodBankcard int = 10 //支付方式： 对公转账
	PayMethodWechat   int = 20 //支付方式： 微信
	PayMethodLicense  int = 60 //支付方式：导入证书

	OrderStatusUnpaid    int = 1  //订单状态:待支付
	OrderStatusPaid      int = 5  //订单状态:已支付
	OrderStatusCompleted int = 8  //订单状态:已完成
	OrderStatusCancel    int = 33 //订单状态:已取消

	ActionAdd           int = 1 //行为：添加
	ActionPayBankcard   int = 2 //行为：对公转账方式付款
	ActionCancel        int = 3 //行为：取消授权
	ActionDecline       int = 4 //行为：修改支付金额
	ActionConfirmPaid   int = 5 //行为：确认支付成功
	ActionConfirmUnpaid int = 6 //行为：确认未付款
	ActionAutoExpire    int = 7 //行为：自动过期
	ActionGrant2Agent   int = 8 //行为：授权给终端
	ActionImportLicense int = 9 //行为：导入证书

	OrderNoPrefixCloud   string = "Yhj" //订单编号前缀: 云幻甲
	OrderNoPrefixOffline string = "Lhj" //订单编号前缀: 离线幻甲
)

var EnumMap = map[string]map[int]string{

	"order_status": {
		OrderStatusUnpaid:    "待支付",
		OrderStatusPaid:      "已支付",
		OrderStatusCompleted: "已完成",
		OrderStatusCancel:    "已取消",
	},
	"order_status_en": {
		OrderStatusUnpaid:    "Pending Payment",
		OrderStatusPaid:      "Paid",
		OrderStatusCompleted: "Completed",
		OrderStatusCancel:    "Cancelled",
	},
	"pay_method": {
		PayMethodBankcard: "对公转账",
		PayMethodLicense:  "授权许可",
	},
	"pay_method_en": {
		PayMethodBankcard: "Bank Transfer",
		PayMethodLicense:  "License",
	},
	"grant_duration": {
		1:   "1个月",
		2:   "2个月",
		3:   "3个月",
		4:   "4个月",
		5:   "5个月",
		6:   "6个月",
		12:  "1年",
		24:  "2年",
		36:  "3年",
		120: "10年",
	},
	"grant_duration_en": {
		1:   "1 Month",
		2:   "2 Months",
		3:   "3 Months",
		4:   "4 Months",
		5:   "5 Months",
		6:   "6 Months",
		12:  "1 Year",
		24:  "2 Years",
		36:  "3 Years",
		120: "10 Years",
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
		ActionImportLicense: "导入授权证书",
	},
	"action_en": {
		ActionAdd:           "Purchase License",
		ActionPayBankcard:   "Pay by Bank Transfer",
		ActionCancel:        "Cancel License",
		ActionDecline:       "Modify Payment Amount",
		ActionConfirmPaid:   "Confirm Payment Success",
		ActionConfirmUnpaid: "Confirm Unpaid",
		ActionAutoExpire:    "Auto Expired",
		ActionGrant2Agent:   "Grant to Agent",
		ActionImportLicense: "Import License Certificate",
	},
	"expire_status": {
		enums.IsNot: "有效",
		enums.IsYes: "无效",
	},
	"expire_status_en": {
		enums.IsNot: "Valid",
		enums.IsYes: "Invalid",
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
		120,
	},
	"expire_status": {
		enums.IsNot,
		enums.IsYes,
	},
	"pay_method": {
		PayMethodBankcard,
	},
}

func FormatPayBankcard(m map[string]string) map[string]string {
	result := make(map[string]string)
	if val, isOk := m[enum_dictionary.OrderBankcardName]; isOk {
		result["card_name"] = val
	} else {
		result["card_name"] = ""
	}
	if val, isOk := m[enum_dictionary.OrderBankcardNo]; isOk {
		result["card_no"] = val
	} else {
		result["card_no"] = ""
	}
	if val, isOk := m[enum_dictionary.OrderBankcardBank]; isOk {
		result["bank_name"] = val
	} else {
		result["bank_name"] = ""
	}
	if val, isOk := m[enum_dictionary.OrderPayingTel]; isOk {
		result["tel"] = val
	} else {
		result["tel"] = ""
	}

	return result
}

// GetPrefixByOrderNo  通过订单编号获取前缀
func GetPrefixByOrderNo(orderNo string) string {
	if strings.Contains(orderNo, OrderNoPrefixCloud) {
		return OrderNoPrefixCloud
	} else {
		return OrderNoPrefixOffline
	}
}

// GetGrantTypeByOrderNo  通过订单编号获取授权类型
func GetGrantTypeByOrderNo(orderNo string) int {
	return string_helper.ToInt(orderNo[3:4])
}
