package port_service

import (
	"goskeleton/app/model/agent/port_models"
)

type PortBase struct {
	masterModel *port_models.VirtualGlobalModel
}

// 校验虚开端口
func (crd *PortBase) checkSourcePort(port string) error {

	return nil
}

// 校验端口范围
func (crd *PortBase) checkPortRange(port string) bool {

	return true
}

// 校验唯一性
func (crd *PortBase) checkPortUnique(port string, existsList []string) error {

	return nil
}
