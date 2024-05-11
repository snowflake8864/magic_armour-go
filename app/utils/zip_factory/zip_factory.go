package zip_factory

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func ZipDir(srcPath string, zipFileName string) error {
	// 创建压缩文件
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建 zip.Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历文件夹
	err = filepath.Walk(srcPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断文件是否为目录
		if info.IsDir() {
			return nil
		}

		// 创建压缩文件的路径
		relativePath, err := filepath.Rel(srcPath, filePath)
		if err != nil {
			return err
		}

		// 创建 zip 文件的路径
		zipPath := filepath.Join(filepath.Base(srcPath), relativePath)

		// 创建压缩文件的头部信息
		zipHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		zipHeader.Name = zipPath

		// 打开待压缩文件
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 写入压缩文件
		writer, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

// UnzipFile 解压zip文件
func UnzipFile(zipPath string, destPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddFileToZip 新增文件到zip包，并重命名
func AddFileToZip(srcPath, destPath, newFileName, newFileContent string) error {
	// 打开原始的zip文件
	srcZipFile, err := os.OpenFile(srcPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer srcZipFile.Close()

	// 创建一个新的zip文件，用于写入新的文件
	newZipFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	// 创建一个zip写入器
	netZipWriter := zip.NewWriter(newZipFile)

	// 将原始zip文件的内容复制到新的zip文件中
	originalReader, err := zip.OpenReader(srcZipFile.Name())
	if err != nil {
		return err
	}
	defer originalReader.Close()

	for _, file := range originalReader.File {
		// 创建一个新的文件头
		header := &zip.FileHeader{
			Name:   file.Name,
			Method: file.Method,
		}

		// 打开原始文件
		originalFile, err := file.Open()
		if err != nil {
			return err
		}
		defer originalFile.Close()

		// 将原始文件复制到新的zip文件中
		newFile, err := netZipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(newFile, originalFile)
		if err != nil {
			return err
		}
	}

	// 创建一个新的文件并写入内容
	newFile, err := netZipWriter.Create(newFileName)
	if err != nil {
		return err
	}
	_, err = newFile.Write([]byte(newFileContent))
	if err != nil {
		return err
	}

	// 关闭zip写入器
	err = netZipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
