# project-app-inventory-restapi-golang-azwin

Mini Project Bootcamp Backend Golang Batch 3 (Inventory System RESTful API)

## Features

- RESTful API with Chi Router
- PostgreSQL Database with pgx driver
- Daily Log Rotation with Zap Logger
- Request Validation
- Middleware (Logging, Authentication)
- Clean Architecture (Handler → Service → Repository)
- CRUD Operations for:
  - Users
  - Categories
  - Items
  - Racks
  - Warehouses
  - Sales
- Reporting System:
  - Items Report (Total barang & stock)
  - Sales Report (Total transaksi & penjualan)
  - Revenue Report (Pendapatan & rata-rata)
- **Low Stock Alert**: Monitor barang dengan stock di bawah threshold minimum

## Quick Start

### 1. Setup Environment

```bash
# Copy environment file
cp .env.example .env

# Edit .env sesuai kebutuhan
nano .env
```

### 2. Run Application

```bash
# Install dependencies
go mod download

# Run application
go run main.go
```

Application will start on `http://localhost:8080`

## Logging System

Aplikasi ini memiliki sistem logging otomatis yang mencatat semua aktivitas ke file log harian.

### Features:

- **Daily Log Files** - File dibuat otomatis per hari: `logs/app-2026-01-04.log`
- **Auto Rotation** - File > 10MB otomatis di-rotate
- **Compression** - File lama di-compress dengan gzip
- **30 Days Retention** - File > 30 hari dihapus otomatis
- **Dual Output** - Log ke file (JSON) dan console (readable)

### Quick View Logs:

```bash
# Real-time tail
tail -f logs/app-2026-01-04.log

# View with jq (pretty print)
cat logs/app-2026-01-04.log | jq '.'
```

## Project Structure

```
.
├── database/           # Database connection & migrations
├── dto/               # Data Transfer Objects
├── handler/           # HTTP handlers (controllers)
├── middleware/        # HTTP middlewares
├── model/            # Database models
├── repository/       # Database operations
├── router/           # Route definitions
├── service/          # Business logic
├── utils/            # Utilities (logger, validator, etc)
├── logs/             # Log files (auto-created)
├── docs/             # Documentation
├── main.go           # Application entry point
└── .env.example      # Environment template
```

## API Endpoints

### Reports

- `GET /reports/items` - Total barang & stock
- `GET /reports/sales` - Total penjualan & transaksi
- `GET /reports/revenue` - Total pendapatan & rata-rata

### Items

- `GET /items` - Get all items (with pagination)
- `GET /items/{id}` - Get item by ID
- `GET /items/low-stock` - Get items dengan stock rendah
- `POST /items` - Create item
- `PUT /items/{id}` - Update item
- `DELETE /items/{id}` - Delete item

### Users

- `GET /users` - Get all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create user
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Categories, Racks, Warehouses, Sales

Similar CRUD operations for each resource.

## Configuration

Edit `.env` file:

```env
# Application
APP_NAME=Inventory REST API
PORT=8080
DEBUG=true              # true = development, false = production
LIMIT=10                # Default pagination limit

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=your_password
DATABASE_NAME=inventory_db
DATABASE_MAX_CONN=10

# Logging
PATH_LOGGING=logs/      # Log directory
```

## Tech Stack

- **Language**: Go 1.21+
- **Router**: Chi
- **Database**: PostgreSQL with pgx/v5
- **Logger**: Zap (uber-go)
- **Validation**: go-playground/validator
- **Config**: Viper
- **Log Rotation**: lumberjack

## Dependencies

```bash
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
go get go.uber.org/zap
go get github.com/go-playground/validator/v10
go get github.com/spf13/viper
go get gopkg.in/natefinch/lumberjack.v2
```

## Development

```bash
# Run with hot reload (air)
air

# Run tests
go test ./...

# Build
go build -o app main.go

# Run binary
./app
```

## License

MIT License - Lumoshive Bootcamp Backend Golang Batch 3
