package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

var db *sql.DB

func main() {
	var err error

	// 환경변수에서 DB 정보 가져오기
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("DB 환경변수 미설정")
	}

	// DSN 구성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB 연결 실패: ", err)
	}

	// DB 연결 테스트
	if err = db.Ping(); err != nil {
		log.Fatal("DB 핑 실패: ", err)
	}

	// 라우터 설정
	r := mux.NewRouter()
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/healthy", healthCheck).Methods("GET")

	log.Println("서버 시작: http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// --- CRUD 핸들러 그대로 ---
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Failed to decode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Age,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := result.LastInsertId()
	user.ID = uint64(id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Failed to encode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
		if err != nil {
			log.Printf("Failed to scan users row: %v", err)
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Printf("Failed to encode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var user User
	err := db.QueryRow(
		"SELECT id, name, email, age FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age)

	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Failed to encode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Failed to decode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		"UPDATE users SET name=?, email=?, age=? WHERE id=?",
		user.Name, user.Email, user.Age, id,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user.ID = uint64(id)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Failed to encode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	type Health struct {
		Status string `json:"status"`
	}

	health := Health{
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.Printf("Failed to encode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
