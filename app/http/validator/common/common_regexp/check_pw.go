package common_regexp

import "regexp"

func CheckPw(pw string) bool {
	re := regexp.MustCompile(`^([a-zA-Z0-9.!@#%^&*()+=$_\-]+)$`)
	return re.MatchString(pw)
}
