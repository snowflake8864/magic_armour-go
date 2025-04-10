package test

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"goskeleton/app/utils/helpers/excel_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestExcelWrite(t *testing.T) {
	handler := excel_helper.NewExcelHandler()

	// 初始化默认样式
	if err := handler.InitDefaultStyles(); err != nil {
		fmt.Println("InitDefaultStyles: ", err.Error())
		return
	}

	// 创建工作表
	sheetName := "Sheet12"
	handler.File.NewSheet(sheetName)

	// 写入标题行并应用样式
	if err := handler.SetRowValues(sheetName, "A1", []interface{}{"姓名", "年龄", "部门"}); err != nil {
		fmt.Println("SetRowValues: ", err.Error())
		return
	}
	if err := handler.ApplyStyle(sheetName, "A1", "C1", "header"); err != nil {
		fmt.Println("ApplyStyle: ", err.Error())
		return
	}

	// 写入数据
	data := [][]interface{}{
		{"张三", 28, "技术部"},
		{"李四", 32, "市场部"},
	}
	for i, row := range data {
		startCell, _ := excelize.CoordinatesToCellName(1, i+2) // 从第二行开始
		_ = handler.SetRowValues(sheetName, startCell, row)
	}

	// 设置列宽
	_ = handler.SetColWidth(sheetName, "A", "C", 15)

	// 合并单元格（示例：合并A4:C4）
	_ = handler.MergeCells(sheetName, 4, 1, 4, 3)
	_ = handler.SetCellValue(sheetName, "A4", "合并单元格示例")

	// 保存文件
	if err := handler.SaveAs("demo.xlsx"); err != nil {
		fmt.Println("InitDefaultStyles: ", err.Error())
	}
}
