package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"BeeCustom/controllers"
	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xmlTemplate"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

//回执解释
func annotationCReturnXmlParse() *toolbox.Task {

	task := toolbox.NewTask("task", "* * * * * *", func() error {
		//o := orm.NewOrm()
		//status12, err := enums.GetSectionWithString("单一已暂存", "annotation_status")
		//if err != nil {
		//	utils.LogDebug(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		//}
		//status11, err := enums.GetSectionWithString("已提交单一", "annotation_status")
		//if err != nil {
		//	utils.LogDebug(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		//}

		parseAnnotationRturns("annotation_return_path_c", "annotation_history_path_c")
		//
		//qs := o.QueryTable(models.AnnotationTBName()).Filter("status", status9)
		//
		//cond := orm.NewCondition()
		//var cond1 *orm.Condition
		//cond1 = cond.And("etps_inner_invt_no__in", sendPathCNames)
		//
		//qs = qs.SetCond(cond1)
		//mun, err := qs.Update(orm.Params{
		//	"status": status11,
		//})
		//
		//if err != nil {
		//	utils.LogDebug(fmt.Sprintf("annotationXmlPasre Update error:%v", err))
		//}
		//
		////ws 自动更新
		//if mun > 0 {
		//	msg := utils.Message{"清单状态更新", true}
		//	utils.Broadcast <- msg
		//}
		//
		return nil
	})

	err := task.Run()
	if err != nil {
		utils.LogDebug(fmt.Sprintf("tk.Run error :%v", err))
		return nil
	}

	return task
}

func parseAnnotationRturns(returnPathCofig, historyPathCofig string) {

	returnPath := beego.AppConfig.String(returnPathCofig)
	historyPath := beego.AppConfig.String(historyPathCofig)
	pathCfiles, err := ioutil.ReadDir(returnPath)
	if err != nil {
		utils.LogDebug(fmt.Sprintf("获取数据列表和总数 error:%v", err))
	}

	for _, f := range pathCfiles {
		fullPath := returnPath + f.Name()
		file, err := os.Open(fullPath) // For read access.
		if err != nil {
			utils.LogDebug(fmt.Sprintf("os.Open :%v", err))
			continue
		}

		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			utils.LogDebug(fmt.Sprintf(" ioutil.ReadAll :%v", err))
			continue
		}

		prefix, ext, failedName := getNameExts(f)
		if len(prefix) > 0 && len(ext) > 0 {
			if ext == "INV" {
				if prefix == "Successed" {
					v := xmlTemplate.CommonResponeMessage{}
					err := xml.Unmarshal(data, &v)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
						continue
					}

					annotation, err := models.GetAnnotationByEtpsInnerInvtNo(v.EtpsPreentNo)
					if err != nil {
						utils.LogDebug(fmt.Sprintf(" models.GetAnnotationByEtpsInnerInvtNo :%v", err))
						continue
					}

					aReturn := models.NewAnnotationReturn(0)
					aReturn.EtpsPreentNo = v.EtpsPreentNo
					aReturn.CheckInfo = v.CheckInfo
					aReturn.DealFlag = v.DealFlag
					aReturn.Annotation = annotation

					if err = models.AnnotationReturnSave(&aReturn); err != nil {
						utils.LogDebug(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
						continue
					}

					//更新状态
					if err = controllers.UpdateAnnotationStatus(annotation, "单一处理中"); err != nil {
						utils.LogDebug(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
					}

					fullHistoryPath := getHistoryPath(historyPath, v.EtpsPreentNo, f)
					err = os.Rename(fullPath, fullHistoryPath)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
						continue
					}

				} else if prefix == "Receipt" {
					fullHistoryPath := getHistoryPath(historyPath, "/Receipt/", f)
					err = os.Rename(fullPath, fullHistoryPath)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
						continue
					}
				}
			} else if ext == "INVT" && prefix == "Receipt" {
				v := xmlTemplate.ReturnPackage{}
				err := xml.Unmarshal(data, &v)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
					continue
				}

				annotation, err := models.GetAnnotationBySeqNo(v.InvPreentNo)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("models.GetAnnotationBySeqNo :%v", err))
					continue
				}

				aReturn := models.NewAnnotationReturn(0)
				aReturn.BusinessId = v.DataInfoBusinessId
				aReturn.ManageResult = v.ManageResult
				aReturn.Reason = v.Reason
				aReturn.Annotation = annotation

				if err = models.AnnotationReturnSave(&aReturn); err != nil {
					utils.LogDebug(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
					continue
				}

				//更新状态
				if err = controllers.UpdateAnnotationStatus(annotation, "已完成"); err != nil {
					utils.LogDebug(fmt.Sprintf("enums.UpdateAnnotationStatus :%v,filename:%v", err, f.Name()))
				}

				fullHistoryPath := getHistoryPath(historyPath, annotation.EtpsInnerInvtNo, f)
				err = os.Rename(fullPath, fullHistoryPath)
				if err != nil {
					utils.LogDebug(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
					continue
				}

			} else if ext == "xml" {
				if prefix == "Failed" {
					v := xmlTemplate.Failed{}
					err := xml.Unmarshal(data, &v)
					if err != nil {
						utils.LogDebug(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
						continue
					}

					annotation, err := models.GetAnnotationByEtpsInnerInvtNo(failedName)
					if err != nil {
						utils.LogDebug(fmt.Sprintf(" models.GetAnnotationByEtpsInnerInvtNo :%v", err))
						continue
					}

					aReturn := models.NewAnnotationReturn(0)
					aReturn.CheckInfo = v.FailInfo
					aReturn.DealFlag = v.FailCode
					aReturn.Annotation = annotation

					if err = models.AnnotationReturnSave(&aReturn); err != nil {
						utils.LogDebug(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
						continue
					}

					continue
				}
			}
		}

		fullHistoryPath := getHistoryPath(historyPath, "/Error/", f)
		err = os.Rename(fullPath, fullHistoryPath)
		if err != nil {
			utils.LogDebug(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
			continue
		}
		continue

	}
}

//历史目录
func getHistoryPath(historyPath, v string, f os.FileInfo) string {
	return historyPath + "/" + time.Now().Format(enums.BaseDateFormat) + "/" + v + "/" + f.Name()
}

//回执文件的前缀，和后缀
func getNameExts(f os.FileInfo) (string, string, string) {
	extNames := strings.Split(f.Name(), ".")
	if len(extNames) > 1 && len(extNames[1]) > 0 && len(extNames[0]) > 0 {
		eNames := strings.Split(extNames[0], "_")
		if len(eNames) > 1 && len(eNames[0]) > 0 {
			if len(eNames[1]) > 0 {
				mNames := strings.Split(extNames[0], "__")
				if len(mNames) > 1 && len(mNames[1]) > 0 {
					return eNames[0], extNames[1], mNames[1]
				} else {
					return eNames[0], extNames[1], ""
				}
			}

		} else {
			return "", extNames[1], ""
		}
	}

	return "", "", ""
}
