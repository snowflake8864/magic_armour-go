package path_helper

import "path/filepath"

// ToStandard  将路径转义为标准格式
func ToStandard(path string) string {
	return filepath.ToSlash(path)

}

// ToWindows 将路径转换为windows默认的路径(斜杠分隔)
func ToWindows(path string) string {
	return filepath.FromSlash(path)
}

// GetFileName 获取路径中的文件名
func GetFileName(path string) string {
	return filepath.Base(path)
}
