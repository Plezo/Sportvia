package data

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// var teams []string = []string{
// 	"ATL",
// 	"BOS",
// 	"BRK",
// 	"CHO",
// 	"CHI",
// 	"CLE",
// 	"DAL",
// 	"DEN",
// 	"DET",
// 	"GSW",
// 	"HOU",
// 	"IND",
// 	"LAC",
// 	"LAL",
// 	"MEM",
// 	"MIA",
// 	"MIL",
// 	"MIN",
// 	"NOP",
// 	"NYK",
// 	"OKC",
// 	"ORL",
// 	"PHI",
// 	"PHO",
// 	"POR",
// 	"SAC",
// 	"SAS",
// 	"TOR",
// 	"UTA",
// 	"WAS",
// }

type Team struct {
	Team        string `json:"team"`
	Conference  string `json:"conference"`
	Division    string `json:"division"`
}

var NW = []string{"DEN", "MIN", "OKC", "POR", "UTA"}
var PAC = []string{"GSW", "LAC", "LAL", "PHO", "SAC"}
var SW = []string{"DAL", "HOU", "MEM", "NOP", "SAS"}
var ATL = []string{"BOS", "BRK", "NYK", "PHI", "TOR"}
var CEN = []string{"CHI", "CLE", "DET", "IND", "MIL"}
var SE = []string{"ATL", "CHO", "MIA", "ORL", "WAS"}

type PlayerScraper struct{}

func (s PlayerScraper) Run() ([]Player, error) {
	// c := colly.NewCollector(
	// 	colly.AllowURLRevisit(),
	// 	colly.MaxDepth(2),
		// colly.Async(true),
	// )

	// c.Limit(&colly.LimitRule{
	// 	DomainGlob: "*",
	// 	Parallelism: 2,
	// })

	// extensions.RandomUserAgent(c)

	var players []Player

	teams := append(NW, PAC...)
	teams = append(teams, SW...)
	teams = append(teams, ATL...)
	teams = append(teams, CEN...)
	teams = append(teams, SE...)

	for _, team := range teams {
		scrapeTeam(team, &players)
		time.Sleep(2 * time.Second)
	}

	if err := savePlayersCSV(players); err != nil {
		return nil, err
	}

	return players, nil
}

func scrapeTeam(team string, players *[]Player) {
	base_site := "https://www.basketball-reference.com"

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(2),
	)

	extensions.RandomUserAgent(c)

	conf, div := getConfDiv(team)

	c.OnHTML("table#roster tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			player := Player{Team: team, Conference: conf, Division: div}

			player.PlayerName = formatName(e.ChildText("td:nth-of-type(1)"))
			player.Age = birthdateToAge(e.ChildText("td:nth-of-type(5)"))
			player.Height = heightToInt(e.ChildText("td:nth-of-type(3)"))
			player.Position = formatPosition(e.ChildText("td:nth-of-type(2)"))

			num, _ := strconv.ParseInt(e.ChildText("th"), 10, 8)
			player.PlayerNumber = int8(num)
			player.PlayerImage = ""

			player_url := e.ChildAttr("td:nth-of-type(1) a", "href")

			time.Sleep(2 * time.Second)

			fmt.Printf("Player saved: %v\n", player.PlayerName)
			*players = append(*players, player)
			e.Request.Visit(fmt.Sprintf("%s%s", base_site, player_url))
		})
	})

	c.OnHTML("div#info.players img", func(e *colly.HTMLElement) {
		(*players)[len(*players)-1].PlayerImage = e.Attr("src")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(fmt.Sprintf("%s/teams/%s/2024.html", base_site, team))

	fmt.Printf("Team saved: %v\n", team)
}

func savePlayersCSV(players []Player) error {
	fName := "../../internal/data/csv/players.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return err
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, player := range players {
		writer.Write([]string{
			player.PlayerName,
			strconv.Itoa(int(player.Age)),
			strconv.Itoa(int(player.Height)),
			player.Team,
			player.Position,
			strconv.Itoa(int(player.PlayerNumber)),
			player.PlayerImage,
		})
	}

	fmt.Printf("Scraped %d players saved to %s\n", len(players), fName)

	return nil
}

func formatName(name string) string {
	name = strings.ReplaceAll(name, "\u00a0", " ")

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, name)
	if e != nil {
		panic(e)
	}

	parts := strings.Split(output, "(")
	output = strings.TrimSpace(parts[0])

	return output
}

func formatPosition(pos string) string {
	if strings.Contains(pos, "G") {
		return "G"
	} else if strings.Contains(pos, "F") {
		return "F"
	} else {
		return "C"
	}
}

func heightToInt(height string) int8 {
	parts := strings.Split(height, "-")
	feet, _ := strconv.Atoi(parts[0])
	inches, _ := strconv.Atoi(parts[1])

	return int8(feet*12 + inches)
}

func birthdateToAge(birth string) int8 {
	year, _ := strconv.Atoi(birth[len(birth)-4:])
	today := time.Now()
	age := today.Year() - year

	return int8(age)
}

func getConfDiv(team string) (string, string) {
	if contains(NW, team) {
		return "WEST", "NW"
	} else if contains(PAC, team) {
		return "WEST", "PAC"
	} else if contains(SW, team) {
		return "WEST", "SW"
	} else if contains(ATL, team) {
		return "EAST", "ATL"
	} else if contains(CEN, team) {
		return "EAST", "CEN"
	} else {
		return "EAST", "SE"
	}
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}