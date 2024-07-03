package order_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *OrderCmd {
	return &OrderCmd{}
}

type OrderCmd struct {
}

// UnpaidExpired  订单未支付且已过期的
func (crd *OrderCmd) UnpaidExpired() {
	variable.ZapLog.Sugar().Debug("UnpaidExpired")
}
