package service

import (
	"context"

	"github.com/reversersed/LitGO-backend/tree/main/api_review/internal/rabbitmq"
	model "github.com/reversersed/LitGO-backend/tree/main/api_review/internal/storage"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
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
	GetUserBookReview(ctx context.Context, bookId string, userId string) (*model.ReviewModel, error)
	GetBookReviews(ctx context.Context, bookId string, page int, count int, sortType string) ([]*model.ReviewModel, error)
	DeleteReview(ctx context.Context, bookId string, reviewId string) error
	CreateBookReview(ctx context.Context, bookId string, text string, authorId string, rating float64) (*model.ReviewModel, error)
	CreateBookReviewReply(ctx context.Context, reviewId string, text string, authorId string) (*model.ReviewModel, error)
	GetBookReviewsData(ctx context.Context, bookId string) (int, float64, error)
}

type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type rabbitService interface {
	SendBookRatingChangedMessage(ctx context.Context, model *rabbitmq.RatingChangeModel) error
}

type reviewServer struct {
	cache       cache
	logger      logger
	storage     storage
	validator   validator
	userService users_pb.UserClient
	rabbit      rabbitService
	reviews_pb.UnimplementedReviewServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator, userService users_pb.UserClient, rabbit rabbitService) *reviewServer {
	return &reviewServer{
		storage:     storage,
		logger:      logger,
		cache:       cache,
		validator:   validator,
		userService: userService,
		rabbit:      rabbit,
	}
}
func (u *reviewServer) Register(s grpc.ServiceRegistrar) {
	reviews_pb.RegisterReviewServer(s, u)
}
