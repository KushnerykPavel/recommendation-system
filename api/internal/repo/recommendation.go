package repo

import (
	"context"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
)

func (r *Repo) SetRecommendationsForUser(ctx context.Context, uid string, recommendations []*models.Recommendation) error {
	tx := r.db.MustBegin()
	tx.MustExec("delete from recommendations where user_id = $1", uid)
	tx.NamedExecContext(ctx, "insert into recommendations (user_id, movie_id, entity_type, created_at) VALUES (:user_id, :movie_id, :entity_type, :created_at)", recommendations)
	return tx.Commit()
}

func (r *Repo) GetMoviesRecommendations(ctx context.Context, uid string) ([]*models.Movie, error) {
	moviesList := make([]*movieDTO, 0)
	err := r.db.SelectContext(ctx, &moviesList, `select id, name, description, rating, votes, duration, duration_unit, link from movies
where id in (select movie_id from recommendations where user_id = $1);`, uid)

	result := make([]*models.Movie, len(moviesList))
	for idx, dto := range moviesList {
		result[idx] = &models.Movie{
			ID:           dto.ID,
			Name:         dto.Name,
			Description:  dto.Description,
			Rating:       dto.Rating,
			Votes:        dto.Votes,
			Duration:     dto.Duration,
			DurationUnit: dto.DurationUnit,
			Link:         dto.Link,
		}
	}

	return result, err
}
