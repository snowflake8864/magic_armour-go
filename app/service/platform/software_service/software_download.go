package software_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/file_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/zip_factory"
	"os"
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
		SoftId: softId,
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

// GetServiceAddr 获取服务地址
// @return err 错误信息, addr 完整地址, addr 域名, port 端口号
func (crd *SoftwareDownload) GetServiceAddr() (err error, addr, url, port string) {
	dictIds := []string{
		enum_dictionary.AgentHttpUrl,
		enum_dictionary.AgentHttpPort,
	}
	mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if string_helper.IsEmpty(mapDict[enum_dictionary.AgentHttpUrl]) || string_helper.IsEmpty(mapDict[enum_dictionary.AgentHttpPort]) {
		return errors.New("没有配置终端访问地址"), addr, url, port
	}

	url = strings.Trim(mapDict[enum_dictionary.AgentHttpUrl], "/")
	port = strings.Trim(mapDict[enum_dictionary.AgentHttpPort], " ")
	addr = url + ":" + port
	return nil, addr, url, port
}

// GetDownloadInfo 获取下载信息
func (crd *SoftwareDownload) GetDownloadInfo(param *software_models.DownloadModel) (error, *software_models.DownloadModel) {
	//获取服务器地址
	err, serviceIp, agentHttpUrl, agentHttpPort := crd.GetServiceAddr()
	if err != nil {
		return err, nil
	}
	param.ServiceIp = serviceIp

	//验证是否存在下载记录, 存在直接返回
	if downloadModel, _ := crd.masterModel.GetItemByUnique(param.OrgId, param.SoftwareId, param.CreateId, param.ServiceIp); downloadModel.Id > 0 {
		return nil, downloadModel
	}

	//验证软件
	softModel := CreateCurdFactory().GetItemById(param.SoftwareId)
	if softModel == nil || softModel.Id == 0 {
		return errors.New("数据不存在"), nil
	}
	if !file_helper.ExistsFile(softModel.FullPath) {
		return errors.New("文件不存在"), nil
	}

	//新文件名称
	newFileName := crd.generateFileName(param.OrgId, param.SoftwareId, param.CreateId, param.ServiceIp)
	//创建路径
	newStorageDir, newFullDir := crd.generateYearMonthPath()
	param.StoragePath = newStorageDir + newFileName
	param.FullPath = newFullDir + newFileName

	//复制到新路径
	if err = file_helper.CopyFile(softModel.FullPath, param.FullPath); err != nil {
		return errors.New("复制ZIP失败"), nil
	}

	//追加配置文件
	if !crd.appendConfig(agentHttpUrl, agentHttpPort, softModel, param) {
		return errors.New("追加INI失败"), nil
	}

	//新文件属性
	_, fileHash := file_helper.Md5sum(param.FullPath)
	param.FileHash = fileHash
	fileInfo, _ := file_helper.FileInfo(param.FullPath)
	param.FileSize = fileInfo.Size()
	param.CreatedAt = date_helper.CurDatetime()
	param.Alias = enum_software.DownloadName(softModel.PlatformType, softModel.OsType, softModel.Ver)

	//入库
	crd.masterModel.InsertItem(param)

	return nil, param
}

// 追加配置文件
func (crd *SoftwareDownload) appendConfig(agentHttpUrl, agentHttpPort string, softModel *software_models.SoftwareModel, downloadModel *software_models.DownloadModel) bool {

	var params []string
	params = append(params, "[SERVER]")
	params = append(params, "URL="+agentHttpUrl)
	params = append(params, "PORT="+agentHttpPort)
	params = append(params, "USER_ID="+crd.EncryptConfig(downloadModel.OrgId, downloadModel.CreateId, downloadModel.SoftwareId))

	//Unix系统里，每行结尾只有“<换行>”，即"\n"；Windows系统里面，每行结尾是“<回车><换行>”，即“\r\n”；
	glue := "\n"
	if softModel.PlatformType == enum_software.PlatformTypeWindows {
		glue = "\r\n"
	}
	content := strings.Join(params, glue)

	if err := zip_factory.AddFileToZip(softModel.FullPath, downloadModel.FullPath, "config.ini", content); err != nil {
		variable.ZapLog.Sugar().Error(err.Error())
		return false
	}

	return true
}

// 生成文件名
func (crd *SoftwareDownload) generateFileName(orgId, softwareId, userId int64, serviceIp string) string {
	str := "%d_%d_%d_%s_%s_%s"

	str = md5_encrypt.MD5(fmt.Sprintf(str, orgId, softwareId, userId, serviceIp, date_helper.CurDatetime(), int_helper.IntRandom2Str(1000, 9999)))
	str += ".zip"
	return str
}

// 文件上传可以设置按照 xxx年-xx月 格式存储
// @return 项目路径, 完整路径
func (crd *SoftwareDownload) generateYearMonthPath() (string, string) {
	basePath := variable.BasePath
	savePathPre := variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath")
	returnPathPre := variable.ConfigYml.GetString("FileUploadSetting.UploadFileReturnPath")
	curYearMonth := date_helper.SpecialDatetime("2006_01")

	newSavePathPre := savePathPre + curYearMonth
	newReturnPathPre := basePath + returnPathPre + curYearMonth
	// 相关路径不存在，创建目录
	if _, err := os.Stat(newReturnPathPre); err != nil {
		if err = os.MkdirAll(newReturnPathPre, os.ModePerm); err != nil {
			variable.ZapLog.Error("文件上传创建目录出错" + err.Error())
			return "", ""
		}
	}
	return newSavePathPre + "/", newReturnPathPre + "/"
}
