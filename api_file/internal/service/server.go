package service

import (
	"github.com/reversersed/LitGO-backend/tree/main/api_file/internal/config"
	files_pb "github.com/reversersed/LitGO-proto/gen/go/files"
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
type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type rabbitservice interface {
	Close() error
}

type fileServer struct {
	cache     cache
	logger    logger
	validator validator
	rabbit    rabbitservice
	fileCfg   *config.FileConfig
	files_pb.UnimplementedFileServer
}

func NewServer(logger logger, cache cache, validator validator, rabbit rabbitservice, fileCfg *config.FileConfig) *fileServer {
	return &fileServer{
		logger:    logger,
		cache:     cache,
		validator: validator,
		rabbit:    rabbit,
		fileCfg:   fileCfg,
	}
}
func (u *fileServer) Register(s grpc.ServiceRegistrar) {
	files_pb.RegisterFileServer(s, u)
}
