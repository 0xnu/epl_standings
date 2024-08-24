package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

type entry struct {
	uuid            string
	position        string
	team_name       string
	played          string
	won             string
	drawn           string
	lost            string
	goals_for       string
	goals_against   string
	goal_difference string
	points          string
	form            string
}

func extractForm(e *colly.HTMLElement) string {
	formData := ""
	e.ForEach("ul.ssrcss-5z9wmy-FormContainer li div[data-testid='letter-content']", func(_ int, el *colly.HTMLElement) {
		formData += el.Text
	})
	return formData
}

func collectData() ([]entry, error) {
	var entries []entry
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.IgnoreRobotsTxt(),
	)

	apiKey, exists := os.LookupEnv("SCRAPER_API_KEY")
	if !exists {
		log.Fatal("SCRAPER_API_KEY environment variable is not set")
		return nil, fmt.Errorf("SCRAPER_API_KEY environment variable is not set")
	}

	targetURL := "https://www.bbc.com/sport/football/premier-league/table"
	proxyURL := fmt.Sprintf("http://api.scraperapi.com?api_key=%s&url=%s", apiKey, url.QueryEscape(targetURL))

	c.OnHTML("table.ssrcss-14j0ip6-Table tbody tr", func(e *colly.HTMLElement) {
		var eData entry
		eData.uuid = uuid.New().String()
		eData.position = e.ChildText("td:nth-child(1) span")
		eData.team_name = strings.TrimSpace(e.ChildText("td:nth-child(2) span.ssrcss-1f39n02-VisuallyHidden"))
		eData.played = e.ChildText("td:nth-child(3)")
		eData.won = e.ChildText("td:nth-child(4)")
		eData.drawn = e.ChildText("td:nth-child(5)")
		eData.lost = e.ChildText("td:nth-child(6)")
		eData.goals_for = e.ChildText("td:nth-child(7)")
		eData.goals_against = e.ChildText("td:nth-child(8)")
		eData.goal_difference = e.ChildText("td:nth-child(9)")
		eData.points = e.ChildText("td:nth-child(10) span")
		eData.form = extractForm(e)
		entries = append(entries, eData)
	})

	err := c.Visit(proxyURL)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func writeToCSV(w io.Writer, entries []entry) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	err := writer.Write([]string{"uuid", "position", "team_name", "played", "won", "drawn", "lost", "goals_for", "goals_against", "goal_difference", "points", "form"})
	if err != nil {
		log.Fatal("Error writing CSV header:", err)
		return
	}
	for _, e := range entries {
		err := writer.Write([]string{e.uuid, e.position, e.team_name, e.played, e.won, e.drawn, e.lost, e.goals_for, e.goals_against, e.goal_difference, e.points, e.form})
		if err != nil {
			log.Fatal("Error writing CSV row:", err)
			return
		}
	}
	log.Println("CSV writing completed.")
}

func writeToDB(entries []entry) {
	env1, err1 := os.LookupEnv("FOOTBALL_DBNAME")
	log.Println(env1, err1)

	db, err := sql.Open("mysql", dsn(env1))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE football_table")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		sql := "INSERT INTO football_table (uuid, position, team_name, played, won, drawn, lost, goals_for, goals_against, goal_difference, points, form) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err := db.Exec(sql, e.uuid, e.position, e.team_name, e.played, e.won, e.drawn, e.lost, e.goals_for, e.goals_against, e.goal_difference, e.points, e.form)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func dsn(dbName string) string {
	env2, err2 := os.LookupEnv("FOOTBALL_USER")
	log.Println(env2, err2)
	env3, err3 := os.LookupEnv("FOOTBALL_PASS")
	log.Println(env3, err3)
	env4, err4 := os.LookupEnv("FOOTBALL_HOST")
	log.Println(env4, err4)
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", env2, env3, env4, dbName)
}

func main() {
	fName := "football_table.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	entries, err := collectData()
	if err != nil || len(entries) == 0 {
		log.Fatal("No data collected. Exiting.")
		return
	}

	log.Println("Writing to CSV...")
	writeToCSV(file, entries)
	log.Println("Writing to CSV done.")

	log.Println("Writing to DB...")
	writeToDB(entries)
	log.Println("Writing to DB done.")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
