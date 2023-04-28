package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Импортируем для возможности подключения к MySQL
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

const (
	port         = ":3000"
	dbDriverName = "mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	dbx := sqlx.NewDb(db, dbDriverName)

	router := mux.NewRouter()
	router.HandleFunc("/home", index(dbx))
	router.HandleFunc("/post/{postID}", post(dbx))

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Start server " + port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	// Здесь прописываем соединение к базе данных
	return sql.Open(dbDriverName, "root:gtqju7VrdU@tcp(localhost:3306)/blog?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
