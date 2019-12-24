package tasks

import (
	"fmt"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

//检测报文是否发送成功，并更新货物状态
func orderUpdateStatus() *toolbox.Task {

	task := toolbox.NewTask("task", "* * * * * *", func() error {
		o := orm.NewOrm()
		aStatusS, err := models.GetSettingRValueByKey("orderStatus", false)
		status9, err, _ := enums.TransformCnToInt(aStatusS, "复核通过")
		if err != nil {
			utils.LogError(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		}

		status11, err, _ := enums.TransformCnToInt(aStatusS, "已提交单一")
		if err != nil {
			utils.LogError(fmt.Sprintf("获取数据列表和总数 error:%v", err))
		}

		var sendPathNames []string
		sendPathNames = getXmlNames("order_send_xml_path", sendPathNames)
		qs := o.QueryTable(models.OrderTBName()).Filter("Status", status9)
		if len(sendPathNames) > 0 && len(sendPathNames[0]) > 0 {
			if len(sendPathNames[0]) > 0 {
				qs = qs.Filter("client_seq_no__in", sendPathNames)
			}

			mun, err := qs.Update(orm.Params{
				"status": status11,
			})

			if err != nil {
				utils.LogError(fmt.Sprintf("orderXmlPasre Update error:%v", err))
			}

			//ws 自动更新
			if mun > 0 {
				msg := utils.Message{Message: "货物状态更新", IsUpdated: true}
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
