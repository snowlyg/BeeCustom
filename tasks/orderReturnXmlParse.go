package tasks

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"BeeCustom/controllers"
	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xmlTemplate"
	"github.com/astaxie/beego/toolbox"
)

// 回执解释
func orderCReturnXmlParse() *toolbox.Task {

	task := toolbox.NewTask("task", "* * * * * *", func() error {
		parseOrderReturns("order_return_path_c", "order_history_path_c")
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
func parseOrderReturns(returnPathConfig, historyPathConfig string) {

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
		if ext != "xml" {
			utils.LogError(fmt.Sprintf("openFile :%v ,filename:%v", errors.New("未知文件格式"), f.Name()))
			continue
		}

		if prefix == "Successed" {

			xmlFile, err, data := openFile(fullPath)
			if err != nil {
				utils.LogError(fmt.Sprintf("openFile :%v ,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			v := xmlTemplate.DecImportResponse{}
			err = xml.Unmarshal(data, &v)
			if err != nil {
				utils.LogError(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			order, err := models.GetOrderByClientSeqNo(v.ClientSeqNo)
			if err != nil {
				xmlFile.Close()
				continue
			}

			aReturn := models.NewOrderReturn(0)
			aReturn.Channel = v.ResponseCode
			aReturn.CusCiqNo = v.SeqNo
			aReturn.EntryId = ""
			aReturn.Note = v.ErrorMessage
			aReturn.NoticeDate = time.Now()
			aReturn.Remark = ""
			aReturn.Order = order

			if err = models.OrderReturnSave(&aReturn); err != nil {
				utils.LogError(fmt.Sprintf("models.OrderReturnSave :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			// 更新状态
			order.SeqNo = v.SeqNo
			if err = controllers.UpdateOrderStatus(order, "单一处理中", false); err != nil {
				utils.LogError(fmt.Sprintf("controllers.UpdateOrderStatus :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			if err := models.OrderUpdate(order, []string{"Status", "StatusUpdatedAt", "SeqNo"}); err != nil {
				utils.LogError(fmt.Sprintf("OrderUpdateStatus :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			xmlFile.Close()
			err = moveFile(historyPath, v.SeqNo, fullPath, f)
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

			xmlFile, err, data := openFile(fullPath)
			if err != nil {
				utils.LogError(fmt.Sprintf("openFile :%v ,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			v := xmlTemplate.DecResult{}
			err = xml.Unmarshal(data, &v)
			if err != nil {
				utils.LogError(fmt.Sprintf(" xml.Unmarshal :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			order, err := models.GetOrderBySeqNo(v.CusCiqNo)
			if err != nil {
				xmlFile.Close()
				continue
			}

			note := v.Note
			if v.CustomMaster != "" {
				note = note + "申报地海关:" + v.CustomMaster
			}
			if v.IEDate != "" {
				note = note + "进/出口日期:" + v.IEDate
			}
			if v.DDate != "" {
				note = note + "DDate:" + v.DDate
			}

			aReturn := models.NewOrderReturn(0)
			aReturn.EntryId = v.EntryId
			aReturn.NoticeDate, _ = time.Parse(v.NoticeDate, enums.RFC3339)
			aReturn.Note = note
			aReturn.Channel = v.Channel
			aReturn.Order = order

			if err = models.OrderReturnSave(&aReturn); err != nil {
				utils.LogError(fmt.Sprintf("models.OrderReturnSave :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			// 更新状态
			if err = controllers.UpdateOrderStatus(order, "已完成", false); err != nil {
				utils.LogError(fmt.Sprintf("enums.UpdateOrderStatus :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}
			order.EntryId = v.EntryId
			if err := models.OrderUpdate(order, []string{"Status", "StatusUpdatedAt", "EntryId"}); err != nil {

				utils.LogError(fmt.Sprintf("OrderUpdateStatus :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			xmlFile.Close()
			err = moveFile(historyPath, order.ClientSeqNo, fullPath, f)
			if err != nil {
				utils.LogError(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
			}

			// ws 自动更新
			msg := utils.Message{Message: "清单状态更新", IsUpdated: true}
			utils.Broadcast <- msg

		} else if prefix == "Failed" {

			xmlFile, err, data := openFile(fullPath)
			if err != nil {
				utils.LogError(fmt.Sprintf("openFile :%v ,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			v := xmlTemplate.DecImportResponse{}
			err = xml.Unmarshal(data, &v)
			if err != nil {
				utils.LogError(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}

			order, err := models.GetOrderByClientSeqNo(failedName)
			if err != nil {
				xmlFile.Close()
				continue
			}

			aReturn := models.NewOrderReturn(0)
			aReturn.Channel = v.ResponseCode
			aReturn.EntryId = ""
			aReturn.Note = v.ErrorMessage
			aReturn.NoticeDate = time.Now()
			aReturn.Remark = ""
			aReturn.Order = order

			if err = models.OrderReturnSave(&aReturn); err != nil {
				utils.LogError(fmt.Sprintf("models.OrderReturnSave :%v,filename:%v", err, f.Name()))
				xmlFile.Close()
				continue
			}
		}

		// 错误文件移动
		err = moveFile(historyPath, "Error", fullPath, f)
		if err != nil {
			utils.LogError(fmt.Sprintf("moveFile:%v,filename:%v", err, f.Name()))
		}
	}
}
