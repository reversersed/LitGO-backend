package service

import (
	"context"
	"time"

	"github.com/reversersed/LitGO-backend-pkg/copier"
	"github.com/reversersed/LitGO-backend-pkg/middleware"
	"github.com/reversersed/LitGO-backend/tree/main/api_review/internal/rabbitmq"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *reviewServer) CreateBookReview(c context.Context, r *reviews_pb.CreateBookReviewRequest) (*reviews_pb.CreateBookReviewResponse, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := s.validator.StructValidation(r); err != nil {
		return nil, err
	}

	var credentials *shared_pb.UserCredentials
	var err error
	if credentials, err = middleware.GetCredentialsFromContext(c, s.logger); err != nil {
		return nil, err
	}

	_, err = s.storage.GetUserBookReview(c, r.GetBookId(), credentials.GetId())
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user already has active review")
	}

	review, err := s.storage.CreateBookReview(c, r.GetBookId(), r.GetText(), credentials.GetId(), float64(r.GetRating()))
	if err != nil {
		return nil, err
	}
	model := &reviews_pb.ReviewModel{}

	if err := copier.Copy(model, review, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	model.Creator = &reviews_pb.UserModel{Id: credentials.GetId(), Login: credentials.GetLogin()}

	total, rating, err := s.storage.GetBookReviewsData(c, r.GetBookId())
	if err != nil {
		return nil, err
	}
	err = s.rabbit.SendBookRatingChangedMessage(c, &rabbitmq.RatingChangeModel{BookId: r.GetBookId(), Rating: rating, TotalReviews: total})
	if err != nil {
		return nil, err
	}

	return &reviews_pb.CreateBookReviewResponse{Review: model}, nil
}
func (s *reviewServer) CreateReviewReply(c context.Context, r *reviews_pb.CreateReplyRequest) (*reviews_pb.CreateReplyResponse, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := s.validator.StructValidation(r); err != nil {
		return nil, err
	}

	var credentials *shared_pb.UserCredentials
	var err error
	if credentials, err = middleware.GetCredentialsFromContext(c, s.logger); err != nil {
		return nil, err
	}

	result, err := s.storage.CreateBookReviewReply(c, r.GetReviewId(), r.GetText(), credentials.GetId())
	if err != nil {
		return nil, err
	}

	var review = &reviews_pb.ReviewModel{}
	err = copier.Copy(review, result, copier.WithPrimitiveToStringConverter)
	if err != nil {
		return nil, err
	}

	usr, err := s.userService.GetUser(c, &users_pb.UserRequest{Id: result.CreatorId.Hex()})
	if err != nil {
		return nil, err
	}
	review.Creator = &reviews_pb.UserModel{Login: usr.GetLogin(), Id: usr.GetId()}

	for i, v := range result.Replies {
		usr, err := s.userService.GetUser(c, &users_pb.UserRequest{Id: v.CreatorId.Hex()})
		if err != nil {
			return nil, err
		}
		review.Replies[i].Creator = &reviews_pb.UserModel{Login: usr.GetLogin(), Id: usr.GetId()}
	}

	return &reviews_pb.CreateReplyResponse{Review: review}, nil
}
func (s *reviewServer) DeleteBookReview(c context.Context, r *reviews_pb.DeleteReviewRequest) (*reviews_pb.DeleteReviewResponse, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := s.validator.StructValidation(r); err != nil {
		return nil, err
	}

	if _, err := middleware.GetCredentialsFromContext(c, s.logger); err != nil {
		return nil, err
	}

	if err := s.storage.DeleteReview(c, r.GetBookId(), r.GetReviewId()); err != nil {
		return nil, err
	}

	total, rating, err := s.storage.GetBookReviewsData(c, r.GetBookId())
	if err != nil {
		return nil, err
	}
	err = s.rabbit.SendBookRatingChangedMessage(c, &rabbitmq.RatingChangeModel{BookId: r.GetBookId(), Rating: rating, TotalReviews: total})
	if err != nil {
		return nil, err
	}

	return &reviews_pb.DeleteReviewResponse{DeletedId: r.GetReviewId()}, nil
}
func (s *reviewServer) GetBookReviews(c context.Context, r *reviews_pb.GetBookReviewsRequest) (*reviews_pb.GetBookReviewsResponse, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := s.validator.StructValidation(r); err != nil {
		return nil, err
	}

	response, err := s.storage.GetBookReviews(c, r.GetId(), int(r.GetPage()), int(r.GetPageSize()), r.GetSort())
	if err != nil {
		return nil, err
	}

	data := make([]*reviews_pb.ReviewModel, len(response))
	for i, v := range response {
		var review = &reviews_pb.ReviewModel{}
		err := copier.Copy(review, v, copier.WithPrimitiveToStringConverter)
		if err != nil {
			return nil, err
		}

		usr, err := s.userService.GetUser(c, &users_pb.UserRequest{Id: v.CreatorId.Hex()})
		if err != nil {
			return nil, err
		}
		review.Creator = &reviews_pb.UserModel{Login: usr.GetLogin(), Id: usr.GetId()}

		for i, v := range v.Replies {
			usr, err := s.userService.GetUser(c, &users_pb.UserRequest{Id: v.CreatorId.Hex()})
			if err != nil {
				return nil, err
			}
			review.Replies[i].Creator = &reviews_pb.UserModel{Login: usr.GetLogin(), Id: usr.GetId()}
		}
		data[i] = review
	}
	return &reviews_pb.GetBookReviewsResponse{Reviews: data}, nil
}
func (s *reviewServer) GetCurrentUserBookReview(c context.Context, r *reviews_pb.GetCurrentUserBookReviewRequest) (*reviews_pb.GetCurrentUserReviewResponse, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := s.validator.StructValidation(r); err != nil {
		return nil, err
	}

	var credentials *shared_pb.UserCredentials
	var err error
	if credentials, err = middleware.GetCredentialsFromContext(c, s.logger); err != nil {
		return nil, err
	}

	review, err := s.storage.GetUserBookReview(c, r.GetId(), credentials.GetId())
	if err != nil {
		return nil, err
	}
	var model = &reviews_pb.ReviewModel{}

	if err := copier.Copy(model, review, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	model.Creator = &reviews_pb.UserModel{Id: credentials.GetId(), Login: credentials.GetLogin()}

	for i, v := range review.Replies {
		usr, err := s.userService.GetUser(c, &users_pb.UserRequest{Id: v.CreatorId.Hex()})
		if err != nil {
			return nil, err
		}
		model.Replies[i].Creator = &reviews_pb.UserModel{Login: usr.GetLogin(), Id: usr.GetId()}
	}

	return &reviews_pb.GetCurrentUserReviewResponse{Review: model}, nil
}
