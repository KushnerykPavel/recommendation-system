package recommender

import (
	"context"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
	"time"
)

func (r *Recommender) HandleUserEvent(ctx context.Context, event *models.Event) error {
	ts := time.Now()
	var reward float64
	switch event.EntityType {
	case models.EntityTypeMovie:
		movie, _ := r.repo.GetMovieByID(ctx, event.EntityID)
		switch event.EventType {
		case models.EventTypeClicked:
			reward = 1.0
		case models.EventTypeSeen:
			reward = 0.0
		}
		r.refreshRewards(ctx, event, movie, ts, reward)
	}
	return nil
}

func (r *Recommender) refreshRewards(ctx context.Context, event *models.Event, movie *models.Movie, ts time.Time, reward float64) {
	interaction, _ := r.repo.GetLastInteraction(ctx, "movies", event.UserID, event.EntityID)
	interaction.CreatedAt = ts
	interaction.Alpha += reward
	interaction.Beta += 1.0 - reward

	if err := r.repo.SetInteraction(ctx, "movies", interaction); err != nil {
		r.logger.With("entity", "movies").Error(err)
	}

	for _, genre := range movie.Genres {
		interaction, _ = r.repo.GetLastInteraction(ctx, "genres", event.UserID, genre.ID)

		interaction.CreatedAt = ts
		interaction.Alpha += reward
		interaction.Beta += 1.0 - reward

		if err := r.repo.SetInteraction(ctx, "genres", interaction); err != nil {
			r.logger.With("entity", "genres").Error(err)
		}
	}

	for _, director := range movie.Directors {
		interaction, _ = r.repo.GetLastInteraction(ctx, "directors", event.UserID, director.ID)

		interaction.CreatedAt = ts
		interaction.Alpha += reward
		interaction.Beta += 1.0 - reward

		if err := r.repo.SetInteraction(ctx, "directors", interaction); err != nil {
			r.logger.With("entity", "directors").Error(err)
		}
	}
}
