package apisvc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testBuildAPI struct{}

func (m testBuildAPI) Call() (interface{}, error) {
	return nil, nil
}

type testRegisterAPI struct{}

func (m testRegisterAPI) Call() (interface{}, error) {
	return nil, nil
}

func Test_factory_Build(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		endpoint := "endpoint"
		apiName := "api-name"
		self := make(factory)
		self[endpoint] = map[string]reflect.Type{
			apiName: reflect.TypeOf(testBuildAPI{}),
		}

		res := self.Build(endpoint, apiName)
		assert.Equal(
			t,
			reflect.TypeOf(res).Elem(),
			reflect.TypeOf(testBuildAPI{}),
		)
	})

	t.Run("nil", func(t *testing.T) {
		endpoint := "endpoint"
		apiName := "api-name"
		self := make(factory)
		res := self.Build(endpoint, apiName)
		assert.Equal(t, res, nilApiPtr)
	})
}

func Test_factory_Register(t *testing.T) {
	endpoint := "endpoint"
	apiName := "api-name"
	self := make(factory)
	self.Register(endpoint, apiName, testRegisterAPI{})

	apiTypes, ok := self[endpoint]
	assert.True(t, ok)

	apiType, ok := apiTypes[apiName]
	assert.True(t, ok)
	assert.Equal(
		t,
		apiType,
		reflect.TypeOf(testRegisterAPI{}),
	)
}
