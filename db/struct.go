package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	client *mongo.Client
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserId   string             `bson:"user_id"`
	Name     string             `bson:"username"`
	Password string             `bson:"password"`
	Email    string             `bson:"email"`
}

type JsonRequest struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Comment   string             `bson:"comment"`
	SessionId string             `bson:"session_id"`
}

type PresentationUrl struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserId string             `bson:"user_id"`
	Url    string             `bson:"url"`
}
