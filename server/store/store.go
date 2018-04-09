package store

import "github.com/nickwu241/simply-do/server/models"

// Store provides CRUD methods to the models.
type Store interface {
	GetAll() []models.Item
	Get(id string) models.Item
	Create(item models.Item) models.Item
	Update(id string, item models.Item) models.Item
	Delete(id string) []models.Item
}
