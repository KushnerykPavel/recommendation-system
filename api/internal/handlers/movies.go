package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
)

type MovieRepo interface {
	GetMovieByID(ctx context.Context, id int) (*models.Movie, error)
	GetMovieList(ctx context.Context, limit, offset int) ([]*models.Movie, error)
}

type Queue interface {
	PublishEvent(event *models.Event)
}

type MovieHandler struct {
	l     *zap.SugaredLogger
	repo  MovieRepo
	queue Queue
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func NewMovieHandler(l *zap.SugaredLogger, r MovieRepo, q Queue) *MovieHandler {
	return &MovieHandler{
		l:     l.With("handler", "movie"),
		repo:  r,
		queue: q,
	}
}

func NewMovieListResponse(movies []*models.Movie) []render.Renderer {
	list := []render.Renderer{}
	for _, movie := range movies {
		list = append(list, movie)
	}
	return list
}

func (h *MovieHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPaginationParams(r.URL.Query())

	movies, err := h.repo.GetMovieList(r.Context(), limit, offset)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.RenderList(w, r, NewMovieListResponse(movies))
}

func (h *MovieHandler) Item(w http.ResponseWriter, r *http.Request) {
	rawMovieID := chi.URLParam(r, "movieID")
	if rawMovieID == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("movie-id must not be empty")))
		return
	}
	movieID, err := strconv.Atoi(rawMovieID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("movie-id must be integer %s given", rawMovieID)))
		return
	}

	movie, err := h.repo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		h.l.Info(err)
		render.Render(w, r, ErrNotFound)
		return
	}

	if uxid := r.Header.Get(models.UserHeaderName); uxid != "" {
		h.queue.PublishEvent(models.NewEvent(models.EventTypeClicked, models.EntityTypeMovie, movieID, uxid))
	}

	render.Render(w, r, movie)
}

func getPaginationParams(urlQuery url.Values) (int, int) {
	page := 1
	numPerPage := 10

	rawPage := urlQuery.Get("page")
	if rawPage != "" {
		l, err := strconv.Atoi(rawPage)
		if err == nil && l > 0 {
			page = l
		}
	}

	rawPerPage := urlQuery.Get("per_page")
	if rawPerPage != "" {
		l, err := strconv.Atoi(rawPerPage)
		if err == nil && l > 0 {
			numPerPage = l
		}
	}

	return numPerPage, (page - 1) * numPerPage
}
