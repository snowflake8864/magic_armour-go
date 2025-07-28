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
func (crd *BackupExcel) appendTrustDirGlobal(jsonStr string, userId int64, userName string, now string) trust_dir_models.GlobalModel {
	var tmpModel trust_dir_models.GlobalModel
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
func (crd *BackupExcel) appendTamperGlobal(jsonStr string, userId int64, userName string, now string) tamper_models.GlobalModel {
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

// Export 导出
func (crd *BackupExcel) Export() (string, error) {
	//创建新 Excel 文件
	excelHandler := excel_helper.NewExcelHandler()
	defer excelHandler.Close()

	//行数据
	headerData := []interface{}{"类型", "数据"}

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
	var insTrustDirGlobal []trust_dir_models.GlobalModel
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
			if tmpModel := crd.appendTrustDirGlobal(rowVals[1], userId, userName, now); tmpModel.Id > 0 {
				insTrustDirGlobal = append(insTrustDirGlobal, tmpModel)
			}
		case typeTamperGlobal: //防篡改全局数据
			if tmpModel := crd.appendTamperGlobal(rowVals[1], userId, userName, now); tmpModel.Id > 0 {
				insTamperGlobal = append(insTamperGlobal, tmpModel)
			}
		}
	}

	if len(insTrustDirGlobal) > 0 {
		trust_dir_models.CreateGlobalFactory().InsertBatch(&insTrustDirGlobal)
	}
	if len(insTamperGlobal) > 0 {
		tamper_models.CreateGlobalFactory().InsertBatch(&insTamperGlobal)
	}

	return nil
}
