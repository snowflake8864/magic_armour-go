package software_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/files"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"strings"
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
	SoftId int64 `json:"s"`
	Expire int64 `json:"t"`
}

// EncryptConfig 加密配置  todo: 底层限制长度不能超过256
func (crd *SoftwareDownload) EncryptConfig(orgId, userId, softId int64) string {

	//从数据字典中获取过期时间
	var hour int64 = 366
	agentConfigExpire := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentConfigExpire)
	if !string_helper.IsEmpty(agentConfigExpire) {
		hour = string_helper.ToInt64(agentConfigExpire)
	}

	//登录信息
	loginInfo := AgentConfig{
		OrgId:  orgId,
		UserId: userId,
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

func (crd *SoftwareDownload) GetServiceAddr() (addr string, err error) {
	dictIds := []string{
		enum_dictionary.AgentHttpUrl,
		enum_dictionary.AgentHttpPort,
	}
	mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if string_helper.IsEmpty(mapDict[enum_dictionary.AgentHttpUrl]) || string_helper.IsEmpty(mapDict[enum_dictionary.AgentHttpPort]) {
		return "", errors.New("没有配置终端访问地址")
	}

	addr = strings.Trim(mapDict[enum_dictionary.AgentHttpUrl], "/") + ":" + mapDict[enum_dictionary.AgentHttpPort]
	return addr, nil
}

// GetDownloadInfo 获取下载信息
func (crd *SoftwareDownload) GetDownloadInfo(downloadModel *software_models.DownloadModel) error {
	addr, err := crd.GetServiceAddr()
	if err != nil {
		return err
	}

	//验证是否存在下载记录, 存在直接返回
	downloadModel, err = crd.masterModel.GetItemByUnique(downloadModel.OrgId, downloadModel.SoftwareId, downloadModel.CreateId, addr)
	if err != nil {
		return errors.New("下载错误")
	}
	if downloadModel.Id > 0 {
		return nil
	}

	//验证软件
	softModel := CreateCurdFactory().GetItemById(downloadModel.SoftwareId)
	if softModel == nil || softModel.Id == 0 {
		return errors.New("数据不存在")
	}
	if !files.ExistsFile(softModel.FullPath) {
		return errors.New("文件不存在")
	}

	return nil
}
