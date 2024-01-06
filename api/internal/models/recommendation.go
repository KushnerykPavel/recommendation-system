package models

import "time"

type Recommendation struct {
	UserID     string     `json:"user_id" db:"user_id"`
	MovieID    int        `json:"movie_id" db:"movie_id"`
	EntityType EntityType `json:"entity_type" db:"entity_type"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
