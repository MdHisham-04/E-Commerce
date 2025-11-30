package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MdHisham-04/E-Commerce/internal/auth"
	"github.com/MdHisham-04/E-Commerce/internal/database"
	"github.com/MdHisham-04/E-Commerce/internal/models"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// sendError writes an error response in JSON format
func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"` // buyer or seller
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// Register creates a new user account
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Name == "" || req.Password == "" {
		sendError(w, "Email, name, and password are required", http.StatusBadRequest)
		return
	}

	if req.Role != "buyer" && req.Role != "seller" {
		req.Role = "buyer"
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		sendError(w, "Email already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		sendError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login authenticates a user
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		sendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		sendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
