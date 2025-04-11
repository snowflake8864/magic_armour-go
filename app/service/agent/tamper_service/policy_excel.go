package tamper_service

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/helpers/excel_helper"
	"goskeleton/app/utils/md5_encrypt"
	"strings"
)

func CreatePolicyExcelFactory(orgId int64, agentId string) *PolicyExcel {
	return &PolicyExcel{
		masterModel: tamper_models.CreatePolicyFactory(),
		orgId:       orgId,
		agentId:     agentId,
		sheetName:   "Sheet1",
		sep:         "|",
	}
}

type PolicyExcel struct {
	masterModel *tamper_models.PolicyModel
	orgId       int64
	agentId     string
	sheetName   string
	sep         string
}

func (crd *PolicyExcel) generateExportPath() string {

	fullDir, _ := upload_file.GenerateYearMonthPath(upload_file.GetDefaultPath())
	name := md5_encrypt.MD5(fmt.Sprintf("%s_%d", crd.agentId, variable.SnowFlake.GetId()))

	return fullDir + name + ".xlsx"
}

// Export 导出
func (crd *PolicyExcel) Export() (string, error) {
	//创建新 Excel 文件
	excelHandler := excel_helper.NewExcelHandler()
	defer excelHandler.Close()
	// 初始化默认样式
	_ = excelHandler.InitDefaultStyles()

	//顶部
	top := "" +
		"类型：单个文件或文件夹 如：文件夹\n" +
		"名称：当前文件或文件夹路径 如：D:\\Backup\\Downloads\n" +
		"保护中：读取|写入|删除|重命名|新建 如：读取|写入\n" +
		"包含文件类型：多个以竖线隔开 如: gif|jpg\n" +
		"排除文件类型：多个以竖线隔开 如: gif|jpg\n" +
		"排除文件夹：多个以竖线隔开 如: www|test\n" +
		"继承关系：0否 1是 如：是\n" +
		"信任进程：多个hash以竖线隔开如: aa|bb"
	if err := excelHandler.MergeCells(crd.sheetName, 1, 1, 1, 9); err != nil {
		return "", fmt.Errorf("merge cells fail: %v", err)
	}
	if err := excelHandler.SetCellValue(crd.sheetName, "A1", top); err != nil {
		return "", fmt.Errorf("set top fail: %v", err)
	}
	_ = excelHandler.SetRowHeight(crd.sheetName, 1, 160)

	//头部
	header := []interface{}{
		"类型",
		"名称",
		"保护中",
		"包含文件类型",
		"排除文件类型",
		"排除文件夹",
		"继承关系",
		"信任进程",
		"备注",
	}
	if err := excelHandler.SetHeader(crd.sheetName, "A2", "I2", header); err != nil {
		return "", fmt.Errorf("set header fail: %v", err)
	}
	_ = excelHandler.SetColWidth(crd.sheetName, "A", "A", 20)
	_ = excelHandler.SetColWidth(crd.sheetName, "B", "B", 50)
	_ = excelHandler.SetColWidth(crd.sheetName, "C", "F", 30)
	_ = excelHandler.SetColWidth(crd.sheetName, "G", "G", 15)
	_ = excelHandler.SetColWidth(crd.sheetName, "H", "H", 50)
	_ = excelHandler.SetColWidth(crd.sheetName, "I", "I", 30)

	//策略数据
	condition := tamper_models.SearchConditions{
		OrgId:   crd.orgId,
		AgentId: crd.agentId,
	}
	models := tamper_models.CreatePolicyFactory().GetListByConditions(condition)
	if len(models) >= 0 {
		var rowsData [][]interface{}
		for _, model := range models {

			rowsData = append(rowsData, []interface{}{
				enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_tamper.EnumMapInt),
				model.Dir,
				strings.Join(enum_tamper.GetPowerLabels(model.ProtectRw), crd.sep),
				model.IncludeSuffix,
				model.ExcludeSuffix,
				model.ExcludeFolder,
				enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend", model.IsExtend, enum_tamper.EnumMapInt),
				"",
				model.Remark,
			})
		}

		if err := excelHandler.SetRowsValues(crd.sheetName, 1, 3, rowsData); err != nil {
			return "", fmt.Errorf("set rows fail: %v", err)
		}
	}

	// 保存文件
	path := crd.generateExportPath()
	fmt.Println("path:", path)
	if err := excelHandler.SaveAs(path); err != nil {
		return "", fmt.Errorf("save path fail: %v", err)
	}

	return path, nil
}
