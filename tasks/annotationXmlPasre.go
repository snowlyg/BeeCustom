package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

//检测报文是否发送成功，并更新清单状态
func annotationUpdateAnnotationStatus() *toolbox.Task {
	task := toolbox.NewTask("task", "* * * * * *", func() error {
		o := orm.NewOrm()
		status9, err := enums.GetSectionWithString("复核通过", "annotation_status")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		}
		status11, err := enums.GetSectionWithString("已提交单一", "annotation_status")
		if err != nil {
			utils.LogDebug(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		}

		sendPathCNames, sendPathENames := []string{}, []string{}
		sendPathCNames = getAnnotationXmlNames(err, "annotation_send_xml_path_c", sendPathCNames)
		sendPathENames = getAnnotationXmlNames(err, "annotation_send_xml_path_e", sendPathENames)
		qs := o.QueryTable(models.AnnotationTBName()).Filter("status", status9)
		if len(sendPathCNames) > 0 || len(sendPathENames) > 0 {
			cond := orm.NewCondition()
			var cond1 *orm.Condition
			if len(sendPathCNames) > 0 {
				cond1 = cond.And("etps_inner_invt_no__in", sendPathCNames)
			} else {
				cond1 = cond.And("etps_inner_invt_no__in", sendPathENames)
			}

			qs = qs.SetCond(cond1)
		}

		_, err = qs.Update(orm.Params{
			"status": status11,
		})

		if err != nil {
			utils.LogDebug(fmt.Sprintf("annotationXmlPasre Update error:%v", err))
		}

		return nil
	})

	err := task.Run()
	if err != nil {
		utils.LogDebug(fmt.Sprintf("tk.Run error :%v", err))
		return nil
	}

	return task
}

func getAnnotationXmlNames(err error, pathConfig string, pathNames []string) []string {
	path := beego.AppConfig.String(pathConfig)
	pathCfiles, err := ioutil.ReadDir(path)
	for i, f := range pathCfiles {
		name := getAnnotationXmlName(f)
		if len(name) > 0 {
			pathNames[i] = name
		}
	}

	return pathNames
}

func getAnnotationXmlName(f os.FileInfo) string {
	names := strings.Split(f.Name(), "__")
	if len(names) > 2 && len(names[1]) > 0 {
		return names[1]
	}

	return ""
}
