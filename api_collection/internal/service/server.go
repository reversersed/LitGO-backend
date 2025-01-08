package service

import (
	collections_pb "github.com/reversersed/LitGO-proto/gen/go/collections"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=server.go -destination=mocks/server.go

type validator interface {
	StructValidation(any) error
}
type logger interface {
	Infof(string, ...any)
	Info(...any)
	Errorf(string, ...any)
	Error(...any)
	Warnf(string, ...any)
	Warn(...any)
}
type storage interface {
}
type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type collectionServer struct {
	cache     cache
	logger    logger
	storage   storage
	validator validator
	collections_pb.UnimplementedCollectionServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator) *collectionServer {
	return &collectionServer{
		storage:   storage,
		logger:    logger,
		cache:     cache,
		validator: validator,
	}
}
func (u *collectionServer) Register(s grpc.ServiceRegistrar) {
	collections_pb.RegisterCollectionServer(s, u)
}
