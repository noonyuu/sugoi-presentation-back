package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	client *mongo.Client
}

type Comment struct {
	ID        string `bson:"_id"`
	Comment   string `bson:"comment"`
	SessionId string `bson:"session_id"`
}
