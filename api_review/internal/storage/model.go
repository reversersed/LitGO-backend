package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	Oldest = "old"
	Newest = "new"
)

type ReviewReplyModel struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text      string             `json:"text" bson:"text"`
	Created   int64              `json:"created" bson:"created"`
	CreatorId primitive.ObjectID `json:"author" bson:"author"`
}
type ReviewModel struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text      string             `json:"text" bson:"text"`
	Rating    float64            `json:"rating" bson:"rating"`
	Created   int64              `json:"created" bson:"created"`
	CreatorId primitive.ObjectID `json:"author" bson:"author"`
	Replies   []ReviewReplyModel `json:"replies" bson:"replies"`
	BookId    primitive.ObjectID `json:"-" bson:"book"`
}
