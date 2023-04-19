package upload

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"threads-service/global"
	"threads-service/pkg/util"
)

type FileType int

const (
	TypeImage FileType = iota + 1
	TypeExcel
	TypeTxt
)

func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

func getFileExt(file string) string {
	return path.Ext(file)
}

func GetFileName(file string) string {
	ext := getFileExt(file)
	fileName := strings.TrimSuffix(file, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

/*检查函数*/

func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

func CheckContainExt(t FileType, name string) bool {
	ext := getFileExt(name)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
				return true
			}
		}
	}
	return false
}

// 检查文件大小是否超出限制。
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := io.ReadAll(f)
	size := len(content)
	global.Logger.Infof(context.TODO(), "上传的图片大小:%v", size)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}
	return false
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

/*写入读取函数*/
func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}
