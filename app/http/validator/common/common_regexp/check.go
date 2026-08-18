package common_regexp

import (
	"errors"
	"goskeleton/app/global/my_errors"
	"regexp"
)

func CheckPw(str string) error {

	re := regexp.MustCompile(`^[a-zA-Z0-9.!@#%^&*()+=$_\-]{6,20}$`)
	if !re.MatchString(str) {
		return errors.New(my_errors.CurdPasswordFormatIncorrect)
	}
	return nil
}

func CheckOrgName(str string) error {

	if length := len(str); length < 5 || length > 128 {
		return errors.New(my_errors.CurdOrgNameFormatIncorrect)
	}

	return nil
}

func CheckUserName(str string) error {

	re := regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z_]{3,24}$`)
	if !re.MatchString(str) {
		return errors.New(my_errors.CurdUsernameFormatIncorrect)
	}

	return nil
}

func CheckPhone(str string) error {

	re := regexp.MustCompile(`^1(3|4|5|6|7|8|9)\d{9}$`)
	if !re.MatchString(str) {
		return errors.New(my_errors.CurdPhoneFormatIncorrect)
	}
	return nil
}

func CheckEmail(str string) error {

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(str) {
		return errors.New(my_errors.CurdEmailFormatIncorrect)
	}
	return nil
}

func CheckIpWithMask(str string) error {

	re := regexp.MustCompile(`^(?:(?:\d{1,3}\.){3}\d{1,3}(?:\/(?:[1-2]?[0-9]|3[0-2]))?)|(?:(?:(?:(?:[0-9a-fA-F]{1,4}):){7}(?:[0-9a-fA-F]{1,4}))(?:\/(?:[1-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?)$`)
	if !re.MatchString(str) {
		return errors.New(my_errors.CurdIpFormatIncorrect)
	}
	return nil
}
