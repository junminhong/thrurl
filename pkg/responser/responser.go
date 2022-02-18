package responser

import (
	"time"
)

var (
	Ok                 = add(0, "ok")
	ErrRequest         = add(400, "")
	ErrNotFind         = add(404, "")
	ErrForbidden       = add(403, "")
	ErrNoPermission    = add(405, "")
	ErrServer          = add(500, "")
	ReqBindErr         = add(1000, "請依照API文件進行請求")
	StoreShortenUrlErr = add(1001, "shorten url儲存失敗")
	StoreShortenUrlOk  = add(1002, "shorten url儲存成功")
)

func New(code int, msg string) ResponseFlag {
	if code < 1000 {
		panic("error code must be greater than 1000")
	}
	return add(code, msg)
}

func add(code int, msg string) ResponseFlag {
	return ResponseFlag{
		code: code, message: msg,
	}
}

func (responseFlag *ResponseFlag) Error() string {
	return responseFlag.message
}

func (responseFlag ResponseFlag) Message() string {
	return responseFlag.message
}

func (responseFlag ResponseFlag) Reload(message string) ResponseFlag {
	responseFlag.message = message
	return responseFlag
}

func (responseFlag ResponseFlag) Code() int {
	return responseFlag.code
}

func FormatResponse(resultCode int, message string, data interface{}, timeStamp time.Time) Response {
	return Response{
		ResultCode: resultCode,
		Message:    message,
		Data:       data,
		TimeStamp:  timeStamp,
	}
}

type Response struct {
	ResultCode int         `json:"result_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	TimeStamp  time.Time   `json:"time_stamp"`
}

type ResponseFlag struct {
	code    int
	message string
}

type ResponseFunc interface {
	Error() string
	Code() int
	Message() string
	Reload(string) ResponseFlag
}