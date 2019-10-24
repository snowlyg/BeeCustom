package validations

func GetMessageTmpls() map[string]string {
	return map[string]string{
		"Required":     "不能为空",
		"Min":          "最小为 %d",
		"Max":          "最大为 %d",
		"Range":        "范围在 %d 至 %d",
		"MinSize":      "最小长度为 %d",
		"MaxSize":      "最大长度为 %d",
		"Length":       "长度必须是 %d",
		"Alpha":        "必须是有效的字母字符",
		"Numeric":      "必须是有效的数字字符",
		"AlphaNumeric": "必须是有效的字母或数字字符",
		"Match":        "必须匹配格式 %s",
		"NoMatch":      "必须不匹配格式 %s",
		"AlphaDash":    "必须是有效的字母或数字或破折号(-_)字符",
		"Email":        "必须是有效的邮件地址",
		"IP":           "必须是有效的IP地址",
		"Base64":       "必须是有效的base64字符",
		"Mobile":       "必须是有效手机号码",
		"Tel":          "必须是有效电话号码",
		"Phone":        "必须是有效的电话号码或者手机号码",
		"ZipCode":      "必须是有效的邮政编码",
	}
}
