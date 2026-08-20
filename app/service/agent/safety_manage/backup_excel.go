package safety_manage

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/agent/trust_dir_models"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/excel_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateBackupExcelFactory(orgId int64) *BackupExcel {
	return &BackupExcel{
		orgId:     orgId,
		sheetName: "Sheet1",
		colsLen:   2,
	}
}

type BackupExcel struct {
	orgId     int64
	sheetName string
	colsLen   int
	rowsData  [][]interface{}
}

var typeTrustDirGlobal = "TrDiGl"
var typeProcGlobal = "PrGl"
var typeTamperGlobal = "TaGl"
var typeBlackmailGlobal = "BlGl"
var typeBlackmailGroup = "BlGr"
var typePortVirtualGlobal = "PoViGl"

// var typePlugDynamic = "PlDy"
var typePlugIpGlobal = "PlIpGl"
var typeUsbGlobal = "UsGl"
var typeOutreachGlobal = "OuGl"
var typeOutreachDetect = "OuDe"

func (crd *BackupExcel) generateExportPath() string {

	fullDir, _ := upload_file.GenerateYearMonthPath(upload_file.GetDefaultPath())
	name := md5_encrypt.MD5(fmt.Sprintf("%d_%d", crd.orgId, variable.SnowFlake.GetId()))

	return fullDir + name + ".xlsx"
}

// 填充信任目录全局数据
func (crd *BackupExcel) fillTrustDirGlobal() {
	condition := trust_dir_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := trust_dir_models.CreateGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeTrustDirGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式信任目录全局数据
func (crd *BackupExcel) formatTrustDirGlobal(jsonStr string, userId int64, userName string, now string) trust_dir_models.GlobalModel {
	var tmpModel trust_dir_models.GlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.OpenSwitch = enums.SwitchOn
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充可信防护全局全局数据
func (crd *BackupExcel) fillProcGlobal() {
	condition := proc_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := proc_models.CreateProcGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeProcGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式可信防护全局数据
func (crd *BackupExcel) formatProcGlobal(jsonStr string, userId int64, userName string, now string) proc_models.ProcGlobalModel {
	var tmpModel proc_models.ProcGlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && !string_helper.IsEmpty(tmpModel.Id) {
		tmpModel.OrgId = crd.orgId
		//tmpModel.OpenSwitch = enums.SwitchOn
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充防篡改全局数据
func (crd *BackupExcel) fillTamperGlobal() {
	condition := tamper_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := tamper_models.CreateGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeTamperGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式防篡改全局数据
func (crd *BackupExcel) formatTamperGlobal(jsonStr string, userId int64, userName string, now string) tamper_models.GlobalModel {
	var tmpModel tamper_models.GlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.OpenSwitch = enums.SwitchOff
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充勒索全局数据
func (crd *BackupExcel) fillBlackmailGlobal() {
	condition := blackmail_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := blackmail_models.CreateGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeBlackmailGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式勒索全局数据
func (crd *BackupExcel) formatBlackmailGlobal(jsonStr string, userId int64, userName string, now string) blackmail_models.GlobalModel {
	var tmpModel blackmail_models.GlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.OpenSwitch = enums.SwitchOn
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充勒索分组数据
func (crd *BackupExcel) fillBlackmailGroup() {
	condition := blackmail_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := blackmail_models.CreateGroupFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeBlackmailGroup,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式勒索分组数据
func (crd *BackupExcel) formatBlackmailGroup(jsonStr string, userId int64, userName string, now string) blackmail_models.GroupModel {
	var tmpModel blackmail_models.GroupModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充端口虚开全局数据
func (crd *BackupExcel) fillPortVirtualGlobal() {
	condition := port_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := port_models.CreateVirtualGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typePortVirtualGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式端口虚开全局数据
func (crd *BackupExcel) formatPortVirtualGlobal(jsonStr string, userId int64, userName string, now string) port_models.VirtualGlobalModel {
	var tmpModel port_models.VirtualGlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.OpenSwitch = enums.SwitchOn
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充封堵IP全局数据
func (crd *BackupExcel) fillPlugIpGlobal() {
	condition := plug_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := plug_models.CreateIpGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typePlugIpGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式封堵IP全局数据
func (crd *BackupExcel) formatPlugIpGlobal(jsonStr string, userId int64, userName string, now string) plug_models.IpGlobalModel {
	var tmpModel plug_models.IpGlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.OpenSwitch = enums.SwitchOn
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
		if string_helper.IsEmpty(tmpModel.EffectStartAt) {
			tmpModel.EffectStartAt = "null"
		}
		if string_helper.IsEmpty(tmpModel.EffectEndAt) {
			tmpModel.EffectEndAt = "null"
		}
	}
	return tmpModel
}

// 填充外设全局数据
func (crd *BackupExcel) fillUsbGlobal() {
	condition := usb_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := usb_models.CreateGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeUsbGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式外设全局数据
func (crd *BackupExcel) formatUsbGlobal(jsonStr string, userId int64, userName string, now string) usb_models.GlobalModel {
	var tmpModel usb_models.GlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && !string_helper.IsEmpty(tmpModel.Id) {
		tmpModel.OrgId = crd.orgId
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充外联全局数据
func (crd *BackupExcel) fillOutreachGlobal() {
	condition := outreach_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := outreach_models.CreateGlobalFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeOutreachGlobal,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式外联全局数据
func (crd *BackupExcel) formatOutreachGlobal(jsonStr string, userId int64, userName string, now string) outreach_models.GlobalModel {
	var tmpModel outreach_models.GlobalModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.OpenSwitch = enums.SwitchOn
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// 填充外联探测数据
func (crd *BackupExcel) fillOutreachDetect() {
	condition := outreach_models.SearchConditions{
		OrgId: crd.orgId,
	}
	models := outreach_models.CreateDetectFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		for _, model := range models {
			model.OrgId = 0
			crd.rowsData = append(crd.rowsData, []interface{}{
				typeOutreachDetect,
				json_helper.Encode(model),
			})
		}
	}
}

// 格式外联探测数据
func (crd *BackupExcel) formatOutreachDetect(jsonStr string, userId int64, userName string, now string) outreach_models.DetectModel {
	var tmpModel outreach_models.DetectModel
	if err := json_helper.ToAny(jsonStr, &tmpModel); err == nil && tmpModel.Id > 0 {
		tmpModel.OrgId = crd.orgId
		tmpModel.OpenSwitch = enums.SwitchOn
		tmpModel.UpdateId = userId
		tmpModel.UpdateName = userName
		tmpModel.UpdatedAt = now
		tmpModel.CreateId = userId
		tmpModel.CreateName = userName
		tmpModel.UpdatedAt = now
	}
	return tmpModel
}

// Export 导出
func (crd *BackupExcel) Export() (string, error) {
	//创建新 Excel 文件
	excelHandler := excel_helper.NewExcelHandler()
	defer excelHandler.Close()

	//行数据
	headerData := []interface{}{"Type", "Data"}

	//信任目录
	crd.fillTrustDirGlobal()
	//可信防护数据
	crd.fillProcGlobal()
	//防篡改数据
	crd.fillTamperGlobal()
	//勒索数据
	crd.fillBlackmailGroup()
	crd.fillBlackmailGlobal()
	//端口虚开数据
	crd.fillPortVirtualGlobal()
	//封堵数据
	crd.fillPlugIpGlobal()
	//外设数据
	crd.fillUsbGlobal()
	//外联数据
	crd.fillOutreachGlobal()
	crd.fillOutreachDetect()

	//流式写入sheet
	if err := excelHandler.StreamWriteRows(crd.sheetName, headerData, crd.rowsData, 1, 1); err != nil {
		return "", fmt.Errorf("StreamWriteRows fail: %v", err)
	}

	// 保存文件
	path := crd.generateExportPath()
	if err := excelHandler.SaveAs(path); err != nil {
		return "", fmt.Errorf("save path fail: %v", err)
	}

	return path, nil
}

// Import 导入
func (crd *BackupExcel) Import(fullPath string, userId int64, userName string) error {
	//创建新 Excel 文件
	excelHandler, err := excel_helper.NewExcelHandlerFromFile(fullPath)
	if err != nil {
		return fmt.Errorf("read file: %v", err)
	}
	defer excelHandler.Close()

	//获取所有行数据
	rowsValues, err2 := excelHandler.GetRowsValues(crd.sheetName)
	if err2 != nil {
		return fmt.Errorf("get rows: %v", err)
	}
	if len(rowsValues) < 2 {
		return nil
	}
	now := date_helper.CurDatetime()
	var insTamperGlobal []tamper_models.GlobalModel
	var insProcGlobal []proc_models.ProcGlobalModel
	var insTrustDirGlobal []trust_dir_models.GlobalModel
	var insBlackmailGroup []blackmail_models.GroupModel
	var insBlackmailGlobal []blackmail_models.GlobalModel
	var insPortVirtualGlobal []port_models.VirtualGlobalModel
	var insPlugIpGlobal []plug_models.IpGlobalModel
	var insUsbGlobal []usb_models.GlobalModel
	var insOutreachGlobal []outreach_models.GlobalModel
	var insOutreachDetect []outreach_models.DetectModel

	//处理行数据
	for i, rowVals := range rowsValues {
		//忽略前1行
		if i < 1 {
			continue
		}
		//补足
		if colRealLen := len(rowVals); colRealLen < crd.colsLen {
			tmp := make([]string, crd.colsLen-colRealLen)
			rowVals = append(rowVals, tmp...)
		}
		switch rowVals[0] {
		case typeTrustDirGlobal: //信任目录全局数据
			if tmpModel := crd.formatTrustDirGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insTrustDirGlobal = append(insTrustDirGlobal, tmpModel)
			}
		case typeProcGlobal: //可信防护全局数据
			if tmpModel := crd.formatProcGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insProcGlobal = append(insProcGlobal, tmpModel)
			}
		case typeTamperGlobal: //防篡改全局数据
			if tmpModel := crd.formatTamperGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insTamperGlobal = append(insTamperGlobal, tmpModel)
			}
		case typeBlackmailGroup: //勒索分组数据
			if tmpModel := crd.formatBlackmailGroup(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insBlackmailGroup = append(insBlackmailGroup, tmpModel)
			}
		case typeBlackmailGlobal: //勒索全局数据
			if tmpModel := crd.formatBlackmailGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insBlackmailGlobal = append(insBlackmailGlobal, tmpModel)
			}
		case typePortVirtualGlobal: //端口虚开全局数据
			if tmpModel := crd.formatPortVirtualGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insPortVirtualGlobal = append(insPortVirtualGlobal, tmpModel)
			}
		case typePlugIpGlobal: //封堵ip全局数据
			if tmpModel := crd.formatPlugIpGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insPlugIpGlobal = append(insPlugIpGlobal, tmpModel)
			}
		case typeUsbGlobal: //外设全局数据
			if tmpModel := crd.formatUsbGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insUsbGlobal = append(insUsbGlobal, tmpModel)
			}
		case typeOutreachGlobal: //外联全局数据
			if tmpModel := crd.formatOutreachGlobal(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insOutreachGlobal = append(insOutreachGlobal, tmpModel)
			}
		case typeOutreachDetect: //外联探测名单
			if tmpModel := crd.formatOutreachDetect(rowVals[1], userId, userName, now); tmpModel.OrgId > 0 {
				insOutreachDetect = append(insOutreachDetect, tmpModel)
			}
		}
	}

	if len(insTrustDirGlobal) > 0 {
		trust_dir_models.CreateGlobalFactory().InsertBatch(&insTrustDirGlobal)
	}

	if len(insProcGlobal) > 0 {
		conditions := proc_models.SearchConditions{OrgId: crd.orgId}
		proc_models.CreateProcessFactory().DeleteByConditions(conditions)
		proc_models.CreateDllFactory().DeleteByConditions(conditions)
		proc_models.CreateProcGlobalFactory().InsertBatch(&insProcGlobal)
	}

	if len(insTamperGlobal) > 0 {
		tamper_models.CreateGlobalFactory().InsertBatch(&insTamperGlobal)
	}

	if len(insBlackmailGroup) > 0 {
		blackmail_models.CreateGroupFactory().InsertBatch(&insBlackmailGroup)
	}
	if len(insBlackmailGlobal) > 0 {
		blackmail_models.CreateGlobalFactory().InsertBatch(&insBlackmailGlobal)
	}

	if len(insPortVirtualGlobal) > 0 {
		port_models.CreateVirtualGlobalFactory().InsertBatch(&insPortVirtualGlobal)
	}

	if len(insPlugIpGlobal) > 0 {
		plug_models.CreateIpGlobalFactory().InsertBatch(&insPlugIpGlobal)
	}

	if len(insUsbGlobal) > 0 {
		usb_models.CreateGlobalFactory().InsertBatch(&insUsbGlobal)
	}

	if len(insOutreachGlobal) > 0 {
		outreach_models.CreateGlobalFactory().InsertBatch(&insOutreachGlobal)
	}
	if len(insOutreachDetect) > 0 {
		outreach_models.CreateDetectFactory().InsertBatch(&insOutreachDetect)
	}
	return nil
}
