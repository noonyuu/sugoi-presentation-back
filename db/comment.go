package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func (db *Database) InsertComment(comment *Comment) error {
	fmt.Printf("insert:セッションID: %s, コメント: %s\n", comment.SessionId, comment.Comment)
	collection := db.client.Database("comments").Collection("comment")
	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}

func (db *Database) GetComments(sessionID string) ([]Comment, error) {
	collection := db.client.Database("comments").Collection("comment")

	// filter
	filter := bson.D{{Key: "session_id", Value: sessionID}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var comments []Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		panic(err)
	}
	return comments, nil
}
