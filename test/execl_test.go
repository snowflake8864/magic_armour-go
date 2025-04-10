package test

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"goskeleton/app/utils/helpers/excel_helper"
	"goskeleton/app/utils/helpers/int_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"runtime"
	"testing"
	"time"
)

func TestExcelWrite(t *testing.T) {
	handler := excel_helper.NewExcelHandler()
	defer handler.Close()

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
		fmt.Println("SaveAs: ", err.Error())
	}
}

func TestCreateSheets(t *testing.T) {
	handler := excel_helper.NewExcelHandler()
	defer handler.Close()

	// 初始化默认样式
	if err := handler.InitDefaultStyles(); err != nil {
		fmt.Println("InitDefaultStyles: ", err.Error())
		return
	}

	for i := 1; i <= 10; i++ {
		sheet := fmt.Sprintf("测试%d", i)
		if err := handler.CreateSheet(sheet); err != nil {
			fmt.Println("CreateSheet: ", err.Error())
		}

		if err := handler.SetHeader(sheet, "A1", "C1", []interface{}{"姓名", "年龄", "部门"}); err != nil {
			fmt.Println("SetHeader: ", err.Error())
		}
	}

	// 保存文件
	if err := handler.SaveAs("demo1.xlsx"); err != nil {
		fmt.Println("SaveAs: ", err.Error())
	}
}

// 获取当前内存占用（单位：MB）
func getMemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}

func TestStreamWriteRows(t *testing.T) {
	handler := excel_helper.NewExcelHandler()
	defer handler.Close()

	// 初始化默认样式
	if err := handler.InitDefaultStyles(); err != nil {
		fmt.Println("InitDefaultStyles: ", err.Error())
		return
	}
	sheet := "Sheet-5"

	// 批量写入10万行数据
	startTime := time.Now()
	var rowsData [][]interface{}
	headerData := []interface{}{"姓名", "年龄", "性别"}

	for i := 0; i < 10; i++ {
		rowsData = append(rowsData, []interface{}{
			fmt.Sprintf("项目%d", i),
			int_helper.IntRandom(20, 55),
			fmt.Sprintf("部门%d", i),
		})
	}

	//流式写入
	if err := handler.StreamWriteRows(sheet, headerData, rowsData, 1, 1); err != nil {
		fmt.Println("CreateSheetAndWriteHeader: ", err.Error())
	}

	// 性能统计
	fmt.Printf("写入完成 总行数: %d\n", len(rowsData))
	fmt.Printf("耗时: %.2f秒\n", time.Since(startTime).Seconds())
	fmt.Printf("内存占用: ≈%.2f MB\n", getMemoryUsage())

	// 保存文件
	if err := handler.SaveAs("10w行.xlsx"); err != nil {
		fmt.Println("SaveAs: ", err.Error())
	}
}
