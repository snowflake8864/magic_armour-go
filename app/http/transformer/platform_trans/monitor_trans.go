package platform_trans

import (
	"goskeleton/app/model/platform/server_monitor_models"
	"strings"
)

type Monitor struct {
}

// TransList 转换成列表
func (t *Monitor) TransList(models []*server_monitor_models.ServerMonitorModel) interface{} {
	result := make(map[string][]string)

	if models != nil {
		var categories, mem, hd, cpu []string
		for _, model := range models {
			tmp := strings.Split(model.CreatedAt, " ")
			categories = append(categories, tmp[1])
			mem = append(mem, model.MemUsage)
			hd = append(hd, model.HdUsage)
			cpu = append(cpu, model.CpuUsage)
		}

		result["categories"] = categories
		result["cpu_usage"] = cpu
		result["mem_usage"] = mem
		result["hd_usage"] = hd
	}

	return result
}
