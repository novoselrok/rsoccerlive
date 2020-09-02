package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // db driver
)

type Datastore interface {
	GetDayHighlights(timestamp time.Time) ([]Highlight, error)
	GetHighlightMirrors(highlightID string) ([]HighlightMirror, error)
	GetHighlightsAfterTimestamp(timestamp time.Time) ([]Highlight, error)
	GetHighlight(highlightID string) (Highlight, error)
	SaveHighlights(highlights []Highlight) ([]string, error)
	SaveHighlightMirrors(highlightMirrors []HighlightMirror) ([]string, error)
}

type DB struct {
	*sqlx.DB
}

func InitDB(dataSourceName string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
