package data

import (
	"database/sql"
)

type Models struct {
	Game interface {
		Insert(game *Game) error
		Get(id int64) (*Game, error)
		Update(game *Game) error
		GenerateGame(userID int64) (*Game, error)
	}

	Player interface {
		BulkUpsert(players *[]Player) error
		Insert(player *Player) error
		GetPlayerByName(name string) (*Player, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Game: GameModel{DB: db},
		Player: PlayerModel{DB: db},
	}
}
