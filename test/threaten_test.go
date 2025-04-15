package test

import (
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/jobs/threaten_job"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/service/veda_other/threaten_platform"
	"goskeleton/app/utils/helpers/json_helper"
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

func TestThreatenRepairProcForCache(t *testing.T) {
	ids := []string{
		"857007c337953c7ea3682e5b159222dd", "912a6e1b7492a9bf32ca6cf4f5982ba0", "aa94e5e3b47af0b801cea20185dc735e",
		"ff65a732903e85f367a64838e54dcd9a", "3c103508bdbc8eaa2a322b76a41686c6", "797616cc1dbfe20c38c3ff59d873d5f4",
		"001786fe94ec78ebaf6206cc066143bf", "20d31264c2a8d891076b4b4fd808f871", "2ce82bd239a502e5317d78978fce3671",
		"fd2208f6280e43fc80c4157cb368d1ba", "422071881bd768642aa94e2961929773", "d37fa1ddf0f1f8c4c7e5583fc3c729ac",
		"2d1322e5ce4586980675602c6e73bb01", "6a8019f7eb33d32c0063b5c7ac2bb657", "2f8dab5347012542de116ada7c762d19",
		"9e8ba45b07f55c4c08c9ba76e34f5ed7", "cb10057e0b7fe1c25fc2c88e04faddef", "3529b35c68acd034bc2dbe901db90ed1",
		"061d72516de2efc2cd49729b365d25d7", "e32b2ac3e643b99d99596898cec53114", "0cb8fb1bf59ad48cc5cffd00145d649e",
		"fa54cdb4a74ea60138dd68fb6e6f539d", "daee0684e83ecd6f03172d480b5bf56a", "f1ac5fcdf2a974e81edf0b14a9f648ac",
		"2f704fda1649614fceffb75baa12bdd5", "4ced7d7af81ea78efb25c8bf23f9fffc", "1731955e6d867f9b386980727227bf4d",
		"34101244442027b4541e57cf78e5f7a6", "147340be9f8c3ffaf1df9dc62fd9f31b", "fd8fa2a2fb19d2439b5e64a8c77b5e60",
		"d3365c4a2e25e5488b1981284bddc470", "7cab268a3a4194aef12bb6789f327e27", "882edc048eb7000045cee806a0d480b2",
		"a4c523b3320ffda27bb2353adc4729e1", "31d2edaf400e9a53a95d23e282f3f765", "138be4f918828d5f07097355aed511e2",
		"841fef896da6602dc0a520e6f09b2f0c", "369252441d6b354eb013ac688489aa35", "c0fb27b2c758e77f4a8bab4c6c519cc3",
		"7a2f0e784489dc1ff58106fe7dc7de13", "877ba8dc2344881e331332dfb0ba63ae", "f174a5f15e762f5d77fdedadb9e4505f",
		"91e044d2a078ac2f0897eab0841a3b97", "cb84f7a25a1c383fa7e772ae73c55f7e", "7ae1624bc3d7424e5b37713d2daf8ef1",
		"8f16c44259859a0f3401760346386384", "791cd3c43b91398811e5fc106fdbdd15", "24e854179ca5b82cc75ed5579c643acc",
		"02352e70ae6c88095911a6df8e9b3ccf", "2451173ff2f99db896c1cd6c69b93eab", "cf7e118b1deeca687c8fe6984b5f7dd7",
		"b30e285ceb0b5b709c66a48cf8cf9be8", "492aab8dd394f7ad848f12860d2cc4e5", "79e987568f1ccccf993ee05834a87e12",
		"319db6b5b14de91a6b94b6f4ad6351fc", "b60d6df741bd0debab70404c4de0e89e", "505d349264b832496448d1333f9ce21b",
		"02a7c69bfe5784fcd82e7f49ac1d3d50", "5eddac67d28ca992ef5e727fa34c1b28", "98aa9a2232238603b9b5bbcd91a32afe"}
	threaten_job.RepairProc(2, 1, json_helper.Encode(ids))
}

func TestThreatenRewrite(t *testing.T) {
	param := threaten_models.CacheModel{
		Id:       "ffffff",
		Level:    enum_threaten.LevelWhite,
		DataFrom: enum_threaten.DataFromDll,
	}

	threaten_job.HandelRewire(json_helper.Encode(param))
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
func TestThreatenBatchAdd(t *testing.T) {
	var list []threaten_models.CacheModel

	list = append(list,
		threaten_models.CacheModel{
			Id:        "0023f8f263615f53ff29795c1b194459",
			Source:    enum_threaten.SourceAgentUpload,
			Dir:       "ggg",
			ParentDir: "",
			Os:        "CentOS\"7 (core)\"_kernel:3.10.0-1160.119.1.el7.x86_64",
			DataFrom:  enum_threaten.DataFromProcess,
		}, threaten_models.CacheModel{
			Id:        "7e20c61405e25702e08b83478d0ef76f",
			Source:    enum_threaten.SourceAgentUpload,
			Dir:       "ggg",
			ParentDir: "",
			Os:        "CentOS\"7 (core)\"_kernel:3.10.0-1160.119.1.el7.x86_64",
			DataFrom:  enum_threaten.DataFromProcess,
		})
	threaten_service.CreateCacheFactory().BatchAdd(&list, false)
}
