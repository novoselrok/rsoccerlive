package models

import (
	"fmt"
	"strings"
	"time"
)

type HighlightMirror struct {
	ID              string    `json:"id" db:"id"`
	HighlightID     string    `json:"highlightId" db:"highlight_id"`
	URL             string    `json:"url" db:"url"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	RedditPermalink string    `json:"redditPermalink" db:"reddit_permalink"`
	RedditAuthor    string    `json:"redditAuthor" db:"reddit_author"`
	RedditCreatedAt time.Time `json:"redditCreatedAt" db:"reddit_created_at"`
}

func (db *DB) GetHighlightMirrors(highlightID string) ([]HighlightMirror, error) {
	highlightMirrors := []HighlightMirror{}
	dayHighlightsQuery := `
		SELECT id, highlight_id, url, created_at, reddit_permalink, reddit_author, reddit_created_at
		FROM highlight_mirrors
		WHERE highlight_id = $1`
	err := db.Select(&highlightMirrors, dayHighlightsQuery, highlightID)
	return highlightMirrors, err
}

func (db *DB) SaveHighlightMirrors(highlightMirrors []HighlightMirror) ([]string, error) {
	if len(highlightMirrors) == 0 {
		return nil, nil
	}

	valuePlaceholderNumber := 1
	valuePlaceholders := []string{}
	valueArgs := []interface{}{}
	for _, highlightMirror := range highlightMirrors {
		valuePlaceholder := fmt.Sprintf("($%d, $%d, NOW(), $%d, $%d, $%d)", valuePlaceholderNumber, valuePlaceholderNumber+1, valuePlaceholderNumber+2, valuePlaceholderNumber+3, valuePlaceholderNumber+4)
		valuePlaceholders = append(valuePlaceholders, valuePlaceholder)
		valueArgs = append(valueArgs, highlightMirror.HighlightID, highlightMirror.URL, highlightMirror.RedditPermalink, highlightMirror.RedditAuthor, highlightMirror.RedditCreatedAt)
		valuePlaceholderNumber = valuePlaceholderNumber + 5
	}

	insertStatement := fmt.Sprintf(
		"INSERT INTO highlight_mirrors (highlight_id, url, created_at, reddit_permalink, reddit_author, reddit_created_at) VALUES %s ON CONFLICT ON CONSTRAINT highlight_mirrors_highlight_id_url_unique DO NOTHING RETURNING id",
		strings.Join(valuePlaceholders, ","),
	)
	rows, err := db.Queryx(insertStatement, valueArgs...)
	if err != nil {
		return nil, err
	}

	insertedIds := []string{}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		insertedIds = append(insertedIds, id)
	}
	return insertedIds, nil
}
