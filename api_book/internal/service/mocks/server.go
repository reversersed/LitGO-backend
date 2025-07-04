// Code generated by MockGen. DO NOT EDIT.
// Source: server.go
//
// Generated by this command:
//
//	mockgen -source=server.go -destination=mocks/server.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	storage "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/storage"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

// Mockvalidator is a mock of validator interface.
type Mockvalidator struct {
	ctrl     *gomock.Controller
	recorder *MockvalidatorMockRecorder
	isgomock struct{}
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
func (mr *MockvalidatorMockRecorder) StructValidation(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StructValidation", reflect.TypeOf((*Mockvalidator)(nil).StructValidation), arg0)
}

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
	isgomock struct{}
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
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockloggerMockRecorder) Error(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mocklogger)(nil).Error), arg0...)
}

// Errorf mocks base method.
func (m *Mocklogger) Errorf(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockloggerMockRecorder) Errorf(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*Mocklogger)(nil).Errorf), varargs...)
}

// Info mocks base method.
func (m *Mocklogger) Info(arg0 ...any) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockloggerMockRecorder) Info(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklogger)(nil).Info), arg0...)
}

// Infof mocks base method.
func (m *Mocklogger) Infof(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockloggerMockRecorder) Infof(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*Mocklogger)(nil).Infof), varargs...)
}

// Warn mocks base method.
func (m *Mocklogger) Warn(arg0 ...any) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockloggerMockRecorder) Warn(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*Mocklogger)(nil).Warn), arg0...)
}

// Warnf mocks base method.
func (m *Mocklogger) Warnf(arg0 string, arg1 ...any) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warnf", varargs...)
}

// Warnf indicates an expected call of Warnf.
func (mr *MockloggerMockRecorder) Warnf(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warnf", reflect.TypeOf((*Mocklogger)(nil).Warnf), varargs...)
}

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
	isgomock struct{}
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

// CreateBook mocks base method.
func (m *Mockstorage) CreateBook(arg0 context.Context, arg1 *storage.Book) (*storage.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBook", arg0, arg1)
	ret0, _ := ret[0].(*storage.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBook indicates an expected call of CreateBook.
func (mr *MockstorageMockRecorder) CreateBook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBook", reflect.TypeOf((*Mockstorage)(nil).CreateBook), arg0, arg1)
}

// Find mocks base method.
func (m *Mockstorage) Find(arg0 context.Context, arg1 string, arg2, arg3 int, arg4 float32, arg5 storage.SortType) ([]*storage.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]*storage.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockstorageMockRecorder) Find(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*Mockstorage)(nil).Find), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetBook mocks base method.
func (m *Mockstorage) GetBook(arg0 context.Context, arg1 string) (*storage.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBook", arg0, arg1)
	ret0, _ := ret[0].(*storage.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBook indicates an expected call of GetBook.
func (mr *MockstorageMockRecorder) GetBook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBook", reflect.TypeOf((*Mockstorage)(nil).GetBook), arg0, arg1)
}

// GetBookByAuthor mocks base method.
func (m *Mockstorage) GetBookByAuthor(ctx context.Context, authorId primitive.ObjectID, page, limit int) ([]*storage.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookByAuthor", ctx, authorId, page, limit)
	ret0, _ := ret[0].([]*storage.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookByAuthor indicates an expected call of GetBookByAuthor.
func (mr *MockstorageMockRecorder) GetBookByAuthor(ctx, authorId, page, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookByAuthor", reflect.TypeOf((*Mockstorage)(nil).GetBookByAuthor), ctx, authorId, page, limit)
}

// GetBookByGenre mocks base method.
func (m *Mockstorage) GetBookByGenre(arg0 context.Context, arg1 []primitive.ObjectID, arg2 storage.SortType, arg3 bool, arg4, arg5 int) ([]*storage.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookByGenre", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]*storage.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookByGenre indicates an expected call of GetBookByGenre.
func (mr *MockstorageMockRecorder) GetBookByGenre(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookByGenre", reflect.TypeOf((*Mockstorage)(nil).GetBookByGenre), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetBookList mocks base method.
func (m *Mockstorage) GetBookList(arg0 context.Context, arg1 []primitive.ObjectID, arg2 []string) ([]*storage.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookList", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*storage.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookList indicates an expected call of GetBookList.
func (mr *MockstorageMockRecorder) GetBookList(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookList", reflect.TypeOf((*Mockstorage)(nil).GetBookList), arg0, arg1, arg2)
}

// Mockcache is a mock of cache interface.
type Mockcache struct {
	ctrl     *gomock.Controller
	recorder *MockcacheMockRecorder
	isgomock struct{}
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
func (mr *MockcacheMockRecorder) Delete(arg0 any) *gomock.Call {
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
func (mr *MockcacheMockRecorder) Get(arg0 any) *gomock.Call {
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
func (mr *MockcacheMockRecorder) Set(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*Mockcache)(nil).Set), arg0, arg1, arg2)
}

// Mockrabbitservice is a mock of rabbitservice interface.
type Mockrabbitservice struct {
	ctrl     *gomock.Controller
	recorder *MockrabbitserviceMockRecorder
	isgomock struct{}
}

// MockrabbitserviceMockRecorder is the mock recorder for Mockrabbitservice.
type MockrabbitserviceMockRecorder struct {
	mock *Mockrabbitservice
}

// NewMockrabbitservice creates a new mock instance.
func NewMockrabbitservice(ctrl *gomock.Controller) *Mockrabbitservice {
	mock := &Mockrabbitservice{ctrl: ctrl}
	mock.recorder = &MockrabbitserviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrabbitservice) EXPECT() *MockrabbitserviceMockRecorder {
	return m.recorder
}

// SendBookCreatedMessage mocks base method.
func (m *Mockrabbitservice) SendBookCreatedMessage(arg0 context.Context, arg1 *books_pb.BookModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendBookCreatedMessage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendBookCreatedMessage indicates an expected call of SendBookCreatedMessage.
func (mr *MockrabbitserviceMockRecorder) SendBookCreatedMessage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendBookCreatedMessage", reflect.TypeOf((*Mockrabbitservice)(nil).SendBookCreatedMessage), arg0, arg1)
}
