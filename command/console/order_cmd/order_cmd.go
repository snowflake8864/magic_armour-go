package order_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *OrderCmd {
	return &OrderCmd{}
}

type OrderCmd struct {
}

// UnpaidExpire  订单未支付过期
func (crd *OrderCmd) UnpaidExpire() {
	variable.ZapLog.Sugar().Debug("UnpaidExpire")
}
