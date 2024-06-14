package port_service

import (
	"errors"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/port_whether_occupy"
	"regexp"
	"strings"
)

type PortBase struct {
}

// 校验虚开端口
// @param portStr 支持单个、多个(10,20)英文逗号隔开、范围(90-99)、混合(10,20-25,80)
func (crd *PortBase) checkSourcePort(portStr string) error {
	portArr := strings.Split(portStr, ",")

	occupyFact := port_whether_occupy.CreateFactory()
	for _, port := range portArr {
		if strings.Contains(port, "-") {
			currPort := strings.Split(portStr, "-")
			if !crd.checkPortRange(currPort[0]) || !crd.checkPortRange(currPort[1]) {
				return errors.New("无效的端口【" + port + "】")
			}
			minPort := string_helper.ToInt(currPort[0])
			maxPort := string_helper.ToInt(currPort[1])
			if minPort >= maxPort {
				return errors.New("端口范围无效【" + port + "】")
			}
			for i := minPort; i <= maxPort; i++ {
				if !occupyFact.CheckPort(i) {
					return errors.New("本次端口有重复添加【" + port + "】")
				}
				occupyFact.SetPort(1)
			}
		} else {
			if !crd.checkPortRange(port) {
				return errors.New("端口无效【" + port + "】")
			}
			intPort := string_helper.ToInt(port)

			if !occupyFact.CheckPort(intPort) {
				return errors.New("本次端口有重复添加【" + port + "】")
			}
			occupyFact.SetPort(intPort)
		}
	}

	return nil
}

// 校验端口范围
func (crd *PortBase) checkPortRange(portStr string) bool {
	re := regexp.MustCompile(`^[0-9]{1,5}$`)
	if !re.MatchString(portStr) {
		return false
	}
	if portInt := string_helper.ToInt(portStr); portInt >= 1 && portInt <= 65535 {
		return true
	}

	return false
}

// 校验唯一性
// @param portStr 支持单个(10)、范围(90-99)
func (crd *PortBase) checkPortUnique(portStr string, existsList []string) error {
	if existsList == nil || len(existsList) == 0 {
		return nil
	}
	return nil
}
