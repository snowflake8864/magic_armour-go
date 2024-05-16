package job_callback_alarm

import (
	"strings"
)

type ProcessBase struct {
}

// ColumnCalcOfUsage 字段： 统计使用数
func (crd *ProcessBase) ColumnCalcOfUsage() string {
	var columns []string
	columns = append(columns, "SUM(buy_number) AS buy_number")
	columns = append(columns, "SUM(usage_number) AS usage_number")
	return strings.Join(columns, ", ")
}
