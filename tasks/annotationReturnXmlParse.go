package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"BeeCustom/controllers"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xmlTemplate"
	"github.com/astaxie/beego/toolbox"
)

// 回执解释
func annotationCReturnXmlParse() *toolbox.Task {

	task := toolbox.NewTask("task", "* * * * * *", func() error {
		parseAnnotationReturns("annotation_return_path_c", "annotation_history_path_c")
		return nil
	})

	err := task.Run()
	if err != nil {
		utils.LogError(fmt.Sprintf("tk.Run error :%v", err))
		return nil
	}

	return task
}

// 解析回执
func parseAnnotationReturns(returnPathConfig, historyPathConfig string) {

	returnPath, err := models.GetSettingValueByKey(returnPathConfig)
	historyPath, err := models.GetSettingValueByKey(historyPathConfig)
	pathCFiles, err := ioutil.ReadDir(returnPath)
	if err != nil {
		return
	}

	for _, f := range pathCFiles {
		fullPath := returnPath + f.Name()

		// 文件前缀和后缀
		prefix, ext, failedName := getNameExt(f)
		if len(prefix) == 0 || len(ext) == 0 {
			continue
		}

		// 首个回执
		if ext == "INV" {

			if prefix == "Successed" {

				xmlFile, err, data := openFile(fullPath)
				if err != nil {
					utils.LogError(fmt.Sprintf("openFile :%v ,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}

				v := xmlTemplate.CommonResponeMessage{}
				err = xml.Unmarshal(data, &v)
				if err != nil {
					utils.LogError(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}

				annotation, err := models.GetAnnotationByEtpsInnerInvtNo(v.EtpsPreentNo)
				if err != nil {
					xmlFile.Close()
					continue
				}

				aReturn := models.NewAnnotationReturn(0)
				aReturn.EtpsPreentNo = v.EtpsPreentNo
				aReturn.CheckInfo = v.CheckInfo
				aReturn.DealFlag = v.DealFlag
				aReturn.SeqNo = v.SeqNo
				aReturn.Annotation = annotation

				if err = models.AnnotationReturnSave(&aReturn); err != nil {
					utils.LogError(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}

				// 更新状态
				annotation.SeqNo = v.SeqNo
				if err = controllers.UpdateAnnotationStatus(annotation, "单一处理中", false); err != nil {
					utils.LogError(fmt.Sprintf("controllers.UpdateAnnotationStatus :%v,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}

				if err := models.AnnotationUpdate(annotation, []string{"Status", "StatusUpdatedAt", "SeqNo"}); err != nil {
					utils.LogError(fmt.Sprintf("AnnotationUpdateStatus :%v,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}

				xmlFile.Close()
				err = moveFile(historyPath, v.EtpsPreentNo, fullPath, f)
				if err != nil {
					utils.LogError(fmt.Sprintf("moveFile:%v,filename:%v", err, f.Name()))
					continue
				}

				// ws 自动更新
				wsPush()

			} else if prefix == "Receipt" {
				err = moveFile(historyPath, "Receipt", fullPath, f)
				if err != nil {
					utils.LogError(fmt.Sprintf("moveFile :%v,filename:%v", err, f.Name()))
				}
			}

		} else if ext == "INVT" && prefix == "Receipt" {

			xmlFile, err, data := openFile(fullPath)
			if err != nil {
				utils.LogError(fmt.Sprintf("openFile :%v ,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			v := xmlTemplate.ReturnPackage{}
			err = xml.Unmarshal(data, &v)
			if err != nil {
				utils.LogError(fmt.Sprintf(" xml.Unmarshal :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			annotation, err := models.GetAnnotationBySeqNo(v.InvPreentNo)
			if err != nil {
				// utils.LogError(fmt.Sprintf("models.GetAnnotationBySeqNo :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			aReturn := models.NewAnnotationReturn(0)
			aReturn.BusinessId = v.DataInfoBusinessId
			aReturn.ManageResult = v.ManageResult
			aReturn.Reason = v.Reason
			aReturn.Annotation = annotation

			if err = models.AnnotationReturnSave(&aReturn); err != nil {
				utils.LogError(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			// 更新状态
			if err = controllers.UpdateAnnotationStatus(annotation, "已完成", false); err != nil {
				utils.LogError(fmt.Sprintf("enums.UpdateAnnotationStatus :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}
			annotation.BondInvtNo = v.DataInfoBusinessId
			if err := models.AnnotationUpdate(annotation, []string{"Status", "StatusUpdatedAt", "BondInvtNo"}); err != nil {

				utils.LogError(fmt.Sprintf("AnnotationUpdateStatus :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			xmlFile.Close()
			err = moveFile(historyPath, annotation.EtpsInnerInvtNo, fullPath, f)
			if err != nil {
				utils.LogError(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
			}

			// ws 自动更新
			msg := utils.Message{Message: "清单状态更新", IsUpdated: true}
			utils.Broadcast <- msg

		} else if ext == "xml" {
			if prefix == "Failed" {
				xmlFile, err, data := openFile(fullPath)
				if err != nil {
					utils.LogError(fmt.Sprintf("openFile :%v ,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}

				v := xmlTemplate.Failed{}
				err = xml.Unmarshal(data, &v)
				if err != nil {
					utils.LogError(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}

				annotation, err := models.GetAnnotationByEtpsInnerInvtNo(failedName)
				if err != nil {
					xmlFile.Close()
					continue
				}

				aReturn := models.NewAnnotationReturn(0)
				aReturn.CheckInfo = v.FailInfo
				aReturn.DealFlag = v.FailCode
				aReturn.Annotation = annotation

				if err = models.AnnotationReturnSave(&aReturn); err != nil {
					utils.LogError(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
					xmlFile.Close()
					continue
				}
			}
		}

		// 错误文件移动
		err = moveFile(historyPath, "Error", fullPath, f)
		if err != nil {
			utils.LogError(fmt.Sprintf("moveFile:%v,filename:%v", err, f.Name()))
		}
	}
}

// 回执文件的前缀，和后缀
func getNameExt(f os.FileInfo) (string, string, string) {
	extNames := strings.Split(f.Name(), ".")
	if len(extNames) > 1 && len(extNames[1]) > 0 && len(extNames[0]) > 0 {
		eNames := strings.Split(extNames[0], "_")
		if len(eNames) > 1 && len(eNames[0]) > 0 {
			if len(eNames) > 3 && len(eNames[3]) > 0 {
				return eNames[0], extNames[1], eNames[3]
			} else {
				return eNames[0], extNames[1], ""
			}
		}

	} else {
		return "", extNames[1], ""
	}

	return "", "", ""
}
