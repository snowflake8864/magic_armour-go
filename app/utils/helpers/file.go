package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func FileMd5sum(dir string) (error, string) {
	// 打开文件
	file, err := os.Open(dir)
	if err != nil {
		return errors.New("文件不存在：" + err.Error()), ""
	}
	defer file.Close()

	// 使用MD5哈希对象
	hash := md5.New()

	// 将文件内容拷贝给哈希对象
	if _, err := io.Copy(hash, file); err != nil {
		return errors.New("文件不存在：" + err.Error()), ""
	}

	// 计算MD5值
	md5Hash := hash.Sum(nil)

	// 将MD5值转换为字符串
	md5String := hex.EncodeToString(md5Hash)

	return nil, md5String
}
