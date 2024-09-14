// Code generated by MockGen. DO NOT EDIT.
// Source: book_service_grpc.pb.go

// Package mock_books_pb is a generated GoMock package.
package mock_books_pb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	books_pb "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/proto/books"
	grpc "google.golang.org/grpc"
)

// MockBookClient is a mock of BookClient interface.
type MockBookClient struct {
	ctrl     *gomock.Controller
	recorder *MockBookClientMockRecorder
}

// MockBookClientMockRecorder is the mock recorder for MockBookClient.
type MockBookClientMockRecorder struct {
	mock *MockBookClient
}

// NewMockBookClient creates a new mock instance.
func NewMockBookClient(ctrl *gomock.Controller) *MockBookClient {
	mock := &MockBookClient{ctrl: ctrl}
	mock.recorder = &MockBookClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBookClient) EXPECT() *MockBookClientMockRecorder {
	return m.recorder
}

// GetBookSuggestions mocks base method.
func (m *MockBookClient) GetBookSuggestions(ctx context.Context, in *books_pb.GetSuggestionRequest, opts ...grpc.CallOption) (*books_pb.GetBooksResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBookSuggestions", varargs...)
	ret0, _ := ret[0].(*books_pb.GetBooksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookSuggestions indicates an expected call of GetBookSuggestions.
func (mr *MockBookClientMockRecorder) GetBookSuggestions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookSuggestions", reflect.TypeOf((*MockBookClient)(nil).GetBookSuggestions), varargs...)
}

// MockBookServer is a mock of BookServer interface.
type MockBookServer struct {
	ctrl     *gomock.Controller
	recorder *MockBookServerMockRecorder
}

// MockBookServerMockRecorder is the mock recorder for MockBookServer.
type MockBookServerMockRecorder struct {
	mock *MockBookServer
}

// NewMockBookServer creates a new mock instance.
func NewMockBookServer(ctrl *gomock.Controller) *MockBookServer {
	mock := &MockBookServer{ctrl: ctrl}
	mock.recorder = &MockBookServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBookServer) EXPECT() *MockBookServerMockRecorder {
	return m.recorder
}

// GetBookSuggestions mocks base method.
func (m *MockBookServer) GetBookSuggestions(arg0 context.Context, arg1 *books_pb.GetSuggestionRequest) (*books_pb.GetBooksResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookSuggestions", arg0, arg1)
	ret0, _ := ret[0].(*books_pb.GetBooksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookSuggestions indicates an expected call of GetBookSuggestions.
func (mr *MockBookServerMockRecorder) GetBookSuggestions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookSuggestions", reflect.TypeOf((*MockBookServer)(nil).GetBookSuggestions), arg0, arg1)
}

// mustEmbedUnimplementedBookServer mocks base method.
func (m *MockBookServer) mustEmbedUnimplementedBookServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedBookServer")
}

// mustEmbedUnimplementedBookServer indicates an expected call of mustEmbedUnimplementedBookServer.
func (mr *MockBookServerMockRecorder) mustEmbedUnimplementedBookServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedBookServer", reflect.TypeOf((*MockBookServer)(nil).mustEmbedUnimplementedBookServer))
}

// MockUnsafeBookServer is a mock of UnsafeBookServer interface.
type MockUnsafeBookServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeBookServerMockRecorder
}

// MockUnsafeBookServerMockRecorder is the mock recorder for MockUnsafeBookServer.
type MockUnsafeBookServerMockRecorder struct {
	mock *MockUnsafeBookServer
}

// NewMockUnsafeBookServer creates a new mock instance.
func NewMockUnsafeBookServer(ctrl *gomock.Controller) *MockUnsafeBookServer {
	mock := &MockUnsafeBookServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeBookServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeBookServer) EXPECT() *MockUnsafeBookServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedBookServer mocks base method.
func (m *MockUnsafeBookServer) mustEmbedUnimplementedBookServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedBookServer")
}

// mustEmbedUnimplementedBookServer indicates an expected call of mustEmbedUnimplementedBookServer.
func (mr *MockUnsafeBookServerMockRecorder) mustEmbedUnimplementedBookServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedBookServer", reflect.TypeOf((*MockUnsafeBookServer)(nil).mustEmbedUnimplementedBookServer))
}
