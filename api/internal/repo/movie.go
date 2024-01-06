package repo

import (
	"context"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
)

type movieDTO struct {
	ID           int     `db:"id"`
	Name         string  `db:"name"`
	Description  string  `db:"description"`
	Rating       float64 `db:"rating"`
	Votes        int     `db:"votes"`
	Duration     int     `db:"duration"`
	DurationUnit string  `db:"duration_unit"`
	Link         string  `db:"link"`
}

func (r *Repo) GetMovieByID(ctx context.Context, id int) (*models.Movie, error) {
	var movie movieDTO
	err := r.db.GetContext(ctx, &movie, `select id, name, description, rating, votes, duration, duration_unit, link from movies where id = $1;`, id)
	if err != nil {
		return nil, err
	}

	genres, err := r.GetRelation(ctx, "genres", id)
	if err != nil {
		return nil, err
	}
	actors, err := r.GetRelation(ctx, "actors", id)
	if err != nil {
		return nil, err
	}
	directors, err := r.GetRelation(ctx, "directors", id)
	if err != nil {
		return nil, err
	}

	return &models.Movie{
		ID:           movie.ID,
		Name:         movie.Name,
		Description:  movie.Description,
		Rating:       movie.Rating,
		Votes:        movie.Votes,
		Duration:     movie.Duration,
		DurationUnit: movie.DurationUnit,
		Link:         movie.Link,

		Genres:    genres,
		Actors:    actors,
		Directors: directors,
	}, nil
}

func (r *Repo) GetMovieList(ctx context.Context, limit, offset int) ([]*models.Movie, error) {
	moviesList := make([]*movieDTO, 0)
	err := r.db.SelectContext(ctx, &moviesList, `select id, name, description, rating, votes, duration, duration_unit, link from movies limit $1 offset $2;`, limit, offset)
	if err != nil {
		return nil, err
	}

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
