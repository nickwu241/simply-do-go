package models

import "time"

// Item is the model for an item in a check list.
type Item struct {
	ID          string    `json:"id"`
	Checked     bool      `json:"checked"`
	Text        string    `json:"text"`
	TimeCreated time.Time `json:"time_created"`
}
