package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MdHisham-04/E-Commerce/internal/database"
	"github.com/MdHisham-04/E-Commerce/internal/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateOrder creates an order from cart items
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get cart items
	var cartItems []models.CartItem
	if err := tx.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(cartItems) == 0 {
		tx.Rollback()
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// Calculate total and check stock
	var total float64
	for _, item := range cartItems {
		if item.Product.Stock < item.Quantity {
			tx.Rollback()
			http.Error(w, "Insufficient stock for "+item.Product.Name, http.StatusBadRequest)
			return
		}
		total += item.Product.Price * float64(item.Quantity)
	}

	// Create order
	order := models.Order{
		UserID: userID,
		Total:  total,
		Status: "pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create order items and update stock
	for _, item := range cartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
			Status:    "pending",
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update product stock
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Clear cart
	if err := tx.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load order with items
	database.DB.Preload("OrderItems.Product").First(&order, order.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrders returns all orders for a user
func GetOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var orders []models.Order
	result := database.DB.Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrder returns a single order
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["order_id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.Order
	result := database.DB.Preload("OrderItems.Product").First(&order, orderID)

	if result.Error != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
