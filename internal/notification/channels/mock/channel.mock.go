// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zitadel/zitadel/v2/internal/notification/channels (interfaces: NotificationChannel)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	channels "github.com/zitadel/zitadel/v2/internal/notification/channels"
	reflect "reflect"
)

// MockNotificationChannel is a mock of NotificationChannel interface
type MockNotificationChannel struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationChannelMockRecorder
}

// MockNotificationChannelMockRecorder is the mock recorder for MockNotificationChannel
type MockNotificationChannelMockRecorder struct {
	mock *MockNotificationChannel
}

// NewMockNotificationChannel creates a new mock instance
func NewMockNotificationChannel(ctrl *gomock.Controller) *MockNotificationChannel {
	mock := &MockNotificationChannel{ctrl: ctrl}
	mock.recorder = &MockNotificationChannelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNotificationChannel) EXPECT() *MockNotificationChannelMockRecorder {
	return m.recorder
}

// HandleMessage mocks base method
func (m *MockNotificationChannel) HandleMessage(arg0 channels.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleMessage indicates an expected call of HandleMessage
func (mr *MockNotificationChannelMockRecorder) HandleMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMessage", reflect.TypeOf((*MockNotificationChannel)(nil).HandleMessage), arg0)
}
