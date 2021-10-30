package apisvc

import (
	"reflect"

	"github.com/ahl5esoft/go-skill/easy-api/contract"
	errorcode "github.com/ahl5esoft/go-skill/easy-api/model/enum/error-code"
	"github.com/ahl5esoft/go-skill/easy-api/service/errorsvc"
)

var (
	errNilApi = errorsvc.Newf(errorcode.API, "")
	nilApiPtr = &nilApi{}
)

type factory map[string]map[string]reflect.Type

func (m factory) Build(endpoint, api string) contract.IApi {
	if apiTypes, ok := m[endpoint]; ok {
		if apiType, ok := apiTypes[api]; ok {
			return reflect.New(apiType).Interface().(contract.IApi)
		}
	}

	return nilApiPtr
}

func (m factory) Register(endpoint, api string, apiInstance contract.IApi) {
	if _, ok := m[endpoint]; !ok {
		m[endpoint] = make(map[string]reflect.Type)
	}

	apiType := reflect.TypeOf(apiInstance)
	if apiType.Kind() == reflect.Ptr {
		apiType = apiType.Elem()
	}
	m[endpoint][api] = apiType
}

type nilApi struct{}

func (m nilApi) Call() (interface{}, error) {
	return nil, errNilApi
}

func NewFactory() contract.IApiFactory {
	return make(factory)
}
