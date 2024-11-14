// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/fanglinhan/fork-quic-go (interfaces: StreamGetter)
//
// Generated by this command:
//
//	mockgen -typed -build_flags=-tags=gomock -package quic -self_package github.com/fanglinhan/fork-quic-go -destination mock_stream_getter_test.go github.com/fanglinhan/fork-quic-go StreamGetter
//

// Package quic is a generated GoMock package.
package quic

import (
	reflect "reflect"

	protocol "github.com/fanglinhan/fork-quic-go/internal/protocol"
	gomock "go.uber.org/mock/gomock"
)

// MockStreamGetter is a mock of StreamGetter interface.
type MockStreamGetter struct {
	ctrl     *gomock.Controller
	recorder *MockStreamGetterMockRecorder
}

// MockStreamGetterMockRecorder is the mock recorder for MockStreamGetter.
type MockStreamGetterMockRecorder struct {
	mock *MockStreamGetter
}

// NewMockStreamGetter creates a new mock instance.
func NewMockStreamGetter(ctrl *gomock.Controller) *MockStreamGetter {
	mock := &MockStreamGetter{ctrl: ctrl}
	mock.recorder = &MockStreamGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStreamGetter) EXPECT() *MockStreamGetterMockRecorder {
	return m.recorder
}

// GetOrOpenReceiveStream mocks base method.
func (m *MockStreamGetter) GetOrOpenReceiveStream(arg0 protocol.StreamID) (receiveStreamI, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrOpenReceiveStream", arg0)
	ret0, _ := ret[0].(receiveStreamI)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrOpenReceiveStream indicates an expected call of GetOrOpenReceiveStream.
func (mr *MockStreamGetterMockRecorder) GetOrOpenReceiveStream(arg0 any) *MockStreamGetterGetOrOpenReceiveStreamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrOpenReceiveStream", reflect.TypeOf((*MockStreamGetter)(nil).GetOrOpenReceiveStream), arg0)
	return &MockStreamGetterGetOrOpenReceiveStreamCall{Call: call}
}

// MockStreamGetterGetOrOpenReceiveStreamCall wrap *gomock.Call
type MockStreamGetterGetOrOpenReceiveStreamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStreamGetterGetOrOpenReceiveStreamCall) Return(arg0 receiveStreamI, arg1 error) *MockStreamGetterGetOrOpenReceiveStreamCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStreamGetterGetOrOpenReceiveStreamCall) Do(f func(protocol.StreamID) (receiveStreamI, error)) *MockStreamGetterGetOrOpenReceiveStreamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStreamGetterGetOrOpenReceiveStreamCall) DoAndReturn(f func(protocol.StreamID) (receiveStreamI, error)) *MockStreamGetterGetOrOpenReceiveStreamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetOrOpenSendStream mocks base method.
func (m *MockStreamGetter) GetOrOpenSendStream(arg0 protocol.StreamID) (sendStreamI, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrOpenSendStream", arg0)
	ret0, _ := ret[0].(sendStreamI)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrOpenSendStream indicates an expected call of GetOrOpenSendStream.
func (mr *MockStreamGetterMockRecorder) GetOrOpenSendStream(arg0 any) *MockStreamGetterGetOrOpenSendStreamCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrOpenSendStream", reflect.TypeOf((*MockStreamGetter)(nil).GetOrOpenSendStream), arg0)
	return &MockStreamGetterGetOrOpenSendStreamCall{Call: call}
}

// MockStreamGetterGetOrOpenSendStreamCall wrap *gomock.Call
type MockStreamGetterGetOrOpenSendStreamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStreamGetterGetOrOpenSendStreamCall) Return(arg0 sendStreamI, arg1 error) *MockStreamGetterGetOrOpenSendStreamCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStreamGetterGetOrOpenSendStreamCall) Do(f func(protocol.StreamID) (sendStreamI, error)) *MockStreamGetterGetOrOpenSendStreamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStreamGetterGetOrOpenSendStreamCall) DoAndReturn(f func(protocol.StreamID) (sendStreamI, error)) *MockStreamGetterGetOrOpenSendStreamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
