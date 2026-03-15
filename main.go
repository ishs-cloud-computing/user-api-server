package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	db = initDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("DB 닫기 실패: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthCheck)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users", getUsers)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("PUT /users/{id}", updateUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	log.Println("서버 시작: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
