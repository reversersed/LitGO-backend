// Code generated by MockGen. DO NOT EDIT.
// Source: server.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
)

// Mockvalidator is a mock of validator interface.
type Mockvalidator struct {
	ctrl     *gomock.Controller
	recorder *MockvalidatorMockRecorder
}

// MockvalidatorMockRecorder is the mock recorder for Mockvalidator.
type MockvalidatorMockRecorder struct {
	mock *Mockvalidator
}

// NewMockvalidator creates a new mock instance.
func NewMockvalidator(ctrl *gomock.Controller) *Mockvalidator {
	mock := &Mockvalidator{ctrl: ctrl}
	mock.recorder = &MockvalidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockvalidator) EXPECT() *MockvalidatorMockRecorder {
	return m.recorder
}

// StructValidation mocks base method.
func (m *Mockvalidator) StructValidation(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StructValidation", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StructValidation indicates an expected call of StructValidation.
func (mr *MockvalidatorMockRecorder) StructValidation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StructValidation", reflect.TypeOf((*Mockvalidator)(nil).StructValidation), arg0)
}

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *Mocklogger) Error(arg0 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockloggerMockRecorder) Error(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mocklogger)(nil).Error), arg0...)
}

// Errorf mocks base method.
func (m *Mocklogger) Errorf(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockloggerMockRecorder) Errorf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*Mocklogger)(nil).Errorf), varargs...)
}

// Info mocks base method.
func (m *Mocklogger) Info(arg0 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockloggerMockRecorder) Info(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklogger)(nil).Info), arg0...)
}

// Infof mocks base method.
func (m *Mocklogger) Infof(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockloggerMockRecorder) Infof(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*Mocklogger)(nil).Infof), varargs...)
}

// Warn mocks base method.
func (m *Mocklogger) Warn(arg0 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockloggerMockRecorder) Warn(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*Mocklogger)(nil).Warn), arg0...)
}

// Warnf mocks base method.
func (m *Mocklogger) Warnf(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warnf", varargs...)
}

// Warnf indicates an expected call of Warnf.
func (mr *MockloggerMockRecorder) Warnf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warnf", reflect.TypeOf((*Mocklogger)(nil).Warnf), varargs...)
}

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage.
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance.
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return m.recorder
}

// GetSuggestions mocks base method.
func (m *Mockstorage) GetSuggestions(ctx context.Context, regex string, limit int64) ([]*storage.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSuggestions", ctx, regex, limit)
	ret0, _ := ret[0].([]*storage.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSuggestions indicates an expected call of GetSuggestions.
func (mr *MockstorageMockRecorder) GetSuggestions(ctx, regex, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSuggestions", reflect.TypeOf((*Mockstorage)(nil).GetSuggestions), ctx, regex, limit)
}

// Mockcache is a mock of cache interface.
type Mockcache struct {
	ctrl     *gomock.Controller
	recorder *MockcacheMockRecorder
}

// MockcacheMockRecorder is the mock recorder for Mockcache.
type MockcacheMockRecorder struct {
	mock *Mockcache
}

// NewMockcache creates a new mock instance.
func NewMockcache(ctrl *gomock.Controller) *Mockcache {
	mock := &Mockcache{ctrl: ctrl}
	mock.recorder = &MockcacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockcache) EXPECT() *MockcacheMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *Mockcache) Delete(arg0 []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockcacheMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*Mockcache)(nil).Delete), arg0)
}

// Get mocks base method.
func (m *Mockcache) Get(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockcacheMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*Mockcache)(nil).Get), arg0)
}

// Set mocks base method.
func (m *Mockcache) Set(arg0, arg1 []byte, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockcacheMockRecorder) Set(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*Mockcache)(nil).Set), arg0, arg1, arg2)
}
