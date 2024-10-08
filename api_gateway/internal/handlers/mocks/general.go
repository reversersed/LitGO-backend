// Code generated by MockGen. DO NOT EDIT.
// Source: general.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Info mocks base method.
func (m *MockLogger) Info(arg0 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), arg0...)
}

// Infof mocks base method.
func (m *MockLogger) Infof(format string, args ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockLoggerMockRecorder) Infof(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockLogger)(nil).Infof), varargs...)
}

// MockJwtMiddleware is a mock of JwtMiddleware interface.
type MockJwtMiddleware struct {
	ctrl     *gomock.Controller
	recorder *MockJwtMiddlewareMockRecorder
}

// MockJwtMiddlewareMockRecorder is the mock recorder for MockJwtMiddleware.
type MockJwtMiddlewareMockRecorder struct {
	mock *MockJwtMiddleware
}

// NewMockJwtMiddleware creates a new mock instance.
func NewMockJwtMiddleware(ctrl *gomock.Controller) *MockJwtMiddleware {
	mock := &MockJwtMiddleware{ctrl: ctrl}
	mock.recorder = &MockJwtMiddlewareMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJwtMiddleware) EXPECT() *MockJwtMiddlewareMockRecorder {
	return m.recorder
}

// GetCredentialsFromContext mocks base method.
func (m *MockJwtMiddleware) GetCredentialsFromContext(c *gin.Context) (*shared_pb.UserCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentialsFromContext", c)
	ret0, _ := ret[0].(*shared_pb.UserCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCredentialsFromContext indicates an expected call of GetCredentialsFromContext.
func (mr *MockJwtMiddlewareMockRecorder) GetCredentialsFromContext(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentialsFromContext", reflect.TypeOf((*MockJwtMiddleware)(nil).GetCredentialsFromContext), c)
}

// Middleware mocks base method.
func (m *MockJwtMiddleware) Middleware(arg0 ...string) gin.HandlerFunc {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Middleware", varargs...)
	ret0, _ := ret[0].(gin.HandlerFunc)
	return ret0
}

// Middleware indicates an expected call of Middleware.
func (mr *MockJwtMiddlewareMockRecorder) Middleware(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Middleware", reflect.TypeOf((*MockJwtMiddleware)(nil).Middleware), arg0...)
}
