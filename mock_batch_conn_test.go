// Code generated by MockGen. DO NOT EDIT.
// Source: sys_conn_oob.go
//
// Generated by this command:
//
//	mockgen -typed -package quic -self_package github.com/fanglinhan/fork-quic-go -source sys_conn_oob.go -destination mock_batch_conn_test.go -mock_names batchConn=MockBatchConn
//

// Package quic is a generated GoMock package.
package quic

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	ipv4 "golang.org/x/net/ipv4"
)

// MockBatchConn is a mock of batchConn interface.
type MockBatchConn struct {
	ctrl     *gomock.Controller
	recorder *MockBatchConnMockRecorder
}

// MockBatchConnMockRecorder is the mock recorder for MockBatchConn.
type MockBatchConnMockRecorder struct {
	mock *MockBatchConn
}

// NewMockBatchConn creates a new mock instance.
func NewMockBatchConn(ctrl *gomock.Controller) *MockBatchConn {
	mock := &MockBatchConn{ctrl: ctrl}
	mock.recorder = &MockBatchConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBatchConn) EXPECT() *MockBatchConnMockRecorder {
	return m.recorder
}

// ReadBatch mocks base method.
func (m *MockBatchConn) ReadBatch(ms []ipv4.Message, flags int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadBatch", ms, flags)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadBatch indicates an expected call of ReadBatch.
func (mr *MockBatchConnMockRecorder) ReadBatch(ms, flags any) *MockBatchConnReadBatchCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadBatch", reflect.TypeOf((*MockBatchConn)(nil).ReadBatch), ms, flags)
	return &MockBatchConnReadBatchCall{Call: call}
}

// MockBatchConnReadBatchCall wrap *gomock.Call
type MockBatchConnReadBatchCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBatchConnReadBatchCall) Return(arg0 int, arg1 error) *MockBatchConnReadBatchCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBatchConnReadBatchCall) Do(f func([]ipv4.Message, int) (int, error)) *MockBatchConnReadBatchCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBatchConnReadBatchCall) DoAndReturn(f func([]ipv4.Message, int) (int, error)) *MockBatchConnReadBatchCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
