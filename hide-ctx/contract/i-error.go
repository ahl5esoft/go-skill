package contract

import errorcode "github.com/ahl5esoft/go-skill/hide-ctx/model/enum/error-code"

type IError interface {
	error

	GetCode() errorcode.Value
	GetData() interface{}
}
