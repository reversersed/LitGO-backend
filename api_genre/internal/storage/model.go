package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type Genre struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translit"`
	BookCount    int64              `bson:"bookCount"`
}
type Category struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translit"`
	Genres       []*Genre           `bson:"genres"`
}
