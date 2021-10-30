package errorsvc

import (
	"fmt"

	"github.com/ahl5esoft/go-skill/easy-api/contract"
	errorcode "github.com/ahl5esoft/go-skill/easy-api/model/enum/error-code"
)

type custom struct {
	error

	code errorcode.Value
	data interface{}
}

func (m custom) Error() string {
	return fmt.Sprintf("[err: %v, code: %v, data: %v]", m.error, m.code, m.data)
}

func (m custom) GetCode() errorcode.Value {
	return m.code
}

func (m custom) GetData() interface{} {
	return m.data
}

func New(code errorcode.Value, data interface{}) contract.IError {
	return custom{
		code: code,
		data: data,
	}
}

func Newf(code errorcode.Value, format string, args ...interface{}) contract.IError {
	return custom{
		error: fmt.Errorf(format, args...),
		code:  code,
		data:  fmt.Sprintf(format, args...),
	}
}

func NewError(code errorcode.Value, err error) contract.IError {
	return custom{
		error: err,
		code:  code,
	}
}
