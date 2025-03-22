package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *Database) CreateUser(user *User) error {
	collection := db.client.Database("users").Collection("user")
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

func (db *Database) GetUser(userId string) (*User, error) {
	collection := db.client.Database("users").Collection("user")

	// filter
	filter := bson.D{{Key: "user_id", Value: userId}}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (db *Database) GetUserInfo(userId string, password string) (*User, error) {
	collection := db.client.Database("users").Collection("user")

	// filter
	filter := bson.D{{Key: "user_id", Value: userId}}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
