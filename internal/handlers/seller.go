package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MdHisham-04/E-Commerce/internal/database"
	"github.com/MdHisham-04/E-Commerce/internal/middleware"
	"github.com/MdHisham-04/E-Commerce/internal/models"
	"github.com/gorilla/mux"
)

// GetSellerProducts returns all products belonging to the authenticated seller
func GetSellerProducts(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var products []models.Product
	result := database.DB.Where("seller_id = ?", claims.UserID).Order("created_at DESC").Find(&products)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetAllOrders retrieves all orders containing the seller's products with filtered order items
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var orders []models.Order
	result := database.DB.
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ?", claims.UserID).
		Preload("User").
		Preload("OrderItems", "product_id IN (SELECT id FROM products WHERE seller_id = ?)", claims.UserID).
		Preload("OrderItems.Product").
		Group("orders.id").
		Order("orders.created_at DESC").
		Find(&orders)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetPendingOrders retrieves orders with pending status containing the seller's products
func GetPendingOrders(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var orders []models.Order
	result := database.DB.
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ? AND orders.status = ?", claims.UserID, "pending").
		Preload("User").
		Preload("OrderItems", "product_id IN (SELECT id FROM products WHERE seller_id = ?)", claims.UserID).
		Preload("OrderItems.Product").
		Group("orders.id").
		Order("orders.created_at ASC").
		Find(&orders)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderItemStatus allows sellers to update the fulfillment status of their order items
func UpdateOrderItemStatus(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	vars := mux.Vars(r)
	orderItemID, err := strconv.Atoi(vars["item_id"])
	if err != nil {
		http.Error(w, "Invalid order item ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validStatuses := map[string]bool{"pending": true, "completed": true}
	if !validStatuses[req.Status] {
		http.Error(w, "Invalid status. Only 'pending' or 'completed' allowed", http.StatusBadRequest)
		return
	}

	// Verify this order item belongs to a product owned by this seller
	var orderItem models.OrderItem
	result := database.DB.
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("order_items.id = ? AND products.seller_id = ?", orderItemID, claims.UserID).
		First(&orderItem)

	if result.Error != nil {
		http.Error(w, "Order item not found or access denied", http.StatusNotFound)
		return
	}

	// Update order item status
	database.DB.Model(&orderItem).Update("status", req.Status)

	// Load updated order item with product
	database.DB.Preload("Product").First(&orderItem, orderItemID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderItem)
}

// CreateProduct creates a new product for the authenticated seller
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.SellerID = claims.UserID

	if err := database.DB.Create(&product).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// UpdateProductStock updates the stock quantity of a seller's product
func UpdateProductStock(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Stock int `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Stock < 0 {
		http.Error(w, "Stock cannot be negative", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND seller_id = ?", productID, claims.UserID).First(&product).Error; err != nil {
		http.Error(w, "Product not found or access denied", http.StatusNotFound)
		return
	}

	product.Stock = req.Stock
	database.DB.Save(&product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct updates product details for the authenticated seller
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND seller_id = ?", productID, claims.UserID).First(&product).Error; err != nil {
		http.Error(w, "Product not found or access denied", http.StatusNotFound)
		return
	}

	var updates models.Product
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updates.Name != "" {
		product.Name = updates.Name
	}
	if updates.Description != "" {
		product.Description = updates.Description
	}
	if updates.Price > 0 {
		product.Price = updates.Price
	}
	if updates.Stock >= 0 {
		product.Stock = updates.Stock
	}

	database.DB.Save(&product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct removes a product owned by the authenticated seller
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	result := database.DB.Where("id = ? AND seller_id = ?", productID, claims.UserID).Delete(&models.Product{})

	if result.RowsAffected == 0 {
		http.Error(w, "Product not found or access denied", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetDashboardStats returns statistics for the seller's dashboard including products, orders, and revenue
func GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var stats struct {
		TotalProducts       int64   `json:"total_products"`
		TotalOrderItems     int64   `json:"total_order_items"`
		PendingOrderItems   int64   `json:"pending_order_items"`
		CompletedOrderItems int64   `json:"completed_order_items"`
		TotalRevenue        float64 `json:"total_revenue"`
		LowStockProducts    int64   `json:"low_stock_products"`
	}

	database.DB.Model(&models.Product{}).Where("seller_id = ?", claims.UserID).Count(&stats.TotalProducts)

	// Count all order items for seller's products
	database.DB.Model(&models.OrderItem{}).
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ?", claims.UserID).
		Count(&stats.TotalOrderItems)

	// Count pending order items
	database.DB.Model(&models.OrderItem{}).
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ? AND order_items.status = ?", claims.UserID, "pending").
		Count(&stats.PendingOrderItems)

	// Count completed order items
	database.DB.Model(&models.OrderItem{}).
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ? AND order_items.status = ?", claims.UserID, "completed").
		Count(&stats.CompletedOrderItems)

	database.DB.Model(&models.Product{}).
		Where("seller_id = ? AND stock < ?", claims.UserID, 5).
		Count(&stats.LowStockProducts)

	// Calculate revenue from completed order items
	database.DB.Table("order_items").
		Select("COALESCE(SUM(order_items.price * order_items.quantity), 0)").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ? AND order_items.status = ?", claims.UserID, "completed").
		Scan(&stats.TotalRevenue)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
