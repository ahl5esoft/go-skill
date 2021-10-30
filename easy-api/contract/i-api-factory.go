package contract

type IApiFactory interface {
	Build(endpoint, apiName string) IApi
	Register(endpoint, apiName string, apiInstance IApi)
}
