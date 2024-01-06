package main

import (
	"context"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/handlers"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/queue"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/recommender"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	NatsURI      string        `envconfig:"NATS_URI"`
	RedisURI     string        `envconfig:"REDIS_URI"`
	CacheTimeout time.Duration `envconfig:"CACHE_TIMEOUT"`

	DbURL string `envconfig:"DB_URL"`

	Address string `envconfig:"ADDR"`
}

func main() {
	ctx := context.Background()
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	l, _ := config.Build()
	defer l.Sync()
	logger := l.Sugar().With("service", "api")

	var cfg Config
	err := envconfig.Process("api", &cfg)
	if err != nil {
		logger.Fatal(err.Error())
	}

	nc, err := nats.Connect(cfg.NatsURI)
	if err != nil {
		logger.Fatal(err.Error())
	}

	db, err := sqlx.Connect("postgres", cfg.DbURL)
	if err != nil {
		logger.Fatal(err.Error())
	}

	repository := repo.New(db, logger)
	rec := recommender.New(logger, repository)

	appQueue := queue.New(logger, nc, rec)
	movieHandler := handlers.NewMovieHandler(logger, repository, appQueue)
	recommenderHandler := handlers.NewRecommenderHandler(logger, repository, appQueue)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api", func(r chi.Router) {
		r.Route("/movies", func(r chi.Router) {
			r.Get("/", movieHandler.List)
			r.Get("/{movieID}", movieHandler.Item)
		})
		r.Route("/recommendations", func(r chi.Router) {
			r.Get("/movies", recommenderHandler.Movies)
		})
	})

	appQueue.EventQueueReceiver(ctx)
	appQueue.CandidatesQueueReceiver(ctx)

	logger.Infof("app running on address: %s", cfg.Address)
	http.ListenAndServe(cfg.Address, r)
}
