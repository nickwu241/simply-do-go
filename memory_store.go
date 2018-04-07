package main

import "fmt"

type MemoryStore struct {
	globalID int
	items    []Item
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		globalID: 0,
		items:    []Item{},
	}
}

func (m *MemoryStore) GetAll() []Item {
	return m.items
}

func (m *MemoryStore) Get(id string) Item {
	for _, item := range m.items {
		if item.ID == id {
			return item
		}
	}
	return Item{}
}

func (m *MemoryStore) Create(item Item) Item {
	item.ID = m.nextID()
	m.items = append(m.items, item)
	return item
}

func (m *MemoryStore) Update(id string, item Item) Item {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items[i].Text = item.Text
			m.items[i].Checked = item.Checked
			return m.items[i]
		}
	}
	return Item{}
}

func (m *MemoryStore) Delete(id string) []Item {
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
