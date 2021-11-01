package errorsvc

import (
	"fmt"
	"testing"

	errorcode "github.com/ahl5esoft/go-skill/hide-ctx/model/enum/error-code"

	"github.com/stretchr/testify/assert"
)

func Test_custom_Error(t *testing.T) {
	self := new(custom)
	self.error = fmt.Errorf("test")

	res := self.Error()
	assert.Equal(t, res, "[err: test, code: 0, data: <nil>]")
}

func Test_custom_GetCode(t *testing.T) {
	self := new(custom)

	res := self.GetCode()
	assert.Equal(t, res, errorcode.Null)
}

func Test_New(t *testing.T) {
	self := New(errorcode.API, "test")
	assert.Equal(
		t,
		self.GetCode(),
		errorcode.API,
	)
	assert.Equal(
		t,
		self.GetData(),
		"test",
	)
}

func Test_Newf(t *testing.T) {
	self := Newf(errorcode.API, "%s-%s", "a", "b")
	assert.Equal(
		t,
		self.GetData(),
		"a-b",
	)
	assert.Equal(
		t,
		self.GetCode(),
		errorcode.API,
	)
}

func Test_NewError(t *testing.T) {
	self := NewError(
		errorcode.API,
		fmt.Errorf("err"),
	)
	assert.Equal(
		t,
		self.GetCode(),
		errorcode.API,
	)
	assert.Nil(
		t,
		self.GetData(),
	)
}
