package main

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

type Store interface {
	GetAll() []Item
	Get(id string) Item
	Create(item Item) Item
	Update(id string, item Item) Item
	Delete(id string) []Item
}

type FirebaseStore struct {
	globalID int
	db       *db.Client
}

func NewFirebaseStore() (*FirebaseStore, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	var globalID int
	if err := db.NewRef("/global_id").Get(context.Background(), &globalID); err != nil {
		return nil, err
	}
	return &FirebaseStore{
		globalID: globalID,
		db:       db,
	}, nil
}

func (f *FirebaseStore) GetAll() []Item {
	var data map[string]Item
	if err := f.db.NewRef("/items").Get(context.Background(), &data); err != nil {
		fmt.Printf("error fetching items: %v\n", err)
	}
	items := []Item{}
	for _, item := range data {
		items = append(items, item)
	}
	return items
}

func (f *FirebaseStore) Get(id string) Item {
	return Item{}
}

func (f *FirebaseStore) Create(item Item) Item {
	item.ID = f.nextID()
	if err := f.db.NewRef("/items/"+item.ID).Set(context.Background(), item); err != nil {
		fmt.Printf("error creating item: %v\n", err)
		return Item{}
	}
	return item
}

func (f *FirebaseStore) Update(id string, item Item) Item {
	item.ID = id
	if err := f.db.NewRef("/items/"+id).Set(context.Background(), item); err != nil {
		fmt.Printf("error updating item: %v\n", err)
		return Item{}
	}
	return item
}

func (f *FirebaseStore) Delete(id string) []Item {
	if err := f.db.NewRef("/items/" + id).Delete(context.Background()); err != nil {
		fmt.Printf("error deleting item: %v\n", err)
	}
	return f.GetAll()
}

func (f *FirebaseStore) nextID() string {
	id := fmt.Sprintf("id_%d", f.globalID)
	f.globalID++
	if err := f.db.NewRef("/global_id").Set(context.Background(), f.globalID); err != nil {
		fmt.Printf("error setting global_id: %v\n", err)
	}
	return id
}

func getDB() (*db.Client, error) {
	conf := &firebase.Config{
		DatabaseURL: "https://simply-do.firebaseio.com",
	}
	opt := option.WithCredentialsFile("simply-do-firebase-adminsdk.json")
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
