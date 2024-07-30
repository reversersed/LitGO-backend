package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type Genre struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translitName"`
}
type Category struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translitName"`
	Genres       []*Genre           `bson:"genres"`
}
