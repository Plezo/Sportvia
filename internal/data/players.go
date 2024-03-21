package data

import (
	"context"
	"database/sql"
	"time"
)

type Player struct {
	ID           string    `json:"id"`
	PlayerName   string    `json:"playerName"`
	Age          int8      `json:"age"`
	Height       int8      `json:"height"`
	Team         string    `json:"team"`
	Conference   string    `json:"conference"`
	Division     string    `json:"division"`
	Position     string    `json:"position"`
	PlayerNumber int8      `json:"playerNumber"`
	PlayerImage  string    `json:"playerImage"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PlayerModel struct {
	DB *sql.DB
}

func (m PlayerModel) BulkUpsert(players *[]Player) error {
	query := `
		INSERT INTO player (playerName, age, height, team, conference, division, position, playerNumber, playerImage, updatedAt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (playerName) DO UPDATE SET
			age = EXCLUDED.age,
			height = EXCLUDED.height,
			team = EXCLUDED.team,
			conference = EXCLUDED.conference,
			division = EXCLUDED.division,
			position = EXCLUDED.position,
			playerNumber = EXCLUDED.playerNumber,
			playerImage = EXCLUDED.playerImage,
			updatedAt = EXCLUDED.updatedAt`
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	for _, player := range *players {
		args := []interface{}{
			player.PlayerName,
			player.Age,
			player.Height,
			player.Team,
			player.Conference,
			player.Division,
			player.Position,
			player.PlayerNumber,
			player.PlayerImage,
			time.Now(),
		}

		_, err := m.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// func (m PlayerModel) BulkUpsert(players *[]Player) error {
// 	query := `
// 	INSERT INTO player (playerName, age, height, team, position, playerNumber, playerImage)
// 	SELECT
// 		unnest($1),
// 		unnest($2),
// 		unnest($3),
// 		unnest($4),
// 		unnest($5),
// 		unnest($6),
// 		unnest($7)
// 		ON CONFLICT (playerName) DO UPDATE SET
// 			age = EXCLUDED.age,
// 			height = EXCLUDED.height,
// 			team = EXCLUDED.team,
// 			position = EXCLUDED.position,
// 			playerNumber = EXCLUDED.playerNumber,
// 			playerImage = EXCLUDED.playerImage,
// 			updatedAt = $8`

// 	playerNames := make([]string, len(*players))
// 	ages := make([]int8, len(*players))
// 	heights := make([]int8, len(*players))
// 	teams := make([]string, len(*players))
// 	positions := make([]string, len(*players))
// 	playerNumbers := make([]int8, len(*players))
// 	playerImages := make([]string, len(*players))

// 	for _, player := range *players {
// 		playerNames = append(playerNames, player.PlayerName)
// 		ages = append(ages, player.Age)
// 		heights = append(heights, player.Height)
// 		teams = append(teams, player.Team)
// 		positions = append(positions, player.Position)
// 		playerNumbers = append(playerNumbers, player.PlayerNumber)
// 		playerImages = append(playerImages, player.PlayerImage)
// 	}

// 	args := []interface{}{
// 		pq.Array(playerNames), 
// 		pq.Array(ages), 
// 		pq.Array(heights), 
// 		pq.Array(teams), 
// 		pq.Array(positions), 
// 		pq.Array(playerNumbers), 
// 		pq.Array(playerImages),
// 		time.Now(),
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := m.DB.ExecContext(ctx, query, args...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (m PlayerModel) Insert(player *Player) error {
	query := `
		INSERT INTO player (playerName, age, height, team, conference, division, position, playerNumber, playerImage)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	args := []interface{}{player.PlayerName, player.Age, player.Height, player.Team, player.Conference, player.Division, player.Position, player.PlayerNumber, player.PlayerImage}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&player.ID)
}

// TODO: add index in db for name
func (m PlayerModel) GetPlayerByName(name string) (*Player, error) {
	query := `
		SELECT *
		FROM player
		WHERE playerName = $1`

	args := []interface{}{name}

	var player Player

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
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

func (m PlayerModel) GetAllPlayerNames() (*[]string, error) {
	query := `
		SELECT playerName
		FROM player`

	var playerNames []string

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var playerName string

		if err := rows.Scan(&playerName); err != nil {
			return nil, err
		}

		playerNames = append(playerNames, playerName)	
	}

	return &playerNames, nil
}