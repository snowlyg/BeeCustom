package validations

import (
	"strings"

	"github.com/astaxie/beego/validation"
)

func RoleValid(permIds string) string {

	valid := validation.Validation{}
	valid.Required(permIds, "权限集合")

	if valid.HasErrors() {
		if len(valid.Errors) > 0 {
			return strings.Split(valid.Errors[0].Key, ".")[0] + ":" + valid.Errors[0].Message
		}
	}
	return ""
}
