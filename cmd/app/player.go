package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Plezo/Sportvia/internal/data"
)

// func (app *application) getPlayerHandler(w http.ResponseWriter, r *http.Request) {

// 	players := data.Players

// 	player := data.Players[rand.Intn(len(players))]

// 	err := app.writeJSON(w, http.StatusOK, player, nil)
// 	if err != nil {
// 		app.logger.PrintError(err, nil)
// 		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
// 	}
// }

func (app *application) createPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerName   string `json:"playerName"`
		Age          int8   `json:"age"`
		Height       int8   `json:"height"`
		Team         string `json:"team"`
		Position     string `json:"position"`
		PlayerNumber int8   `json:"playerNumber"`
		PlayerImage  string `json:"playerImage"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusBadRequest)
	}
	
	player := &data.Player{
		PlayerName:   req.PlayerName,
		Age:          req.Age,
		Height:       req.Height,
		Team:         req.Team,
		Position:     req.Position,
		PlayerNumber: req.PlayerNumber,
		PlayerImage:  req.PlayerImage,
	}

	err = app.models.Player.Insert(player)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	err = app.writeJSON(w, http.StatusOK, player, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) scrapePlayersHandler(w http.ResponseWriter, r *http.Request) {
	players, err := app.scrapers.PlayerScraper.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
	
	// track time taken

	err = app.models.Player.BulkUpsert(&players)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	fmt.Printf("Upserted %d players to db\n", len(players))

	// for _, player := range players {
	// 	err = app.models.Player.Insert(&player)
	// 	if err != nil {
	// 		fmt.Printf("Error: %v\n", err)
	// 		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	// 	}
	// }

	data := map[string]string{
		"message": fmt.Sprintf("Scraped %d players successfully", len(players)),
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}