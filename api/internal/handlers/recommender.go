package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type RecommenderRepository interface {
	GetMoviesRecommendations(ctx context.Context, uid string) ([]*models.Movie, error)
}

type RecommenderQueue interface {
	PublishEvent(event *models.Event)
}

type RecommenderHandler struct {
	l     *zap.SugaredLogger
	repo  RecommenderRepository
	queue RecommenderQueue
}

func NewRecommenderHandler(
	l *zap.SugaredLogger,
	r RecommenderRepository,
	q RecommenderQueue,
) *RecommenderHandler {
	return &RecommenderHandler{
		l:     l.With("handler", "recommender"),
		repo:  r,
		queue: q,
	}
}

func (h *RecommenderHandler) Movies(w http.ResponseWriter, r *http.Request) {
	uxid := r.Header.Get(models.UserHeaderName)
	if uxid == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("uxid required for recommendations")))
		return
	}
	recommendations, err := h.repo.GetMoviesRecommendations(r.Context(), uxid)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	for _, recom := range recommendations {
		h.queue.PublishEvent(models.NewEvent(models.EventTypeSeen, models.EntityTypeMovie, recom.ID, uxid))
	}
	render.RenderList(w, r, NewMovieListResponse(recommendations))
}
