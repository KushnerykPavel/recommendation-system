package recommender

import (
	"context"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
	"go.uber.org/zap"
)

var (
	genresLimit    = 2
	directorsLimit = 5
	moviesLimit    = 5
)

type Repo interface {
	GetMovieByID(ctx context.Context, id int) (*models.Movie, error)

	GetLastInteraction(ctx context.Context, relationName string, userID string, entityID int) (*models.Interaction, error)
	GetLastInteractions(ctx context.Context, relationName models.EntityType, userID string) ([]*models.Interaction, error)
	SetInteraction(ctx context.Context, relationName string, interaction *models.Interaction) error
	GetMovieCandidatesForEntities(ctx context.Context, relationName models.EntityType, userID string, entityIDList []int) ([]*models.Interaction, error)

	SetRecommendationsForUser(ctx context.Context, uid string, recommendations []*models.Recommendation) error
}

type Recommender struct {
	logger *zap.SugaredLogger
	repo   Repo
}

type candidate struct {
	ID   int
	Prob float64
}

func New(l *zap.SugaredLogger, r Repo) *Recommender {
	return &Recommender{
		logger: l.With("module", "recommender"),
		repo:   r,
	}
}
