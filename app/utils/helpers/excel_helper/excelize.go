package excel_helper

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

const (
	StyleHeader    string = "header"
	StyleHighlight string = "highlight"
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

// CreateSheet 创建Sheet，存在不在新建
func (h *ExcelHandler) CreateSheet(sheet string) error {
	// 确保工作表存在
	if index, _ := h.File.GetSheetIndex(sheet); index > -1 {
		return nil
	}
	// 创建工作表
	if _, err := h.File.NewSheet(sheet); err != nil {
		return fmt.Errorf("NewSheet fail: %s", err.Error())
	}
	return nil
}

// SetCellValue 设置单元格值
// sheet: 工作表名称, cell: 单元格坐标(如"A1"), value: 要写入的值
func (h *ExcelHandler) SetCellValue(sheet, cell string, value interface{}) error {
	return h.File.SetCellValue(sheet, cell, value)
}

// GetCellValue 获取单元格值
func (h *ExcelHandler) GetCellValue(sheet, cell string) (string, error) {
	return h.File.GetCellValue(sheet, cell)
}

// SetHeader 设置标题
func (h *ExcelHandler) SetHeader(sheet, startCell, endCell string, values []interface{}) error {
	// 写入标题行
	if err := h.SetRowValues(sheet, startCell, values); err != nil {
		return fmt.Errorf("WriteHeader fail: %s", err.Error())
	}
	//应用样式
	_ = h.ApplyStyle(sheet, startCell, endCell, StyleHeader)

	return nil
}

// SetRowValues 按行连续写入数据
// 示例：SetRowValues("Sheet1", "B1", []interface{}{"姓名", "年龄", "性别"})
func (h *ExcelHandler) SetRowValues(sheet, startCell string, values []interface{}) error {
	axis := startCell
	for i, value := range values {
		// 写入当前单元格
		err := h.File.SetCellValue(sheet, axis, value)
		if err != nil {
			return fmt.Errorf("写入当前单元格失败: %v", err)
		}

		// 生成新坐标
		axis, err = h.GetNextCell(axis)
		if err != nil {
			return fmt.Errorf("生成新坐标失败: %v", err)
		}

		if i >= 1000 { // 安全保护
			break
		}
	}
	return nil
}

// ======================= 流式操作 =======================

//流式写入器不支持以下操作
//MergeCell() 合并单元格
//SetCellStyle() 单独设置样式
//AddComment()    添加批注
//如需样式处理，需在 SetRow 时通过 excelize.Cell{Value: ..., StyleID: ...} 指定

// StreamSetHeader 流式设置标题
func (h *ExcelHandler) StreamSetHeader(sw *excelize.StreamWriter, values []interface{}, startCol, startRow int) error {
	startCell, err := excelize.CoordinatesToCellName(startCol, startRow)
	if err != nil {
		return fmt.Errorf("生成坐标失败: %v", err)
	}

	cells := make([]interface{}, len(values))
	if styleID, ok := h.GetStyleID(StyleHeader); ok {
		for _, value := range values {
			cells = append(cells, excelize.Cell{
				Value:   value,
				StyleID: styleID,
			})
		}
	} else {
		for _, value := range values {
			cells = append(cells, excelize.Cell{
				Value: value,
			})
		}
	}

	// 流式写入标题（需直接操作原始接口）
	err = sw.SetRow(startCell, cells)
	if err != nil {
		return fmt.Errorf("流式写入标题失败: %v", err)
	}

	return nil
}

// StreamWriteRows 使用流式写入器提升大批量数据性能
// 示例：StreamWriteRows("Sheet1", []interface{}{"姓名", "年龄", "性别"}, [][]interface{}{{"张三", 28, "技术部"},{"李四", 32, "市场部"}}, 1, 1)
func (h *ExcelHandler) StreamWriteRows(sheet string, headerData []interface{}, rowsData [][]interface{}, startCol, startRow int) error {
	// 创建流式写入器
	streamWriter, err := h.File.NewStreamWriter(sheet)
	if err != nil {
		return err
	}

	// 流式写入标题（需直接操作原始接口）
	_ = h.StreamSetHeader(streamWriter, headerData, startCol, startRow)

	//起始行增加
	startRow++

	/*for i, rowData := range rowsData {
		//生成起始单元格坐标
		startCell, _ := excelize.CoordinatesToCellName(startCol, startRow+i)
		fmt.Println("startCell:", startCell)

		// 创建单元格数据
		cells := make([]interface{}, len(rowData))
		for _, value := range rowData {
			cells = append(cells, excelize.Cell{
				Value: value,
			})
		}

		// 写入行数据
		if err1 := streamWriter.SetRow(startCell, rowData); err1 != nil {
			//if err1 := streamWriter.SetRow(startCell, cells); err1 != nil {
			return fmt.Errorf("写入第(%d)行失败: %v", i, err1)
		}

	}*/
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

// GetStyleID 获取样式变化
func (h *ExcelHandler) GetStyleID(styleName string) (int, bool) {
	styleID, ok := h.Styles[styleName]

	return styleID, ok
}

// ApplyStyle 应用已存在的样式到单元格
func (h *ExcelHandler) ApplyStyle(sheet, topLeftCell, bottomRightCell, styleName string) error {
	styleID, ok := h.GetStyleID(styleName)
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
	if err := h.AddStyle(StyleHeader, headerStyle); err != nil {
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
	return h.AddStyle(StyleHighlight, highlightStyle)
}

// ======================= 单元格 =======================

// GetNextCell 获取右侧相邻单元格（辅助函数）
func (h *ExcelHandler) GetNextCell(cell string) (string, error) {
	//解析当前坐标
	col, row, err := excelize.CellNameToCoordinates(cell)
	if err != nil {
		return "", err
	}

	// 计算下一个单元格坐标（列+1，行不变）
	nextCol := col + 1
	nextRow := row

	// 生成新坐标
	cell, err = excelize.CoordinatesToCellName(nextCol, nextRow)
	if err != nil {
		return "", fmt.Errorf("生成坐标失败: %v", err)
	}
	return cell, nil
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
