package tasks

import (
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

//检测报文是否发送成功，并更新清单状态
func annotationUpdateAnnotationStatus() *toolbox.Task {

	task := toolbox.NewTask("task", "* * * * * *", func() error {
		o := orm.NewOrm()
		aStatusS, err := models.GetSettingRValueByKey("annotationStatus", false)
		status9, err, _ := enums.TransformCnToInt(aStatusS, "复核通过")
		if err != nil {
			utils.LogError(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		}

		status11, err, _ := enums.TransformCnToInt(aStatusS, "已提交单一")
		if err != nil {
			utils.LogError(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		}

		sendPathCNames, sendPathENames := []string{}, []string{}
		sendPathCNames = getAnnotationXmlNames("annotation_send_xml_path_c", sendPathCNames)
		sendPathENames = getAnnotationXmlNames("annotation_send_xml_path_e", sendPathENames)
		qs := o.QueryTable(models.AnnotationTBName()).Filter("Status", status9)
		if (len(sendPathCNames) > 0 && len(sendPathCNames[0]) > 0) || (len(sendPathENames) > 0 && len(sendPathENames[0]) > 0) {

			if len(sendPathCNames[0]) > 0 {
				qs = qs.Filter("etps_inner_invt_no__in", sendPathCNames)
			} else {
				qs = qs.Filter("etps_inner_invt_no__in", sendPathENames)
			}

			mun, err := qs.Update(orm.Params{
				"status": status11,
			})

			if err != nil {
				utils.LogError(fmt.Sprintf("annotationXmlPasre Update error:%v", err))
			}

			//ws 自动更新
			if mun > 0 {
				msg := utils.Message{"清单状态更新", true}
				utils.Broadcast <- msg
			}
		}

		return nil
	})

	err := task.Run()
	if err != nil {
		utils.LogError(fmt.Sprintf("tk.Run error :%v", err))
		return nil
	}

	return task
}
