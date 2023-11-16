package errcode

import "encoding/json"

// 错误码结构体
type ErrCode struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func NewErrCode(code int32, msg string) *ErrCode {
	return &ErrCode{
		Code: code,
		Msg:  msg,
	}
}

func (err *ErrCode) ToJSON() []byte {
	str, _ := json.Marshal(err)
	return str
}

func (err *ErrCode) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, err)
}

func (err *ErrCode) ISOK() bool {
	if err.Code == OK.Code {
		return true
	} else {
		return false
	}
}
