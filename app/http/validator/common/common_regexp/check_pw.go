package common_regexp

import "regexp"

func CheckPw(pw string) bool {

	re := regexp.MustCompile(`^([a-zA-Z0-9.!@#%^&*()+=$_\-]+)$`)
	return re.MatchString(pw)
}

func CheckIpWithMask(pw string) bool {

	re := regexp.MustCompile(`^(?:(?:\d{1,3}\.){3}\d{1,3}(?:\/(?:[1-2]?[0-9]|3[0-2]))?)|(?:(?:(?:(?:[0-9a-fA-F]{1,4}):){7}(?:[0-9a-fA-F]{1,4}))(?:\/(?:[1-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?)$`)
	return re.MatchString(pw)
}
