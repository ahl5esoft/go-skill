// Code generated by MockGen. DO NOT EDIT.
// Source: one\contract\i-io-factory.go

// Package contract is a generated GoMock package.
package contract

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIIOFactory is a mock of IIOFactory interface.
type MockIIOFactory struct {
	ctrl     *gomock.Controller
	recorder *MockIIOFactoryMockRecorder
}

// MockIIOFactoryMockRecorder is the mock recorder for MockIIOFactory.
type MockIIOFactoryMockRecorder struct {
	mock *MockIIOFactory
}

// NewMockIIOFactory creates a new mock instance.
func NewMockIIOFactory(ctrl *gomock.Controller) *MockIIOFactory {
	mock := &MockIIOFactory{ctrl: ctrl}
	mock.recorder = &MockIIOFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIIOFactory) EXPECT() *MockIIOFactoryMockRecorder {
	return m.recorder
}

// BuildDirectory mocks base method.
func (m *MockIIOFactory) BuildDirectory(pathArgs ...string) IIODirectory {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range pathArgs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BuildDirectory", varargs...)
	ret0, _ := ret[0].(IIODirectory)
	return ret0
}

// BuildDirectory indicates an expected call of BuildDirectory.
func (mr *MockIIOFactoryMockRecorder) BuildDirectory(pathArgs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildDirectory", reflect.TypeOf((*MockIIOFactory)(nil).BuildDirectory), pathArgs...)
}

// BuildFile mocks base method.
func (m *MockIIOFactory) BuildFile(pathArgs ...string) IIOFile {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range pathArgs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BuildFile", varargs...)
	ret0, _ := ret[0].(IIOFile)
	return ret0
}

// BuildFile indicates an expected call of BuildFile.
func (mr *MockIIOFactoryMockRecorder) BuildFile(pathArgs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildFile", reflect.TypeOf((*MockIIOFactory)(nil).BuildFile), pathArgs...)
}
