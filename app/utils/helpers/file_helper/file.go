package file_helper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Md5sum 文件hash值
func Md5sum(dir string) (error, string) {
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

// ExistsFile 文件是否存在
func ExistsFile(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CopyFile 复制文件
func CopyFile(srcFile, destFile string) error {
	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	err = destination.Sync()
	if err != nil {
		return err
	}

	return nil
}

// FileInfo 文件信息
func FileInfo(filePath string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

// ListDir 获取目录中的所有文件
func ListDir(dirname string) ([]string, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names, nil
}

// GetAllFiles 递归获取目录下所有文件路径
func GetAllFiles(dir string) ([]string, error) {
	var files []string

	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// 遍历目录条目
	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// 递归处理子目录
			subFiles, err := GetAllFiles(fullPath)
			if err != nil {
				return nil, err // 遇到错误立即返回
			}
			files = append(files, subFiles...)
		} else {
			// 添加文件路径
			files = append(files, fullPath)
		}
	}

	return files, nil
}
