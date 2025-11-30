package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MdHisham-04/E-Commerce/internal/database"
	"github.com/MdHisham-04/E-Commerce/internal/models"
	"github.com/gorilla/mux"
)

// GetProducts returns all products available in the store
func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	result := database.DB.Preload("Seller").Find(&products)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetProduct returns a single product by ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	result := database.DB.Preload("Seller").First(&product, id)

	if result.Error != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
