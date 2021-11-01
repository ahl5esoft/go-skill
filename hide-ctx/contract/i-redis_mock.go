// Code generated by MockGen. DO NOT EDIT.
// Source: hide-ctx\contract\i-redis.go

// Package contract is a generated GoMock package.
package contract

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockIRedis is a mock of IRedis interface.
type MockIRedis struct {
	ctrl     *gomock.Controller
	recorder *MockIRedisMockRecorder
}

// MockIRedisMockRecorder is the mock recorder for MockIRedis.
type MockIRedisMockRecorder struct {
	mock *MockIRedis
}

// NewMockIRedis creates a new mock instance.
func NewMockIRedis(ctrl *gomock.Controller) *MockIRedis {
	mock := &MockIRedis{ctrl: ctrl}
	mock.recorder = &MockIRedisMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRedis) EXPECT() *MockIRedisMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockIRedis) Del(arg0 ...string) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Del", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Del indicates an expected call of Del.
func (mr *MockIRedisMockRecorder) Del(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockIRedis)(nil).Del), arg0...)
}

// Get mocks base method.
func (m *MockIRedis) Get(k string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", k)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIRedisMockRecorder) Get(k interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIRedis)(nil).Get), k)
}

// Set mocks base method.
func (m *MockIRedis) Set(k, v string, expires time.Duration) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", k, v, expires)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockIRedisMockRecorder) Set(k, v, expires interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockIRedis)(nil).Set), k, v, expires)
}
