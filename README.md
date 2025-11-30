# E-Commerce Platform

A Go-based REST API designed for a e-commerce platform.

## Features

### Authentication & Authorization
- User registration and login
- JWT-based authentication
- Role-based access control (buyer/seller)
- Password encryption with bcrypt

### Product Management
- Create, read, update, and delete products
- Seller-specific product listings
- Stock management
- Low-stock alerts (threshold: 5 units)

### Shopping Cart
- Add items to cart
- Update quantities
- Remove items
- Cart validation with stock checking

### Order Management
- Create orders from cart items
- Automatic stock deduction
- Order history for buyers
- Per-product status tracking (pending/completed)

### Multi-Seller Support
- Independent seller dashboards
- Sellers see only their products in orders
- Each seller can fulfill their items independently
- Revenue tracking per seller

### Seller Dashboard
- Product analytics
- Order item statistics
- Revenue calculation (completed items only)
- Pending and completed order counts

## Tech Stack

- **Language:** Go 1.21
- **Framework:** Gorilla Mux
- **Database:** PostgreSQL with GORM
- **Authentication:** JWT
- **Security:** bcrypt password hashing

## Run Locally

1. **Install Prerequisites**
   - Go 1.21+
   - PostgreSQL

2. **Setup Database**
   ```bash
   createdb ecommerce
   ```
    Database tables are created automatically on startup.
3. **Clone and Install**
   ```bash
   git clone <repo-url>
   cd e-commerce
   go mod download
   ```

4. **Configure (Optional)**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=ecommerce
   ```

5. **Run**
   ```bash
   go run app/main.go
   ```
