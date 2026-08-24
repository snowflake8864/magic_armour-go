package org_trans

import (
	"goskeleton/app/model/agent/agent_monitor_models"
	"strings"
)

type MonitorTrans struct {
}

// TransList 转换成列表
func (t *MonitorTrans) TransList(models []*agent_monitor_models.MonitorModel) interface{} {
	result := make(map[string][]interface{})

	if models != nil {
		var categories, mem, hd, cpu []interface{}
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
