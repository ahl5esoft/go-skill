package response

import errorcode "github.com/ahl5esoft/go-skill/hide-ctx/model/enum/error-code"

type Api struct {
	Data  interface{}     `json:"data"`
	Error errorcode.Value `json:"err"`
}
