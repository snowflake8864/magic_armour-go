package excel_helper

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

// ExcelHandler 封装 Excel 操作句柄
type ExcelHandler struct {
	File   *excelize.File
	Styles map[string]int // 样式缓存（例如："header", "highlight"）
}

// NewExcelHandler 创建新 Excel 文件
func NewExcelHandler() *ExcelHandler {
	return &ExcelHandler{
		File:   excelize.NewFile(),
		Styles: make(map[string]int),
	}
}

// NewExcelHandlerFromFile 打开现有 Excel 文件
func NewExcelHandlerFromFile(filename string) (*ExcelHandler, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	return &ExcelHandler{
		File:   f,
		Styles: make(map[string]int),
	}, nil
}

// ======================= 基础操作 =======================

// SetCellValue 设置单元格值
// sheet: 工作表名称, cell: 单元格坐标(如"A1"), value: 要写入的值
func (h *ExcelHandler) SetCellValue(sheet, cell string, value interface{}) error {
	return h.File.SetCellValue(sheet, cell, value)
}

// GetCellValue 获取单元格值
func (h *ExcelHandler) GetCellValue(sheet, cell string) (string, error) {
	return h.File.GetCellValue(sheet, cell)
}

// SetRowValues 按行连续写入数据
// 示例：SetRowValues("Sheet1", "A1", []interface{}{"姓名", "年龄", "性别"})
func (h *ExcelHandler) SetRowValues(sheet, startCell string, values []interface{}) error {
	axis := startCell
	for i, value := range values {
		// 写入当前单元格
		if err := h.File.SetCellValue(sheet, axis, value); err != nil {
			return err
		}

		// 解析当前坐标
		col, row, err := excelize.CellNameToCoordinates(axis)
		if err != nil {
			return fmt.Errorf("解析坐标失败: %v", err)
		}

		// 计算下一个单元格坐标（列+1，行不变）
		nextCol := col + 1
		nextRow := row

		// 生成新坐标
		axis, err = excelize.CoordinatesToCellName(nextCol, nextRow)
		if err != nil {
			return fmt.Errorf("生成坐标失败: %v", err)
		}

		if i >= 1000 { // 安全保护
			break
		}
	}
	return nil
}

// StreamWriteRows 使用流式写入器提升大批量数据性能
// 示例：SetRowValues("Sheet1", "A1", []interface{}{"姓名", "年龄", "性别"})
func (h *ExcelHandler) StreamWriteRows(sheet string, rows [][]interface{}) error {
	streamWriter, err := h.File.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	for i, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		if err := streamWriter.SetRow(cell, row); err != nil {
			return err
		}
	}
	return streamWriter.Flush()
}

// ======================= 样式管理 =======================

// AddStyle 创建自定义样式
// 示例样式配置见下方代码
func (h *ExcelHandler) AddStyle(styleName string, style *excelize.Style) error {
	styleID, err := h.File.NewStyle(style)
	if err != nil {
		return err
	}
	h.Styles[styleName] = styleID
	return nil
}

// ApplyStyle 应用已存在的样式到单元格
func (h *ExcelHandler) ApplyStyle(sheet, topLeftCell, bottomRightCell, styleName string) error {
	styleID, ok := h.Styles[styleName]
	if !ok {
		return fmt.Errorf("样式未定义: %s", styleName)
	}
	return h.File.SetCellStyle(sheet, topLeftCell, bottomRightCell, styleID)
}

// ======================= 常用预设样式 =======================

// InitDefaultStyles 初始化常用样式
func (h *ExcelHandler) InitDefaultStyles() error {
	// 标题样式
	headerStyle := &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"2c7be5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	}
	if err := h.AddStyle("header", headerStyle); err != nil {
		return err
	}

	// 高亮样式
	highlightStyle := &excelize.Style{
		Font: &excelize.Font{Color: "FF0000"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
		},
	}
	return h.AddStyle("highlight", highlightStyle)
}

// ======================= 单元格 =======================

// GetNextCell 获取右侧相邻单元格（辅助函数）
func GetNextCell(cell string) (string, error) {
	//解析当前坐标
	col, row, err := excelize.CellNameToCoordinates(cell)
	if err != nil {
		return "", err
	}
	//计算下一个单元格坐标（列+1，行不变）
	return excelize.CoordinatesToCellName(col+1, row)
}

// MergeCells 合并单元格
func (h *ExcelHandler) MergeCells(sheet string, startRow, startCol, endRow, endCol int) error {
	startCell, _ := excelize.CoordinatesToCellName(startCol, startRow)
	endCell, _ := excelize.CoordinatesToCellName(endCol, endRow)
	return h.File.MergeCell(sheet, startCell, endCell)
}

// ======================= 文件操作 =======================

// SaveAs 保存文件
func (h *ExcelHandler) SaveAs(filename string) error {
	if err := h.File.SaveAs(filename); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}
	return nil
}

// Close 关闭文件（如果需要处理流式写入）
func (h *ExcelHandler) Close() error {
	return h.File.Close()
}

// ======================= 高级功能 =======================

// SetColWidth 设置列宽
func (h *ExcelHandler) SetColWidth(sheet, startCol, endCol string, width float64) error {
	return h.File.SetColWidth(sheet, startCol, endCol, width)
}
