package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reportIfMail/core"
	"strings"
	"text/template"

	"github.com/hpcloud/tail"
	_ "github.com/mattn/go-sqlite3"
)

type Email struct {
	ID     int
	Date   string
	Agency string
}

func main() {
	db, err := sql.Open("sqlite3", "./locale/emails.db")
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS emails (id INTEGER PRIMARY KEY, date TEXT, agency TEXT UNIQUE)")
	if err != nil {
		log.Fatalf("Could not create table: %v", err)
	}

	go func(db *sql.DB) {
		t, _ := tail.TailFile(core.ConfigEnv("pathLog"), tail.Config{Follow: true})
		var toId, from string
		for line := range t.Lines {
			segment := strings.Fields(line.Text)
			//-- Write DB
			if toId != segment[5] {
				toId, from = segment[5], segment[6]
				continue
			}
			if core.ExtractEmail(segment[6]) == core.ConfigEnv("toEmail") {
				date := segment[0] + " " + segment[1] + " " + segment[2]
				_, err = db.Exec("INSERT INTO emails (date, agency) VALUES (?, ?) ON CONFLICT(agency) DO UPDATE SET date = ?", date, core.GetAgency(core.ExtactSubDomain(core.ExtractEmail(from))), date)

				if err != nil {
					fmt.Println(err)
				}
			}
		}
		fmt.Println(toId, from)
	}(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, date, agency FROM emails")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		emails := []Email{}
		for rows.Next() {
			email := Email{}
			err := rows.Scan(&email.ID, &email.Date, &email.Agency)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			emails = append(emails, email)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("./locale/template.html"))
		err = tmpl.Execute(w, emails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server started at http://localhost:" + core.ConfigEnv("port"))
	http.ListenAndServe(":"+core.ConfigEnv("port"), nil)
}
