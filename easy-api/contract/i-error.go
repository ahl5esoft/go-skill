package contract

import errorcode "github.com/ahl5esoft/go-skill/easy-api/model/enum/error-code"

type IError interface {
	error

	GetCode() errorcode.Value
	GetData() interface{}
}
