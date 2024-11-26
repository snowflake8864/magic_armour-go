package order_service

import (
	"bytes"
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/file_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"os/exec"
	"strings"
)

func CreateLicenseFactory() *License {
	return &License{}
}

type License struct {
}

type LicenseResp struct {
	Serial      string                          `json:"serial"`
	Type        int                             `json:"type"`
	StartAt     string                          `json:"startAt"`
	ExpireAt    string                          `json:"expireAt"`
	Duration    int                             `json:"duration"`
	Number      int                             `json:"number"`
	Modules     enum_product.ProductGrantConfig `json:"modules"`
	ModuleLabel []map[string]string             `json:"moduleLabel"`
}

// Parse 解析
func (crd *License) Parse(path string) (error, *LicenseResp) {
	if !file_helper.ExistsFile(path) {
		return errors.New("授权文件不存在"), nil
	}

	shellVedaLicence := cur_app.ShellVedaLicence()
	cmd := exec.Command(shellVedaLicence, "2", path)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	if err != nil {
		return errors.New("exec command Run failed: " + err.Error()), nil
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if string_helper.IsEmpty(outStr) || !string_helper.IsEmpty(errStr) {
		return errors.New("exec command output failed： " + errStr), nil
	}
	if string_helper.IsEmpty(outStr) {
		return errors.New("exec command output empty"), nil
	}
	/*	outStr = `
		load_license success
		auth_factory:5Lit6IO96J6N5ZCI5pm65oWn56eR5oqA5pyJ6ZmQ5YWs5Y+4
		copy ascription:5Lit6IO96J6N5ZCI5pm65oWn56eR5oqA5pyJ6ZmQ5YWs5Y+4
		client name:5YyX5Lqs6Ieq5p2l5rC0
		auth type:0
		hd type:UFJBLU5CLUI3NS1DLTA4LTUtNQ==
		language:chinese
		hd_iden:43c006e91799cd13d4968920
		work mode:2
		expire date:1734710399
		host num:500
		auth time:1703088000
		serial:43c006e91799cd13d4968920
		status:1
		alive time:2102089
		module:1|2|3|2|3|1
		reserved:
		product_id:1
		`*/
	parseRes := make(map[string]string)
	outList := strings.Split(outStr, "\n")
	for _, line := range outList {
		arr := strings.Split(line, ":")

		if len(arr) == 2 {
			key := strings.ReplaceAll(arr[0], " ", "_")
			parseRes[key] = strings.TrimSpace(arr[1])
		}
	}

	//证书结果
	licenseResp := &LicenseResp{}
	if val, isOK := parseRes["serial"]; isOK && !string_helper.IsEmpty(val) {
		licenseResp.Serial = val
	} else {
		return errors.New("授权文件缺少设备标识码"), nil
	}

	if val, isOK := parseRes["auth_type"]; isOK && !string_helper.IsEmpty(val) {
		//0:正式, 1:测试, 2样本
		licenseResp.Type = string_helper.ToInt(val) + 1
		if label := enums.Get1stIsStrAnd2ndIsIntReturnStr("type", licenseResp.Type, enum_product.EnumIntMap); string_helper.IsEmpty(label) {
			return errors.New("未知的授权类型: " + val), nil
		}
	} else {
		return errors.New("授权文件缺少授权类型"), nil
	}

	if val, isOK := parseRes["host_num"]; isOK && !string_helper.IsEmpty(val) {
		licenseResp.Number = string_helper.ToInt(val)
		if licenseResp.Number <= 0 {
			return errors.New("授权文件终端数目不正确"), nil
		}
	} else {
		return errors.New("授权文件缺少终端数目"), nil
	}

	if val, isOK := parseRes["alive_time"]; isOK && !string_helper.IsEmpty(val) {
		timestamp := string_helper.ToInt64(val)
		if timestamp <= 0 {
			return errors.New("授权文件已过期"), nil
		}
	} else {
		return errors.New("授权文件缺少有效时间"), nil
	}

	if val, isOK := parseRes["auth_time"]; isOK && !string_helper.IsEmpty(val) {
		timestamp := string_helper.ToInt64(val)
		if timestamp == 0 {
			return errors.New("授权开始时间错误"), nil
		}
		licenseResp.StartAt = date_helper.TimestampToDatetimeByFormat(timestamp, variable.DayFormat)
	} else {
		return errors.New("授权文件缺少开始时间"), nil
	}

	if val, isOK := parseRes["expire_date"]; isOK && !string_helper.IsEmpty(val) {
		timestamp := string_helper.ToInt64(val)
		if timestamp == 0 {
			return errors.New("授权过期时间错误"), nil
		}
		licenseResp.ExpireAt = date_helper.TimestampToDatetimeByFormat(timestamp, variable.DayFormat)
	} else {
		return errors.New("授权文件缺少过期时间"), nil
	}

	//计算授权时长
	startAtArr := strings.Split(licenseResp.StartAt, "-")
	expireAtArr := strings.Split(licenseResp.ExpireAt, "-")
	startAtYear := string_helper.ToInt(startAtArr[0])
	startAtMonth := string_helper.ToInt(startAtArr[1])
	expireAtYear := string_helper.ToInt(expireAtArr[0])
	expireAtMonth := string_helper.ToInt(expireAtArr[1])
	if diffYear := expireAtYear - startAtYear; diffYear <= 0 {
		if diffMonth := expireAtMonth - startAtMonth; diffMonth <= 0 {
			licenseResp.Duration = 1
		} else {
			licenseResp.Duration = diffMonth
		}
	} else {
		licenseResp.Duration = diffYear * 12
	}

	if val, isOK := parseRes["module"]; isOK && !string_helper.IsEmpty(val) {
		//@todo 模块顺序: 可信防护、勒索、防篡改、端口虚开、动态封堵、外设
		moduleList := strings.Split(val, "|")
		//补足模块
		moduleList = append(moduleList, "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0")

		licenseResp.Modules.Process = string_helper.ToInt(moduleList[0])
		licenseResp.Modules.Blackmail = string_helper.ToInt(moduleList[1])
		licenseResp.Modules.Tamper = string_helper.ToInt(moduleList[2])
		licenseResp.Modules.Port = string_helper.ToInt(moduleList[3])
		licenseResp.Modules.Plug = string_helper.ToInt(moduleList[4])
		licenseResp.Modules.Usb = string_helper.ToInt(moduleList[5])
		//其他模块暂时与可信防护数目一致
		licenseResp.Modules.Syslog = licenseResp.Modules.Process
		licenseResp.Modules.Outreach = licenseResp.Modules.Process
		licenseResp.Modules.Login = licenseResp.Modules.Process
		licenseResp.Modules.Threaten = licenseResp.Modules.Process
		//展示
		licenseResp.ModuleLabel = append(licenseResp.ModuleLabel,
			map[string]string{
				"label": enums.Get1stIsStrAnd2ndIsStrReturnStr("module", enum_product.ModuleProcess, enum_product.EnumStrMap),
				"val":   moduleList[0],
			}, map[string]string{
				"label": enums.Get1stIsStrAnd2ndIsStrReturnStr("module", enum_product.ModuleBlackmail, enum_product.EnumStrMap),
				"val":   moduleList[1],
			}, map[string]string{
				"label": enums.Get1stIsStrAnd2ndIsStrReturnStr("module", enum_product.ModuleTamper, enum_product.EnumStrMap),
				"val":   moduleList[2],
			}, map[string]string{
				"label": enums.Get1stIsStrAnd2ndIsStrReturnStr("module", enum_product.ModulePort, enum_product.EnumStrMap),
				"val":   moduleList[3],
			}, map[string]string{
				"label": enums.Get1stIsStrAnd2ndIsStrReturnStr("module", enum_product.ModulePlug, enum_product.EnumStrMap),
				"val":   moduleList[4],
			}, map[string]string{
				"label": enums.Get1stIsStrAnd2ndIsStrReturnStr("module", enum_product.ModuleUsb, enum_product.EnumStrMap),
				"val":   moduleList[5],
			})

	} else {
		return errors.New("授权文件缺少模块数量"), nil
	}

	return nil, licenseResp
}
