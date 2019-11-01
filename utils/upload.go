package utils

import (
	"errors"
	"fmt"
	"mime/multipart"

	uuid "github.com/iris-contrib/go.uuid"
)

//获取上传文件唯一地址
func GetUploadFileUPath(f multipart.File, h *multipart.FileHeader, filetype string) (string, error) {

	if f != nil {
		defer f.Close()
	} else {
		LogDebug("上传失败")
		return "", errors.New("上传失败")
	}

	uid, _ := uuid.NewV4()
	if h != nil {

		filepath := "static/upload/" + filetype
		if err := CreateFile(filepath); err != nil {
			LogDebug(fmt.Sprintf("文件夹创建失败:%v", err))
			return "", err
		}

		fileNamePath := filepath + uid.String() + "_" + h.Filename

		return fileNamePath, nil

	} else {
		LogDebug("上传失败")
		return "", errors.New("上传失败")
	}

}
