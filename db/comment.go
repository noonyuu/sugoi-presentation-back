package db

import (
	"context"
	"fmt"
)

func (db *Database) InsertComment(comment *Comment) error {
	fmt.Printf("insert:セッションID: %s, コメント: %s\n", comment.SessionId, comment.Comment)
	collection := db.client.Database("comments").Collection("comment")
	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}

func (db *Database) GetComments(sessionID string) ([]Comment, error) {
	collection := db.client.Database("comment").Collection("comments")
	cursor, err := collection.Find(context.TODO(), map[string]string{"session_id": sessionID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var comments []Comment
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
