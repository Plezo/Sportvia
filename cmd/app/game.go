package main

import (
	"net/http"
	"strconv"

	"github.com/Plezo/Sportvia/internal/utils"
)

type Attempt struct {
	PlayerName string `json:"playerName"`
	Age bool `json:"age"`
	Height bool `json:"height"`
	Team bool `json:"team"`
	Conference bool `json:"conference"`
	Division bool `json:"division"`
	Position bool `json:"position"`
	PlayerNumber bool `json:"playerNumber"`
	Attempt int `json:"attempt"`
	MaxAttempts int `json:"maxAttempts"`
	PlayerImage string `json:"playerImage"`
}

type PlayerFormatted struct {
	PlayerName string `json:"playerName"`
	Age int `json:"age"`
	Height int `json:"height"`
	HeightFormatted string `json:"heightFormatted"`
	Team string `json:"team"`
	Conference string `json:"conference"`
	Division string `json:"division"`
	Position string `json:"position"`
	PlayerNumber int `json:"playerNumber"`
	PlayerImage string `json:"playerImage"`
}

func (app *application) gameView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if !utils.IsValidUUID(id) {
		app.logger.PrintError(nil, map[string]string{
			"error": "Error getting game id",
			"function": "gameView",
		})

		return
	}

	game, err := app.models.Game.Get(id)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error getting game",
			"function": "gameView",
		})

		return
	}

	playerNames, err := app.models.Player.GetAllPlayerNames()
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error getting all players",
			"function": "gameView",
		})

		return
	}

	data := map[string]interface{} {
		"Game": game,
		"PlayerNames": playerNames,
	}

	err = app.templates.ExecuteTemplate(w, "game.html", data)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error executing template",
			"function": "gameView",
		})
		
		return
	}

	// tmp1 := template.Must(template.ParseFiles("../../ui/html/partials/view.html"))
	// tmp1.Execute(w, attempts)
}

func (app *application) createGameHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"userID"`
	}

	req.UserID = "ad24cf83-38e3-408b-a153-3a778b021db4"

	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	app.logger.PrintError(err, map[string]string{
	// 		"error": "Error decoding request",
	// 		"function": "createGameHandler",
	// 	})

	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// }

	// look into adding generateGame into helpers
	game, err := app.models.Game.GenerateGame(req.UserID)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error generating game",
			"function": "createGameHandler",
		})

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = app.models.Game.Insert(game); err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error inserting game",
			"function": "createGameHandler",
		})

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("HX-Location", "/games?id=" + game.ID)

	// r.Header.Set("HX-Redirect", "http://localhost:4000/games?id=" + strconv.FormatInt(game.ID, 10))

	// http.Redirect(w, r, "/games?id=" + strconv.FormatInt(game.ID, 10), http.StatusSeeOther)

	// if err = app.writeJSON(w, http.StatusOK, game, nil); err != nil {
	// 	app.logger.PrintError(err, map[string]string{
	// 		"error": "Error writing JSON",
	// 		"function": "createGameHandler",
	// 	})

	// 	http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	// }
}

// checks if input is correct for specified game id
// TODO: add 1 to attempt for specified game each time this is called
func (app *application) patchGameHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
		PlayerNameGuess string `json:"playerNameGuess"`
	}

	req.ID = r.URL.Query().Get("id")
	req.PlayerNameGuess = r.PostFormValue("playerNameGuess")

	/*

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error decoding request",
			"function": "patchGameHandler",
		})

		http.Error(w, "The server encountered a problem and could not process your request",  http.StatusBadRequest)
	}

	*/

	game, err := app.models.Game.Get(req.ID)

	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error getting game",
			"function": "patchGameHandler",
		})

		http.Error(w, "The server encountered a problem and could not process your request",  http.StatusInternalServerError)
	}

	// TODO: not sure what to do here, feel theres a better way to implement this
	if game.Attempt >= game.MaxAttempts || game.Win {
		http.Error(w, "The server encountered a problem and could not process your request",  http.StatusForbidden)
		app.templates.ExecuteTemplate(w, "gameStatus.html", game)
		return
	}

	// use lev distance to compare player name
	playerGuess, err := app.models.Player.GetPlayerByName(req.PlayerNameGuess)

	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error getting player",
			"function": "patchGameHandler",
		})

		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	// instead of lev dist, load up all player names with ssr
	// and have a autofill name feature in input field on
	// client side
	guessCorrect := lev(game.PlayerName, playerGuess.PlayerName)
	game.Win = guessCorrect

	// this is current attempt
	game.Attempt += 1

	// bools indicate which field is correct
	guessData := map[string]interface{}{
		"PlayerName": guessCorrect,
		"Age": game.Age == playerGuess.Age,
		"Height": game.Height == playerGuess.Height,
		"Team": game.Team == playerGuess.Team,
		"Conference": game.Conference == playerGuess.Conference,
		"Division": game.Division == playerGuess.Division,
		"Position": game.Position == playerGuess.Position,
		"PlayerNumber": game.PlayerNumber == playerGuess.PlayerNumber,
		"Attempt": game.Attempt,
		"MaxAttempts": game.MaxAttempts,
		"PlayerImage": playerGuess.PlayerImage,
	}

	data := map[string]interface{} {
		"PlayerGuess": playerGuess,
		"GuessData": guessData,
		"Game": game,
		"InToHeight": func(inches int8) string {
			return strconv.Itoa(int(inches / 12)) + "'" + strconv.Itoa(int(inches % 12)) + "\""
		},
	}

	err = app.models.Game.Update(game)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error updating game",
			"function": "patchGameHandler",
		})
		
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	app.templates.ExecuteTemplate(w, "attemptRow.html", data)

	if (game.Attempt >= game.MaxAttempts) || game.Win {
		app.templates.ExecuteTemplate(w, "gameStatus.html", data)
	}


	/*
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"error": "Error writing JSON",
			"function": "patchGameHandler",
		})

		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	*/
}