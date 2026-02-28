package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON 인코딩 실패: %v", err)
	}
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func validateUser(user User) string {
	if strings.TrimSpace(user.Name) == "" {
		return "name은 필수입니다."
	}
	if strings.TrimSpace(user.Email) == "" {
		return "email은 필수입니다."
	}
	if !strings.Contains(user.Email, "@") {
		return "email 형식이 올바르지 않습니다."
	}
	return ""
}

// POST /users
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("JSON 디코딩 실패: %v", err)
		http.Error(w, "잘못된 요청 형식입니다.", http.StatusBadRequest)
		return
	}

	if msg := validateUser(user); msg != "" {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		user.Name, user.Email, user.Age,
	)
	if err != nil {
		log.Printf("유저 생성 실패: %v", err)
		http.Error(w, "유저 생성에 실패했습니다.", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	user.ID = uint64(id)

	writeJSON(w, http.StatusCreated, user)
}

// GET /users
func getUsers(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		log.Printf("유저 목록 조회 실패: %v", err)
		http.Error(w, "유저 목록 조회에 실패했습니다", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows 닫기 실패: %v", err)
		}
	}()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
			log.Printf("rows 스캔 실패: %v", err)
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows 순회 중 에러: %v", err)
		http.Error(w, "데이터 처리 중 오류가 발생했습니다", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// GET /users/{id}
func getUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var user User
	err := db.QueryRow(
		"SELECT id, name, email, age FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age)

	if err == sql.ErrNoRows {
		http.Error(w, "유저를 찾을 수 없습니다", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("유저 조회 실패 (id=%d): %v", id, err)
		http.Error(w, "유저 조회에 실패했습니다", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// PUT /users/{id}
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("JSON 디코딩 실패: %v", err)
		http.Error(w, "잘못된 요청 형식입니다", http.StatusBadRequest)
		return
	}

	if msg := validateUser(user); msg != "" {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"UPDATE users SET name=?, email=?, age=? WHERE id=?",
		user.Name, user.Email, user.Age, id,
	)
	if err != nil {
		log.Printf("유저 업데이트 실패 (id=%d): %v", id, err)
		http.Error(w, "유저 업데이트에 실패했습니다", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "유저를 찾을 수 없습니다", http.StatusNotFound)
		return
	}

	user.ID = uint64(id)
	writeJSON(w, http.StatusOK, user)
}

// DELETE /users/{id}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		log.Printf("유저 삭제 실패 (id=%d): %v", id, err)
		http.Error(w, "유저 삭제에 실패했습니다", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "유저를 찾을 수 없습니다", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /health
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	status := "OK"

	// DB 상태도 함께 체크
	if err := db.Ping(); err != nil {
		log.Printf("헬스체크 DB 핑 실패: %v", err)
		status = "DB_UNAVAILABLE"
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": status})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": status})
}
