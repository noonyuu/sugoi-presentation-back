package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	client *mongo.Client
}
