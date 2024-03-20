package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// for handling static files
	fileServer := http.FileServer(http.Dir("ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/games", app.gameView)
	// router.HandlerFunc(http.MethodGet, "/players", app.getPlayerHandler)
	router.HandlerFunc(http.MethodPost, "/players", app.createPlayerHandler)
	router.HandlerFunc(http.MethodPost, "/games", app.createGameHandler)
	router.HandlerFunc(http.MethodPatch, "/games", app.patchGameHandler)

	router.HandlerFunc(http.MethodPost, "/scrape-players", app.scrapePlayersHandler)

	return router
}
