package port_service

import (
	"errors"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/port_whether_occupy"
	"regexp"
	"strings"
)

type PortBase struct {
	Lang string
}

// 校验虚开端口
// @param portStr 支持单个、多个(10,20)英文逗号隔开、范围(90-99)、混合(10,20-25,80)
func (crd *PortBase) checkSourcePort(portStr string) error {
	occupyFact := port_whether_occupy.CreateFactory()
	portArr := strings.Split(portStr, ",")
	for _, portRange := range portArr {
		if strings.Contains(portRange, "-") {
			currPort := strings.Split(portRange, "-")
			if !crd.checkPortRange(currPort[0]) || !crd.checkPortRange(currPort[1]) {
				return errors.New(my_errors.GetMessagef(crd.Lang, my_errors.CurdPortInvalid, portRange))
			}
			minPort := string_helper.ToInt64(currPort[0])
			maxPort := string_helper.ToInt64(currPort[1])
			if minPort >= maxPort {
				return errors.New(my_errors.GetMessagef(crd.Lang, my_errors.CurdPortRangeInvalid, portRange))
			}
			for i := minPort; i <= maxPort; i++ {
				if !occupyFact.CheckPort(i) {
					return errors.New(my_errors.GetMessagef(crd.Lang, my_errors.CurdPortRepeatedAdd, portRange))
				}
				occupyFact.SetPort(i)
			}
		} else {
			if !crd.checkPortRange(portRange) {
				return errors.New(my_errors.GetMessagef(crd.Lang, my_errors.CurdPortInvalid, portRange))
			}
			intPort := string_helper.ToInt64(portRange)

			if !occupyFact.CheckPort(intPort) {
				return errors.New(my_errors.GetMessagef(crd.Lang, my_errors.CurdPortRepeatedAdd, portRange))
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
	//
	occupyFact := port_whether_occupy.CreateFactory()

	//设置已占用的
	for _, existsStr := range existsList {
		portArr := strings.Split(existsStr, ",")
		for _, portRange := range portArr {
			if strings.Contains(portRange, "-") {
				currPort := strings.Split(portRange, "-")
				minPort := string_helper.ToInt64(currPort[0])
				maxPort := string_helper.ToInt64(currPort[1])
				for i := minPort; i <= maxPort; i++ {
					occupyFact.SetPort(i)
				}
			} else {
				occupyFact.SetPort(string_helper.ToInt64(portRange))
			}
		}
	}

	//校验是否占用
	portArr := strings.Split(portStr, ",")
	for _, portRange := range portArr {
		if strings.Contains(portRange, "-") {
			currPort := strings.Split(portRange, "-")
			minPort := string_helper.ToInt64(currPort[0])
			maxPort := string_helper.ToInt64(currPort[1])
			for i := minPort; i <= maxPort; i++ {
				if !occupyFact.CheckPort(i) {
					return errors.New(my_errors.GetMessagef(crd.Lang, my_errors.CurdPortRepeatedWithExistPolicy, portRange))
				}
			}
		} else {
			intPort := string_helper.ToInt64(portRange)
			if !occupyFact.CheckPort(intPort) {
				return errors.New(my_errors.GetMessagef(crd.Lang, my_errors.CurdPortRepeatedWithExistPolicy, portRange))
			}
		}
	}
	return nil
}
