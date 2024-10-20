package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id           primitive.ObjectID   `bson:"_id,omitempty"`
	Name         string               `bson:"name"`
	TranslitName string               `bson:"translit"`
	Description  string               `bson:"description"`
	Picture      string               `bson:"picture"`
	Filepath     string               `bson:"filepath"`
	Genre        primitive.ObjectID   `bson:"genre"`
	Authors      []primitive.ObjectID `bson:"authors"`
	Rating       float64              `bson:"rating"`
	Reviews      int                  `bson:"reviews"`
	Price        int                  `bson:"price"`
	Published    int64                `bson:"published"`
}
