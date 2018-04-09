package store

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/nickwu241/simply-do/server/models"
	"google.golang.org/api/option"
)

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

func (f *FirebaseStore) GetAll() []models.Item {
	var data map[string]models.Item
	if err := f.db.NewRef("/items").Get(context.Background(), &data); err != nil {
		fmt.Printf("error fetching items: %v\n", err)
	}
	items := []models.Item{}
	for _, item := range data {
		items = append(items, item)
	}
	return items
}

func (f *FirebaseStore) Get(id string) models.Item {
	return models.Item{}
}

func (f *FirebaseStore) Create(item models.Item) models.Item {
	item.ID = f.nextID()
	if err := f.db.NewRef("/items/"+item.ID).Set(context.Background(), item); err != nil {
		fmt.Printf("error creating item: %v\n", err)
		return models.Item{}
	}
	return item
}

func (f *FirebaseStore) Update(id string, item models.Item) models.Item {
	item.ID = id
	if err := f.db.NewRef("/items/"+id).Set(context.Background(), item); err != nil {
		fmt.Printf("error updating item: %v\n", err)
		return models.Item{}
	}
	return item
}

func (f *FirebaseStore) Delete(id string) []models.Item {
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