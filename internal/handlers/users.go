package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MdHisham-04/E-Commerce/internal/database"
	"github.com/MdHisham-04/E-Commerce/internal/models"
	"github.com/gorilla/mux"
)

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	result := database.DB.Find(&users)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser returns a single user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	result := database.DB.First(&user, id)

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if user.Email == "" || user.Name == "" {
		http.Error(w, "Email and name are required", http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
