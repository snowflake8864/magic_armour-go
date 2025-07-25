package enum_upload

const (
	TypeSoftware        string = "software"         //类型:终端软件
	TypeUpgrade         string = "upgrade"          //类型:系统升级
	TypeLicense         string = "license"          //类型:授权证书
	TypeProcSample      string = "proc_sample"      //类型:可信防护-样本
	TypeAgentError      string = "agent_error"      //类型:终端错误
	TypeTamperPolicy    string = "tamper_policy"    //类型: 防篡改策略
	TypeThreatenCaches  string = "threaten_caches"  //类型: 威胁情报缓存
	TypeImportantBackup string = "important_backup" //类型: 重要数据备份
)
