package queue

import (
	"context"
	"encoding/json"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	eventsQueue     = "events"
	candidatesQueue = "candidates"
)

type Recommender interface {
	CalculateCandidatesForUser(ctx context.Context, uid string) error
	HandleUserEvent(ctx context.Context, event *models.Event) error
}

type Queue struct {
	logger *zap.SugaredLogger
	conn   *nats.Conn
	rec    Recommender
}

func New(l *zap.SugaredLogger, c *nats.Conn, rec Recommender) *Queue {
	return &Queue{
		logger: l.With("module", "queue"),
		conn:   c,
		rec:    rec,
	}
}

func (q *Queue) PublishEvent(event *models.Event) {
	body, err := json.Marshal(event)
	if err != nil {
		q.logger.Error(err)
		return
	}

	if err := q.conn.Publish(eventsQueue, body); err != nil {
		q.logger.Error(err)
		return
	}
}

func (q *Queue) EventQueueReceiver(ctx context.Context) {
	q.conn.Subscribe(eventsQueue, func(msg *nats.Msg) {
		var event models.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			return
		}

		if err := q.rec.HandleUserEvent(ctx, &event); err != nil {
			q.logger.Error(err)
		}

		if err := q.conn.Publish(candidatesQueue, []byte(event.UserID)); err != nil {
			q.logger.Error(err)
		}
	})
}

func (q *Queue) CandidatesQueueReceiver(ctx context.Context) {
	q.conn.Subscribe(candidatesQueue, func(msg *nats.Msg) {
		if err := q.rec.CalculateCandidatesForUser(ctx, string(msg.Data)); err != nil {
			q.logger.Error(err)
		}
	})
}
