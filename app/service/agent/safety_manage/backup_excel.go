package safety_manage

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/excel_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateBackupExcelFactory(orgId int64) *BackupExcel {
	return &BackupExcel{
		tamperGlobal: tamper_models.CreateGlobalFactory(),
		orgId:        orgId,
		sheetName:    "Sheet1",
		colsLen:      2,
	}
}

type BackupExcel struct {
	tamperGlobal *tamper_models.GlobalModel
	orgId        int64
	sheetName    string
	colsLen      int
}

var backupTypeTamperGlobal = "TaG"

func (crd *BackupExcel) generateExportPath() string {

	fullDir, _ := upload_file.GenerateYearMonthPath(upload_file.GetDefaultPath())
	name := md5_encrypt.MD5(fmt.Sprintf("%d_%d", crd.orgId, variable.SnowFlake.GetId()))

	return fullDir + name + ".xlsx"
}

func (crd *BackupExcel) getTamperGlobal() []*tamper_models.GlobalModel {
	condition := tamper_models.SearchConditions{
		OrgId: crd.orgId,
	}
	return crd.tamperGlobal.GetListByConditions(condition)
}

// Export 导出
func (crd *BackupExcel) Export() (string, error) {
	//创建新 Excel 文件
	excelHandler := excel_helper.NewExcelHandler()
	defer excelHandler.Close()

	//行数据
	headerData := []interface{}{"类型", "数据"}
	var rowsData [][]interface{}

	//防篡改全局数据
	tamperGlobalModels := crd.getTamperGlobal()
	if len(tamperGlobalModels) >= 0 {
		for _, model := range tamperGlobalModels {
			model.OrgId = 0
			rowsData = append(rowsData, []interface{}{
				backupTypeTamperGlobal,
				json_helper.Encode(model),
			})
		}
	}

	//流式写入sheet
	if err := excelHandler.StreamWriteRows(crd.sheetName, headerData, rowsData, 1, 1); err != nil {
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
	//处理行数据
	for i, rowVals := range rowsValues {
		//忽略前2行
		if i < 1 {
			continue
		}
		//补足
		if colRealLen := len(rowVals); colRealLen < crd.colsLen {
			tmp := make([]string, crd.colsLen-colRealLen)
			rowVals = append(rowVals, tmp...)
		}
		switch rowVals[0] {
		case backupTypeTamperGlobal: //防篡改全局数据
			var tmpModel tamper_models.GlobalModel
			if err1 := json_helper.ToAny(rowVals[1], &tmpModel); err1 == nil && tmpModel.Id > 0 {
				tmpModel.OrgId = crd.orgId
				tmpModel.OpenSwitch = enums.SwitchOff
				tmpModel.UpdateId = userId
				tmpModel.UpdateName = userName
				tmpModel.UpdatedAt = now
				tmpModel.CreateId = userId
				tmpModel.CreateName = userName
				tmpModel.UpdatedAt = now

				insTamperGlobal = append(insTamperGlobal, tmpModel)
			}
		}
	}

	if len(insTamperGlobal) > 0 {
		crd.tamperGlobal.InsertBatch(&insTamperGlobal)
	}

	return nil
}
