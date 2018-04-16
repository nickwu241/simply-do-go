package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

const (
	firebaseURL     = "https://simply-do.firebaseio.com"
	credentialsFile = "simply-do-firebase-adminsdk.json"
)

// AdminSnapshot takes a snapshot of the Firebase database.
func AdminSnapshot() error {
	db, err := getDB()
	if err != nil {
		return err
	}

	conf := &firebase.Config{
		StorageBucket: "simply-do.appspot.com",
	}
	opt := option.WithCredentialsFile(credentialsFile)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return err
	}
	client, err := app.Storage(ctx)
	if err != nil {
		return err
	}
	bucket, err := client.DefaultBucket()
	if err != nil {
		return err
	}
	objName := time.Now().UTC().Format("2006-01-02") + "-simply-do.json"
	w := bucket.Object(objName).NewWriter(ctx)
	defer w.Close()

	var data interface{}
	if err := db.NewRef("/").Get(context.Background(), &data); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

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
