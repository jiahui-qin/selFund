package app

var MsgFlags = map[string]string{
	"SUCCESS":             "ok",
	"ERROR":               "fail",
	"DATA_INSERT_ERROR":   "数据插入错误",
	"PARA_ANALYSIS_ERROR": "参数解析错误",
}

// GetMsg get error information based on Code
func GetMsg(errMsg string) string {
	msg, ok := MsgFlags[errMsg]
	if ok {
		return msg
	}

	return MsgFlags[errMsg]
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (*Response) Error() string {
	return "new error response"
}
