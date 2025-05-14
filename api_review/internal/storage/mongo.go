package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=mongo.go -destination=mocks/mongo.go

type logger interface {
	Infof(string, ...any)
	Info(...any)
	Warnf(string, ...any)
	Warn(...any)
	Fatalf(string, ...any)
	Fatal(...any)
}
type db struct {
	logger         logger
	bookCollection *mongodb.Collection
}

func NewStorage(storage *mongodb.Database, collection string, logger logger) *db {
	db := &db{
		bookCollection: storage.Collection(collection),
		logger:         logger,
	}
	return db
}

func (d *db) GetUserBookReview(ctx context.Context, bookId string, userId string) (*ReviewModel, error) {
	primitive_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	book_id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"$and": []bson.M{{"author": primitive_id, "book": book_id}}}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := d.bookCollection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		d.logger.Warnf("error while fetching review from db: %v", err)
		return nil, err
	}
	var u ReviewModel
	if err := result.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *db) GetBookReviews(ctx context.Context, bookId string, page int, count int, sortType string) ([]*ReviewModel, error) {
	book_id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return nil, err
	}

	options := options.Find()

	switch sortType {
	case Oldest:
		options.SetSort(bson.M{"created": 1})
	case Newest:
		options.SetSort(bson.M{"created": -1})
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("not known sort type: %v", sortType))
	}
	options.SetSkip(int64(page * count))
	options.SetLimit(int64(count))

	response, err := d.bookCollection.Find(ctx, bson.M{"book": book_id}, options)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer response.Close(ctx)

	var reviews []*ReviewModel = make([]*ReviewModel, 0)
	if err := response.All(ctx, &reviews); err != nil {
		return nil, status.Error(codes.Internal, "error decoding response: "+err.Error())
	}

	return reviews, nil
}
func (d *db) DeleteReview(ctx context.Context, bookId string, reviewId string) error {
	primitive_id, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return err
	}
	book_id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return err
	}

	filter := bson.M{"$and": []bson.M{{"_id": primitive_id, "book": book_id}}}
	result, err := d.bookCollection.DeleteOne(ctx, filter)
	if err != nil {
		return status.Error(codes.Internal, "error deleting: "+err.Error())
	}
	if result.DeletedCount == 0 {
		return status.Error(codes.NotFound, "nothing was deleted")
	}
	return nil
}
func (d *db) CreateBookReview(ctx context.Context, bookId string, text string, authorId string, rating float64) (*ReviewModel, error) {
	primitive_id, err := primitive.ObjectIDFromHex(authorId)
	if err != nil {
		return nil, err
	}
	book_id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return nil, err
	}

	review := &ReviewModel{
		Text:      text,
		Rating:    rating,
		Created:   time.Now().UTC().Unix(),
		CreatorId: primitive_id,
		Replies:   make([]ReviewReplyModel, 0),
		BookId:    book_id,
	}

	result, err := d.bookCollection.InsertOne(ctx, review)
	if err != nil {
		return nil, status.Error(codes.Internal, "error creating: "+err.Error())
	}
	review.Id, _ = result.InsertedID.(primitive.ObjectID)

	return review, nil
}
func (d *db) CreateBookReviewReply(ctx context.Context, reviewId string, text string, authorId string) (*ReviewModel, error) {
	primitive_id, err := primitive.ObjectIDFromHex(authorId)
	if err != nil {
		return nil, err
	}
	review_id, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return nil, err
	}

	review := &ReviewReplyModel{
		Id:        primitive.NewObjectID(),
		Text:      text,
		Created:   time.Now().UTC().Unix(),
		CreatorId: primitive_id,
	}
	filter := bson.M{"_id": review_id}
	update := bson.M{
		"$push": bson.M{
			"replies": review,
		},
	}

	result, err := d.bookCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, status.Error(codes.Internal, "error creating: "+err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(codes.NotFound, "no book was found")
	}
	if result.ModifiedCount == 0 {
		return nil, status.Error(codes.NotFound, "no book was modified")
	}

	reviewResult := d.bookCollection.FindOne(ctx, bson.M{"_id": review_id})
	var u ReviewModel
	if err := reviewResult.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *db) GetBookReviewsData(ctx context.Context, bookId string) (int, float64, error) {
	book_id, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return 0, 0, err
	}

	pipeline := mongodb.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "book", Value: book_id}}}},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalReviews", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "avgRating", Value: bson.D{{Key: "$avg", Value: "$rating"}}},
			}},
		},
	}

	cursor, err := d.bookCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(context.TODO())

	type ReviewStats struct {
		TotalReviews int     `bson:"totalReviews"`
		AvgRating    float64 `bson:"avgRating"`
	}

	var stats []ReviewStats
	if err = cursor.All(context.TODO(), &stats); err != nil {
		return 0, 0, err
	}

	if len(stats) == 0 {
		return 0, 0, nil
	}

	return stats[0].TotalReviews, stats[0].AvgRating, nil
}
