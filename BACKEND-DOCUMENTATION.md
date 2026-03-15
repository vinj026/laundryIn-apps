# Dokumentasi Teknis Backend LaundryIn

**Single Source of Truth** untuk semua developer dan AI agent yang mengerjakan project LaundryIn.

Dokumentasi ini dibuat berdasarkan pembacaan langsung seluruh kode backend yang ada di repository.

---

## Section 1 — Konfigurasi & Setup

### 1.1 main.go — Entry Point & Inisialisasi

**File:** `cmd/api/main.go`

**Urutan Inisialisasi:**

```
1. Load environment variables (godotenv.Load())
2. Register custom validators (utils.RegisterCustomValidators())
3. Connect to database (database.ConnectDB())
4. Auto-migrate models (db.AutoMigrate(...))
5. Initialize WebSocket Hub (websocket.NewHub(), go hub.RUN())
6. Dependency Injection:
   - Repository layer (NewUserRepository, NewOutletRepository, etc.)
   - Usecase layer (NewAuthUsecase, NewOutletUsecase, etc.)
   - Handler layer (NewAuthHandler, NewOutletHandler, etc.)
7. Router Setup (gin.Default(), middleware, route groups)
8. Server Start (r.Run(":" + port))
```

**Port:**
- Dari env variable `PORT`
- Default: `8080` jika tidak diset

**Cara Env Variable Dibaca:**
```go
_ = godotenv.Load()  // Load .env file (optional untuk production)
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
```

---

### 1.2 Environment Variables

| Variable | Dibaca Di | Default Value | Wajib? | Deskripsi |
|----------|-----------|---------------|--------|-----------|
| `PORT` | `main.go` | `8080` | Tidak | Port server HTTP |
| `DATABASE_URL` | `database/postgres.go` | - | Tidak* | Full connection string PostgreSQL (prioritas utama) |
| `DATABASE_PRIVATE_URL` | `database/postgres.go` | - | Tidak* | Fallback untuk Railway internal |
| `DB_HOST` | `database/postgres.go` | `localhost` (dev) / `REQUIRED_VARIABLE_MISSING` (Railway) | Tidak* | Host PostgreSQL |
| `DB_USER` | `database/postgres.go` | `postgres` | Tidak | Username PostgreSQL |
| `DB_PASSWORD` | `database/postgres.go` | - | Ya* | Password PostgreSQL |
| `DB_NAME` | `database/postgres.go` | - | Ya* | Nama database |
| `DB_PORT` | `database/postgres.go` | `5432` | Tidak | Port PostgreSQL |
| `DB_SSLMODE` | `database/postgres.go` | `require` (Railway) / `disable` (dev) | Tidak | SSL mode PostgreSQL |
| `JWT_SECRET` | `pkg/utils/jwt.go` | - | **Ya** | Secret key untuk JWT signing |
| `GIN_MODE` | `main.go` | - | Tidak | `release` untuk production, empty untuk development |
| `RAILWAY_ENVIRONMENT` | `database/postgres.go` | - | Tidak | Indicator running di Railway |
| `RAILWAY_STATIC_URL` | `database/postgres.go` | - | Tidak | Indicator running di Railway |

\* Wajib jika `DATABASE_URL` tidak diset

---

### 1.3 Database Connection

**File:** `internal/database/postgres.go`

**Cara Connection String Dibuild:**

**Prioritas 1: DATABASE_URL atau DATABASE_PRIVATE_URL**
```go
databaseURL := os.Getenv("DATABASE_URL")
if databaseURL == "" {
    databaseURL = os.Getenv("DATABASE_PRIVATE_URL")
}
if databaseURL != "" {
    dsn = databaseURL
}
```

**Prioritas 2: Individual Variables (Fallback)**
```go
host := getEnvFallback("DB_HOST", "PGHOST", "POSTGRES_HOST")
user := getEnvFallback("DB_USER", "PGUSER", "POSTGRES_USER")
password := getEnvFallback("DB_PASSWORD", "PGPASSWORD", "POSTGRES_PASSWORD")
dbName := getEnvFallback("DB_NAME", "PGDATABASE", "DB_DATABASE", "POSTGRES_DB")
port := getEnvFallback("DB_PORT", "PGPORT", "POSTGRES_PORT")
sslMode := getEnvFallback("DB_SSLMODE", "PGSSLMODE", "POSTGRES_SSLMODE")

dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
    host, user, password, dbName, port, sslMode)
```

**SSL Mode Handling:**
```go
if sslMode == "" {
    if isRailway {
        sslMode = "require"  // Force SSL di Railway
    } else {
        sslMode = "disable"  // No SSL untuk local dev
    }
}
```

**Connection Pool Settings:** Tidak ada konfigurasi eksplisit — menggunakan default GORM/PostgreSQL driver.

**Auto-Migrate Models:**
```go
db.AutoMigrate(
    &models.User{},
    &models.Outlet{},
    &models.Service{},
    &models.Order{},
    &models.OrderItem{},
    &models.OrderLog{},
    &models.Notification{}
)
```

**Extension:**
```go
db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
```

---

## Section 2 — Database Schema

### 2.1 Tabel: users

**File:** `internal/repository/models/models.go`

**Fields:**
| Column | Go Type | GORM Tag | Nullable | Default | Deskripsi |
|--------|---------|----------|----------|---------|-----------|
| `id` | `string` | `primaryKey;type:uuid` | Tidak | - | UUID primary key |
| `name` | `string` | `type:text;not null` | Tidak | - | Nama lengkap user |
| `phone` | `string` | `type:text;uniqueIndex;not null` | Tidak | - | Nomor HP (unique) |
| `email` | `string` | `type:text` | Ya | - | Email (optional) |
| `password` | `string` | `type:text;not null` | Tidak | - | Hashed password (bcrypt) |
| `role` | `string` | `type:text;not null` | Tidak | - | Role: `owner` atau `customer` |
| `created_at` | `time.Time` | (embedded) | Tidak | - | Timestamp created |
| `updated_at` | `time.Time` | (embedded) | Tidak | - | Timestamp updated |
| `deleted_at` | `gorm.DeletedAt` | `index` | Ya | - | Soft delete timestamp |

**Relations:**
- `Outlets`: One-to-Many dengan `outlets` (foreignKey: `UserID`)

**Indexes:**
- Primary: `id` (uuid)
- Unique: `phone`

**Soft Delete:** Ya (`gorm.DeletedAt`)

---

### 2.2 Tabel: outlets

**File:** `internal/repository/models/models.go`

**Fields:**
| Column | Go Type | GORM Tag | Nullable | Default | Deskripsi |
|--------|---------|----------|----------|---------|-----------|
| `id` | `string` | `primaryKey;type:uuid` | Tidak | - | UUID primary key |
| `user_id` | `string` | `type:uuid;not null;index` | Tidak | - | Owner ID (FK ke users) |
| `name` | `string` | `type:varchar(100);not null` | Tidak | - | Nama outlet |
| `address` | `string` | `type:text;not null` | Tidak | - | Alamat lengkap |
| `phone` | `string` | `type:varchar(20);not null` | Tidak | - | Nomor telepon outlet |
| `created_at` | `time.Time` | (embedded) | Tidak | - | Timestamp created |
| `updated_at` | `time.Time` | (embedded) | Tidak | - | Timestamp updated |
| `deleted_at` | `gorm.DeletedAt` | (embedded) | Ya | - | Soft delete timestamp |

**Relations:**
- `User`: Many-to-One dengan `users` (foreignKey: `UserID`, constraint: `OnDelete:CASCADE`)
- `Services`: One-to-Many dengan `services` (foreignKey: `OutletID`, constraint: `OnDelete:CASCADE`)

**Indexes:**
- Primary: `id` (uuid)
- Index: `user_id`

**Soft Delete:** Ya

---

### 2.3 Tabel: services

**File:** `internal/repository/models/models.go`

**Fields:**
| Column | Go Type | GORM Tag | Nullable | Default | Deskripsi |
|--------|---------|----------|----------|---------|-----------|
| `id` | `string` | `primaryKey;type:uuid` | Tidak | - | UUID primary key |
| `outlet_id` | `string` | `type:uuid;not null;index` | Tidak | - | Outlet ID (FK ke outlets) |
| `name` | `string` | `type:varchar(100);not null` | Tidak | - | Nama layanan |
| `price` | `decimal.Decimal` | `type:numeric(10,2);not null` | Tidak | - | Harga (decimal precision) |
| `unit` | `string` | `type:varchar(20);not null` | Tidak | - | Unit: `KG`, `PCS`, `METER` |
| `created_at` | `time.Time` | (embedded) | Tidak | - | Timestamp created |
| `updated_at` | `time.Time` | (embedded) | Tidak | - | Timestamp updated |
| `deleted_at` | `gorm.DeletedAt` | (embedded) | Ya | - | Soft delete timestamp |

**Relations:**
- `Outlet`: Many-to-One dengan `outlets` (foreignKey: `OutletID`, constraint: `OnDelete:CASCADE`)

**Indexes:**
- Primary: `id` (uuid)
- Index: `outlet_id`

**Soft Delete:** Ya

---

### 2.4 Tabel: orders

**File:** `internal/repository/models/models.go`

**Fields:**
| Column | Go Type | GORM Tag | Nullable | Default | Deskripsi |
|--------|---------|----------|----------|---------|-----------|
| `id` | `string` | `primaryKey;type:uuid` | Tidak | - | UUID primary key |
| `user_id` | `string` | `type:uuid;not null;index` | Tidak | - | Customer ID (FK ke users) |
| `outlet_id` | `string` | `type:uuid;not null;index` | Tidak | - | Outlet ID (FK ke outlets) |
| `total_price` | `decimal.Decimal` | `type:numeric(12,2);not null` | Tidak | - | Total harga estimasi |
| `final_total_price` | `*decimal.Decimal` | `type:numeric(12,2)` | Ya | - | Total harga final (setelah timbang) |
| `status` | `string` | `type:varchar(20);not null;default:'pending'` | Tidak | `pending` | Status order |
| `order_date` | `time.Time` | `autoCreateTime` | Tidak | - | Timestamp order dibuat |
| `created_at` | `time.Time` | (embedded) | Tidak | - | Timestamp created |
| `updated_at` | `time.Time` | (embedded) | Tidak | - | Timestamp updated |
| `deleted_at` | `gorm.DeletedAt` | (embedded) | Ya | - | Soft delete timestamp |

**Relations:**
- `User`: Many-to-One dengan `users` (foreignKey: `UserID`)
- `Outlet`: Many-to-One dengan `outlets` (foreignKey: `OutletID`)
- `Items`: One-to-Many dengan `order_items` (foreignKey: `OrderID`, constraint: `OnDelete:CASCADE`)
- `Logs`: One-to-Many dengan `order_logs` (foreignKey: `OrderID`, constraint: `OnDelete:CASCADE`)

**Indexes:**
- Primary: `id` (uuid)
- Index: `user_id`
- Index: `outlet_id`

**Soft Delete:** Ya

---

### 2.5 Tabel: order_items

**File:** `internal/repository/models/models.go`

**Fields:**
| Column | Go Type | GORM Tag | Nullable | Default | Deskripsi |
|--------|---------|----------|----------|---------|-----------|
| `id` | `string` | `primaryKey;type:uuid` | Tidak | - | UUID primary key |
| `order_id` | `string` | `type:uuid;not null;index` | Tidak | - | Order ID (FK ke orders) |
| `service_name` | `string` | `type:varchar(100);not null` | Tidak | - | Nama layanan (snapshot) |
| `service_price` | `decimal.Decimal` | `type:numeric(10,2);not null` | Tidak | - | Harga per unit (snapshot) |
| `qty` | `decimal.Decimal` | `type:numeric(6,2);not null` | Tidak | - | Quantity estimasi |
| `actual_qty` | `*decimal.Decimal` | `type:numeric(6,2)` | Ya | - | Quantity aktual (setelah timbang) |
| `unit` | `string` | `type:varchar(20);not null` | Tidak | - | Unit: `KG`, `PCS`, `METER` |
| `subtotal` | `decimal.Decimal` | `type:numeric(12,2);not null` | Tidak | - | Subtotal estimasi (qty × price) |
| `final_price` | `*decimal.Decimal` | `type:numeric(12,2)` | Ya | - | Harga final (actual_qty × price) |

**Relations:**
- `Order`: Many-to-One dengan `orders` (foreignKey: `OrderID`, constraint: `OnDelete:CASCADE`)

**Indexes:**
- Primary: `id` (uuid)
- Index: `order_id`

**Soft Delete:** Tidak (tidak ada `deleted_at`)

---

### 2.6 Tabel: order_logs

**File:** `internal/repository/models/models.go`

**Fields:**
| Column | Go Type | GORM Tag | Nullable | Default | Deskripsi |
|--------|---------|----------|----------|---------|-----------|
| `id` | `string` | `primaryKey;type:uuid` | Tidak | - | UUID primary key |
| `order_id` | `string` | `type:uuid;not null;index` | Tidak | - | Order ID (FK ke orders) |
| `updated_by` | `string` | `type:uuid;not null` | Tidak | - | User yang update status (FK ke users) |
| `old_status` | `string` | `type:varchar(20)` | Ya | - | Status sebelumnya |
| `new_status` | `string` | `type:varchar(20);not null` | Tidak | - | Status baru |
| `created_at` | `time.Time` | `autoCreateTime` | Tidak | - | Timestamp perubahan |

**Relations:**
- `Order`: Many-to-One dengan `orders` (foreignKey: `OrderID`, constraint: `OnDelete:CASCADE`)
- `User`: Many-to-One dengan `users` (foreignKey: `UpdatedBy`)

**Indexes:**
- Primary: `id` (uuid)
- Index: `order_id`

**Soft Delete:** Tidak

---

### 2.7 Tabel: notifications

**File:** `internal/repository/models/models.go`

**Fields:**
| Column | Go Type | GORM Tag | Nullable | Default | Deskripsi |
|--------|---------|----------|----------|---------|-----------|
| `id` | `string` | `primaryKey;type:uuid` | Tidak | - | UUID primary key |
| `user_id` | `string` | `type:uuid;not null;index:idx_notif_user_read;index:idx_notif_user_created` | Tidak | - | User ID penerima (FK ke users) |
| `type` | `string` | `type:varchar(50);not null` | Tidak | - | Tipe notifikasi |
| `title` | `string` | `type:varchar(200);not null` | Tidak | - | Judul notifikasi |
| `body` | `string` | `type:varchar(500);not null` | Tidak | - | Isi notifikasi |
| `data` | `string` | `type:jsonb` | Ya | - | Data tambahan (JSON) |
| `is_read` | `bool` | `default:false;index:idx_notif_user_read` | Tidak | `false` | Status sudah dibaca |
| `created_at` | `time.Time` | `autoCreateTime;index:idx_notif_user_created,sort:desc` | Tidak | - | Timestamp dibuat |

**Relations:** Tidak ada relasi GORM (notification menyimpan data snapshot)

**Indexes:**
- Primary: `id` (uuid)
- Composite: `idx_notif_user_read` (`user_id`, `is_read`)
- Composite: `idx_notif_user_created` (`user_id`, `created_at DESC`)

**Soft Delete:** Tidak

---

## Section 3 — DTOs

### 3.1 RegisterRequest

**File:** `internal/dto/dto.go`

**Digunakan di:** `POST /api/v1/auth/register`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `Name` | `name` | `string` | `required,min=2,max=100` | Ya | Nama lengkap |
| `Phone` | `phone` | `string` | `required,e164_strict` | Ya | Nomor HP format E.164 |
| `Email` | `email` | `string` | `omitempty,email` | Tidak | Email (optional) |
| `Password` | `password` | `string` | `required,min=8,max=64` | Ya | Password (min 8 char) |
| `Role` | `role` | `string` | `required,oneof=owner customer` | Ya | Role user |

**Contoh JSON:**
```json
{
  "name": "John Doe",
  "phone": "+6281234567890",
  "email": "john@example.com",
  "password": "Password123",
  "role": "customer"
}
```

---

### 3.2 LoginRequest

**File:** `internal/dto/dto.go`

**Digunakan di:** `POST /api/v1/auth/login`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `Phone` | `phone` | `string` | `required` | Ya | Nomor HP |
| `Password` | `password` | `string` | `required` | Ya | Password |

**Contoh JSON:**
```json
{
  "phone": "+6281234567890",
  "password": "Password123"
}
```

---

### 3.3 UserResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** Semua response yang mengandung user data

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `ID` | `id` | `string` | Ya | User ID (UUID) |
| `Name` | `name` | `string` | Ya | Nama lengkap |
| `Phone` | `phone` | `string` | Ya | Nomor HP |
| `Email` | `email,omitempty` | `string` | Tidak | Email (optional) |
| `Role` | `role` | `string` | Ya | Role user |

**⚠️ Penting:** `Password` TIDAK PERNAH termasuk di response.

**Contoh JSON:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "phone": "+6281234567890",
  "email": "john@example.com",
  "role": "customer"
}
```

---

### 3.4 AuthResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** `POST /api/v1/auth/register`, `POST /api/v1/auth/login`

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `Token` | `token` | `string` | Ya | JWT token |
| `User` | `user` | `UserResponse` | Ya | Data user |

**Contoh JSON:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "phone": "+6281234567890",
    "role": "customer"
  }
}
```

---

### 3.5 OutletRequest

**File:** `internal/dto/dto.go`

**Digunakan di:** `POST /api/v1/outlets`, `PUT /api/v1/outlets/:id`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `Name` | `name` | `string` | `required,min=3,max=100` | Ya | Nama outlet |
| `Address` | `address` | `string` | `required,min=10,max=500` | Ya | Alamat lengkap |
| `Phone` | `phone` | `string` | `required,e164_strict` | Ya | Nomor telepon |

**⚠️ Penting:** `UserID` TIDAK PERNAH diterima dari frontend — diinject dari JWT context.

**Contoh JSON:**
```json
{
  "name": "Laundry Express Cabang 1",
  "address": "Jl. Sudirman No. 123, Jakarta",
  "phone": "+6281234567890"
}
```

---

### 3.6 OutletResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** Semua endpoint outlet

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `ID` | `id` | `string` | Ya | Outlet ID (UUID) |
| `Name` | `name` | `string` | Ya | Nama outlet |
| `Address` | `address` | `string` | Ya | Alamat |
| `Phone` | `phone` | `string` | Ya | Nomor telepon |
| `CreatedAt` | `created_at` | `string` | Ya | ISO 8601 timestamp |
| `UpdatedAt` | `updated_at` | `string` | Ya | ISO 8601 timestamp |

**Contoh JSON:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Laundry Express Cabang 1",
  "address": "Jl. Sudirman No. 123, Jakarta",
  "phone": "+6281234567890",
  "created_at": "2025-03-16T10:00:00Z",
  "updated_at": "2025-03-16T10:00:00Z"
}
```

---

### 3.7 PaginationQuery

**File:** `internal/dto/dto.go`

**Digunakan di:** Semua endpoint dengan pagination

**Arah:** Request (Query Params)

**Fields:**
| Field | Form Key | Go Type | Binding/Validation | Default | Deskripsi |
|-------|----------|---------|-------------------|---------|-----------|
| `Page` | `page` | `int` | `min=1` | `1` | Halaman |
| `Limit` | `limit` | `int` | `min=1,max=100` | `10` | Items per halaman (max 100) |

---

### 3.8 PaginatedResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** Semua endpoint dengan pagination

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `Data` | `data` | `interface{}` | Ya | Array data |
| `Page` | `page` | `int` | Ya | Halaman saat ini |
| `Limit` | `limit` | `int` | Ya | Items per halaman |
| `Total` | `total` | `int64` | Ya | Total items |
| `TotalPages` | `total_pages` | `int` | Ya | Total halaman |

**Contoh JSON:**
```json
{
  "data": [...],
  "page": 1,
  "limit": 10,
  "total": 50,
  "total_pages": 5
}
```

---

### 3.9 ServiceRequest

**File:** `internal/dto/dto.go`

**Digunakan di:** `POST /api/v1/services`, `PUT /api/v1/services/:id`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `OutletID` | `outlet_id` | `string` | `required,uuid` | Ya | Outlet ID (UUID) |
| `Name` | `name` | `string` | `required,min=3,max=100` | Ya | Nama layanan |
| `Price` | `price` | `string` | `required` | Ya | Harga (string untuk decimal precision) |
| `Unit` | `unit` | `string` | `required,oneof=KG PCS METER` | Ya | Unit |

**Contoh JSON:**
```json
{
  "outlet_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Cuci Kering Regular",
  "price": "7000",
  "unit": "KG"
}
```

---

### 3.10 ServiceResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** Semua endpoint service

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `ID` | `id` | `string` | Ya | Service ID (UUID) |
| `OutletID` | `outlet_id` | `string` | Ya | Outlet ID |
| `Name` | `name` | `string` | Ya | Nama layanan |
| `Price` | `price` | `string` | Ya | Harga (2 decimal places) |
| `Unit` | `unit` | `string` | Ya | Unit |
| `CreatedAt` | `created_at` | `string` | Ya | ISO 8601 timestamp |
| `UpdatedAt` | `updated_at` | `string` | Ya | ISO 8601 timestamp |

**Contoh JSON:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "outlet_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Cuci Kering Regular",
  "price": "7000.00",
  "unit": "KG",
  "created_at": "2025-03-16T10:00:00Z",
  "updated_at": "2025-03-16T10:00:00Z"
}
```

---

### 3.11 OrderItemRequest

**File:** `internal/dto/dto.go`

**Digunakan di:** `POST /api/v1/orders`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `ServiceID` | `service_id` | `string` | `required,uuid` | Ya | Service ID (UUID) |
| `Qty` | `qty` | `string` | `required` | Ya | Quantity (string untuk decimal precision) |

**Contoh JSON:**
```json
{
  "service_id": "550e8400-e29b-41d4-a716-446655440001",
  "qty": "2.5"
}
```

---

### 3.12 OrderRequest

**File:** `internal/dto/dto.go`

**Digunakan di:** `POST /api/v1/orders`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `OutletID` | `outlet_id` | `string` | `required,uuid` | Ya | Outlet ID (UUID) |
| `Items` | `items` | `[]OrderItemRequest` | `required,min=1,dive` | Ya | List items (min 1) |

**Contoh JSON:**
```json
{
  "outlet_id": "550e8400-e29b-41d4-a716-446655440000",
  "items": [
    {
      "service_id": "550e8400-e29b-41d4-a716-446655440001",
      "qty": "2.5"
    }
  ]
}
```

---

### 3.13 OrderStatusRequest

**File:** `internal/dto/dto.go`

**Digunakan di:** `PATCH /api/v1/orders/:id/status`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `Status` | `status` | `string` | `required,oneof=pending process completed picked_up cancelled` | Ya | Status baru |
| `Items` | `items` | `[]ActualQtyItem` | - | Tidak | List actual qty (untuk status `process`) |

**Contoh JSON (transition ke process):**
```json
{
  "status": "process",
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "actual_qty": "2.3"
    }
  ]
}
```

---

### 3.14 ActualQtyItem

**File:** `internal/dto/dto.go`

**Digunakan di:** `PATCH /api/v1/orders/:id/status`

**Arah:** Request

**Fields:**
| Field | JSON Key | Go Type | Binding/Validation | Required | Deskripsi |
|-------|----------|---------|-------------------|----------|-----------|
| `ID` | `id` | `string` | `required,uuid` | Ya | OrderItem ID (UUID) |
| `ActualQty` | `actual_qty` | `string` | `required` | Ya | Quantity aktual |

---

### 3.15 OrderItemResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** Response order

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `ID` | `id` | `string` | Ya | OrderItem ID |
| `ServiceName` | `service_name` | `string` | Ya | Nama layanan (snapshot) |
| `ServicePrice` | `service_price` | `string` | Ya | Harga per unit |
| `Qty` | `qty` | `string` | Ya | Quantity estimasi |
| `ActualQty` | `actual_qty,omitempty` | `*string` | Tidak | Quantity aktual (jika ada) |
| `Unit` | `unit` | `string` | Ya | Unit |
| `Subtotal` | `subtotal` | `string` | Ya | Subtotal estimasi |
| `FinalPrice` | `final_price,omitempty` | `*string` | Tidak | Harga final (jika ada) |

---

### 3.16 OrderLogResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** Response order (logs)

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `OldStatus` | `old_status` | `string` | Ya | Status sebelumnya |
| `NewStatus` | `new_status` | `string` | Ya | Status baru |
| `UpdatedBy` | `updated_by` | `string` | Ya | User ID yang update |
| `CreatedAt` | `created_at` | `string` | Ya | ISO 8601 timestamp |

---

### 3.17 OrderResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** Semua endpoint order

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `ID` | `id` | `string` | Ya | Order ID |
| `UserID` | `user_id` | `string` | Ya | Customer ID |
| `CustomerName` | `customer_name,omitempty` | `string` | Tidak | Nama customer (dari preload User) |
| `OutletID` | `outlet_id` | `string` | Ya | Outlet ID |
| `TotalPrice` | `total_price` | `string` | Ya | Total estimasi |
| `FinalTotalPrice` | `final_total_price,omitempty` | `*string` | Tidak | Total final |
| `Status` | `status` | `string` | Ya | Status order |
| `OrderDate` | `order_date` | `string` | Ya | ISO 8601 timestamp |
| `Items` | `items,omitempty` | `[]OrderItemResponse` | Tidak | List items |
| `Logs` | `logs,omitempty` | `[]OrderLogResponse` | Tidak | List logs |

---

### 3.18 NotificationResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** `GET /api/v1/notifications`

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `ID` | `id` | `string` | Ya | Notification ID |
| `Type` | `type` | `string` | Ya | Tipe notifikasi |
| `Title` | `title` | `string` | Ya | Judul |
| `Body` | `body` | `string` | Ya | Isi |
| `Data` | `data` | `interface{}` | Ya | Data tambahan (parsed dari JSONB) |
| `IsRead` | `is_read` | `bool` | Ya | Status sudah dibaca |
| `CreatedAt` | `created_at` | `string` | Ya | ISO 8601 timestamp |

---

### 3.19 UnreadCountResponse

**File:** `internal/dto/dto.go`

**Digunakan di:** `GET /api/v1/notifications/unread-count`

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `Count` | `count` | `int64` | Ya | Jumlah unread |

---

### 3.20 ReportQuery

**File:** `internal/dto/report_dto.go`

**Digunakan di:** Semua endpoint reports

**Arah:** Request (Query Params)

**Fields:**
| Field | Form Key | Go Type | Binding/Validation | Default | Deskripsi |
|-------|----------|---------|-------------------|---------|-----------|
| `StartDate` | `start_date` | `string` | - | - | YYYY-MM-DD |
| `EndDate` | `end_date` | `string` | - | - | YYYY-MM-DD |
| `OutletID` | `outlet_id` | `string` | `omitempty,uuid` | - | Filter by outlet |

---

### 3.21 OmzetResponse

**File:** `internal/dto/report_dto.go`

**Digunakan di:** `GET /api/v1/reports/omzet`

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `TotalOmzet` | `total_omzet` | `string` | Ya | Total revenue (2 decimal places) |

---

### 3.22 OrderStatusSummaryResponse

**File:** `internal/dto/report_dto.go`

**Digunakan di:** `GET /api/v1/reports/orders/summary`

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `Pending` | `pending` | `int64` | Ya | Jumlah orders pending |
| `Process` | `process` | `int64` | Ya | Jumlah orders process |
| `Completed` | `completed` | `int64` | Ya | Jumlah orders completed |
| `PickedUp` | `picked_up` | `int64` | Ya | Jumlah orders picked_up |
| `Cancelled` | `cancelled` | `int64` | Ya | Jumlah orders cancelled |

---

### 3.23 TopServiceResponse

**File:** `internal/dto/report_dto.go`

**Digunakan di:** `GET /api/v1/reports/services/top`

**Arah:** Response

**Fields:**
| Field | JSON Key | Go Type | Required | Deskripsi |
|-------|----------|---------|----------|-----------|
| `ServiceName` | `service_name` | `string` | Ya | Nama layanan |
| `OutletName` | `outlet_name` | `string` | Ya | Nama outlet |
| `TotalQty` | `total_qty` | `string` | Ya | Total quantity (2 decimal) |
| `TotalRevenue` | `total_revenue` | `string` | Ya | Total revenue (2 decimal) |

---

## Section 4 — API Endpoints

### 4.1 Public Endpoints (Tanpa Auth)

#### GET /ping

**Handler:** Inline di `main.go`

**Middleware:** Tidak ada

**Auth Required:** Tidak

**Response Sukses:**
```json
{
  "message": "LaundryIn API is running! 🚀"
}
```

---

#### POST /api/v1/auth/register

**Handler:** `auth_handler.Register`

**Middleware:** PayloadLimit, RateLimiter

**Auth Required:** Tidak

**Request Body:**
```json
{
  "name": "string",      // required, min=2, max=100
  "phone": "string",     // required, e164_strict
  "email": "string",     // optional, email format
  "password": "string",  // required, min=8, max=64
  "role": "string"       // required, oneof=owner customer
}
```

**Response Sukses:**
- Status: `201 Created`
- Body: `AuthResponse`

**Response Error:**
- `400 Bad Request`: Validasi format data gagal / Password lemah
- `408 Request Timeout`: Proses terlalu lama
- `409 Conflict`: Nomor HP sudah terdaftar
- `500 Internal Server Error`: Error internal

**Business Logic:**
1. Sanitize input (name, phone, email)
2. Validasi password complexity (regex: uppercase, lowercase, digit)
3. Cek duplicate phone number
4. Hash password dengan bcrypt
5. Create user dengan UUID
6. Generate JWT token
7. Return token + user data

---

#### POST /api/v1/auth/login

**Handler:** `auth_handler.Login`

**Middleware:** PayloadLimit, RateLimiter

**Auth Required:** Tidak

**Request Body:**
```json
{
  "phone": "string",     // required
  "password": "string"   // required
}
```

**Response Sukses:**
- Status: `200 OK`
- Body: `AuthResponse`

**Response Error:**
- `400 Bad Request`: Validasi format data gagal
- `401 Unauthorized`: Nomor HP atau password salah
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Error internal

**Business Logic:**
1. Sanitize phone
2. Find user by phone
3. Check password dengan bcrypt
4. Generate JWT token
5. Return token + user data

---

#### GET /api/v1/public/outlets

**Handler:** `outlet_handler.GetAllOutletsPublic`

**Middleware:** PayloadLimit, RateLimiter

**Auth Required:** Tidak

**Query Params:**
- `page`: int (default: 1)
- `limit`: int (default: 10, max: 100)

**Response Sukses:**
- Status: `200 OK`
- Body: `PaginatedResponse<OutletResponse>`

**Response Error:**
- `400 Bad Request`: Format pagination tidak valid
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Error internal

---

#### GET /api/v1/public/outlets/:id

**Handler:** `outlet_handler.GetOutletByIDPublic`

**Middleware:** PayloadLimit, RateLimiter

**Auth Required:** Tidak

**Path Params:**
- `id`: UUID — Outlet ID

**Response Sukses:**
- Status: `200 OK`
- Body: `OutletResponse`

**Response Error:**
- `400 Bad Request`: Format ID tidak valid
- `404 Not Found`: Outlet tidak ditemukan
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Error internal

---

#### GET /api/v1/public/outlets/:id/services

**Handler:** `service_handler.GetAllByOutletIDPublic`

**Middleware:** PayloadLimit, RateLimiter

**Auth Required:** Tidak

**Path Params:**
- `id`: UUID — Outlet ID

**Response Sukses:**
- Status: `200 OK`
- Body: `[]ServiceResponse`

**Response Error:**
- `400 Bad Request`: Format Outlet ID tidak valid
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Error internal

---

### 4.2 Customer Endpoints

#### POST /api/v1/orders

**Handler:** `order_handler.CreateOrder`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("customer")

**Auth Required:** Ya

**Role Required:** customer

**Request Headers:**
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "outlet_id": "uuid",
  "items": [
    {
      "service_id": "uuid",
      "qty": "string"
    }
  ]
}
```

**Response Sukses:**
- Status: `201 Created`
- Body: `OrderResponse`

**Response Error:**
- `400 Bad Request`: Validasi format / outlet tidak ditemukan / layanan tidak ditemukan / qty tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role customer
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal memproses pesanan

**Business Logic:**
1. Verify outlet exists
2. For each item: verify service exists AND belongs to the outlet (Anti-IDOR)
3. Parse qty to decimal, validate > 0
4. Calculate subtotal = qty × service_price
5. Calculate grand_total = sum(subtotal)
6. Create order with status "pending"
7. Create order log entry
8. Insert order + items in transaction
9. Fire notification to owner (background)

**Catatan Penting:**
- Zero-trust pricing: frontend tidak bisa set harga, harga diambil dari database
- Deep Anti-IDOR: service harus belong to outlet yang sama

---

#### GET /api/v1/orders

**Handler:** `order_handler.GetAllByUserID`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("customer")

**Auth Required:** Ya

**Role Required:** customer

**Query Params:**
- `page`: int (default: 1)
- `limit`: int (default: 10)

**Response Sukses:**
- Status: `200 OK`
- Body: `PaginatedResponse<OrderResponse>`

**Response Error:**
- `400 Bad Request`: Format pagination tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role customer
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal mengambil daftar pesanan

---

### 4.3 Owner — Outlets

#### POST /api/v1/outlets

**Handler:** `outlet_handler.CreateOutlet`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Request Headers:**
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Request Body:** `OutletRequest`

**Response Sukses:**
- Status: `201 Created`
- Body: `OutletResponse`

**Response Error:**
- `400 Bad Request`: Validasi format data gagal
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal membuat outlet

**Business Logic:**
1. Get user_id dari JWT context
2. Sanitize inputs
3. Create outlet dengan UUID
4. Return outlet data

---

#### GET /api/v1/outlets

**Handler:** `outlet_handler.GetAllOutlets`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Query Params:**
- `page`: int (default: 1)
- `limit`: int (default: 10)

**Response Sukses:**
- Status: `200 OK`
- Body: `PaginatedResponse<OutletResponse>`

**Response Error:**
- `400 Bad Request`: Format pagination tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal mengambil data outlet

---

#### GET /api/v1/outlets/:id

**Handler:** `outlet_handler.GetOutletByID`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Outlet ID

**Response Sukses:**
- Status: `200 OK`
- Body: `OutletResponse`

**Response Error:**
- `400 Bad Request`: Format ID tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner / Bukan pemilik outlet
- `404 Not Found`: Outlet tidak ditemukan
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Error internal

**Business Logic:**
- Anti-IDOR: FindByIDAndUserID — verify outlet belongs to user

---

#### PUT /api/v1/outlets/:id

**Handler:** `outlet_handler.UpdateOutlet`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Outlet ID

**Request Body:** `OutletRequest`

**Response Sukses:**
- Status: `200 OK`
- Body: `OutletResponse`

**Response Error:**
- `400 Bad Request`: Validasi format / Format ID tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner / Bukan pemilik outlet
- `404 Not Found`: Outlet tidak ditemukan
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal mengupdate outlet

---

#### DELETE /api/v1/outlets/:id

**Handler:** `outlet_handler.DeleteOutlet`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Outlet ID

**Response Sukses:**
- Status: `200 OK`
- Body: Success message

**Response Error:**
- `400 Bad Request`: Format ID tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner / Bukan pemilik outlet
- `404 Not Found`: Outlet tidak ditemukan
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal menghapus outlet

**Business Logic:**
- Soft delete (GORM DeletedAt)
- Anti-IDOR: Verify ownership before delete

---

### 4.4 Owner — Services

#### POST /api/v1/services

**Handler:** `service_handler.CreateService`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Request Body:** `ServiceRequest`

**Response Sukses:**
- Status: `201 Created`
- Body: `ServiceResponse`

**Response Error:**
- `400 Bad Request`: Validasi format / Harga tidak valid / Gagal validasi outlet
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner
- `404 Not Found`: Outlet tidak ditemukan / Bukan pemilik outlet
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal membuat layanan

**Business Logic:**
1. Sanitize inputs
2. Parse price dari string ke decimal, validate > 0
3. Verify outlet belongs to user (Anti-IDOR Create)
4. Create service dengan UUID
5. Return service data

---

#### GET /api/v1/outlets/:id/services

**Handler:** `service_handler.GetAllByOutletID`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Outlet ID

**Response Sukses:**
- Status: `200 OK`
- Body: `[]ServiceResponse`

**Response Error:**
- `400 Bad Request`: Format Outlet ID tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner / Bukan pemilik outlet
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal mengambil daftar layanan

**Business Logic:**
- Anti-IDOR: JOIN outlets untuk verify ownership

---

#### PUT /api/v1/services/:id

**Handler:** `service_handler.UpdateService`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Service ID

**Request Body:** `ServiceRequest`

**Response Sukses:**
- Status: `200 OK`
- Body: `ServiceResponse`

**Response Error:**
- `400 Bad Request`: Validasi format / Harga tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner
- `404 Not Found`: Layanan tidak ditemukan / Outlet tujuan tidak ditemukan
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal mengupdate layanan

**Business Logic:**
1. Verify service belongs to user (Anti-IDOR)
2. If OutletID changed, verify new outlet also belongs to user
3. Update fields
4. Return updated service

---

#### DELETE /api/v1/services/:id

**Handler:** `service_handler.DeleteService`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Service ID

**Response Sukses:**
- Status: `200 OK`
- Body: Success message

**Response Error:**
- `400 Bad Request`: Format Service ID tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner
- `404 Not Found`: Layanan tidak ditemukan
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal menghapus layanan

---

### 4.5 Owner — Orders

#### GET /api/v1/outlets/:id/orders

**Handler:** `order_handler.GetAllByOutletID`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Outlet ID

**Query Params:**
- `page`: int (default: 1)
- `limit`: int (default: 10)

**Response Sukses:**
- Status: `200 OK`
- Body: `PaginatedResponse<OrderResponse>`

**Response Error:**
- `400 Bad Request`: Format Outlet ID / pagination tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner / Bukan pemilik outlet
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal mengambil data pesanan

**Business Logic:**
- Anti-IDOR: JOIN outlets untuk verify ownership

---

#### PATCH /api/v1/orders/:id/status

**Handler:** `order_handler.UpdateStatus`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Path Params:**
- `id`: UUID — Order ID

**Request Body:** `OrderStatusRequest`

**Response Sukses:**
- Status: `200 OK`
- Body: `OrderResponse`

**Response Error:**
- `400 Bad Request`: Validasi format / Transisi status tidak valid / Berat aktual tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner / Bukan pemilik outlet
- `404 Not Found`: Pesanan tidak ditemukan
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal memperbarui status

**Business Logic:**
1. Get order + verify ownership (Anti-IDOR)
2. Validate FSM transition
3. If status = "process":
   - Parse actual_qty for KG items
   - Calculate final_price = actual_qty × service_price
   - Calculate final_total_price
4. Create order log entry
5. Update order status + prices in transaction
6. Fire notification (background)

---

### 4.6 Owner — Reports

#### GET /api/v1/reports/omzet

**Handler:** `report_handler.GetOmzet`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Query Params:**
- `start_date`: YYYY-MM-DD (optional)
- `end_date`: YYYY-MM-DD (optional)
- `outlet_id`: UUID (optional)

**Response Sukses:**
- Status: `200 OK`
- Body: `OmzetResponse`

**Response Error:**
- `400 Bad Request`: Format query parameter tidak valid
- `401 Unauthorized`: Token tidak valid
- `403 Forbidden`: Bukan role owner
- `408 Request Timeout`: Proses terlalu lama
- `500 Internal Server Error`: Gagal memproses data omzet

---

#### GET /api/v1/reports/orders/summary

**Handler:** `report_handler.GetOrderSummary`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Query Params:**
- `start_date`: YYYY-MM-DD (optional)
- `end_date`: YYYY-MM-DD (optional)
- `outlet_id`: UUID (optional)

**Response Sukses:**
- Status: `200 OK`
- Body: `OrderStatusSummaryResponse`

---

#### GET /api/v1/reports/services/top

**Handler:** `report_handler.GetTopServices`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware, RoleMiddleware("owner")

**Auth Required:** Ya

**Role Required:** owner

**Query Params:**
- `start_date`: YYYY-MM-DD (optional)
- `end_date`: YYYY-MM-DD (optional)
- `outlet_id`: UUID (optional)

**Response Sukses:**
- Status: `200 OK`
- Body: `[]TopServiceResponse`

---

### 4.7 Notifications

#### GET /api/v1/notifications

**Handler:** `notification_handler.GetNotifications`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware

**Auth Required:** Ya

**Query Params:**
- `page`: int (default: 1)
- `limit`: int (default: 20)

**Response Sukses:**
- Status: `200 OK`
- Body: Custom format dengan `unread_count`

---

#### GET /api/v1/notifications/unread-count

**Handler:** `notification_handler.GetUnreadCount`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware

**Auth Required:** Ya

**Response Sukses:**
- Status: `200 OK`
- Body: `{"count": int64}`

---

#### PATCH /api/v1/notifications/:id/read

**Handler:** `notification_handler.MarkAsRead`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware

**Auth Required:** Ya

**Path Params:**
- `id`: UUID — Notification ID

**Response Sukses:**
- Status: `200 OK`

---

#### PATCH /api/v1/notifications/read-all

**Handler:** `notification_handler.MarkAllAsRead`

**Middleware:** PayloadLimit, RateLimiter, AuthMiddleware

**Auth Required:** Ya

**Response Sukses:**
- Status: `200 OK`

---

### 4.8 WebSocket

#### GET /api/v1/ws/connect

**Handler:** `notification_handler.Connect`

**Middleware:** AuthMiddleware

**Auth Required:** Ya

**Query Params:**
- `token`: JWT token (fallback jika tidak ada Authorization header)

**Upgrade:** WebSocket upgrade via gorilla/websocket.Upgrader

**Response:**
- Upgrade ke WebSocket connection
- Client registered ke hub

**Business Logic:**
1. Get user_id dan role dari context (setelah AuthMiddleware)
2. Upgrade connection ke WebSocket
3. Create Client struct
4. Register ke hub
5. Start ReadPump dan WritePump goroutines

---

## Section 5 — Middleware

### 5.1 AuthMiddleware

**File:** `internal/delivery/http/middleware.go`

**Fungsi:** Validasi JWT token dari Authorization header atau query param

**Diapply di:** Semua protected routes

**Logic:**
1. Get Authorization header
2. Fallback ke query param `token` (untuk WebSocket)
3. Validate format "Bearer <token>"
4. Validate token dengan `utils.ValidateToken()`
5. Set `user_id` dan `role` ke context

**Set ke context:**
- `c.Set("user_id", claims.UserID)` — User ID dari JWT
- `c.Set("role", claims.Role)` — Role user dari JWT

**On fail:**
- `401 Unauthorized`: "Token tidak ditemukan" / "Format token tidak valid" / "Token tidak valid atau sudah kadaluarsa"

---

### 5.2 RoleMiddleware

**File:** `internal/delivery/http/middleware.go`

**Fungsi:** Cek apakah user memiliki role yang diizinkan

**Diapply di:** Customer routes, Owner routes

**Logic:**
1. Get role dari context
2. Check role exists dan valid type
3. Iterate allowed roles, check match
4. If match, continue; else abort

**On fail:**
- `403 Forbidden`: "Akses ditolak" / "Anda tidak memiliki izin untuk mengakses resource ini"

---

### 5.3 CORSMiddleware

**File:** `internal/delivery/http/middleware.go`

**Fungsi:** Handle Cross-Origin Resource Sharing

**Diapply di:** Global (semua routes)

**Allowed Origins:**
```go
allowedOrigins := map[string]bool{
    "https://laundryin.vercel.app":      true,
    "https://www.laundryin.vercel.app":  true,
    "https://laundryinapps.vinjo.me":    true,
    "https://www.laundryinapps.vinjo.me": true,
    "http://localhost:3000":             true,
    "http://localhost:3001":             true,
}
```

**Logic:**
1. Get Origin header
2. Check if origin is allowed atau subdomain .vercel.app
3. In production (`GIN_MODE=release`): reject unknown origins dengan 403
4. In development: fallback to `*`
5. Set CORS headers
6. Handle OPTIONS preflight

**Headers Set:**
- `Access-Control-Allow-Origin`: Dynamic atau `*`
- `Access-Control-Allow-Methods`: GET, POST, PUT, PATCH, DELETE, OPTIONS
- `Access-Control-Allow-Headers`: Origin, Content-Type, Accept, Authorization
- `Access-Control-Allow-Credentials`: true

**⚠️ Inconsistency:** Di production mode, unknown origin di-reject dengan 403, tapi ada fallback ke `*` untuk local dev yang bisa bypass jika GIN_MODE tidak diset.

---

### 5.4 RateLimiter

**File:** `internal/delivery/http/middleware.go`

**Fungsi:** IP-based rate limiting dengan token bucket algorithm

**Diapply di:** Semua API routes (v1 group)

**Logic:**
1. Get client IP
2. Get or create rate limiter for IP
3. Check if token available
4. If not available, return 429
5. Background cleaner: remove inactive IPs setelah 3 menit

**Rate Limit:**
- Rate: 1.66 tokens/second (~100 requests/minute)
- Burst: 100 tokens

**On fail:**
- `429 Too Many Requests`: "Terlalu banyak request, silakan coba lagi nanti"

**🐛 Potential Bug:** Background cleaner goroutine tidak ada graceful shutdown — bisa leak jika server restart.

---

### 5.5 PayloadLimit

**File:** `internal/delivery/http/middleware.go`

**Fungsi:** Limit ukuran request body untuk mencegah DoS

**Diapply di:** API v1 group

**Limit:** 1MB (1024 * 1024 bytes)

**Logic:**
```go
c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
```

---

### 5.6 SecurityHeaders

**File:** `internal/delivery/http/middleware.go`

**Fungsi:** Add security headers ke setiap response

**Diapply di:** Global

**Headers Set:**
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Cache-Control: no-store`
- `Content-Security-Policy: default-src 'none'`
- `Referrer-Policy: no-referrer`

---

## Section 6 — WebSocket

### 6.1 Hub Structure

**File:** `internal/websocket/hub.go`

**Fields:**
```go
type Hub struct {
    clients    map[string]*Client  // Map userID -> Client
    register   chan *Client        // Channel untuk register client baru
    unregister chan *Client        // Channel untuk unregister client
    broadcast  chan *Message       // Channel untuk broadcast message
    mu         sync.RWMutex        // Mutex untuk thread safety
}
```

**Functions:**
- `NewHub()`: Create new hub instance
- `Run()`: Main loop untuk handle register/unregister/broadcast
- `SendToUser(userID, msg)`: Send message ke specific user
- `BroadcastToOwners(msg)`: Broadcast ke semua owner clients

---

### 6.2 Client Structure

**Fields:**
```go
type Client struct {
    UserID string               // User ID dari JWT
    Role   string               // Role user
    Conn   *websocket.Conn      // WebSocket connection
    Send   chan []byte          // Channel untuk outgoing messages (buffer 256)
}
```

---

### 6.3 Message Format

```go
type Message struct {
    ID        string      `json:"id,omitempty"`
    Type      string      `json:"type"`
    Title     string      `json:"title"`
    Body      string      `json:"body"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
}
```

**Contoh JSON:**
```json
{
  "id": "notif-uuid",
  "type": "new_order",
  "title": "Pesanan Baru Masuk",
  "body": "John Doe memesan Cuci Kering di Laundry Express",
  "data": {
    "order_id": "order-uuid",
    "outlet_id": "outlet-uuid",
    "customer_name": "John Doe",
    "total_price": "50000.00"
  },
  "timestamp": "2025-03-16T10:00:00Z"
}
```

---

### 6.4 Client Lifecycle

**1. Connect:**
```
Client → GET /api/v1/ws/connect?token=<jwt> → AuthMiddleware → notification_handler.Connect
→ upgrader.Upgrade() → Create Client → hub.Register() <- client
→ go client.WritePump()
→ go client.ReadPump()
```

**2. Register:**
```go
case client := <-h.register:
    h.mu.Lock()
    if old, exists := h.clients[client.UserID]; exists {
        old.Conn.Close()  // Close old connection
        close(old.Send)
    }
    h.clients[client.UserID] = client
    h.mu.Unlock()
```

**3. WritePump (Hub → Client):**
- Ticker: Ping setiap 30 detik
- Write deadline: 10 detik
- Handle messages dari `client.Send` channel
- Send PingMessage untuk keepalive

**4. ReadPump (Client → Hub):**
- Read limit: 512 bytes
- Read deadline: 60 detik (reset on pong)
- Pong handler: Reset read deadline
- On error/close: `hub.unregister <- client`

**5. Disconnect:**
```go
case client := <-h.unregister:
    h.mu.Lock()
    if _, ok := h.clients[client.UserID]; ok {
        delete(h.clients, client.UserID)
        close(client.Send)
    }
    h.mu.Unlock()
```

---

### 6.5 Notification Triggers

| Event | Dipanggil Di | Penerima | Type | Title | Body Template | Data |
|-------|-------------|----------|------|-------|---------------|------|
| Order Created | `order_usecase.Create` (background) | Owner | `new_order` | "Pesanan Baru Masuk" | "{customer_name} memesan {services} di {outlet_name}" | `order_id`, `outlet_id`, `customer_name`, `total_price` |
| Status → Process | `order_usecase.UpdateStatus` (background) | Customer | `order_status` | "Pesananmu Sedang Diproses" | "Outlet {outlet_name} mulai memproses pesanan #{order_id_short}" | `order_id`, `new_status` |
| Status → Process (with final price) | `order_usecase.UpdateStatus` (background) | Customer | `price_updated` | "Harga Final Pesananmu Sudah Diketahui" | "Total pembayaran pesanan #{order_id_short} adalah Rp {final_total_price}" | `order_id`, `estimated_price`, `final_price` |
| Status → Completed | `order_usecase.UpdateStatus` (background) | Customer | `order_status` | "Pesananmu Siap Diambil! 🎉" | "Cucian kamu sudah selesai, silakan ambil di {outlet_name}" | `order_id`, `new_status` |
| Status → Cancelled | `order_usecase.UpdateStatus` (background) | Customer | `order_cancelled` | "Pesananmu Dibatalkan" | "Pesanan #{order_id_short} di {outlet_name} telah dibatalkan" | `order_id` |

---

### 6.6 Ping/Pong

**Ping Interval:** 30 detik (WritePump ticker)

**Pong Timeout:** 60 detik (Read deadline)

**Implementation:**
```go
// WritePump
ticker := time.NewTicker(30 * time.Second)
c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
c.Conn.WriteMessage(websocket.PingMessage, nil)

// ReadPump
c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
c.Conn.SetPongHandler(func(string) error {
    c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    return nil
})
```

---

## Section 7 — Business Logic Penting

### 7.1 Order FSM (Finite State Machine)

**Valid Transitions:**

```
┌─────────┐
│ pending │
└────┬────┘
     │
     ├─────────────┐
     │             │
     ▼             ▼
┌─────────┐   ┌───────────┐
│ process │   │ cancelled │
└────┬────┘   └───────────┘
     │
     ▼
┌───────────┐
│ completed │
└─────┬─────┘
      │
      ▼
┌───────────┐
│ picked_up │
└───────────┘
```

**Implementation:**
```go
func isValidTransition(current, next string) bool {
    switch current {
    case "pending":
        return next == "process" || next == "cancelled"
    case "process":
        return next == "completed"
    case "completed":
        return next == "picked_up"
    default:
        return false
    }
}
```

**Validasi di usecase:**
```go
if !isValidTransition(order.Status, req.Status) {
    return nil, ErrStateInvalid
}
```

---

### 7.2 Actual Qty & Final Price Calculation

**Kapan Dihitung:** Saat owner transition order dari `pending` → `process`

**Logic:**
```go
if req.Status == "process" {
    // Parse actual_qty dari request
    reqItemsMap := make(map[string]decimal.Decimal)
    for _, ri := range req.Items {
        qty, _ := decimal.NewFromString(ri.ActualQty)
        reqItemsMap[ri.ID] = qty.Round(2)
    }

    var finalTotalPrice decimal.Decimal

    for i, item := range order.Items {
        if item.Unit == "KG" {
            if actualQty, exists := reqItemsMap[item.ID]; exists {
                order.Items[i].ActualQty = &actualQty
                finalPrice := actualQty.Mul(item.ServicePrice).Round(2)
                order.Items[i].FinalPrice = &finalPrice
                finalTotalPrice = finalTotalPrice.Add(finalPrice)
            } else {
                return nil, errors.New("berat aktual wajib diisi untuk layanan per KG")
            }
        } else {
            // PCS: finalPrice = Subtotal (tidak ada timbang ulang)
            finalPrice := item.Subtotal
            order.Items[i].FinalPrice = &finalPrice
            finalTotalPrice = finalTotalPrice.Add(finalPrice)
        }
    }

    order.FinalTotalPrice = &finalTotalPrice
}
```

**Catatan:**
- Hanya item dengan `unit = "KG"` yang butuh `actual_qty`
- Item `PCS` langsung pakai `subtotal` sebagai `final_price`
- `final_total_price` = sum dari semua `final_price`

---

### 7.3 Password Validation

**Rules:**
- Min 8 karakter, max 64
- Harus mengandung huruf besar (A-Z)
- Harus mengandung huruf kecil (a-z)
- Harus mengandung angka (0-9)

**Implementation:**
```go
hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(req.Password)
hasLower := regexp.MustCompile(`[a-z]`).MatchString(req.Password)
hasDigit := regexp.MustCompile(`[0-9]`).MatchString(req.Password)

if !hasUpper || !hasLower || !hasDigit {
    return nil, ErrWeakPassword
}
```

---

### 7.4 Phone Format Validation

**Format Diterima:** E.164 strict

**Regex:** `^\+[1-9]\d{6,14}$`

**Contoh Valid:**
- `+6281234567890` ✅
- `+12025551234` ✅

**Contoh Invalid:**
- `081234567890` ❌ (tidak ada +)
- `6281234567890` ❌ (tidak ada +)
- `+62 812-345-6789` ❌ (ada spasi/dash)

**Custom Validator:**
```go
func validateE164Strict(fl validator.FieldLevel) bool {
    return e164Regex.MatchString(fl.Field().String())
}
```

**Sanitization:**
```go
req.Phone = utils.Sanitize(req.Phone)  // Trim whitespace, remove null bytes
```

---

### 7.5 JWT Claims

**Structure:**
```go
type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
```

**Expiry Time:** 24 jam

**Issuer:** `laundryin-api`

**Signing Method:** HS256

**Generate Token:**
```go
claims := Claims{
    UserID: userID,
    Role:   role,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        Issuer:    "laundryin-api",
    },
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString([]byte(secret))
```

**Validation:**
```go
token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, errors.New("unexpected signing method")
    }
    return []byte(secret), nil
})
```

---

### 7.6 Anti-IDOR (Insecure Direct Object Reference)

**Pattern:** Setiap akses ke resource (outlet, service, order) harus verify ownership via JOIN.

**Implementation Examples:**

**Outlet (FindByIDAndUserID):**
```go
err := r.db.WithContext(ctx).
    Where("id = ? AND user_id = ?", outletID, userID).
    First(&outlet).Error
```

**Service (FindByIDAndOwner):**
```go
err := r.db.WithContext(ctx).
    Joins("JOIN outlets ON outlets.id = services.outlet_id").
    Where("services.id = ? AND outlets.user_id = ?", serviceID, userID).
    First(&service).Error
```

**Order (FindByIDAndOwner):**
```go
err := r.db.WithContext(ctx).
    Joins("JOIN outlets ON outlets.id = orders.outlet_id").
    Preload("Items").
    Preload("User").
    Where("orders.id = ? AND outlets.user_id = ?", orderID, userID).
    First(&order).Error
```

**Create Service (Anti-IDOR Create):**
```go
// Step 1: Verify the outlet belongs to the current user
_, err = u.outletRepo.FindByIDAndUserID(ctx, req.OutletID, userID)
if err != nil {
    return nil, ErrOutletNotFound
}
// Step 2: Create service
```

**Create Order (Deep Anti-IDOR & Zero-Trust Pricing):**
```go
for _, itemReq := range req.Items {
    // Verify service exists AND belongs to the outlet
    s, err := u.serviceRepo.FindByIDAndOutletID(ctx, itemReq.ServiceID, req.OutletID)
    if err != nil {
        return nil, errors.New("satu atau lebih layanan tidak ditemukan atau bukan milik outlet ini")
    }
    // Use service.Price from DB, NOT from frontend
}
```

---

## Section 8 — Repository Layer

### 8.1 UserRepository

**File:** `internal/repository/user_repository.go`

**Interface:**
```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByPhone(ctx context.Context, phone string) (*models.User, error)
    FindByID(ctx context.Context, id string) (*models.User, error)
}
```

**Methods:**
- `Create(ctx, user)`: Insert user baru
- `FindByPhone(ctx, phone)`: Find user by phone number
- `FindByID(ctx, id)`: Find user by ID

---

### 8.2 OutletRepository

**File:** `internal/repository/outlet_repository.go`

**Interface:**
```go
type OutletRepository interface {
    Create(ctx context.Context, outlet *models.Outlet) error
    FindAll(ctx context.Context, limit, offset int) ([]models.Outlet, int64, error)
    FindAllByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Outlet, int64, error)
    FindByID(ctx context.Context, outletID string) (*models.Outlet, error)
    FindByIDAndUserID(ctx context.Context, outletID, userID string) (*models.Outlet, error)
    Update(ctx context.Context, outlet *models.Outlet) error
    Delete(ctx context.Context, outletID, userID string) error
}
```

**Methods:**
- `Create(ctx, outlet)`: Insert outlet baru
- `FindAll(ctx, limit, offset)`: Get all outlets (public) dengan pagination
- `FindAllByUserID(ctx, userID, limit, offset)`: Get outlets by owner dengan pagination
- `FindByID(ctx, outletID)`: Find outlet by ID (public)
- `FindByIDAndUserID(ctx, outletID, userID)`: Find outlet by ID + verify ownership (Anti-IDOR)
- `Update(ctx, outlet)`: Update outlet
- `Delete(ctx, outletID, userID)`: Soft delete outlet dengan ownership check

---

### 8.3 ServiceRepository

**File:** `internal/repository/service_repository.go`

**Interface:**
```go
type ServiceRepository interface {
    Create(ctx context.Context, service *models.Service) error
    FindAllByOutletID(ctx context.Context, outletID, userID string) ([]models.Service, error)
    FindAllByOutletIDPublic(ctx context.Context, outletID string) ([]models.Service, error)
    FindByIDAndOwner(ctx context.Context, serviceID, userID string) (*models.Service, error)
    FindByIDAndOutletID(ctx context.Context, serviceID, outletID string) (*models.Service, error)
    Update(ctx context.Context, service *models.Service) error
    Delete(ctx context.Context, serviceID, userID string) error
}
```

**Methods:**
- `Create(ctx, service)`: Insert service baru
- `FindAllByOutletID(ctx, outletID, userID)`: Get services by outlet + verify ownership (Anti-IDOR)
  - Query: JOIN outlets ON outlets.id = services.outlet_id WHERE services.outlet_id = ? AND outlets.user_id = ?
- `FindAllByOutletIDPublic(ctx, outletID)`: Get services by outlet (public, no auth)
- `FindByIDAndOwner(ctx, serviceID, userID)`: Find service by ID + verify ownership (Anti-IDOR)
  - Query: JOIN outlets ON outlets.id = services.outlet_id WHERE services.id = ? AND outlets.user_id = ?
- `FindByIDAndOutletID(ctx, serviceID, outletID)`: Find service by ID + verify belongs to outlet
- `Update(ctx, service)`: Update service (GORM Save)
- `Delete(ctx, serviceID, userID)`: Soft delete dengan Anti-IDOR check

---

### 8.4 OrderRepository

**File:** `internal/repository/order_repository.go`

**Interface:**
```go
type OrderRepository interface {
    CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error
    FindAllByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Order, int64, error)
    FindAllByOutletIDAndOwner(ctx context.Context, outletID, userID string, limit, offset int) ([]models.Order, int64, error)
    FindByIDAndOwner(ctx context.Context, orderID, userID string) (*models.Order, error)
    UpdateOrderStatus(ctx context.Context, order *models.Order, logEntry *models.OrderLog) error
}
```

**Methods:**
- `CreateOrderWithItems(ctx, order, items)`: Insert order + items atomically dalam transaction
  - Logic: Transaction → Create order → Loop create items → Commit
- `FindAllByUserID(ctx, userID, limit, offset)`: Get orders by customer dengan pagination
  - Preload: Items, User
- `FindAllByOutletIDAndOwner(ctx, outletID, userID, limit, offset)`: Get orders by outlet + verify ownership (Anti-IDOR)
  - Query: JOIN outlets ON outlets.id = orders.outlet_id WHERE orders.outlet_id = ? AND outlets.user_id = ?
- `FindByIDAndOwner(ctx, orderID, userID)`: Find order by ID + verify ownership (Anti-IDOR)
  - Query: JOIN outlets ON outlets.id = orders.outlet_id WHERE orders.id = ? AND outlets.user_id = ?
- `UpdateOrderStatus(ctx, order, logEntry)`: Update order status + final prices + create log atomically
  - Transaction: Update order → Update items (actual_qty, final_price) → Create log

---

### 8.5 ReportRepository

**File:** `internal/repository/report_repository.go`

**Interface:**
```go
type ReportRepository interface {
    GetTotalOmzet(ctx context.Context, userID string, req dto.ReportQuery) (decimal.Decimal, error)
    GetOrderStatusSummary(ctx context.Context, userID string, req dto.ReportQuery) (map[string]int64, error)
    GetTopServices(ctx context.Context, userID string, req dto.ReportQuery) ([]TopServiceRow, error)
}
```

**Methods:**
- `GetTotalOmzet(ctx, userID, req)`: Get total revenue untuk owner dalam periode
  - Query: JOIN outlets, SUM(orders.total_price) WHERE outlets.user_id = ? AND status != 'cancelled'
- `GetOrderStatusSummary(ctx, userID, req)`: Get count per status order
  - Query: JOIN outlets, GROUP BY orders.status COUNT(orders.id)
- `GetTopServices(ctx, userID, req)`: Get top 5 services by revenue
  - Query: JOIN orders, JOIN outlets, GROUP BY service_name, outlets.name ORDER BY total_revenue DESC LIMIT 5

---

### 8.6 NotificationRepository

**File:** `internal/repository/notification_repository.go`

**Interface:**
```go
type NotificationRepository interface {
    Create(ctx context.Context, notif *models.Notification) error
    GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error)
    GetUnreadCount(ctx context.Context, userID string) (int64, error)
    GetTotalCount(ctx context.Context, userID string) (int64, error)
    MarkAsRead(ctx context.Context, notifID string, userID string) error
    MarkAllAsRead(ctx context.Context, userID string) error
}
```

**Methods:**
- `Create(ctx, notif)`: Insert notification baru
- `GetByUserID(ctx, userID, limit, offset)`: Get notifications by user dengan pagination
- `GetUnreadCount(ctx, userID)`: Get count unread notifications
- `GetTotalCount(ctx, userID)`: Get total notifications
- `MarkAsRead(ctx, notifID, userID)`: Mark single notification as read
- `MarkAllAsRead(ctx, userID)`: Mark all notifications as read

---

## Section 9 — Error Handling Patterns

### 9.1 Custom Errors

**auth_usecase.go:**
```go
var ErrDuplicatePhone = errors.New("nomor HP sudah terdaftar")
var ErrInvalidCredentials = errors.New("nomor HP atau password salah")
var ErrWeakPassword = errors.New("Password harus mengandung huruf besar, huruf kecil, dan angka")
```

**outlet_usecase.go:**
```go
var ErrOutletNotFound = errors.New("outlet tidak ditemukan")
```

**service_usecase.go:**
```go
var ErrServiceNotFound = errors.New("layanan tidak ditemukan atau akses ditolak")
```

**order_usecase.go:**
```go
var ErrOrderNotFound = errors.New("pesanan tidak ditemukan atau akses ditolak")
var ErrStateInvalid = errors.New("transisi status pesanan tidak valid")
```

---

### 9.2 HTTP Status Code Mapping

| Status Code | Kondisi |
|-------------|---------|
| `200 OK` | Success (GET, PUT, PATCH) |
| `201 Created` | Success create (POST) |
| `400 Bad Request` | Validasi gagal / Format tidak valid / Transisi FSM invalid |
| `401 Unauthorized` | Token tidak ada / Token tidak valid / Token expired / Credentials salah |
| `403 Forbidden` | Role tidak match / Akses resource orang lain (Anti-IDOR) |
| `404 Not Found` | Resource tidak ditemukan |
| `408 Request Timeout` | Context deadline exceeded |
| `409 Conflict` | Duplicate phone number |
| `422 Unprocessable Entity` | Binding/validation error dari Gin |
| `429 Too Many Requests` | Rate limit exceeded |
| `500 Internal Server Error` | Error internal server |

---

### 9.3 Error Response Format

**Success:**
```json
{
  "status": "success",
  "message": "Berhasil",
  "data": { ... }
}
```

**Error:**
```json
{
  "status": "error",
  "message": "Pesan error",
  "errors": "Detail error (optional)"
}
```

**Implementation:**
```go
func ErrorResponse(c *gin.Context, statusCode int, message string, errs interface{}) {
    if statusCode >= 500 {
        fmt.Printf("🔴 SERVER ERROR (%d): %s | Details: %+v\n", statusCode, message, errs)
    }

    response := gin.H{
        "status":  "error",
        "message": message,
    }
    if errs != nil {
        response["errors"] = errs
    }
    c.JSON(statusCode, response)
}
```

---

## Section 10 — Known Issues & Inconsistencies

### 10.1 Security Issues

| Issue | File | Severity | Deskripsi | Rekomendasi |
|-------|------|----------|-----------|-------------|
| CORS wildcard fallback | `middleware.go` | 🐛 Medium | Di development, fallback ke `*` bisa bypass origin check jika GIN_MODE tidak diset | Hardcode strict origin check untuk semua environment, atau gunakan env variable untuk allowed origins list |
| CheckOrigin WebSocket | `notification_handler.go` | ⚠️ Low | Whitelist hardcoded, tapi ada fallback `return true` untuk empty origin | Tambah validasi origin yang lebih strict, reject empty origin di production |
| JWT_SECRET validation | `database/postgres.go` | ⚠️ Medium | Hanya warning jika JWT_SECRET tidak diset, tapi app tetap jalan | Fail fast di startup jika JWT_SECRET tidak diset di production |

---

### 10.2 Missing Validation

| Issue | File | Deskripsi | Rekomendasi |
|-------|------|-----------|-------------|
| No sanitize on OrderRequest fields | `order_handler.go` | Hanya sanitize OutletID dan ServiceID, tapi tidak semua string fields | Tambah sanitize untuk semua user input |
| No max limit enforcement on some endpoints | Beberapa handler | Default limit 10, max 100, tapi tidak semua endpoint enforce max | Tambah enforcement `if limit > 100 { limit = 100 }` di semua endpoint |

---

### 10.3 Inconsistent Patterns

| Issue | File 1 | File 2 | Rekomendasi |
|-------|--------|--------|-------------|
| Error message exposure | Beberapa handler expose `err.Error()` ke client | Beberapa handler hide detail error | Standardisasi: hanya expose error message yang aman, log detail error di server |
| Context timeout | Auth: 15 detik | Handler lain: 5 detik | Standardisasi timeout berdasarkan operasi (auth boleh lebih lama, CRUD 5 detik cukup) |

---

### 10.4 TODO/FIXME Comments

**Tidak ada TODO/FIXME comments yang ditemukan di kode.**

---

### 10.5 Potential Bugs

| Issue | File | Severity | Deskripsi | Rekomendasi |
|-------|------|----------|-----------|-------------|
| Rate limiter goroutine leak | `middleware.go` | 🐛 Medium | Background cleaner goroutine tidak ada graceful shutdown | Tie context to server shutdown, call cancel() saat server stop |
| WebSocket broadcast blocking | `hub.go` | 🐛 Low | Broadcast loop bisa block jika client.Send channel full | Drop message atau disconnect slow client |
| Notification goroutine panic | `order_usecase.go` | ⚠️ Low | `go u.notifUsecase.NotifyOrderCreated(...)` tidak ada recover | Wrap dengan `defer recover()` untuk prevent panic crash |
| Decimal conversion error ignored | Beberapa usecase | ⚠️ Medium | `decimal.NewFromString()` error kadang di-ignore dengan `_` | Handle error dengan return 400 Bad Request |

---

### 10.6 Missing Error Handling

| Issue | File | Deskripsi | Rekomendasi |
|-------|------|-----------|-------------|
| Notification creation error logged but not returned | `notification_usecase.go` | Error hanya di-print, notifikasi bisa fail silent | Tambah retry mechanism atau queue untuk reliability |
| Outlet/User lookup fail silent | `notification_usecase.go` | `outlet, _ := u.outletRepo.FindByID(...)` ignore error | Log error untuk debugging, tapi continue karena notification adalah feature tambahan |
| JSON marshal error ignored | `notification_usecase.go` | `notifData, _ := json.Marshal(...)` | JSON marshal seharusnya tidak fail untuk valid input, tapi bisa ditambahkan fallback |

---

## Appendix A — File Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── database/
│   │   └── postgres.go          # Database connection
│   ├── delivery/
│   │   └── http/
│   │       ├── auth_handler.go
│   │       ├── middleware.go
│   │       ├── notification_handler.go
│   │       ├── order_handler.go
│   │       ├── outlet_handler.go
│   │       ├── report_handler.go
│   │       └── service_handler.go
│   ├── dto/
│   │   ├── dto.go               # Main DTOs
│   │   └── report_dto.go        # Report DTOs
│   ├── repository/
│   │   ├── models/
│   │   │   └── models.go        # GORM models
│   │   ├── notification_repository.go
│   │   ├── order_repository.go
│   │   ├── outlet_repository.go
│   │   ├── report_repository.go
│   │   ├── service_repository.go
│   │   └── user_repository.go
│   ├── usecase/
│   │   ├── auth_usecase.go
│   │   ├── notification_usecase.go
│   │   ├── order_usecase.go
│   │   ├── outlet_usecase.go
│   │   ├── report_usecase.go
│   │   └── service_usecase.go
│   └── websocket/
│       └── hub.go               # WebSocket hub
├── pkg/
│   └── utils/
│       ├── jwt.go               # JWT utils
│       ├── password.go          # Password hashing
│       ├── response.go          # Response helpers
│       ├── string.go            # String utils (sanitize)
│       └── validator.go         # Custom validators
├── go.mod
├── go.sum
└── Dockerfile
```

---

## Appendix B — Dependency Graph

```
┌─────────────────────────────────────────────────────────────────┐
│                         MAIN.GO                                  │
│  (Init: DB → Migrate → WebSocket → Repos → Usecases → Handlers) │
└─────────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│    Database     │  │   WebSocket     │  │     Utils       │
│   (PostgreSQL)  │  │     (Hub)       │  │ (JWT, Password) │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                    │                    │
         ▼                    ▼                    │
┌─────────────────┐          │                    │
│   Repository    │          │                    │
│     Layer       │          │                    │
└────────┬────────┘          │                    │
         │                   │                    │
         ▼                   │                    │
┌─────────────────┐          │                    │
│    Usecase      │◀─────────┘                    │
│     Layer       │                               │
└────────┬────────┘                               │
         │                                        │
         ▼                                        │
┌─────────────────┐                               │
│     Handler     │◀──────────────────────────────┘
│     Layer       │         (Middleware)
└─────────────────┘
```

---

## Appendix C — Authentication Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                      AUTHENTICATION FLOW                         │
└─────────────────────────────────────────────────────────────────┘

1. REGISTER
   Client → POST /api/v1/auth/register → Handler → Usecase
   → Repo.Create() → Generate JWT → Return token + user

2. LOGIN
   Client → POST /api/v1/auth/login → Handler → Usecase
   → Repo.FindByPhone() → CheckPassword() → Generate JWT → Return token + user

3. AUTHENTICATED REQUEST
   Client → GET /api/v1/orders (Authorization: Bearer <token>)
   → AuthMiddleware: ValidateToken() → c.Set("user_id", claims.UserID)
   → Handler: userID := c.MustGet("user_id").(string)
   → Usecase: Check ownership (Anti-IDOR)
   → Repo: Query with user_id filter

4. WebSocket CONNECTION
   Client → GET /api/v1/ws/connect?token=<jwt>
   → AuthMiddleware: ValidateToken() → c.Set("user_id", ...)
   → Handler: Upgrade to WebSocket → Create Client → hub.Register()
   → hub.SendToUser(userID, message) → client.Send channel → WritePump → Conn.WriteMessage()
```

---

**Dokumentasi ini selesai dibuat pada:** March 16, 2026

**Versi:** 1.0.0

**Catatan:** Dokumentasi ini adalah living document. Update setiap ada perubahan signifikan di kode backend.
