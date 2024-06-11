package software_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
)

func CreateDownloadFactory() *SoftwareDownload {
	return &SoftwareDownload{masterModel: software_models.CreateDownloadFactory()}
}

type SoftwareDownload struct {
	masterModel *software_models.DownloadModel
}

type AgentConfig struct {
	OrgId  int64 `json:"o"`
	UserId int64 `json:"u"`
	Expire int64 `json:"t"`
}

// EncryptConfig 加密配置
// 底层限制长度不能超过260
func (crd *SoftwareDownload) EncryptConfig(OrgId, UserId int64) string {

	//从数据字典中获取过期时间
	var hour int64 = 366
	agentConfigExpire := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentConfigExpire)
	if !string_helper.IsEmpty(agentConfigExpire) {
		hour = string_helper.ToInt64(agentConfigExpire)
	}

	//登录信息
	loginInfo := AgentConfig{
		OrgId:  OrgId,
		UserId: UserId,
		Expire: date_helper.AddTimestampByHour(hour),
	}
	jsonStr := json_helper.Encode(loginInfo)
	//加密
	cryptStr, _ := openssl_encrypt.CreateAesFactory().CbcEncryptToBase64(jsonStr)
	return cryptStr
}

// DecryptConfig 解密配置
func (crd *SoftwareDownload) DecryptConfig(cryptStr string) (loginInfo AgentConfig, err error) {
	//解密成json
	jsonStr, err := openssl_encrypt.CreateAesFactory().CbcDecryptByBase64(cryptStr)
	if err != nil {
		return AgentConfig{}, errors.New("userid decrypt fail")
	}

	//转成登录结构体
	err = json_helper.ToStruct(jsonStr, &loginInfo)
	if err != nil {
		return AgentConfig{}, errors.New("userid json fail")
	}
	return loginInfo, nil
}

// Clear 清空
func (crd *SoftwareDownload) Clear() bool {
	return crd.masterModel.TruncateTable()
}
