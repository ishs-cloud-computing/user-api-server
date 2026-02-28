package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	db = initDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("DB 닫기 실패: %v", err)
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/users", createUser).Methods(http.MethodPost)
	r.HandleFunc("/users", getUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", getUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", updateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", deleteUser).Methods(http.MethodDelete)

	r.HandleFunc("/health", healthCheck).Methods(http.MethodGet)

	log.Println("서버 시작: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
