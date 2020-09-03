package highlights

import (
	"time"

	"github.com/novoselrok/rsoccerlive/src/models"
	"github.com/novoselrok/rsoccerlive/src/websockethub"

	"github.com/novoselrok/rsoccerlive/src/redditclient"
	log "github.com/sirupsen/logrus"
)

const (
	pollingDelay = 30 * time.Second

	newHighlightsEventType       = "NEW_HIGHLIGHTS"
	newHighlightMirrorsEventType = "NEW_HIGHLIGHT_MIRRORS"
)

type event struct {
	Type string   `json:"type"`
	IDs  []string `json:"ids"`
}

func HighlightUpdater(db models.Datastore, hub *websockethub.Hub, env map[string]string) {
	config := redditclient.Config{
		env["RSOCCERLIVE_REDDIT_USERNAME"],
		env["RSOCCERLIVE_REDDIT_PASSWORD"],
		env["RSOCCERLIVE_REDDIT_CLIENT_ID"],
		env["RSOCCERLIVE_REDDIT_CLIENT_SECRET"],
		env["RSOCCERLIVE_REDDIT_USER_AGENT"],
	}
	redditClient := redditclient.NewRedditClient(config)

	for {
		time.Sleep(pollingDelay)

		dayAgo := time.Now().Add(-24 * time.Hour)
		existingDayOldHighlights, err := db.GetHighlightsAfterTimestamp(dayAgo)
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

		newHighlightMirrorIds, err := db.SaveHighlightMirrors(highlightMirrors)
		if err != nil {
			log.Error("Failed to save new highlight mirrors ", err)
		} else {
			log.Infof("Saved %d new highlight mirrors", len(newHighlightMirrorIds))

			if len(newHighlightMirrorIds) > 0 {
				hub.BroadcastJSON(event{newHighlightMirrorsEventType, newHighlightMirrorIds})
			}
		}

		latestHighlightSubmissions := GetLatestHighlightSubmissions(redditClient)
		highlights := ConvertSubmissionsToHighlights(latestHighlightSubmissions)

		newHighlightIds, err := db.SaveHighlights(highlights)
		if err != nil {
			log.Error("Failed to save new highlights ", err)
		} else {
			log.Infof("Saved %d new highlights", len(newHighlightIds))

			if len(newHighlightIds) > 0 {
				hub.BroadcastJSON(event{newHighlightsEventType, newHighlightIds})
			}
		}
	}
}
