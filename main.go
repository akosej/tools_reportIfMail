package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reportIfMail/core"
	"strings"
	"text/template"
	"time"

	"github.com/hpcloud/tail"
	_ "github.com/mattn/go-sqlite3"
)

// Email es una estructura que representa una entrada en la tabla de correos electrónicos.
type Email struct {
	ID     int
	Date   string
	Agency string
}

func main() {
	// Abrir la base de datos SQLite3
	db, err := sql.Open("sqlite3", "./locale/emails.db")
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer db.Close()

	// Crear la tabla de correos electrónicos si no existe
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS emails (id INTEGER PRIMARY KEY, date TEXT, agency TEXT UNIQUE)")
	if err != nil {
		log.Fatalf("Could not create table: %v", err)
	}

	// Crear un mapa de grupos de líneas por valor de segment[5]
	lineGroups := make(map[string][]string)

	// Monitorear el archivo de registro de Postfix en segundo plano
	go func() {
		t, _ := tail.TailFile(core.ConfigEnv("pathLog"), tail.Config{Follow: true})
		for line := range t.Lines {
			// Si la línea contiene "from=<" o "to=<", agréguela al grupo correspondiente
			if strings.Contains(line.Text, "from=<") || strings.Contains(line.Text, "to=<") {
				segment := strings.Fields(line.Text)
				key := segment[5]
				if _, ok := lineGroups[key]; !ok {
					// si la clave no existe, crear una nueva slice para ese grupo
					lineGroups[key] = make([]string, 0)
				}
				// agregar la línea actual a la slice correspondiente
				lineGroups[key] = append(lineGroups[key], line.Text)
			}
		}
	}()

	// Procesar grupos de líneas cada 10 segundos
	go func(db *sql.DB) {
		for {
			for key, _ := range lineGroups {
				var from string
				for _, line := range lineGroups[key] {

					segment := strings.Fields(line)
					if strings.Contains(line, "from=<") {
						from = segment[6]
					}
					// Si el correo es para toEmail, actualizar entrada en la tabla de correos electrónicos
					if core.ExtractEmail(segment[6]) == core.ConfigEnv("toEmail") {
						date := segment[0] + " " + segment[1] + " " + segment[2]
						_, err = db.Exec("INSERT INTO emails (date, agency) VALUES (?, ?) ON CONFLICT(agency) DO UPDATE SET date = ?", date, core.GetAgency(core.ExtactSubDomain(core.ExtractEmail(from))), date)

						if err != nil {
							fmt.Println(err)
						}
					}
				}
				// Eliminar el grupo de entradas procesadas
				delete(lineGroups, key)
			}
			// Esperar 10 segundos para volver a procesar
			time.Sleep(10 * time.Second)
		}
	}(db)

	// Manejar las solicitudes HTTP en la raíz del servidor
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Consultar la tabla de correos electrónicos
		rows, err := db.Query("SELECT id, date, agency FROM emails")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Construir una lista de correos electrónicos a partir de los resultados de la consulta
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

		// Renderizar la plantilla HTML con la lista de correos electrónicos
		tmpl := template.Must(template.ParseFiles("./locale/template.html"))
		err = tmpl.Execute(w, emails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Iniciar el servidor HTTP
	fmt.Println("Server started at http://localhost:" + core.ConfigEnv("port"))
	http.ListenAndServe(":"+core.ConfigEnv("port"), nil)
}
