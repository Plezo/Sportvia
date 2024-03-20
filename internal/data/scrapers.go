package data

type Scrapers struct {
	PlayerScraper interface {
		Run() ([]Player, error)
	}
}

func NewScrapers() Scrapers {
	return Scrapers{
		PlayerScraper: PlayerScraper{},
	}
}