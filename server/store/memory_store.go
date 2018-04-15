package store

import (
	"fmt"

	"github.com/nickwu241/simply-do/server/models"
)

// MemoryStore implemnts Store in memory.
// The data will not persist after application restart.
type MemoryStore struct {
	globalID int
	items    []models.Item
}

// NewMemoryStore returns an instance of MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		globalID: 0,
		items:    []models.Item{},
	}
}

// SetUser is a noop, here to satisfy Store inteface.
func (m *MemoryStore) SetUser(uid string) error {
	return nil
}

// GetAll returns all the items.
func (m *MemoryStore) GetAll() []models.Item {
	return m.items
}

// Get returns an item by id.
func (m *MemoryStore) Get(id string) models.Item {
	for _, item := range m.items {
		if item.ID == id {
			return item
		}
	}
	return models.Item{}
}

// Create returns the created item.
func (m *MemoryStore) Create(item models.Item) models.Item {
	item.ID = m.nextID()
	m.items = append(m.items, item)
	return item
}

// Update returns the updated item if the id exists, otherwise an an empty item.
func (m *MemoryStore) Update(id string, item models.Item) models.Item {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items[i].Text = item.Text
			m.items[i].Checked = item.Checked
			return m.items[i]
		}
	}
	return models.Item{}
}

// Delete returns the list after the operation.
func (m *MemoryStore) Delete(id string) []models.Item {
	for i, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			break
		}
	}
	return m.items
}

func (m *MemoryStore) nextID() string {
	id := fmt.Sprintf("%d", m.globalID)
	m.globalID++
	return id
}
