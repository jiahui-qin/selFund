package service

type Exception struct {
	Code      int    `json:"-"`
	ErrorCode int    `json:"error_code"`
	Msg       string `json:"msg"`
	Request   string `json:"request"`
}

func (e *Exception) Error() string {
	return e.Msg
}

func newException(code int, errorCode int, msg string) *Exception {
	return &Exception{
		Code:      code,
		ErrorCode: errorCode,
		Msg:       msg,
	}
}
