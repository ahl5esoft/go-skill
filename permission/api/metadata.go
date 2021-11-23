package api

import (
    "github.com/ahl5esoft/lite-go/contract"
    mobile "github.com/ahl5esoft/go-skill/permission/api/mobile"
)

func Register(apiFactory contract.IApiFactory) {
    apiFactory.Register("mobile", "enter", mobile.EnterApi{})
    apiFactory.Register("mobile", "play", mobile.PlayApi{})
    apiFactory.Register("mobile", "quick-play", mobile.QuickPlayApi{})
}