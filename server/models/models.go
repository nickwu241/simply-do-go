package models

// Item is the model for an item in a check list.
type Item struct {
	ID      string `json:"id"`
	Checked bool   `json:"checked"`
	Text    string `json:"text"`
}
