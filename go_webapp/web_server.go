package main

import (
	"database/sql"
	// "errors"
	"fmt"
	"io"
	"net/http"
	"os"
	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "test"
	password = "test"
	dbname   = "postgres"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func getDB(w http.ResponseWriter, r *http.Request) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	rows, err := db.Query(`SELECT version();`)
	CheckError(err)
	defer db.Close()
	for rows.Next() {
		var db_version string
		err = rows.Scan(&db_version)
		CheckError(err)

		io.WriteString(w, db_version)
		db.Close()
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	name, err := os.Hostname()
	CheckError(err)

	// fmt.Fprintf(w, n ame)
	io.WriteString(w, name)

	// io.WriteString(w, "This is my website!\n")
}
func getImage(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := os.ReadFile("/srv/static/index.png")
	CheckError(err)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	// return
}

func getALL(w http.ResponseWriter, r *http.Request) {
	name, err := os.Hostname()
	CheckError(err)

	io.WriteString(w, name)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	rows, err := db.Query(`SELECT version();`)
	CheckError(err)
	defer db.Close()
	for rows.Next() {
		var db_version string
		err = rows.Scan(&db_version)
		CheckError(err)

		io.WriteString(w, db_version)
		db.Close()
	}

	fileBytes, err := os.ReadFile("/srv/static/index.png")
	CheckError(err)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)


}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/db", getDB)
	http.HandleFunc("/image", getImage)
	http.HandleFunc("/all", getALL)

	

	http.ListenAndServe(":8090", nil)
}