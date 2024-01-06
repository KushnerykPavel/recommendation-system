package models

import "net/http"

type Relation struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (rel *Relation) Bind(r *http.Request) error {
	return nil
}

func (rel *Relation) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
