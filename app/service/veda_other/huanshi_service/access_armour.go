package huanshi_service

type accessArmour struct {
	paramTimestamp int64
	paramSecret    string
}

var accessConfig accessArmour

func (aa *accessArmour) config(timestamp int64, secret string) {
	accessConfig.paramTimestamp = timestamp
	accessConfig.paramSecret = secret
}

func (aa *accessArmour) checkSecret(timestamp int64, secret string) {
	accessConfig.paramTimestamp = timestamp
	accessConfig.paramSecret = secret
}
