package common_regexp

import (
	"errors"
	"regexp"
)

func CheckPw(str string) error {

	re := regexp.MustCompile(`^[a-zA-Z0-9.!@#%^&*()+=$_\-]{6,20}$`)
	if !re.MatchString(str) {
		return errors.New("密码格式不正确")
	}
	return nil
}

func CheckOrgName(str string) error {

	if length := len(str); length < 5 || length > 128 {
		return errors.New("企业名称格式不正确")
	}

	return nil
}

func CheckUserName(str string) error {

	re := regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z_]{3,24}$`)
	if !re.MatchString(str) {
		return errors.New("用户名格式不正确")
	}

	return nil
}

func CheckPhone(str string) error {

	re := regexp.MustCompile(`^1(3|4|5|6|7|8|9)\d{9}$`)
	if !re.MatchString(str) {
		return errors.New("手机号码格式不正确")
	}
	return nil
}

func CheckEmail(str string) error {

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(str) {
		return errors.New("邮箱格式不正确")
	}
	return nil
}

func CheckIpWithMask(str string) error {

	re := regexp.MustCompile(`^(?:(?:\d{1,3}\.){3}\d{1,3}(?:\/(?:[1-2]?[0-9]|3[0-2]))?)|(?:(?:(?:(?:[0-9a-fA-F]{1,4}):){7}(?:[0-9a-fA-F]{1,4}))(?:\/(?:[1-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?)$`)
	if !re.MatchString(str) {
		return errors.New("IP格式不正确")
	}
	return nil
}
