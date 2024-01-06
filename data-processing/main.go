package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Config struct {
	DbURL string `envconfig:"DB_URL"`
}

type Record struct {
	ID        string  `csv:"ID"`
	Name      string  `csv:"Movie Name"`
	Rating    float64 `csv:"Rating"`
	Runtime   string  `csv:"Runtime"`
	Genre     string  `csv:"Genre"`
	Metascore int     `csv:"Metascore"`
	Plot      string  `csv:"Plot"`
	Directors string  `csv:"Directors"`
	Stars     string  `csv:"Stars"`
	Votes     int     `csv:"Votes"`
	Gross     int     `csv:"Gross"`
	Link      string  `csv:"Link"`
}

func main() {
	var cfg Config
	err := envconfig.Process("crawler", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	db, err := sqlx.Connect("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()
	genreMap := map[string]int{}
	actorsMap := map[string]int{}
	directorsMap := map[string]int{}

	var records []Record
	if err := gocsv.UnmarshalFile(file, &records); err != nil {
		panic(err)
	}

	for _, record := range records {
		for _, genre := range strings.Split(record.Genre, ",") {
			genreMap[strings.TrimSpace(genre)]++
		}

		stars, _ := parseQuotedArray(record.Stars)
		directors, _ := parseQuotedArray(record.Directors)

		for _, star := range stars {
			actorsMap[strings.TrimSpace(star)]++
		}
		for _, director := range directors {
			directorsMap[strings.TrimSpace(director)]++
		}
	}
	fmt.Println("num genres: ", len(genreMap))
	fmt.Println()
	fmt.Println("num actors: ", len(actorsMap))
	fmt.Println()
	fmt.Println("num directors: ", len(directorsMap))

	actorID := 1
	for actor := range actorsMap {
		db.Exec("insert into actors (id, name) values ($1, $2) on conflict DO NOTHING;", actorID, actor)
		actorID++
	}

	directorID := 1
	for director := range directorsMap {
		db.Exec("insert into directors (id, name) values ($1, $2) on conflict DO NOTHING;", directorID, director)
		directorID++
	}

	genreID := 1
	for genre := range genreMap {
		db.Exec("insert into genres (id, name) values ($1, $2) on conflict DO NOTHING;", genreID, genre)
		genreID++
	}
	movieID := 1
	for _, record := range records {
		duration, _ := strconv.Atoi(strings.Split(record.Runtime, " ")[0])
		durationUnit := strings.Split(record.Runtime, " ")[1]
		db.Exec("insert into movies (id, name, description, rating, votes, duration, duration_unit, link) values ($1, $2, $3, $4, $5, $6, $7, $8) on conflict DO NOTHING;",
			movieID, record.Name, record.Plot, record.Rating, record.Votes, duration, durationUnit, record.Link)

		stars, _ := parseQuotedArray(record.Stars)
		directors, _ := parseQuotedArray(record.Directors)

		var id int
		for _, star := range stars {
			_ = db.Get(&id, "select id from actors where name = $1;", strings.TrimSpace(star))
			db.Exec("insert into movies_actors_relation (source_id, destination_id) values ($1, $2) on conflict DO NOTHING;", movieID, id)
		}
		for _, director := range directors {
			_ = db.Get(&id, "select id from directors where name = $1;", strings.TrimSpace(director))
			db.Exec("insert into movies_directors_relation (source_id, destination_id) values ($1, $2) on conflict DO NOTHING;", movieID, id)
		}
		for _, genre := range strings.Split(record.Genre, ",") {
			_ = db.Get(&id, "select id from genres where name = $1;", strings.TrimSpace(genre))
			db.Exec("insert into movies_genres_relation (source_id, destination_id) values ($1, $2) on conflict DO NOTHING;", movieID, id)
		}
		movieID++
	}
}

func parseQuotedArray(input string) ([]string, error) {

	input = strings.ReplaceAll(input, "[", "")
	input = strings.ReplaceAll(input, "]", "")

	elements := strings.Split(input, ",")

	var result []string
	for _, elem := range elements {
		elem = strings.TrimSpace(elem)
		elem = strings.TrimLeft(elem, "'")
		elem = strings.TrimRight(elem, "'")

		result = append(result, elem)
	}

	return result, nil
}
