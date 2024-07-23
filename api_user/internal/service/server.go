package service

import (
	users_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto/users"
	"google.golang.org/grpc"
)

type userServer struct {
	users_pb.UnimplementedUserServer
}

func NewServer() *userServer {
	return &userServer{}
}
func (u *userServer) Register(s grpc.ServiceRegistrar) {
	users_pb.RegisterUserServer(s, u)
}
