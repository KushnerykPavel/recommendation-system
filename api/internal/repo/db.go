package repo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repo struct {
	db *sqlx.DB
	l  *zap.SugaredLogger
}

func New(d *sqlx.DB, l *zap.SugaredLogger) *Repo {
	return &Repo{
		db: d,
		l:  l,
	}
}
