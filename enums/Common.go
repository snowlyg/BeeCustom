package enums

type JsonResultCode int

const (
	JRCodeFailed JsonResultCode = iota
	JRCodeSucc
	JRCode302 = 302 //跳转至地址
	JRCode401 = 401 //未授权访问
)

const (
	Deleted = iota - 1
	Disabled
	Enabled
)
