package api

import (
	mobile "github.com/ahl5esoft/go-skill/easy-api/api/mobile"
	"github.com/ahl5esoft/go-skill/easy-api/contract"
)

func Register(apiFactory contract.IApiFactory) {
	apiFactory.Register("mobile", "one", mobile.OneApi{})
	apiFactory.Register("mobile", "two", mobile.TwoApi{})
}
