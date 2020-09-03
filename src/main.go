package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/novoselrok/rsoccerlive/src/highlights"
	"github.com/novoselrok/rsoccerlive/src/models"
	"github.com/novoselrok/rsoccerlive/src/websockethub"
	log "github.com/sirupsen/logrus"
)

type App struct {
	db         models.Datastore
	wsUpgrader websocket.Upgrader
	hub        *websockethub.Hub
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var env map[string]string
	env, err := godotenv.Read()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		env["RSOCCERLIVE_DB_HOST"],
		env["RSOCCERLIVE_DB_PORT"],
		env["RSOCCERLIVE_DB_USER"],
		env["RSOCCERLIVE_DB_PASSWORD"],
		env["RSOCCERLIVE_DB_NAME"],
	)
	db, err := models.InitDB(psqlInfo)
	if err != nil {
		log.Fatal("Failed to create the database ", err)
	}

	allowedOrigins := strings.Split(env["RSOCCERLIVE_ALLOWED_ORIGINS"], ",")
	wsUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == origin {
					return true
				}
			}
			return false
		},
	}
	websocketHub := websockethub.NewWebsocketHub()
	app := &App{db, wsUpgrader, websocketHub}
	go highlights.HighlightUpdater(app.db, app.hub, env)

	log.Info("Starting web server")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins(allowedOrigins)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "OPTIONS"})

	router := mux.NewRouter()
	router.HandleFunc("/ws", app.websocketUpgradeHandler)
	router.HandleFunc("/api/highlights/{id}", app.highlightAPIHandler).Methods("GET")
	router.HandleFunc("/api/highlights", app.dayHighlightsAPIHandler).Queries("day", "{day}").Methods("GET")
	router.HandleFunc("/api/highlightMirrors", app.highlightMirrorsAPIHandler).Queries("highlightId", "{highlightId}").Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
