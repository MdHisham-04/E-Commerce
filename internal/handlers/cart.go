package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/MdHisham-04/E-Commerce/internal/database"
	"github.com/MdHisham-04/E-Commerce/internal/models"
)

type AddToCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// GetCart returns all items in a user's cart
func GetCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var cartItems []models.CartItem
	result := database.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItems)
}

// AddToCart adds an item to the cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if product exists and has enough stock
	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if product.Stock < req.Quantity {
		http.Error(w, "Insufficient stock", http.StatusBadRequest)
		return
	}

	// Check if item already in cart
	var existingItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingItem)

	if result.Error == nil {
		// Update quantity - but check total doesn't exceed stock
		newQuantity := existingItem.Quantity + req.Quantity
		if product.Stock < newQuantity {
			http.Error(w, "Insufficient stock", http.StatusBadRequest)
			return
		}

		existingItem.Quantity = newQuantity
		database.DB.Save(&existingItem)

		// Reload with product data
		database.DB.Preload("Product").First(&existingItem, existingItem.ID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingItem)
		return
	}

	// Create new cart item
	cartItem := models.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := database.DB.Create(&cartItem).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load product relation
	database.DB.Preload("Product").First(&cartItem, cartItem.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cartItem)
}

// RemoveFromCart removes an item from the cart
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cartItemID, err := strconv.Atoi(vars["item_id"])
	if err != nil {
		http.Error(w, "Invalid cart item ID", http.StatusBadRequest)
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", cartItemID, userID).Delete(&models.CartItem{})
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Cart item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
