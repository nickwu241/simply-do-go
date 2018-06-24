package store

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"firebase.google.com/go/db"
	"github.com/nickwu241/simply-do/server/models"
	"github.com/pkg/errors"
)

// FirebaseStore implements Store using Firebase as persistent storage.
type FirebaseStore struct {
	globalID int
	userRoot *db.Ref
	userList *db.Ref
	db       *db.Client
}

// NewFirebaseStore returns an instance of FirebaseStore.
func NewFirebaseStore(uid string, lid string) (*FirebaseStore, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	store := &FirebaseStore{
		userRoot: nil,
		globalID: 0,
		db:       db,
	}

	if err := store.SetUserList(uid, lid); err != nil {
		return nil, err
	}
	return store, nil
}

// SetUserList initializes the Store to use the UID and LID for all subsequent operations.
// If UID is empty, "default" will be used.
// If the UID doesn't exist, it will be created.
func (f *FirebaseStore) SetUserList(uid string, lid string) error {
	if uid == "" {
		uid = "default"
	}
	uid = strings.ToLower(uid)

	f.userRoot = f.db.NewRef("/" + uid)
	f.userList = f.userRoot.Child(lid)
	var userData interface{}
	if err := f.userRoot.Get(context.Background(), &userData); err != nil {
		return errors.Wrap(err, "getting uid")
	}
	if userData == nil {
		if err := f.userRoot.Set(context.Background(), uid); err != nil {
			return errors.Wrap(err, "setting up uid")
		}
	}

	globalIDNode := f.userRoot.Child("global_id")
	var globalID int
	if err := globalIDNode.Get(context.Background(), &globalID); err != nil {
		return errors.Wrap(err, "getting global id")
	}
	if err := globalIDNode.Set(context.Background(), globalID); err != nil {
		return errors.Wrap(err, "setting up global id")
	}
	f.globalID = globalID
	return nil
}

// GetAll returns all the items or an empty slice of items if it fails.
func (f *FirebaseStore) GetAll() []models.Item {
	f.updateAndGetLastAccessTime()
	var data map[string]models.Item
	if err := f.userRoot.Child("items").Get(context.Background(), &data); err != nil {
		fmt.Printf("error fetching items: %v\n", err)
	}
	items := []models.Item{}
	for _, item := range data {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].TimeCreated.Before(items[j].TimeCreated)
	})
	return items
}

// Get returns an item by id.
func (f *FirebaseStore) Get(id string) models.Item {
	f.updateAndGetLastAccessTime()
	var item models.Item
	if err := f.userList.Child("items/"+id).Get(context.Background(), item); err != nil {
		fmt.Printf("error fetching item with id %q: %v\n", id, err)
	}
	return models.Item{}
}

// Create returns the created item or an empty item if it failed.
func (f *FirebaseStore) Create(item models.Item) models.Item {
	item.ID = f.nextID()
	item.TimeCreated = f.updateAndGetLastAccessTime()
	item.TimeUpdated = item.TimeCreated
	if err := f.userList.Child("items/"+item.ID).Set(context.Background(), item); err != nil {
		fmt.Printf("error creating item: %v\n", err)
		return models.Item{}
	}
	return item
}

// Update returns the updated item if the id exists, otherwise an an empty item.
func (f *FirebaseStore) Update(id string, item models.Item) models.Item {
	item.ID = id
	item.TimeUpdated = f.updateAndGetLastAccessTime()
	if err := f.userList.Child("items/"+id).Set(context.Background(), item); err != nil {
		fmt.Printf("error updating item: %v\n", err)
		return models.Item{}
	}
	return item
}

// Delete returns the list after the operation.
func (f *FirebaseStore) Delete(id string) []models.Item {
	f.updateAndGetLastAccessTime()
	if err := f.userList.Child("items/" + id).Delete(context.Background()); err != nil {
		fmt.Printf("error deleting item: %v\n", err)
	}
	return f.GetAll()
}

func (f *FirebaseStore) nextID() string {
	id := fmt.Sprintf("id_%d", f.globalID)
	f.globalID++
	if err := f.userList.Child("global_id").Set(context.Background(), f.globalID); err != nil {
		fmt.Printf("error setting global_id: %v\n", err)
	}
	return id
}

func (f *FirebaseStore) updateAndGetLastAccessTime() time.Time {
	now := time.Now()
	if err := f.userRoot.Child("last_accessed_time").Set(context.Background(), now); err != nil {
		fmt.Printf("error updating user last_accessed_time with %q: %v\n", now, err)
	}
	return now
}
