package models

type Todo struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Active bool   `json:"active"`
	Done   bool   `json:"done"`
}

type App struct {
	Name    string
	Version string
}
