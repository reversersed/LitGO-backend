package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

// TOD add authors and category+genre
type Book struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translit"`
	Description  string             `bson:"description"`
	Picture      string             `bson:"picture"`
	Filepath     string             `bson:"filepath"`
}

var mocked_books []*Book = []*Book{
	{
		Name:        "Книга о книгопечатании",
		Description: "Эта книга должна была быть о чем-то хорошем, но в итоге...",
	},
}
