package store

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// AdminCopyData copies items and global_id from src to dst.
func AdminCopyData(src, dst string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	srcNode := db.NewRef(src)
	dstNode := db.NewRef(dst)
	var srcData interface{}
	var srcID interface{}
	if err := srcNode.Child("items").Get(context.Background(), &srcData); err != nil {
		return err
	}
	if err := srcNode.Child("global_id").Get(context.Background(), &srcID); err != nil {
		return err
	}
	if err := dstNode.Child("items").Set(context.Background(), srcData); err != nil {
		return err
	}
	if err := dstNode.Child("global_id").Set(context.Background(), srcID); err != nil {
		return err
	}
	return nil
}

func getDB() (*db.Client, error) {
	conf := &firebase.Config{
		DatabaseURL: firebaseURL,
	}
	opt := option.WithCredentialsFile(credentialsFile)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing database client: %v", err)
	}
	return client, nil
}
