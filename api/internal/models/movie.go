package models

import (
	"net/http"
)

type Movie struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Rating       float64 `json:"rating"`
	Votes        int     `json:"votes"`
	Duration     int     `json:"duration"`
	DurationUnit string  `json:"duration_unit"`
	Link         string  `json:"link"`

	Genres    []*Relation `json:"genres,omitempty"`
	Actors    []*Relation `json:"actors,omitempty"`
	Directors []*Relation `json:"directors,omitempty"`
}

func (m *Movie) Bind(r *http.Request) error {
	return nil
}

func (m *Movie) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
