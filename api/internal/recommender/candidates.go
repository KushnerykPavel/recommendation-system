package recommender

import (
	"context"
	"sort"
	"time"

	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func (r *Recommender) CalculateCandidatesForUser(ctx context.Context, uid string) error {
	ts := time.Now()
	recommendationDeduplicateMap := make(map[int]models.EntityType)
	movieGenresCandidates, err := r.getCandidatesForRelation(ctx, models.EntityTypeGenre, uid, genresLimit, moviesLimit)
	if err != nil {
		return err
	}
	for _, movieID := range movieGenresCandidates {
		recommendationDeduplicateMap[movieID] = models.EntityTypeGenre
	}

	movieDirectorsCandidates, err := r.getCandidatesForRelation(ctx, models.EntityTypeDirector, uid, directorsLimit, moviesLimit)
	if err != nil {
		return err
	}
	for _, movieID := range movieDirectorsCandidates {
		recommendationDeduplicateMap[movieID] = models.EntityTypeDirector
	}

	recommendations := make([]*models.Recommendation, 0)
	for movieID, entityType := range recommendationDeduplicateMap {
		recommendations = append(recommendations, &models.Recommendation{
			UserID:     uid,
			MovieID:    movieID,
			EntityType: entityType,
			CreatedAt:  ts,
		})
	}
	return r.repo.SetRecommendationsForUser(ctx, uid, recommendations)
}

func (r *Recommender) getCandidatesForRelation(ctx context.Context, relation models.EntityType, uid string, entityLimit, movieLimit int) ([]int, error) {
	interactions, err := r.repo.GetLastInteractions(ctx, relation, uid)
	if err != nil {
		return nil, err
	}

	topEntities := r.getTopCandidatesList(interactions, entityLimit)
	movieCandidates, err := r.repo.GetMovieCandidatesForEntities(ctx, relation, uid, topEntities)
	if err != nil {
		return nil, err
	}

	return r.getTopCandidatesList(movieCandidates, movieLimit), nil
}

func (r *Recommender) getTopCandidatesList(interactions []*models.Interaction, limit int) []int {
	entityCandidates := make([]*candidate, len(interactions))
	for idx, itr := range interactions {
		distribution := distuv.Beta{
			Alpha: itr.Alpha,
			Beta:  itr.Beta,
			Src:   rand.New(rand.NewSource(rand.Uint64())),
		}

		entityCandidates[idx] = &candidate{
			ID:   itr.EntityID,
			Prob: distribution.Rand(),
		}
	}

	sort.Slice(entityCandidates[:], func(i, j int) bool {
		return entityCandidates[i].Prob > entityCandidates[j].Prob
	})

	if len(entityCandidates) > limit {
		entityCandidates = entityCandidates[:limit]
	}
	topGenres := make([]int, min(len(entityCandidates), limit))
	for id, cd := range entityCandidates {
		topGenres[id] = cd.ID
	}
	return topGenres
}
