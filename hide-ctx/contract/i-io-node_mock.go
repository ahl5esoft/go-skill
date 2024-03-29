// Code generated by MockGen. DO NOT EDIT.
// Source: hide-ctx\contract\i-io-node.go

// Package contract is a generated GoMock package.
package contract

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIIONode is a mock of IIONode interface.
type MockIIONode struct {
	ctrl     *gomock.Controller
	recorder *MockIIONodeMockRecorder
}

// MockIIONodeMockRecorder is the mock recorder for MockIIONode.
type MockIIONodeMockRecorder struct {
	mock *MockIIONode
}

// NewMockIIONode creates a new mock instance.
func NewMockIIONode(ctrl *gomock.Controller) *MockIIONode {
	mock := &MockIIONode{ctrl: ctrl}
	mock.recorder = &MockIIONodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIIONode) EXPECT() *MockIIONodeMockRecorder {
	return m.recorder
}

// GetName mocks base method.
func (m *MockIIONode) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockIIONodeMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockIIONode)(nil).GetName))
}

// GetParent mocks base method.
func (m *MockIIONode) GetParent() IIODirectory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParent")
	ret0, _ := ret[0].(IIODirectory)
	return ret0
}

// GetParent indicates an expected call of GetParent.
func (mr *MockIIONodeMockRecorder) GetParent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParent", reflect.TypeOf((*MockIIONode)(nil).GetParent))
}

// GetPath mocks base method.
func (m *MockIIONode) GetPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPath indicates an expected call of GetPath.
func (mr *MockIIONodeMockRecorder) GetPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPath", reflect.TypeOf((*MockIIONode)(nil).GetPath))
}

// IsExist mocks base method.
func (m *MockIIONode) IsExist() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExist")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsExist indicates an expected call of IsExist.
func (mr *MockIIONodeMockRecorder) IsExist() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExist", reflect.TypeOf((*MockIIONode)(nil).IsExist))
}
