package validation

import (
	"strings"

	"BeeCustom/models"
	"github.com/astaxie/beego/validation"
)

func freezeValid(m *models.BackendUser) string {

	valid := validation.Validation{}
	valid.Required(m.Status, "Status")

	if valid.HasErrors() {
		if len(valid.Errors) > 0 {
			return strings.Split(valid.Errors[0].Key, ".")[0] + ":" + valid.Errors[0].Message

		}

	}

	return ""

}
