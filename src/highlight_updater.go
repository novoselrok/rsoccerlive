package main

import (
	"time"

	"github.com/novoselrok/rsoccerlive/src/models"

	"github.com/novoselrok/rsoccerlive/src/redditclient"
	log "github.com/sirupsen/logrus"
)

func highlightUpdater(app *App) {
	config := redditclient.Config{
		app.env["RSOCCERLIVE_REDDIT_USERNAME"],
		app.env["RSOCCERLIVE_REDDIT_PASSWORD"],
		app.env["RSOCCERLIVE_REDDIT_CLIENT_ID"],
		app.env["RSOCCERLIVE_REDDIT_CLIENT_SECRET"],
		app.env["RSOCCERLIVE_REDDIT_USER_AGENT"],
	}
	redditClient := redditclient.NewRedditClient(config)

	for {
		time.Sleep(30 * time.Second)

		dayAgo := time.Now().Add(-24 * time.Hour)
		existingDayOldHighlights, err := app.db.GetHighlightsAfterTimestamp(dayAgo)
		if err != nil {
			log.Error("Failed to retrieve existing day old highlights ", err)
			continue
		}

		highlightMirrors := []models.HighlightMirror{}
		for _, highlight := range existingDayOldHighlights {
			mirrors, err := GetHighlightMirrors(redditClient, highlight)
			if err != nil {
				continue
			}
			highlightMirrors = append(highlightMirrors, mirrors...)
		}

		newHighlightMirrorIds, err := app.db.SaveHighlightMirrors(highlightMirrors)
		if err != nil {
			log.Error("Failed to save new highlight mirrors ", err)
		} else {
			log.Infof("Saved %d new highlight mirrors", len(newHighlightMirrorIds))
		}

		latestHighlightSubmissions := GetLatestHighlightSubmissions(redditClient)
		highlights := ConvertSubmissionsToHighlights(latestHighlightSubmissions)

		newHighlightIds, err := app.db.SaveHighlights(highlights)
		if err != nil {
			log.Error("Failed to save new highlights ", err)
		} else {
			log.Infof("Saved %d new highlights", len(newHighlightIds))
		}
	}
}
