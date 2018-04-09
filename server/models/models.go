package models

type Item struct {
	ID      string `json:"id"`
	Checked bool   `json:"checked"`
	Text    string `json:"text"`
}
