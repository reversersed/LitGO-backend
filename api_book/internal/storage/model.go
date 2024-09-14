package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translit"`
	Description  string             `bson:"description"`
	Picture      string             `bson:"picture"`
	FilePath     string             `bson:"filepath"`
}
