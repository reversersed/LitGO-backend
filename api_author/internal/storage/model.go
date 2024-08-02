package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type Author struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	TranslitName   string             `bson:"translit"`
	About          string             `bson:"about"`
	ProfilePicture string             `bson:"profilepic"`
	Rating         float32            `bson:"rating"`
}
