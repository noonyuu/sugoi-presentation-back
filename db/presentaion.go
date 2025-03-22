package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *Database) InsertPresentationUrl(presentation *PresentationUrl) error {
	collection := db.client.Database("users").Collection("presentation")
	_, err := collection.InsertOne(context.TODO(), presentation)
	return err
}

func (db *Database) GetAllPresentationUrls(userId string) ([]PresentationUrl, error) {
	collection := db.client.Database("users").Collection("presentation")

	// filter
	filter := bson.D{{Key: "user_id", Value: userId}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var presentationUrls []PresentationUrl
	if err = cursor.All(context.TODO(), &presentationUrls); err != nil {
		panic(err)
	}
	return presentationUrls, nil
}

func (db *Database) GetSelectedPresentationUrl(userId string, url string) (*PresentationUrl, error) {
	collection := db.client.Database("users").Collection("presentation")

	// filter
	filter := bson.D{{Key: "user_id", Value: userId}, {Key: "url", Value: url}}
	
	var presentationUrl PresentationUrl
	err := collection.FindOne(context.TODO(), filter).Decode(&presentationUrl)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &presentationUrl, nil
}
