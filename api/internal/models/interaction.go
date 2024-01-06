package models

import "time"

type Interaction struct {
	UserID    string    `db:"user_id"`
	EntityID  int       `db:"entity_id"`
	Alpha     float64   `db:"alpha"`
	Beta      float64   `db:"beta"`
	CreatedAt time.Time `db:"created_at"`
}
