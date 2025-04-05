# Loyalty-Points-System-API

The Loyalty Point System is a robust backend API built with Golang, Gin, and PostgreSQL for managing loyalty points in an e-commerce platform. Users earn points with each purchase and can redeem them as discounts during future orders. The system intelligently ensures a maximum of 10% of the order value can be redeemed through points and automatically utilizes available points for optimal discounts. It includes secure authentication and a cron-based expiry system for unused points.
The project is built with a clean architecture approach.

## Key Features

- **Automatic Points Calculation:** Earn and redeem points based on purchase transactions with 10% redemption cap logic.

- **Point Expiry Scheduler:** Background job that expires unused points based on custom business logic (e.g., after 1 year).

- **Audit Logging:** Logs every event of earning, redeeming, and expiring points to maintain data traceability.

- **Data Consistency:** Uses atomic DB transactions and validations to ensure accuracy in point updates.

- **Secure Auth System:** JWT-based authentication with role-specific access for different user types.

- **Clean Architecture:** Structured into handlers, use cases, and repositories for high maintainability.

## Installation

To set up the project locally, follow these steps:

1. **Clone the Repository:**

     ```bash
    git clone https://github.com/abin-saji-2003/Loyalty-Points-System-API.git
    cd Loyalty-Points-System-API
    ```
2. **Set Up the Environment Variables:**

    Create a `.env` file in the root directory and add the following variables:

    ```bash
    DB_HOST=localhost
    DB_USER=your_database_username
    DB_PASSWORD=your_database_password
    DB_NAME=your_database_name
    DB_PORT=5432
    DB_SSLMODE=your_database_sslmode
    PORT=your_port
    ACCESS_SECRET=your_access_secrete
    REFRESH_SECRET=your_refresh_secrete 
    ```

3. **Install Dependencies:**

    ```bash
    go mod tidy
    ```

4. **Run the Application:**

    ```bash
    go run cmd/api/main.go
    ```

5. **Running Cron Jobs:**

    ```bash
    go run cmd/cron/main.go
    ```

## Example API Usage

### 1. Login (User Authentication)

**Endpoint:**  
`POST /api/auth/login`

**Request Body:**

```json
{
  "email": "abin@gmail.com",
  "password": "abin123"
}

```

### 2. Refresh Token

**Endpoint:**  
`POST /api/auth/refresh`

**Request Body:**

```json
{
  "refresh_token": "your_refresh_token_here"
}

```

### 3. Record a Transaction

**Endpoint:**  
`POST /api/transaction`

**Request Body:**

```json
{
  "transaction_id": "be0a5dfc-6b27-4ef5-8d46-5d70f5e6e0d1",
  "user_id": 1,
  "transaction_amount": 1100,
  "category": "groceries",
  "transaction_date": "2025-04-04",
  "product_code": "ELEC-X102",
  "use_points": true
}

```

### 4. View Point History (Paginated)

**Endpoint:**  
`GET /api/points/1/history?page=1&limit=5`

### 5. Filter Point History by Type and Date Range

**Endpoint:**  
`GET /api/points/1/history/filter?tx_type=redeem&start_date=2025-01-01&end_date=2025-05-01`
