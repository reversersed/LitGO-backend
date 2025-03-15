package endpoint

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/config"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterEndpoints(ctx context.Context, cfg *config.UrlConfig) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := users_pb.RegisterUserHandlerFromEndpoint(ctx, mux, cfg.UserServiceUrl, opts)
	if err != nil {
		return nil, err
	}
	err = genres_pb.RegisterGenreHandlerFromEndpoint(ctx, mux, cfg.GenreServiceUrl, opts)
	if err != nil {
		return nil, err
	}
	err = authors_pb.RegisterAuthorHandlerFromEndpoint(ctx, mux, cfg.AuthorServiceUrl, opts)
	if err != nil {
		return nil, err
	}
	err = books_pb.RegisterBookHandlerFromEndpoint(ctx, mux, cfg.BookServiceUrl, opts)
	if err != nil {
		return nil, err
	}
	err = reviews_pb.RegisterReviewHandlerFromEndpoint(ctx, mux, cfg.ReviewServiceUrl, opts)
	if err != nil {
		return nil, err
	}

	return mux, nil
}
