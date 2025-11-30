package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MdHisham-04/E-Commerce/internal/database"
	"github.com/MdHisham-04/E-Commerce/internal/handlers"
	"github.com/MdHisham-04/E-Commerce/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// main function
func main() {
	config := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "ecommerce"),
	}

	if err := database.Connect(config); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	api.HandleFunc("/auth/register", handlers.Register).Methods("POST")
	api.HandleFunc("/auth/login", handlers.Login).Methods("POST")

	api.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	api.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/users/{user_id}/cart", handlers.GetCart).Methods("GET")
	protected.HandleFunc("/users/{user_id}/cart", handlers.AddToCart).Methods("POST")
	protected.HandleFunc("/users/{user_id}/cart/{item_id}", handlers.RemoveFromCart).Methods("DELETE")
	protected.HandleFunc("/users/{user_id}/orders", handlers.GetOrders).Methods("GET")
	protected.HandleFunc("/users/{user_id}/orders", handlers.CreateOrder).Methods("POST")

	seller := protected.PathPrefix("/seller").Subrouter()
	seller.Use(middleware.RequireRole("seller"))

	seller.HandleFunc("/products", handlers.GetSellerProducts).Methods("GET")
	seller.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	seller.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")
	seller.HandleFunc("/products/{id}/stock", handlers.UpdateProductStock).Methods("PATCH")
	seller.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	seller.HandleFunc("/orders", handlers.GetAllOrders).Methods("GET")
	seller.HandleFunc("/orders/pending", handlers.GetPendingOrders).Methods("GET")
	seller.HandleFunc("/order-items/{item_id}/status", handlers.UpdateOrderItemStatus).Methods("PATCH")

	seller.HandleFunc("/dashboard/stats", handlers.GetDashboardStats).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets")))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// getEnv retrieves environment variables with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
