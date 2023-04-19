package errcode

import (
	"fmt"
	"net/http"
)

// 抽象错误响应体 （数据抽象）
type Err struct {
	code    int      `json:"code"`
	message string   `json:"message"`
	details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Err {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码[%d]已存在，请勿重复定义", code))
	}
	codes[code] = msg
	return &Err{
		code:    code,
		message: msg,
	}
}

func (e *Err) ErrMsg() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.code, e.message)
}

func (e *Err) GetCode() int {
	return e.code
}

func (e *Err) GetMsg() string {
	return e.message
}

func (e *Err) MsgF(args []interface{}) string {
	return fmt.Sprintf(e.message, args...)
}

func (e *Err) GetDetails() []string {
	return e.details
}

func (e *Err) WithDetails(details ...string) *Err {
	e.details = []string{}
	// for _, d := range details {
	// 	e.details = append(e.details, d)
	// }
	e.details = append(e.details, details...)
	return e
}

func (e *Err) StatusCode() int {
	switch e.code {
	case Success.code:
		return http.StatusOK
	case ServerError.code:
		return http.StatusInternalServerError
	case InvalidParams.code:
		return http.StatusBadRequest
	case NotFound.code:
		return http.StatusNotFound
	case UnauthorizedAuthNotExist.code:
		fallthrough
	case UnauthorizedTokenError.code:
		fallthrough
	case UnauthorizedTokenGenerate.code:
		fallthrough
	case UnauthorizedTokenTimeout.code:
		return http.StatusUnauthorized
	case TooManyRequests.code:
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
