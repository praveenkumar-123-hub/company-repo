package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"companyintersala/internal/api/auth"
	"companyintersala/internal/api/db"
	"companyintersala/internal/api/middleware"
	"companyintersala/internal/api/models"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}
	if req.Role == "" {
		req.Role = models.RoleUser
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password, role) VALUES (?,?,?)", req.Username, hashedPassword, string(req.Role))
	if err != nil {
		http.Error(w, "username already exists or database error", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", req.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.Username, string(user.Role))
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := db.DB.QueryRow("SELECT id, username, role, created_at FROM users WHERE username = ?", claims.Username).Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, username, role, created_at FROM users")
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt); err != nil {
			http.Error(w, "error scanning users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}
