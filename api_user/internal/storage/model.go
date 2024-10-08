package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login          string             `json:"login" bson:"login"`
	Password       []byte             `json:"-" bson:"password"`
	Roles          []string           `json:"roles" bson:"roles"`
	Email          string             `json:"email" bson:"email"`
	EmailConfirmed bool               `json:"emailconfirmed" bson:"emailconfirmed"`
}

var seedAdmin = &User{
	Id:             primitive.NewObjectID(),
	Login:          "admin",
	Roles:          []string{"user", "admin"},
	Email:          "admin@example.com",
	EmailConfirmed: true,
}
