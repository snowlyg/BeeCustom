package controllers

import (
	"fmt"
	"time"

	"BeeCustom/enums"
	"BeeCustom/models"
	"BeeCustom/utils"
)

// 更新状态和状态更新时间
func UpdateAnnotationStatus(m *models.Annotation, StatusString string) error {
	aStatus, err := enums.GetSectionWithString(StatusString, "annotation_status")
	if err != nil {
		utils.LogDebug(fmt.Sprintf("转换清单状态出错:%v", err))
		return err
	}

	// 禁止状态回退
	if m.Status < aStatus {
		m.Status = aStatus
		m.StatusUpdatedAt = time.Now()
	}

	return nil
}
