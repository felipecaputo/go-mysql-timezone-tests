package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type model struct {
	ID        string    `json:"id,omitempty"`
	LocalTime time.Time `json:"local_time,omitempty"`
	DBTime    time.Time `json:"db_time,omitempty"`
}

// ***************   DataBase   ***********
func migrateDB(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS times_table (
			id VARCHAR(40),
			local_time TIMESTAMP NULL, 
			db_time TIMESTAMP NULL
		) ENGINE=MyISAM  DEFAULT CHARSET=latin1
	`)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database OK")
	}
}

func saveItemToDb(db *sql.DB) (string, error) {
	id := uuid.New().String()
	date := time.Now()
	fmt.Println(date.Format(time.RFC3339))
	if _, err := db.Exec("INSERT INTO times_table (id, local_time, db_time) VALUES (?,?,NOW())", id, date); err != nil {
		return "", errors.New(fmt.Sprintf("Error: %s", err))
	}

	return id, nil
}

func getItemFromDb(db *sql.DB, id string) (*model, error) {
	row := db.QueryRow("SELECT id, local_time, db_time FROM times_table WHERE id = ?", id)

	if row == nil {
		return nil, nil
	}

	var localTime, dbTime time.Time
	if err := row.Scan(&id, &localTime, &dbTime); err != nil {
		return nil, errors.New(fmt.Sprintf("Error: %s", err))
	}

	return &model{id, localTime, dbTime}, nil
}

func returnDbItemAsJSON(db *sql.DB, id string, w http.ResponseWriter) {
	item, err := getItemFromDb(db, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, `
		Localtime: %v,
		DB Time: %v,
		Difference Local to DB: %v
	`, item.LocalTime, item.DBTime, item.LocalTime.Sub(item.DBTime))
}

// ********** Handlers *************
func handlePost(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id, err := saveItemToDb(db)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	returnDbItemAsJSON(db, id, w)
}

func handleGet(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("id")

	returnDbItemAsJSON(db, id, w)
}

// ***************   Main   ***********
func main() {
	conString := os.Getenv("CON_STRING")
	if conString == "" {
		conString = "root:password@/timesdb?parseTime=true&loc=America%2FSao_Paulo&time_zone=%27America%2FSao_Paulo%27"
	}
	db, err := sql.Open("mysql", conString)
	if err != nil {
		panic(err)
	}

	migrateDB(db)

	if err != nil {
		fmt.Printf("Error ocurred: %s", err)
	}

	http.HandleFunc("/tz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePost(db, w, r)
		} else if r.Method == http.MethodGet {
			handleGet(db, w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", http.DefaultServeMux)
	fmt.Println("Server started at 8080")
}
