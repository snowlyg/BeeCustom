package file

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"BeeCustom/utils"
	uuid "github.com/iris-contrib/go.uuid"
)

//调用os.MkdirAll递归创建文件夹
func CreateFile(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//获取上传文件唯一地址
func GetUploadFileUPath(f multipart.File, h *multipart.FileHeader, fileType string) (string, error) {

	if f != nil {
		defer f.Close()
	} else {
		utils.LogDebug("上传失败")
		return "", errors.New("上传失败")
	}

	uid, _ := uuid.NewV4()
	if h != nil {

		filepath := "static/upload/" + fileType
		if err := CreateFile(filepath); err != nil {
			utils.LogDebug(fmt.Sprintf("文件夹创建失败:%v", err))
			return "", err
		}

		fileNamePath := filepath + uid.String() + "_" + h.Filename

		return fileNamePath, nil

	} else {
		utils.LogDebug("上传失败")
		return "", errors.New("上传失败")
	}

}
