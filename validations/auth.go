package validations

import (
	"strings"

	"github.com/astaxie/beego/validation"
)

func LoginValid(UserName, UserPwd string) string {

	valid := validation.Validation{}
	validation.MessageTmpls = GetMessageTmpls()

	valid.Required(UserName, "用户名")
	valid.MinSize(UserName, 5, "用户名")
	valid.MaxSize(UserName, 24, "用户名")

	valid.Required(UserPwd, "密码")
	valid.MinSize(UserPwd, 6, "密码")
	valid.MaxSize(UserPwd, 16, "密码")
	valid.AlphaNumeric(UserPwd, "密码")

	if valid.HasErrors() {
		if len(valid.Errors) > 0 {
			return strings.Split(valid.Errors[0].Key, ".")[0] + ":" + valid.Errors[0].Message
		}
	}

	return ""
}
