package mocks

//go:generate mockgen -source=../users/user_service_grpc.pb.go -destination=users/user_service_mock.go
//go:generate mockgen -source=../genres/genre_service_grpc.pb.go -destination=genres/genre_service_mock.go
//go:generate mockgen -source=../authors/author_service_grpc.pb.go -destination=authors/author_service_mock.go
