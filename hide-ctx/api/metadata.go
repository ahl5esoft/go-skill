package api

import (
    "github.com/ahl5esoft/go-skill/hide-ctx/contract"
    mobile "github.com/ahl5esoft/go-skill/hide-ctx/api/mobile"
)

func Register(apiFactory contract.IApiFactory) {
    apiFactory.Register("mobile", "demo", mobile.DemoApi{})
}