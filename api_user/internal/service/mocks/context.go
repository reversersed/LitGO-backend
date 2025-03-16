package mock_service

import "google.golang.org/grpc/metadata"

type MockServerTransportStream struct{}

func (m *MockServerTransportStream) Method() string {
	return "foo"
}

func (m *MockServerTransportStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockServerTransportStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockServerTransportStream) SetTrailer(md metadata.MD) error {
	return nil
}
