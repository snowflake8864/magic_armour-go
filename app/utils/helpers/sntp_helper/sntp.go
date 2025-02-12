package sntp_helper

import (
	"bytes"
	"errors"
	"goskeleton/app/utils/helpers/string_helper"
	"os/exec"
)

// Exec 执行命令
func Exec(command, addr string) (error, string) {
	cmd := exec.Command("/bin/bash", "-c", command+" "+addr)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	if err != nil {
		return errors.New("exec command Run failed: " + err.Error()), ""
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if string_helper.IsEmpty(outStr) {
		return errors.New("exec command output empty"), ""
	}
	if !string_helper.IsEmpty(errStr) {
		return errors.New("exec command error not empty"), ""
	}

	return nil, outStr
}
