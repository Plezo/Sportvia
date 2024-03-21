package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Plezo/Sportvia/internal/utils"
)

/*
This game will work like wordle but for guessing
NBA players instead

One difficult feature to implement
will be to show a list of players as you're typing in the input field
*/

// Look into making a Height type that accepts an int and return feet and inches string
type Game struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userID"`
	PlayerName   string    `json:"playerName"`
	Age          int8      `json:"age"`
	Height       int8      `json:"height"`
	Team         string    `json:"team"`
	Conference   string    `json:"conference"`
	Division     string    `json:"division"`
	Position     string    `json:"position"`
	PlayerNumber int8      `json:"playerNumber"`
	PlayerImage  string    `json:"playerImage"`
	Attempt      int8      `json:"attempt"`
	MaxAttempts  int8      `json:"maxAttempts"`
	Win		  	 bool      `json:"win"`
	CreatedAt    time.Time `json:"createdAt"`
}

type GameModel struct {
	DB *sql.DB
}

func (m GameModel) Insert(game *Game) error {
	query := `
		INSERT INTO game (userID, playerName, age, height, team, conference, division, position, playerNumber, playerImage, attempt, maxAttempts)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, createdAt`

	args := []interface{}{game.UserID, game.PlayerName, game.Age, game.Height, game.Team, game.Conference, game.Division, game.Position, game.PlayerNumber, game.PlayerImage, game.Attempt, game.MaxAttempts}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&game.ID, &game.CreatedAt)
}

func (m GameModel) Get(id string) (*Game, error) {
	if !utils.IsValidUUID(id) {
		return nil, errors.New("record not found")
	}

	query := `
		SELECT *
		FROM game
		WHERE id = $1`

	var game Game

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&game.ID,
		&game.UserID,
		&game.PlayerName,
		&game.Age,
		&game.Height,
		&game.Team,
		&game.Conference,
		&game.Division,
		&game.Position,
		&game.PlayerNumber,
		&game.PlayerImage,
		&game.Attempt,
		&game.MaxAttempts,
		&game.Win,
		&game.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (m GameModel) Update(game *Game) error {
	query := `
		UPDATE game
		SET attempt = $1, win = $2
		WHERE id = $3`

	args := []interface{}{game.Attempt, game.Win, game.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

/*

Start a game (Insert)

Check input (Get)
- Return info about the player that was input (Mention win/loss too, so it wont be calculated by the client)

*/

// look into updating pointer to game instead of returning a game pointer
func (m GameModel) GenerateGame(userID string) (*Game, error) {

	// check if user exists

	player, err := m.GetRandomPlayer()

	if err != nil {
		return nil, err
	}

	// consider replacing player properties with a struct
	return &Game{
		UserID:       userID,
		PlayerName:   player.PlayerName,
		Age:          player.Age,
		Height:       player.Height,
		Team:         player.Team,
		Conference:   player.Conference,
		Division:     player.Division,
		Position:     player.Position,
		PlayerNumber: player.PlayerNumber,
		PlayerImage:  player.PlayerImage,
		Attempt:      0,
		MaxAttempts:  10,
		Win:          false,
		CreatedAt:    time.Now().UTC(),
	}, nil
}

func (m GameModel) GetRandomPlayer() (*Player, error) {

	query := `
		SELECT * FROM player
		ORDER BY RANDOM()
		LIMIT 1`

	var player Player

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query).Scan(
		&player.ID,
		&player.PlayerName,
		&player.Age,
		&player.Height,
		&player.Team,
		&player.Conference,
		&player.Division,
		&player.Position,
		&player.PlayerNumber,
		&player.PlayerImage,
		&player.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &player, nil
}