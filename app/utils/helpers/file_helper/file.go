package file_helper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
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

// FormatFileSize 将字节大小转为易读的字符串，单位只使用 KB、MB、GB
// 参数: bytes - 字节大小
// 返回: 格式化后的字符串，如 "1.2 MB"
func FormatFileSize(bytes int64) string {
	if bytes < 0 {
		return "0 B"
	}

	const unit = 1024

	// 小于 1KB 直接显示字节
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	// 计算应该使用哪个单位
	var size float64
	var unitStr string

	switch {
	case bytes < unit*unit: // < 1MB，显示 KB
		size = float64(bytes) / float64(unit)
		unitStr = "KB"
	case bytes < unit*unit*unit: // < 1GB，显示 MB
		size = float64(bytes) / float64(unit*unit)
		unitStr = "MB"
	default: // >= 1GB，显示 GB
		size = float64(bytes) / float64(unit*unit*unit)
		unitStr = "GB"
	}

	// 智能显示小数位数
	// 如果是整数，不显示小数点
	if size == math.Trunc(size) {
		return fmt.Sprintf("%d %s", int64(size), unitStr)
	}
	// 否则保留 1 位小数
	return fmt.Sprintf("%.1f %s", size, unitStr)
}
