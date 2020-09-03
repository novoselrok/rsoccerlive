package models

import (
	"fmt"
	"strings"
	"time"
)

type Highlight struct {
	ID                 string    `json:"id" db:"id"`
	URL                string    `json:"url" db:"url"`
	Title              string    `json:"title" db:"title"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	RedditSubmissionID string    `json:"redditSubmissionId" db:"reddit_submission_id"`
	RedditPermalink    string    `json:"redditPermalink" db:"reddit_permalink"`
	RedditAuthor       string    `json:"redditAuthor" db:"reddit_author"`
	RedditCreatedAt    time.Time `json:"redditCreatedAt" db:"reddit_created_at"`
	NumMirrors         int64     `json:"numMirrors" db:"num_mirrors"`
}

const (
	highlightSelectFields = "id, url, title, created_at, reddit_submission_id, reddit_permalink, reddit_author, reddit_created_at, (select COUNT(*) from highlight_mirrors where highlight_id = highlights.id) as num_mirrors"
)

func (db *DB) GetHighlight(highlightID string) (Highlight, error) {
	highlight := Highlight{}
	highlightQuery := fmt.Sprintf(`
		SELECT %s
		FROM highlights
		WHERE id = $1`, highlightSelectFields)
	err := db.Get(&highlight, highlightQuery, highlightID)
	return highlight, err
}

func (db *DB) GetDayHighlights(timestamp time.Time) ([]Highlight, error) {
	highlights := []Highlight{}
	dayHighlightsQuery := fmt.Sprintf(`
		SELECT %s
		FROM highlights
		WHERE reddit_created_at::date = $1::date`, highlightSelectFields)
	err := db.Select(&highlights, dayHighlightsQuery, timestamp)
	return highlights, err
}

func (db *DB) GetHighlightsAfterTimestamp(timestamp time.Time) ([]Highlight, error) {
	highlights := []Highlight{}
	highlightsAfterTimestampQuery := fmt.Sprintf(`
		SELECT %s
		FROM highlights
		WHERE reddit_created_at > $1`, highlightSelectFields)
	err := db.Select(&highlights, highlightsAfterTimestampQuery, timestamp)
	return highlights, err
}

func (db *DB) SaveHighlights(highlights []Highlight) ([]string, error) {
	if len(highlights) == 0 {
		return nil, nil
	}

	valuePlaceholderNumber := 1
	valuePlaceholders := []string{}
	valueArgs := []interface{}{}
	for _, highlight := range highlights {
		valuePlaceholder := fmt.Sprintf("($%d, $%d, NOW(), $%d, $%d, $%d, $%d)", valuePlaceholderNumber, valuePlaceholderNumber+1, valuePlaceholderNumber+2, valuePlaceholderNumber+3, valuePlaceholderNumber+4, valuePlaceholderNumber+5)
		valuePlaceholders = append(valuePlaceholders, valuePlaceholder)
		valueArgs = append(valueArgs, highlight.URL, highlight.Title, highlight.RedditSubmissionID, highlight.RedditPermalink, highlight.RedditAuthor, highlight.RedditCreatedAt)
		valuePlaceholderNumber = valuePlaceholderNumber + 6
	}

	insertStatement := fmt.Sprintf(
		"INSERT INTO highlights (url, title, created_at, reddit_submission_id, reddit_permalink, reddit_author, reddit_created_at) VALUES %s ON CONFLICT ON CONSTRAINT highlights_reddit_submission_id_unique DO NOTHING RETURNING id",
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
