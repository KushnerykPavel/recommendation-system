package models

var UserHeaderName = "uxid"

type EventType string

var (
	EventTypeSeen    EventType = "seen"
	EventTypeClicked EventType = "clicked"
)

type EntityType string

var (
	EntityTypeMovie    EntityType = "movies"
	EntityTypeGenre    EntityType = "genres"
	EntityTypeActor    EntityType = "actors"
	EntityTypeDirector EntityType = "directors"
)

type Event struct {
	EventType  EventType  `json:"event_type"`
	EntityType EntityType `json:"entity_type"`
	EntityID   int        `json:"entity_id"`
	UserID     string     `json:"user_id"`
}

func NewEvent(event EventType, entity EntityType, eid int, uid string) *Event {
	return &Event{
		EventType:  event,
		EntityType: entity,
		UserID:     uid,
		EntityID:   eid,
	}
}
