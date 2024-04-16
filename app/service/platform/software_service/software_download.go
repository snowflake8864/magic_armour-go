package software_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/openssl_encrypt"
)

func CreateDownloadFactory() *SoftwareDownload {
	return &SoftwareDownload{masterModel: model_platform.CreateSoftwareDownloadFactory()}
}

type SoftwareDownload struct {
	masterModel *model_platform.SoftwareDownloadModel
}

type AgentConfig struct {
	OrgId  int64 `json:"o"`
	UserId int64 `json:"u"`
	Expire int64 `json:"t"`
}

// EncryptConfig 加密配置
func (crd *SoftwareDownload) EncryptConfig(OrgId, UserId int64) string {

	//从数据字典中获取过期时间
	var hour int64 = 366
	agentConfigExpire := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentConfigExpire)
	if !helpers.StringIsEmpty(agentConfigExpire) {
		hour = helpers.String2Int64(agentConfigExpire)
	}

	//登录信息
	loginInfo := AgentConfig{
		OrgId:  OrgId,
		UserId: UserId,
		Expire: helpers.AddTimestampByHour(hour),
	}
	jsonStr := helpers.JsonEncode(loginInfo)

	//加密
	cryptStr, _ := openssl_encrypt.CreateAesFactory().CbcDecryptByHex(jsonStr)
	return cryptStr
}

// DecryptConfig 解密配置
func (crd *SoftwareDownload) DecryptConfig(cryptStr string) (loginInfo AgentConfig, err error) {
	//解密成json
	jsonStr, err := openssl_encrypt.CreateAesFactory().CbcDecryptByHex(cryptStr)
	if err != nil {
		return AgentConfig{}, errors.New("userid decrypt fail")
	}

	//转成登录结构体
	err = helpers.Json2Struct(jsonStr, &loginInfo)
	if err != nil {
		return AgentConfig{}, errors.New("userid json fail")
	}
	return loginInfo, nil
}
