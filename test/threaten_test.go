package test

import (
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/service/veda_other/threaten_platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console/alarm_cmd"
	"goskeleton/command/console/threaten_cmd"
	"testing"
)

func TestAlarmStatistics(t *testing.T) {
	alarm_cmd.CreateCmdFactory().YesterdayStatistics()
	alarm_cmd.CreateCmdFactory().TodayStatistics()
}

func TestThreatenAddFromProc(t *testing.T) {
	threaten_cmd.CreateCmdFactory().AddFromProcForCache()
}

func TestThreatenSyncPushes(t *testing.T) {
	threaten_cmd.CreateCmdFactory().SyncForPushes()
}

func TestThreatenSyncCache(t *testing.T) {
	threaten_cmd.CreateCmdFactory().SyncDefaultForCache()
}

func TestThreatenRewrite(t *testing.T) {
	threaten_service.CreateCacheFactory().TestRewriteAfter("4c7b5857db3f1c52a07d5feda1370107", 1, 1)
}

func TestThreatenCallbackByHash(t *testing.T) {

	proc_service.CreateProcessCurdFactory().ThreatenCallbackByHash("ffffff", 11)
}

func TestThreatenCheckData(t *testing.T) {

	str1 := "{\"code\":1,\"msg\":\"success\",\"data\":[],\"appId\":\"bff1f0668a0c17fc10a68d5c4c517169\"}"
	threaten_platform.CreateRequestFactory(threaten_platform.RequestConfig{}).TestCheckData(str1)

	str2 := `
{"code":1,"msg":"success","data":{"6D220604AA4240303DD8DEAEAB428377":{"name":"6D220604AA4240303DD8DEAEAB428377","type":1,"source":4,"filename":"C:\\Windows\\System32\\ie4uinit.exe","dir":"C:\\Windows\\System32\\ie4uinit.exe","os":"Microsoft Windows Server 2008 "},"685C6E1854C51ACB440F332B46B4AA8A":{"name":"685C6E1854C51ACB440F332B46B4AA8A","type":1,"source":4,"attr":"","filename":"J:\\Windows\\winsxs\\amd64_microsoft-windows-ie-ielowutil_31bf3856ad364e35_11.2.9600.18893_none_e8a9573","dir":"J:\\Windows\\winsxs\\amd64_microsoft-windows-ie-ielowutil_31bf3856ad364e35_11.2.9600.18893_none_e8a9573e8662871f\\ielowutil.exe","os":"Microsoft Windows 11 \\u5bb6\\u5ead\\u7248"},"6A3D3F818FAAD570E026ECCDC6E6D2D8":{"name":"6A3D3F818FAAD570E026ECCDC6E6D2D8","type":1,"source":4,"attr":"","filename":"J:\\Windows\\winsxs\\amd64_microsoft-windows-s..-downlevel.binaries_31bf3856ad364e35_6.3.9600.18893_non","dir":"J:\\Windows\\winsxs\\amd64_microsoft-windows-s..-downlevel.binaries_31bf3856ad364e35_6.3.9600.18893_none_5f5eb1a8ff7ad292\\MsSpellCheckingFacility.exe","os":"Microsoft Windows 11 \\u5bb6\\u5ead\\u7248"}}}
`
	threaten_platform.CreateRequestFactory(threaten_platform.RequestConfig{}).TestCheckData(str2)

	str3 := `
{"code":1,"msg":"success","data":{"6D220604AA4240303DD8DEAEAB428377":{"name":"6D220604AA4240303DD8DEAEAB428377","type":1,"source":4,"attr":"{\"hash\":\"6d220604aa4240303dd8deaeab428377\",\"fileName\":\"ie4uinit.exe\",\"source\":\"network\",\"detections\":[{\"engineName\":\"Bkav\",\"category\":\"undetected\"},{\"engineName\":\"Lionic\",\"category\":\"undetected\"},{\"engineName\":\"Elastic\",\"category\":\"undetected\"},{\"engineName\":\"DrWeb\",\"category\":\"undetected\"},{\"engineName\":\"ClamAV\",\"category\":\"undetected\"},{\"engineName\":\"FireEye\",\"category\":\"undetected\"},{\"engineName\":\"CAT-QuickHeal\",\"category\":\"undetected\"},{\"engineName\":\"McAfee\",\"category\":\"undetected\"},{\"engineName\":\"Cylance\",\"category\":\"undetected\"},{\"engineName\":\"Zillya\",\"category\":\"undetected\"},{\"engineName\":\"Sangfor\",\"category\":\"undetected\"},{\"engineName\":\"K7AntiVirus\",\"category\":\"undetected\"},{\"engineName\":\"Alibaba\",\"category\":\"undetected\"},{\"engineName\":\"K7GW\",\"category\":\"undetected\"},{\"engineName\":\"Cybereason\",\"category\":\"undetected\"},{\"engineName\":\"BitDefenderTheta\",\"category\":\"undetected\"},{\"engineName\":\"VirIT\",\"category\":\"undetected\"},{\"engineName\":\"Cyren\",\"category\":\"undetected\"},{\"engineName\":\"SymantecMobileInsight\",\"category\":\"type-unsupported\"},{\"engineName\":\"Symantec\",\"category\":\"undetected\"},{\"engineName\":\"tehtris\",\"category\":\"undetected\"},{\"engineName\":\"ESET-NOD32\",\"category\":\"undetected\"},{\"engineName\":\"Zoner\",\"category\":\"undetected\"},{\"engineName\":\"TrendMicro-HouseCall\",\"category\":\"undetected\"},{\"engineName\":\"Paloalto\",\"category\":\"undetected\"},{\"engineName\":\"Cynet\",\"category\":\"undetected\"},{\"engineName\":\"Kaspersky\",\"category\":\"undetected\"},{\"engineName\":\"BitDefender\",\"category\":\"undetected\"},{\"engineName\":\"NANO-Antivirus\",\"category\":\"undetected\"},{\"engineName\":\"SUPERAntiSpyware\",\"category\":\"undetected\"},{\"engineName\":\"MicroWorld-eScan\",\"category\":\"undetected\"},{\"engineName\":\"Avast\",\"category\":\"undetected\"},{\"engineName\":\"Tencent\",\"category\":\"undetected\"},{\"engineName\":\"Ad-Aware\",\"category\":\"undetected\"},{\"engineName\":\"Trustlook\",\"category\":\"type-unsupported\"},{\"engineName\":\"Sophos\",\"category\":\"undetected\"},{\"engineName\":\"Comodo\",\"category\":\"undetected\"},{\"engineName\":\"F-Secure\",\"category\":\"undetected\"},{\"engineName\":\"Baidu\",\"category\":\"undetected\"},{\"engineName\":\"VIPRE\",\"category\":\"undetected\"},{\"engineName\":\"TrendMicro\",\"category\":\"undetected\"},{\"engineName\":\"McAfee-GW-Edition\",\"category\":\"timeout\"},{\"engineName\":\"Trapmine\",\"category\":\"undetected\"},{\"engineName\":\"CMC\",\"category\":\"undetected\"},{\"engineName\":\"Emsisoft\",\"category\":\"undetected\"},{\"engineName\":\"SentinelOne\",\"category\":\"undetected\"},{\"engineName\":\"GData\",\"category\":\"undetected\"},{\"engineName\":\"Jiangmin\",\"category\":\"undetected\"},{\"engineName\":\"Webroot\",\"category\":\"undetected\"},{\"engineName\":\"Avira\",\"category\":\"undetected\"},{\"engineName\":\"MAX\",\"category\":\"undetected\"},{\"engineName\":\"Antiy-AVL\",\"category\":\"undetected\"},{\"engineName\":\"Kingsoft\",\"category\":\"undetected\"},{\"engineName\":\"Gridinsoft\",\"category\":\"undetected\"},{\"engineName\":\"Arcabit\",\"category\":\"undetected\"},{\"engineName\":\"ViRobot\",\"category\":\"undetected\"},{\"engineName\":\"ZoneAlarm\",\"category\":\"undetected\"},{\"engineName\":\"Avast-Mobile\",\"category\":\"type-unsupported\"},{\"engineName\":\"Microsoft\",\"category\":\"undetected\"},{\"engineName\":\"Google\",\"category\":\"undetected\"},{\"engineName\":\"BitDefenderFalx\",\"category\":\"type-unsupported\"},{\"engineName\":\"AhnLab-V3\",\"category\":\"undetected\"},{\"engineName\":\"Acronis\",\"category\":\"undetected\"},{\"engineName\":\"VBA32\",\"category\":\"undetected\"},{\"engineName\":\"ALYac\",\"category\":\"undetected\"},{\"engineName\":\"TACHYON\",\"category\":\"undetected\"},{\"engineName\":\"Malwarebytes\",\"category\":\"undetected\"},{\"engineName\":\"APEX\",\"category\":\"undetected\"},{\"engineName\":\"Rising\",\"category\":\"undetected\"},{\"engineName\":\"Yandex\",\"category\":\"undetected\"},{\"engineName\":\"Ikarus\",\"category\":\"undetected\"},{\"engineName\":\"MaxSecure\",\"category\":\"undetected\"},{\"engineName\":\"Fortinet\",\"category\":\"undetected\"},{\"engineName\":\"AVG\",\"category\":\"undetected\"},{\"engineName\":\"Panda\",\"category\":\"undetected\"}],\"undetected\":70,\"malicious\":0}","filename":"C:\\Windows\\System32\\ie4uinit.exe","dir":"C:\\Windows\\System32\\ie4uinit.exe","os":"Microsoft Windows Server 2008 "},"685C6E1854C51ACB440F332B46B4AA8A":{"name":"685C6E1854C51ACB440F332B46B4AA8A","type":1,"source":4,"attr":"","filename":"J:\\Windows\\winsxs\\amd64_microsoft-windows-ie-ielowutil_31bf3856ad364e35_11.2.9600.18893_none_e8a9573","dir":"J:\\Windows\\winsxs\\amd64_microsoft-windows-ie-ielowutil_31bf3856ad364e35_11.2.9600.18893_none_e8a9573e8662871f\\ielowutil.exe","os":"Microsoft Windows 11 \\u5bb6\\u5ead\\u7248"},"6A3D3F818FAAD570E026ECCDC6E6D2D8":{"name":"6A3D3F818FAAD570E026ECCDC6E6D2D8","type":1,"source":4,"attr":"","filename":"J:\\Windows\\winsxs\\amd64_microsoft-windows-s..-downlevel.binaries_31bf3856ad364e35_6.3.9600.18893_non","dir":"J:\\Windows\\winsxs\\amd64_microsoft-windows-s..-downlevel.binaries_31bf3856ad364e35_6.3.9600.18893_none_5f5eb1a8ff7ad292\\MsSpellCheckingFacility.exe","os":"Microsoft Windows 11 \\u5bb6\\u5ead\\u7248"}}}
`
	threaten_platform.CreateRequestFactory(threaten_platform.RequestConfig{}).TestCheckData(str3)
}
