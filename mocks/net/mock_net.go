// Automatically generated by MockGen. DO NOT EDIT!
// Source: net (interfaces: Addr)

package net

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Addr interface
type MockAddr struct {
	ctrl     *gomock.Controller
	recorder *_MockAddrRecorder
}

// Recorder for MockAddr (not exported)
type _MockAddrRecorder struct {
	mock *MockAddr
}

func NewMockAddr(ctrl *gomock.Controller) *MockAddr {
	mock := &MockAddr{ctrl: ctrl}
	mock.recorder = &_MockAddrRecorder{mock}
	return mock
}

func (_m *MockAddr) EXPECT() *_MockAddrRecorder {
	return _m.recorder
}

func (_m *MockAddr) Network() string {
	ret := _m.ctrl.Call(_m, "Network")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockAddrRecorder) Network() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Network")
}

func (_m *MockAddr) String() string {
	ret := _m.ctrl.Call(_m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockAddrRecorder) String() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "String")
}
