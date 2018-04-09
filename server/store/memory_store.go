package store

import (
	"fmt"

	"github.com/nickwu241/simply-do/server/models"
)

type MemoryStore struct {
	globalID int
	items    []models.Item
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		globalID: 0,
		items:    []models.Item{},
	}
}

func (m *MemoryStore) GetAll() []models.Item {
	return m.items
}

func (m *MemoryStore) Get(id string) models.Item {
	for _, item := range m.items {
		if item.ID == id {
			return item
		}
	}
	return models.Item{}
}

func (m *MemoryStore) Create(item models.Item) models.Item {
	item.ID = m.nextID()
	m.items = append(m.items, item)
	return item
}

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
