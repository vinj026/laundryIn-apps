# Dokumentasi Lengkap LaundryIn Platform

> **Versi:** 2.0 (Terakhir diupdate: Maret 2026)  
> **Status:** Production Ready

---

## Daftar Isi

1. [Gambaran Umum](#1-gambaran-umum)
2. [Peta Halaman](#2-peta-halaman)
3. [Flow Customer](#3-flow-customer)
4. [Flow Owner](#4-flow-owner)
5. [Flow Notifikasi (WebSocket)](#5-flow-notifikasi-websocket)
6. [Auth & Guard Rules](#6-auth--guard-rules)
7. [State Management](#7-state-management)
8. [Database Schema](#8-database-schema)
9. [API Routes Summary](#9-api-routes-summary)
10. [Catatan Teknis Penting](#10-catatan-teknis-penting)
11. [File Structure](#11-file-structure)
12. [Environment Variables](#12-environment-variables)

---

## 1. Gambaran Umum

### Tech Stack

#### Frontend
| Teknologi | Versi | Deskripsi |
|-----------|-------|-----------|
| **Framework** | Nuxt 3.x | Vue.js full-stack framework |
| **Language** | TypeScript | Type-safe JavaScript |
| **State Management** | Pinia | Vue 3 store |
| **Styling** | TailwindCSS + Custom CSS | Utility-first CSS |
| **UI Components** | Material Symbols | Google icon font |
| **HTTP Client** | Nuxt $fetch | Built-in fetch wrapper |
| **WebSocket** | Native WebSocket API | Real-time notifications |

#### Backend
| Teknologi | Versi | Deskripsi |
|-----------|-------|-----------|
| **Language** | Go 1.21+ | High-performance backend |
| **Framework** | Gin | HTTP web framework |
| **ORM** | GORM | Go ORM library |
| **Database** | PostgreSQL | Relational database |
| **JWT** | golang-jwt/jwt/v5 | Token authentication |
| **Validation** | go-playground/validator | Request validation |
| **Decimal** | shopspring/decimal | Financial precision |
| **WebSocket** | gorilla/websocket | Real-time communication |

### Arsitektur Aplikasi

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           FRONTEND (Nuxt 3)                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │   Pages     │  │  Components │  │    Stores   │  │  Composables│   │
│  │   (Views)   │  │    (UI)     │  │  (Pinia)    │  │  (Utilities)│   │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │
│         │                │                │                │            │
│         └────────────────┴────────────────┴────────────────┘            │
│                                  │                                      │
│                          Nuxt Proxy (/api)                              │
└──────────────────────────────────┼──────────────────────────────────────┘
                                   │ HTTP/WebSocket
                                   ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           BACKEND (Go + Gin)                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │   Handler   │──▶│  Usecase    │──▶│ Repository  │──▶│  Database   │   │
│  │   (HTTP)    │  │ (Business)  │  │   (Data)    │  │  (Postgres) │   │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │
│         │                │                                              │
│         ▼                ▼                                              │
│  ┌─────────────┐  ┌─────────────┐                                      │
│  │ Middleware  │  │  WebSocket  │                                      │
│  │ (Auth/Rate) │  │    Hub      │                                      │
│  └─────────────┘  └─────────────┘                                      │
└─────────────────────────────────────────────────────────────────────────┘
```

### Cara Kerja Autentikasi JWT

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Client     │     │   Server     │     │   Database   │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       │ POST /auth/login   │                    │
       │ {phone, password}  │                    │
       │───────────────────▶│                    │
       │                    │  Query user by phone│
       │                    │────────────────────▶│
       │                    │◀────────────────────│
       │                    │  Verify password    │
       │                    │  Generate JWT       │
       │                    │                    │
       │  {token, user}     │                    │
       │◀───────────────────│                    │
       │                    │                    │
       │  Store in localStorage & Pinia          │
       │                    │                    │
       │  Subsequent requests with Authorization │
       │  header: "Bearer <token>"               │
       │───────────────────▶│                    │
       │                    │ Validate JWT       │
       │                    │ Set user_id, role  │
       │                    │ in context         │
       │  {protected data}  │                    │
       │◀───────────────────│                    │
       │                    │                    │
```

**Token Structure:**
```json
{
  "user_id": "uuid-string",
  "role": "customer|owner",
  "exp": 1710489600,
  "iat": 1710403200,
  "iss": "laundryin-api"
}
```

**Token Expiry:** 24 jam

### Cara Kerja Nuxt Proxy

```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  routeRules: {
    '/api/**': { proxy: 'http://localhost:8080/api/v1/**' }
  }
})
```

**Flow:**
1. Frontend request ke `/api/orders`
2. Nuxt proxy rewrite ke `http://localhost:8080/api/v1/orders`
3. Backend handle request
4. Response dikembalikan ke frontend

**Keuntungan:**
- Tidak perlu CORS configuration di development
- Single origin untuk production
- Transparent API routing

### Format Response API

#### Response Sukses
```json
{
  "status": "success",
  "message": "Pesanan berhasil dibuat",
  "data": { ... }
}
```

#### Response Error
```json
{
  "status": "error",
  "message": "Nomor HP atau password salah",
  "errors": null
}
```

#### Response Paginated
```json
{
  "status": "success",
  "message": "Daftar outlet berhasil diambil",
  "data": {
    "data": [...],
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

#### Response Notifikasi (dengan unread_count)
```json
{
  "data": [...],
  "total": 20,
  "unread_count": 5
}
```

---

## 2. Peta Halaman

### Customer Routes

| Route | File | Layout | Auth Required | Role | Deskripsi |
|-------|------|--------|---------------|------|-----------|
| `/` | `pages/index.vue` | None | ❌ | Public | Landing page |
| `/customer` | `pages/customer/index.vue` | customer | ❌ | Public/Customer | Explore outlets |
| `/customer/login` | `pages/customer/login.vue` | None | ❌ | Public | Login customer |
| `/customer/register` | `pages/customer/register.vue` | None | ❌ | Public | Register customer |
| `/customer/outlet/[id]` | `pages/customer/outlet/[id].vue` | customer | ❌ | Public | Order booking |
| `/customer/orders` | `pages/customer/orders/index.vue` | customer | ✅ | Customer | Order history |
| `/customer/orders/[id]` | `pages/customer/orders/[id].vue` | customer | ✅ | Customer | Order tracking |
| `/customer/profile` | `pages/customer/profile.vue` | customer | ✅ | Customer | User profile |
| `/customer/tracking` | `pages/customer/tracking.vue` | customer | ✅ | Customer | Live tracking (demo) |

### Owner Routes

| Route | File | Layout | Auth Required | Role | Deskripsi |
|-------|------|--------|---------------|------|-----------|
| `/owner` | `pages/owner/index.vue` | owner | ✅ | Owner | Analytics dashboard |
| `/owner/login` | `pages/owner/login.vue` | None | ❌ | Public | Login owner |
| `/owner/outlets` | `pages/owner/outlets.vue` | owner | ✅ | Owner | Outlet CRUD |
| `/owner/services` | `pages/owner/services.vue` | owner | ✅ | Owner | Service CRUD |
| `/owner/orders` | `pages/owner/orders.vue` | owner | ✅ | Owner | Order pipeline |

---

## 3. Flow Customer

### 3.1 Landing Page (`/`)

**File:** `pages/index.vue`

**UI:**
- Logo LaundryIn
- Hero section dengan judul "Smart Laundry Platform"
- 2 CTA buttons: "Owner Dashboard" → `/owner/login`, "Find a Laundry" → `/customer`

**Redirect Logic:**
```
User sudah login (customer) → /customer
User sudah login (owner) → /owner
```

---

### 3.2 Explore Outlets (`/customer`)

**File:** `pages/customer/index.vue`

**UI:**
- Header "Find a Laundry"
- Search bar (filter by name/address)
- Grid outlet cards (responsive: 1 col mobile, 2 col tablet, 3 col desktop)
- Outlet card: nama, alamat, rating, status "Open"

**API Call:**
```
GET /api/public/outlets
Headers: None (public endpoint)
Response Success:
{
  "status": "success",
  "message": "Daftar outlet public berhasil diambil",
  "data": {
    "data": [
      {
        "id": "uuid",
        "name": "LaundryIn Premium Branch",
        "address": "Jl. Kertajaya No. 123",
        "phone": "+628123456789",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "page": 1,
    "limit": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

**User Actions:**
1. **Search:** Filter outlet by name/address (client-side)
2. **Click outlet card:** Navigate to `/customer/outlet/[id]`
3. **Pull to refresh:** Refresh outlet list

---

### 3.3 Login Customer (`/customer/login`)

**File:** `pages/customer/login.vue`

**UI:**
- Form: Phone Number, Password
- Toggle show/hide password
- Link ke register page

**Guard:**
```
User sudah login (customer) → redirect ke /customer
User sudah login (owner) → toast error + redirect ke /owner
```

**API Call:**
```
POST /api/auth/login
Headers: None
Body:
{
  "phone": "+628123456789",
  "password": "Password123"
}

Response Success:
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "phone": "+628123456789",
      "role": "customer"
    }
  }
}

Response Error (401):
{
  "status": "error",
  "message": "Nomor HP atau password salah"
}
```

**Post-Login Flow:**
1. Simpan token di localStorage (`laundryin_token`)
2. Simpan user di localStorage (`laundryin_user`)
3. Update Pinia auth store
4. Connect WebSocket
5. Redirect ke `/customer` atau `?redirect=` query param

---

### 3.4 Register Customer (`/customer/register`)

**File:** `pages/customer/register.vue`

**UI:**
- Form: Full Name, Phone Number, Password, Confirm Password
- Password hint: min 8 chars, uppercase, lowercase, number
- Toggle show/hide password
- Link ke login page

**Guard:**
```
User sudah login → redirect sesuai role
```

**Validasi Frontend:**
- Nama: required
- Phone: required
- Password: min 8 characters
- Confirm Password: must match

**API Call:**
```
POST /api/auth/register
Headers: None
Body:
{
  "name": "John Doe",
  "phone": "+628123456789",
  "password": "Password123",
  "role": "customer"
}

Response Success (201):
{
  "status": "success",
  "message": "Registrasi berhasil",
  "data": {
    "token": "...",
    "user": { ... }
  }
}

Response Error (409 - Duplicate):
{
  "status": "error",
  "message": "Nomor HP sudah terdaftar, silakan login"
}

Response Error (400 - Weak Password):
{
  "status": "error",
  "message": "Password harus mengandung huruf besar, huruf kecil, dan angka"
}
```

---

### 3.5 Order Booking (`/customer/outlet/[id]`)

**File:** `pages/customer/outlet/[id].vue`

**UI:**
- **Left Panel (Form):**
  - Selected outlet info (with "Change" button)
  - Contact Info: Name, Phone
  - Pickup Address (textarea)
  - Pickup Date (DateCarousel component - 14 days)
  - Pickup Time: Morning (09:00-12:00) / Afternoon (13:00-16:00)
  - Special Instructions (textarea)

- **Right Panel (Services):**
  - Search bar
  - Category filter: All, Clothes (KG), Shoes (PCS), Others
  - Service list dengan add/remove qty control
  - Estimated total preview

- **Modals:**
  - Change Outlet Modal: list semua outlet available

**Guard:**
```
User belum login → redirect ke /customer/login?redirect=/customer/outlet/[id]
```

**API Calls:**

1. **Get Outlet Detail:**
```
GET /api/public/outlets/:id
Response:
{
  "status": "success",
  "message": "Data outlet public berhasil diambil",
  "data": {
    "id": "uuid",
    "name": "LaundryIn Premium Branch",
    "address": "Jl. Kertajaya No. 123",
    "phone": "+628123456789"
  }
}
```

2. **Get Services:**
```
GET /api/public/outlets/:id/services
Response:
{
  "status": "success",
  "message": "Daftar layanan public berhasil diambil",
  "data": [
    {
      "id": "uuid",
      "outlet_id": "uuid",
      "name": "Cuci Kering Regular",
      "price": "8000.00",
      "unit": "KG",
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

3. **Create Order (Checkout):**
```
POST /api/orders
Headers: Authorization: Bearer <token>
Body:
{
  "outlet_id": "uuid",
  "customer_name": "John Doe",
  "customer_phone": "+628123456789",
  "customer_address": "Jl. Test No. 123",
  "pickup_date": "2024-03-15",
  "pickup_time": "afternoon",
  "notes": "Antar ke depan rumah",
  "items": [
    {
      "service_id": "uuid",
      "qty": "2.0"
    }
  ]
}

Response Success (201):
{
  "status": "success",
  "message": "Pesanan berhasil dibuat",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "outlet_id": "uuid",
    "total_price": "16000.00",
    "status": "pending",
    "order_date": "2024-03-14T10:00:00Z",
    "items": [...]
  }
}

Response Error (401 - Unauthorized):
{
  "status": "error",
  "message": "Sesi kamu habis, silakan login ulang"
}

Response Error (400 - Validation):
{
  "status": "error",
  "message": "Data tidak valid"
}
```

**Post-Checkout Flow:**
1. Clear cart store
2. Toast success: "Pesanan berhasil dibuat!"
3. Redirect ke `/customer/orders`
4. Backend fires notification ke owner (WebSocket + DB)

---

### 3.6 Order History (`/customer/orders`)

**File:** `pages/customer/orders/index.vue`

**UI:**
- Header "My Orders"
- List order cards (sorted by order_date DESC)
- Order card: order ID (short), outlet name, date, status badge, total price
- Status badges: Pending, In Process, Ready Pickup, Finished

**Guard:**
```
User belum login → redirect ke /customer/login?redirect=/customer/orders
```

**API Call:**
```
GET /api/orders?page=1&limit=10
Headers: Authorization: Bearer <token>
Response:
{
  "status": "success",
  "message": "Daftar pesanan berhasil diambil",
  "data": {
    "data": [
      {
        "id": "uuid",
        "status": "pending",
        "total_price": "16000.00",
        "final_total_price": null,
        "outlet_id": "uuid",
        "order_date": "2024-03-14T10:00:00Z"
      }
    ],
    "page": 1,
    "limit": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

**User Actions:**
1. **Click order card:** Navigate to `/customer/orders/[id]`
2. **Pull to refresh:** Refresh order list

---

### 3.7 Order Tracking (`/customer/orders/[id]`)

**File:** `pages/customer/orders/[id].vue`

**UI:**
- **Left Section (Tracking):**
  - Breadcrumbs: My Orders > Tracking
  - Order ID header
  - Live Status Banner (icon, label, description)
  - Timeline History (vertical steps dengan timestamp)

- **Right Section (Receipt):**
  - Order Summary (items list)
  - Outlet info card
  - Total Payment (estimated vs final)

**Status Config:**
| Status | Icon | Label | Description |
|--------|------|-------|-------------|
| pending | inventory_2 | Menunggu Konfirmasi | Pesanan sudah masuk dan sedang menunggu konfirmasi |
| process | local_laundry_service | Sedang Dicuci | Mesin sedang membersihkan bajumu |
| completed | check_circle | Siap Diambil | Cucian kamu sudah selesai! |
| picked_up | hail | Selesai | Pesanan sudah diambil |

**Timeline Steps:**
```
1. Pesanan Dibuat (10:30 AM) ✓
2. Diterima di Cabang (11:15 AM) ✓
3. Sedang Diproses (12:45 PM) ● LIVE
4. Siap Diambil (pending)
5. Selesai (pending)
```

**API Call:**
```
GET /api/orders?limit=100&page=1
Headers: Authorization: Bearer <token>
Response: Same as order history, dengan items & logs populated
```

**Price Display Logic:**
```
Jika final_total_price ada:
  - Tampilkan harga final (warna success)
  - Coret estimasi harga (line-through)
  - Badge "Harga Final"

Jika final_total_price null:
  - Tampilkan estimasi (warna primary)
  - Badge "Estimasi"
```

---

### 3.8 Profile (`/customer/profile`)

**File:** `pages/customer/profile.vue`

**UI:**
- Avatar circle (initial nama)
- User info: name, phone, role badge
- Detail rows: Nama, Nomor HP, Role
- Buttons: Edit Profil, Keluar

**Edit Mode:**
- Form: Nama, Nomor HP
- Buttons: Batal, Simpan Perubahan

**Guard:**
```
User belum login → redirect ke /customer/login
User role owner → redirect ke /owner/login
```

**API Call (Update Profile):**
```
PUT /api/users/me
Headers: Authorization: Bearer <token>
Body:
{
  "name": "John Doe Updated",
  "phone": "+628123456789"
}

Response Success:
{
  "status": "success",
  "message": "Profil berhasil diperbarui"
}
```

**Local Update Flow:**
1. API call (fallback jika endpoint tidak ada)
2. Update Pinia auth store
3. Update localStorage
4. Show success alert (3 detik)

---

## 4. Flow Owner

### 4.1 Login Owner (`/owner/login`)

**File:** `pages/owner/login.vue`

**UI:**
- Sama dengan customer login
- Tambahan link "Forgot password?" (non-functional)

**Guard:**
```
User sudah login (owner) → redirect ke /owner
User sudah login (customer) → toast error + redirect ke /customer
```

**API Call:**
```
POST /api/auth/login
Body: { phone, password }
Response: Same as customer login
```

**Post-Login Validation:**
```javascript
if (res.data.user.role !== 'owner') {
  toastError('Akun ini bukan akun owner')
  return
}
```

---

### 4.2 Analytics Dashboard (`/owner`)

**File:** `pages/owner/index.vue`

**UI:**
- **Date Filter:** 7 Hari, 30 Hari, 3 Bulan
- **Revenue Card:** Total omzet periode terpilih
- **Pipeline Status:** 4 cards (Pending, Process, Completed, Picked Up)
- **Top Services:** List 5 layanan teratas (revenue-based)

**Guard:**
```
User belum login → redirect ke /owner/login
User role customer → redirect ke /customer
```

**API Calls:**

1. **Get Omzet:**
```
GET /api/reports/omzet?start_date=2024-02-14&end_date=2024-03-14
Headers: Authorization: Bearer <token>
Response:
{
  "status": "success",
  "message": "Berhasil mengambil data omzet",
  "data": {
    "total_omzet": "1500000.00"
  }
}
```

2. **Get Order Status Summary:**
```
GET /api/reports/orders/summary?start_date=...&end_date=...
Response:
{
  "status": "success",
  "message": "Berhasil mengambil ringkasan status pesanan",
  "data": {
    "pending": 5,
    "process": 3,
    "completed": 10,
    "picked_up": 8,
    "cancelled": 2
  }
}
```

3. **Get Top Services:**
```
GET /api/reports/services/top?start_date=...&end_date=...
Response:
{
  "status": "success",
  "message": "Berhasil mengambil data layanan terlaris",
  "data": [
    {
      "service_name": "Cuci Kering Regular",
      "outlet_name": "LaundryIn Premium Branch",
      "total_qty": "50.00",
      "total_revenue": "400000.00"
    }
  ]
}
```

---

### 4.3 Outlet Management (`/owner/outlets`)

**File:** `pages/owner/outlets.vue`

**UI:**
- Header dengan tombol "Add Outlet" (+)
- Search bar
- Outlet list (name, address, status)
- Action buttons per outlet: Edit, Delete

**Modals:**
1. **Add/Edit Modal:**
   - Form: Nama Outlet, Alamat, Nomor Telepon
   - Validation: semua field required, phone format E.164

2. **Delete Confirmation Modal:**
   - Warning message
   - Buttons: Ya Hapus, Batal

**Guard:**
```
User belum login → redirect ke /owner/login
User role customer → redirect ke /customer
```

**API Calls:**

1. **Get All Outlets:**
```
GET /api/outlets?page=1&limit=10
Headers: Authorization: Bearer <token>
Response:
{
  "status": "success",
  "message": "Daftar outlet berhasil diambil",
  "data": {
    "data": [...]
  }
}
```

2. **Create Outlet:**
```
POST /api/outlets
Headers: Authorization: Bearer <token>
Body:
{
  "name": "LaundryIn Cabang 2",
  "address": "Jl. Baru No. 456",
  "phone": "+628987654321"
}
Response (201):
{
  "status": "success",
  "message": "Outlet berhasil dibuat",
  "data": { ... }
}
```

3. **Update Outlet:**
```
PUT /api/outlets/:id
Headers: Authorization: Bearer <token>
Body: { name, address, phone }
Response (200):
{
  "status": "success",
  "message": "Outlet berhasil diupdate",
  "data": { ... }
}
```

4. **Delete Outlet:**
```
DELETE /api/outlets/:id
Headers: Authorization: Bearer <token>
Response (200):
{
  "status": "success",
  "message": "Outlet berhasil dihapus"
}
```

**Validation:**
```javascript
// Phone format validation
if (!/^\+62\d{8,13}$/.test(phone.trim())) {
  formErrors.value.phone = 'Format harus +62xxx (contoh: +628123456789)'
}
```

---

### 4.4 Service Management (`/owner/services`)

**File:** `pages/owner/services.vue`

**UI:**
- Header dengan tombol "New Service"
- Filter Outlet dropdown
- Service list: icon, name, outlet_name, unit, price
- Action buttons: Edit, Delete

**Modals:**
1. **Add/Edit Modal:**
   - Outlet selection (dropdown, readonly jika edit)
   - Nama Layanan
   - Harga (number input)
   - Unit (KG/PCS dropdown)

2. **Delete Confirmation Modal**

**Guard:** Same as owner routes

**API Calls:**

1. **Get Services by Outlet:**
```
GET /api/outlets/:id/services
Headers: Authorization: Bearer <token>
Response:
{
  "status": "success",
  "message": "Daftar layanan berhasil diambil",
  "data": [
    {
      "id": "uuid",
      "outlet_id": "uuid",
      "name": "Cuci Kering Regular",
      "price": "8000.00",
      "unit": "KG"
    }
  ]
}
```

2. **Create Service:**
```
POST /api/services
Headers: Authorization: Bearer <token>
Body:
{
  "outlet_id": "uuid",
  "name": "Cuci Karpet",
  "price": "15000",
  "unit": "PCS"
}
```

3. **Update Service:**
```
PUT /api/services/:id
Headers: Authorization: Bearer <token>
Body: { name, price, unit }
```

4. **Delete Service:**
```
DELETE /api/services/:id
Headers: Authorization: Bearer <token>
```

---

### 4.5 Order Pipeline (`/owner/orders`)

**File:** `pages/owner/orders.vue`

**UI:**
- Filter Outlet dropdown
- Status Tabs: Semua, PENDING, PROCESS, COMPLETED, PICKED_UP, CANCELLED
- Order cards dengan:
  - Status badge
  - Customer name
  - Items list (qty unit service_name)
  - Total price (estimated/final)
  - Action buttons berdasarkan status

**Input Berat Aktual:**
- Muncul jika status = `pending` dan ada item KG
- Input field per item KG
- Required sebelum bisa "Proses"

**Status Transitions:**
```
pending → process (input actual_qty untuk KG)
process → completed
completed → picked_up
pending → cancelled
process → cancelled
```

**Guard:** Same as owner routes

**API Calls:**

1. **Get Orders by Outlet:**
```
GET /api/outlets/:id/orders?page=1&limit=10
Headers: Authorization: Bearer <token>
Response:
{
  "status": "success",
  "message": "Data pesanan outlet berhasil diambil",
  "data": {
    "data": [
      {
        "id": "uuid",
        "status": "pending",
        "total_price": "16000.00",
        "final_total_price": null,
        "customer_name": "John Doe",
        "items": [
          {
            "id": "uuid",
            "service_name": "Cuci Kering",
            "qty": "2.0",
            "actual_qty": null,
            "unit": "KG",
            "subtotal": "16000.00",
            "final_price": null
          }
        ]
      }
    ]
  }
}
```

2. **Update Order Status:**
```
PATCH /api/orders/:id/status
Headers: Authorization: Bearer <token>
Body:
{
  "status": "process",
  "items": [
    {
      "id": "order_item_id",
      "actual_qty": "2.5"
    }
  ]
}

Response (200):
{
  "status": "success",
  "message": "Status pesanan berhasil diperbarui",
  "data": { ... }
}

Response Error (400 - Invalid Transition):
{
  "status": "error",
  "message": "Transisi status pesanan tidak valid"
}

Response Error (400 - Missing KG Input):
{
  "status": "error",
  "message": "Berat aktual wajib diisi untuk layanan per KG"
}
```

**Post-Status Update Flow:**
1. Refresh order list
2. Toast success
3. Backend fires notification ke customer (WebSocket + DB)
4. Jika status = `process` dengan actual_qty, trigger `NotifyPriceUpdated`

---

## 5. Flow Notifikasi (WebSocket)

### 5.1 Koneksi WebSocket

**Endpoint:** `ws://localhost:8080/api/v1/ws/connect?token=<jwt_token>`

**Flow:**
```
1. User login → token disimpan di Pinia auth store
2. Plugin websocket.client.ts watch authStore.isLoggedIn
3. Jika logged in → connect() dipanggil
4. WebSocket upgrade dengan token di query param
5. Middleware validates token → set user_id, role di context
6. Client registered di Hub
7. WritePump & ReadPump dimulai
```

**Reconnect Logic:**
```javascript
let reconnectDelay = 1000
const MAX_RECONNECT_DELAY = 30000

const scheduleReconnect = () => {
  reconnectTimer = setTimeout(() => {
    reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY)
    connect()
  }, reconnectDelay)
}
```

**Exponential Backoff:**
- Attempt 1: 1s
- Attempt 2: 2s
- Attempt 3: 4s
- ...
- Max: 30s

---

### 5.2 Notification Triggers

#### 5.2.1 Order Created (NotifyOrderCreated)

**Trigger:** Customer creates order (POST /api/orders)

**Penerima:** Owner dari outlet

**Type:** `new_order`

**Title:** "Pesanan Baru Masuk"

**Body:** "{customer_name} memesan {services} di {outlet_name}"

**Data:**
```json
{
  "order_id": "uuid",
  "outlet_id": "uuid",
  "customer_name": "John Doe",
  "total_price": "16000.00"
}
```

**Flow:**
```
1. OrderUsecase.Create() selesai
2. go u.notifUsecase.NotifyOrderCreated(context.Background(), order)
3. Create notification di DB
4. hub.SendToUser(outlet.UserID, message)
5. Frontend WebSocket onmessage → addNotification()
6. Toast info + badge unread count increment
```

---

#### 5.2.2 Status Changed (NotifyStatusChanged)

**Trigger:** Owner updates order status (PATCH /api/orders/:id/status)

**Penerima:** Customer yang make order

**Type:** `order_status` | `order_cancelled`

**Status Transitions:**

| New Status | Type | Title | Body |
|------------|------|-------|------|
| process | order_status | "Pesananmu Sedang Diproses" | "Outlet {name} mulai memproses pesanan #{id}" |
| completed | order_status | "Pesananmu Siap Diambil! 🎉" | "Cucian kamu sudah selesai, silakan ambil di {name}" |
| cancelled | order_cancelled | "Pesananmu Dibatalkan" | "Pesanan #{id} di {name} telah dibatalkan" |

**Data:**
```json
{
  "order_id": "uuid",
  "new_status": "process"
}
```

**Special Case - Price Updated:**
Jika status = `process` dan `FinalTotalPrice` set:
- Trigger additional `notifyPriceUpdated()`
- Type: `price_updated`
- Title: "Harga Final Pesananmu Sudah Diketahui"
- Body: "Total pembayaran pesanan #{id} adalah Rp {final_price}"

---

### 5.3 Frontend WebSocket Handling

**Plugin:** `app/plugins/websocket.client.ts`

```typescript
export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()
  const notifStore = useNotificationStore()
  const { connect, disconnect } = useWebSocket()

  // Initial fetch unread count
  if (authStore.isLoggedIn) {
    notifStore.fetchUnreadCount()
  }

  // Watch login state
  watch(() => authStore.isLoggedIn, (loggedIn) => {
    if (loggedIn) connect()
    else disconnect()
  }, { immediate: true })
})
```

**Composable:** `app/composables/useWebSocket.ts`

```typescript
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data)
  notifStore.addNotification({
    ...msg,
    id: msg.id || Math.random().toString(36).substring(7),
    is_read: false,
    created_at: new Date().toISOString()
  })
  info(msg.title || 'Notifikasi Baru')
}
```

---

### 5.4 Notification Store Structure

**File:** `app/stores/notification.ts`

```typescript
interface Notification {
  id: string
  type: string
  title: string
  body: string
  data: any
  is_read: boolean
  created_at: string
}

interface NotificationStore {
  notifications: Notification[]
  unreadCount: number
  isOpen: boolean
  loading: boolean
}
```

**Actions:**
- `addNotification(notif)` - Add from WebSocket
- `fetchNotifications(page, limit)` - Fetch from API
- `fetchUnreadCount()` - Get unread count
- `markAsRead(id)` - Mark single as read
- `markAllAsRead()` - Mark all as read
- `toggleDropdown(val)` - Toggle UI dropdown

---

## 6. Auth & Guard Rules

### 6.1 Redirect Conditions

| Kondisi | Halaman | Aksi |
|---------|---------|------|
| User sudah login (customer) | `/customer/login` | Redirect ke `/customer` |
| User sudah login (owner) | `/customer/login` | Toast error + redirect ke `/owner` |
| User sudah login (customer) | `/customer/register` | Redirect ke `/customer` |
| User sudah login (owner) | `/customer/register` | Toast error + redirect ke `/owner` |
| User sudah login (owner) | `/owner/login` | Redirect ke `/owner` |
| User sudah login (customer) | `/owner/login` | Toast error + redirect ke `/customer` |
| User belum login | `/customer/orders` | Redirect ke `/customer/login?redirect=/customer/orders` |
| User belum login | `/customer/orders/[id]` | Redirect ke `/customer/login` |
| User belum login | `/customer/profile` | Redirect ke `/customer/login` |
| User role owner | `/customer/profile` | Redirect ke `/owner/login` |
| User belum login | `/owner` | Redirect ke `/owner/login` |
| User role customer | `/owner` | Redirect ke `/customer` |
| User belum login | `/owner/outlets` | Redirect ke `/owner/login` |
| User belum login | `/owner/services` | Redirect ke `/owner/login` |
| User belum login | `/owner/orders` | Redirect ke `/owner/login` |

### 6.2 Guard Implementation Examples

**Customer Login Guard:**
```typescript
watchEffect(() => {
  if (import.meta.client && authStore.isLoggedIn) {
    if (authStore.user?.role === 'owner') {
      toastError('Kamu sudah login sebagai owner, silakan logout terlebih dahulu')
      router.push('/owner')
    } else if (authStore.user?.role === 'customer') {
      router.push('/customer')
    }
  }
})
```

**Owner Page Guard:**
```typescript
watchEffect(() => {
  if (import.meta.client) {
    if (!authStore.isLoggedIn || authStore.user?.role !== 'owner') {
      router.push('/owner/login')
    }
  }
})
```

**Customer Page Guard:**
```typescript
watchEffect(() => {
  if (import.meta.client && !authStore.isLoggedIn) {
    router.push('/customer/login?redirect=/customer/orders')
  }
})
```

---

## 7. State Management

### 7.1 Auth Store

**File:** `app/stores/auth.ts`

```typescript
interface User {
  id: string
  name: string
  phone: string
  email?: string
  role: 'owner' | 'customer'
}

interface AuthState {
  token: string | null
  user: User | null
}
```

**State:**
- `token`: JWT token string
- `user`: User object

**Getters:**
- `isLoggedIn`: boolean - Check if token exists
- `isOwner`: boolean - Check if role === 'owner'
- `isCustomer`: boolean - Check if role === 'customer'
- `authHeader`: string - Format "Bearer <token>" untuk API calls

**Actions:**
- `setAuth(token, user)` - Set token & user, save to localStorage
- `logout()` - Clear token & user, remove from localStorage
- `restoreSession()` - Load from localStorage on app init

**LocalStorage Keys:**
- `laundryin_token`: JWT token
- `laundryin_user`: JSON.stringify(user)

---

### 7.2 Cart Store

**File:** `app/stores/cart.ts`

```typescript
interface CartItem {
  serviceId: string
  name: string
  price: number
  unit: string
  qty: string // zero-trust: string untuk presisi decimal
}

interface CartState {
  items: CartItem[]
  outletId: string | null
}
```

**State:**
- `items`: Array of CartItem
- `outletId`: Current selected outlet ID

**Getters:**
- `totalPreview`: number - Sum of (price * qty)
- `itemCount`: number - Length of items array

**Actions:**
- `setOutlet(id)` - Set outlet, clear items if different
- `addItem(item)` - Add or merge qty
- `updateQty(serviceId, qty)` - Update quantity
- `removeItem(serviceId)` - Remove from cart
- `clearCart()` - Reset all

---

### 7.3 Notification Store

**File:** `app/stores/notification.ts`

(See section 5.4 for structure)

---

### 7.4 Session Restore Flow

```
1. App init (auth.client.ts plugin)
2. authStore.restoreSession() called
3. Read localStorage:
   - laundryin_token
   - laundryin_user
4. Parse and set state
5. If valid → websocket connects automatically
6. If invalid/expired → logout() clears storage
```

---

## 8. Database Schema

### 8.1 Tables Overview

```
┌─────────────────┐     ┌─────────────────┐
│     users       │     │     outlets     │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │◀────│ id (PK)         │
│ name            │     │ user_id (FK)    │
│ phone           │     │ name            │
│ email           │     │ address         │
│ password        │     │ phone           │
│ role            │     │ created_at      │
│ created_at      │     │ updated_at      │
│ updated_at      │     │ deleted_at      │
│ deleted_at      │     └─────────────────┘
└─────────────────┘              │
                                 │ 1:N
                                 ▼
                        ┌─────────────────┐
                        │    services     │
                        ├─────────────────┤
                        │ id (PK)         │
                        │ outlet_id (FK)  │
                        │ name            │
                        │ price (decimal) │
                        │ unit            │
                        │ created_at      │
                        │ updated_at      │
                        │ deleted_at      │
                        └─────────────────┘

┌─────────────────┐     ┌─────────────────┐
│     orders      │     │   order_items   │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │◀────│ id (PK)         │
│ user_id (FK)    │     │ order_id (FK)   │
│ outlet_id (FK)  │     │ service_name    │
│ total_price     │     │ service_price   │
│ final_total_price    │ qty (decimal)   │
│ status          │     │ actual_qty      │
│ order_date      │     │ unit            │
│ created_at      │     │ subtotal        │
│ updated_at      │     │ final_price     │
│ deleted_at      │     └─────────────────┘
└─────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐     ┌─────────────────┐
│   order_logs    │     │  notifications  │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │     │ id (PK)         │
│ order_id (FK)   │     │ user_id (FK)    │
│ updated_by (FK) │     │ type            │
│ old_status      │     │ title           │
│ new_status      │     │ body            │
│ created_at      │     │ data (jsonb)    │
└─────────────────┘     │ is_read         │
                        │ created_at      │
                        └─────────────────┘
```

---

### 8.2 Table Details

#### users
| Column | Type | Constraints | Deskripsi |
|--------|------|-------------|-----------|
| id | UUID | PK, default uuid_generate_v4() | User ID |
| name | TEXT | NOT NULL | Full name |
| phone | TEXT | UNIQUE, NOT NULL | E.164 phone (+62xxx) |
| email | TEXT | NULLABLE | Email (optional) |
| password | TEXT | NOT NULL | Bcrypt hash |
| role | TEXT | NOT NULL | 'owner' or 'customer' |
| created_at | TIMESTAMPTZ | NOT NULL | Created timestamp |
| updated_at | TIMESTAMPTZ | NOT NULL | Updated timestamp |
| deleted_at | TIMESTAMPTZ | NULLABLE | Soft delete timestamp |

#### outlets
| Column | Type | Constraints | Deskripsi |
|--------|------|-------------|-----------|
| id | UUID | PK | Outlet ID |
| user_id | UUID | FK → users(id), NOT NULL, INDEX | Owner ID |
| name | VARCHAR(100) | NOT NULL | Outlet name |
| address | TEXT | NOT NULL | Full address |
| phone | VARCHAR(20) | NOT NULL | E.164 phone |
| created_at | TIMESTAMPTZ | NOT NULL | |
| updated_at | TIMESTAMPTZ | NOT NULL | |
| deleted_at | TIMESTAMPTZ | NULLABLE | Soft delete |

#### services
| Column | Type | Constraints | Deskripsi |
|--------|------|-------------|-----------|
| id | UUID | PK | Service ID |
| outlet_id | UUID | FK → outlets(id), NOT NULL, INDEX | Outlet ID |
| name | VARCHAR(100) | NOT NULL | Service name |
| price | NUMERIC(10,2) | NOT NULL | Price (decimal) |
| unit | VARCHAR(20) | NOT NULL | 'KG', 'PCS', 'METER' |
| created_at | TIMESTAMPTZ | NOT NULL | |
| updated_at | TIMESTAMPTZ | NOT NULL | |
| deleted_at | TIMESTAMPTZ | NULLABLE | Soft delete |

#### orders
| Column | Type | Constraints | Deskripsi |
|--------|------|-------------|-----------|
| id | UUID | PK | Order ID |
| user_id | UUID | FK → users(id), NOT NULL, INDEX | Customer ID |
| outlet_id | UUID | FK → outlets(id), NOT NULL, INDEX | Outlet ID |
| total_price | NUMERIC(12,2) | NOT NULL | Estimated total |
| final_total_price | NUMERIC(12,2) | NULLABLE | Actual total (after weighing) |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'pending' | FSM state |
| order_date | TIMESTAMPTZ | autoCreateTime | Order timestamp |
| created_at | TIMESTAMPTZ | NOT NULL | |
| updated_at | TIMESTAMPTZ | NOT NULL | |
| deleted_at | TIMESTAMPTZ | NULLABLE | Soft delete |

**Status FSM:**
```
pending → process → completed → picked_up
   ↓          ↓
cancelled    cancelled
```

#### order_items
| Column | Type | Constraints | Deskripsi |
|--------|------|-------------|-----------|
| id | UUID | PK | Order Item ID |
| order_id | UUID | FK → orders(id), NOT NULL, INDEX | Order ID |
| service_name | VARCHAR(100) | NOT NULL | Snapshot service name |
| service_price | NUMERIC(10,2) | NOT NULL | Snapshot price |
| qty | NUMERIC(6,2) | NOT NULL | Ordered quantity |
| actual_qty | NUMERIC(6,2) | NULLABLE | Actual quantity (KG only) |
| unit | VARCHAR(20) | NOT NULL | 'KG' or 'PCS' |
| subtotal | NUMERIC(12,2) | NOT NULL | qty * service_price |
| final_price | NUMERIC(12,2) | NULLABLE | actual_qty * service_price |

#### order_logs
| Column | Type | Constraints | Deskripsi |
|--------|------|-------------|-----------|
| id | UUID | PK | Log ID |
| order_id | UUID | FK → orders(id), NOT NULL, INDEX | Order ID |
| updated_by | UUID | FK → users(id), NOT NULL | Who changed status |
| old_status | VARCHAR(20) | NULLABLE | Previous status |
| new_status | VARCHAR(20) | NOT NULL | New status |
| created_at | TIMESTAMPTZ | autoCreateTime | Change timestamp |

#### notifications
| Column | Type | Constraints | Deskripsi |
|--------|------|-------------|-----------|
| id | UUID | PK | Notification ID |
| user_id | UUID | FK → users(id), NOT NULL, INDEX | Recipient ID |
| type | VARCHAR(50) | NOT NULL | 'new_order', 'order_status', 'price_updated', 'order_cancelled' |
| title | VARCHAR(200) | NOT NULL | Notification title |
| body | VARCHAR(500) | NOT NULL | Notification body |
| data | JSONB | NOT NULL | Structured data |
| is_read | BOOLEAN | DEFAULT false | Read status |
| created_at | TIMESTAMPTZ | autoCreateTime | Created timestamp |

---

### 8.3 Relasi

```
users (1) ──────< outlets (N)
  │
  │ (1)
  │
  └──────< orders (N)
              │
              │ (1)
              │
              └──────< order_items (N)
              │
              │ (1)
              │
              └──────< order_logs (N)

outlets (1) ────< services (N)

users (1) ──────< notifications (N)
```

---

## 9. API Routes Summary

### 9.1 Public Routes (No Auth)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/ping` | Health check |
| POST | `/api/v1/auth/register` | Register user |
| POST | `/api/v1/auth/login` | Login |
| GET | `/api/v1/public/outlets` | Get all outlets (public) |
| GET | `/api/v1/public/outlets/:id` | Get outlet by ID (public) |
| GET | `/api/v1/public/outlets/:id/services` | Get services by outlet (public) |

---

### 9.2 Auth Routes (Customer Only)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/orders` | Create order |
| GET | `/api/v1/orders` | Get user's orders |

---

### 9.3 Auth Routes (Owner Only)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/outlets` | Create outlet |
| GET | `/api/v1/outlets` | Get owner's outlets |
| GET | `/api/v1/outlets/:id` | Get outlet by ID |
| PUT | `/api/v1/outlets/:id` | Update outlet |
| DELETE | `/api/v1/outlets/:id` | Delete outlet |
| POST | `/api/v1/services` | Create service |
| GET | `/api/v1/outlets/:id/services` | Get services by outlet |
| PUT | `/api/v1/services/:id` | Update service |
| DELETE | `/api/v1/services/:id` | Delete service |
| GET | `/api/v1/outlets/:id/orders` | Get orders by outlet |
| PATCH | `/api/v1/orders/:id/status` | Update order status |
| GET | `/api/v1/reports/omzet` | Get revenue report |
| GET | `/api/v1/reports/orders/summary` | Get order status summary |
| GET | `/api/v1/reports/services/top` | Get top services |

---

### 9.4 Auth Routes (Both Roles)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/notifications` | Get user notifications |
| GET | `/api/v1/notifications/unread-count` | Get unread count |
| PATCH | `/api/v1/notifications/:id/read` | Mark as read |
| PATCH | `/api/v1/notifications/read-all` | Mark all as read |
| GET | `/api/v1/ws/connect` | WebSocket connection |

---

### 9.5 WebSocket Messages

#### Server → Client

**Order Created:**
```json
{
  "type": "new_order",
  "title": "Pesanan Baru Masuk",
  "body": "John Doe memesan Cuci Kering di LaundryIn Premium Branch",
  "data": {
    "order_id": "uuid",
    "outlet_id": "uuid",
    "customer_name": "John Doe",
    "total_price": "16000.00"
  },
  "timestamp": "2024-03-14T10:00:00Z"
}
```

**Status Changed:**
```json
{
  "type": "order_status",
  "title": "Pesananmu Sedang Diproses",
  "body": "Outlet LaundryIn Premium Branch mulai memproses pesanan #abc12345",
  "data": {
    "order_id": "uuid",
    "new_status": "process"
  },
  "timestamp": "2024-03-14T11:00:00Z"
}
```

**Price Updated:**
```json
{
  "type": "price_updated",
  "title": "Harga Final Pesananmu Sudah Diketahui",
  "body": "Total pembayaran pesanan #abc12345 adalah Rp 20000.00",
  "data": {
    "order_id": "uuid",
    "estimated_price": "16000.00",
    "final_price": "20000.00"
  },
  "timestamp": "2024-03-14T11:00:00Z"
}
```

**Order Cancelled:**
```json
{
  "type": "order_cancelled",
  "title": "Pesananmu Dibatalkan",
  "body": "Pesanan #abc12345 di LaundryIn Premium Branch telah dibatalkan",
  "data": {
    "order_id": "uuid"
  },
  "timestamp": "2024-03-14T12:00:00Z"
}
```

---

## 10. Catatan Teknis Penting

### 10.1 Format Data Khusus

#### Phone Number (E.164 Strict)
```
Format: +[country_code][number]
Contoh: +628123456789 (Indonesia)
Validasi: ^\+[1-9]\d{6,14}$
Frontend auto-format: 08112233445 → +628112233445
```

#### Price (Decimal Precision)
```
Backend: decimal.Decimal (shopspring/decimal)
API Response: string dengan 2 decimal places
Contoh: "8000.00", "15000.00"
Frontend: parseFloat() untuk kalkulasi
```

#### Quantity (Zero-Trust)
```
DTO: string (bukan number)
Alasan: preserve decimal precision
Frontend: qty selalu string "2.0", "0.5"
Backend: parse ke decimal.Decimal
Validasi: must be positive number
```

#### UUID (String-based)
```
Go: plain string (bukan uuid.UUID)
Format: uuid_generate_v4() dari PostgreSQL
Contoh: "550e8400-e29b-41d4-a716-446655440000"
Display: slice 8 chars → "#550E8400"
```

---

### 10.2 Gotcha di API

#### 1. Payload Limit
```
Max: 1MB (1024 * 1024 bytes)
Middleware: PayloadLimit(1024 * 1024)
Error: 413 Payload Too Large
```

#### 2. Rate Limiting
```
Limit: 60 requests/minute per IP
Algorithm: Token bucket (1 token/sec, burst 60)
Cleanup: Background worker setiap 1 menit
Error: 429 Too Many Requests
```

#### 3. Context Timeout
```
Timeout: 5 detik per request
Error: 408 Request Timeout
Message: "Proses terlalu lama, silakan coba lagi"
```

#### 4. Anti-IDOR Pattern
```go
// SALAH: Direct ID lookup
order := FindByID(orderID)

// BENAR: JOIN dengan ownership check
order := FindByIDAndOwner(orderID, userID)
// JOIN orders JOIN outlets WHERE outlets.user_id = userID
```

#### 5. Zero-Trust Pricing
```go
// Frontend hanya kirim service_id dan qty
// Backend fetch service dari DB, hitung subtotal sendiri
service := FindByIDAndOutletID(service_id, outlet_id)
// Anti-IDOR: verify service belongs to outlet
subtotal := qty * service.Price
```

---

### 10.3 Pattern yang Dipakai

#### Repository Pattern
```
Handler → Usecase → Repository → Database
```

#### Dependency Injection
```go
// main.go
userRepo := repository.NewUserRepository(db)
authUsecase := usecase.NewAuthUsecase(userRepo)
authHandler := handler.NewAuthHandler(authUsecase)
```

#### FSM (Finite State Machine)
```go
// Valid transitions only
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

#### ACID Transactions
```go
// Order creation dengan items
db.Transaction(func(tx *gorm.DB) error {
  // 1. Create order header
  tx.Create(order)
  // 2. Create items sequentially
  for _, item := range items {
    tx.Create(&item)
  }
  // Rollback otomatis jika ada error
})
```

#### Soft Delete
```go
// GORM DeletedAt field
type Order struct {
  gorm.DeletedAt `gorm:"index"`
}

// Query otomatis exclude deleted
db.Find(&orders)
// WHERE deleted_at IS NULL

// Include deleted
db.Unscoped().Find(&orders)
```

---

### 10.4 Known Issues

#### 1. WebSocket Origin Check
```go
// Current: Allow all origins
CheckOrigin: func(r *http.Request) bool {
  return true
}
// TODO: Restrict to production domain
```

#### 2. Notification Memory Limit
```go
// Frontend limit: 50 notifications in memory
if (this.notifications.length > 50) {
  this.notifications.pop()
}
```

#### 3. Profile Update Fallback
```typescript
// Frontend update locally jika API endpoint tidak ada
await $fetch('/api/users/me', {...}).catch(err => {
  console.warn('API update failed, updating locally only')
})
```

#### 4. Tracking Page Hardcoded
```
/customer/tracking menggunakan static data
TODO: Connect ke real order tracking API
```

---

## 11. File Structure

### Frontend
```
frontend/
├── app/
│   ├── app.vue                 # Root component
│   ├── assets/
│   │   └── css/
│   │       └── main.css        # Global styles
│   ├── components/
│   │   └── ui/
│   │       ├── DateCarousel.vue
│   │       ├── NotificationDropdown.vue
│   │       └── ToastContainer.vue
│   ├── composables/
│   │   ├── useToast.ts
│   │   └── useWebSocket.ts
│   ├── layouts/
│   │   ├── customer.vue
│   │   └── owner.vue
│   ├── middleware/             # (empty)
│   ├── pages/
│   │   ├── index.vue           # Landing page
│   │   ├── customer/
│   │   │   ├── index.vue       # Explore outlets
│   │   │   ├── login.vue
│   │   │   ├── register.vue
│   │   │   ├── profile.vue
│   │   │   ├── tracking.vue
│   │   │   └── outlet/
│   │   │       └── [id].vue    # Order booking
│   │   ├── customer/orders/
│   │   │   ├── index.vue       # Order history
│   │   │   └── [id].vue        # Order tracking
│   │   └── owner/
│   │       ├── index.vue       # Analytics
│   │       ├── login.vue
│   │       ├── outlets.vue
│   │       ├── services.vue
│   │       └── orders.vue
│   ├── plugins/
│   │   ├── auth.client.ts      # Session restore
│   │   └── websocket.client.ts # WS connection
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── cart.ts
│   │   └── notification.ts
│   └── types/
│       └── api.ts
├── nuxt.config.ts
├── package.json
└── ...
```

### Backend
```
backend/
├── cmd/
│   └── api/
│       └── main.go             # Entry point
├── internal/
│   ├── database/
│   │   └── postgres.go         # DB connection
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
│   │   ├── dto.go              # Main DTOs
│   │   └── report_dto.go       # Report DTOs
│   ├── repository/
│   │   ├── models/
│   │   │   └── models.go       # GORM models
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
│       └── hub.go              # WebSocket hub
├── pkg/
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── password.go
│   │   ├── response.go
│   │   ├── string.go
│   │   └── validator.go
├── tests/
├── .env
├── go.mod
├── go.sum
└── ...
```

---

## 12. Environment Variables

### Backend (.env)

```bash
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password_lu
DB_NAME=laundryin_db
DB_PORT=5432
DB_SSLMODE=disable

# JWT
JWT_SECRET=rahasia_negara_lu

# Server
PORT=8080
```

### Frontend

Tidak ada environment variables khusus. Proxy configured di `nuxt.config.ts`:

```typescript
export default defineNuxtConfig({
  routeRules: {
    '/api/**': { proxy: 'http://localhost:8080/api/v1/**' }
  }
})
```

**Production Deployment:**
- Ganti proxy URL ke production backend
- Atau gunakan environment variable:
```typescript
routeRules: {
  '/api/**': { proxy: `${process.env.BACKEND_URL}/api/v1/**` }
}
```

---

## Lampiran A: Status FSM Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    ORDER STATUS FSM                         │
└─────────────────────────────────────────────────────────────┘

                    ┌─────────────┐
                    │   PENDING   │
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               │
    ┌──────────┐    ┌──────────┐          │
    │CANCELLED │    │ PROCESS  │          │
    │  (END)   │    └────┬─────┘          │
    └──────────┘         │                │
                         │                │
                         ▼                │
                  ┌──────────┐            │
                  │COMPLETED │            │
                  └────┬─────┘            │
                       │                  │
                       ▼                  │
                ┌──────────┐              │
                │PICKED_UP │              │
                │  (END)   │◀─────────────┘
                └──────────┘

Valid Transitions:
- pending → process (input actual_qty for KG)
- pending → cancelled
- process → completed
- process → cancelled
- completed → picked_up
```

---

## Lampiran B: Notification Type Matrix

| Event | Trigger | Recipient | Type | Title | Body Template |
|-------|---------|-----------|------|-------|---------------|
| Order Created | POST /orders | Owner | new_order | "Pesanan Baru Masuk" | "{customer} memesan {services} di {outlet}" |
| Status → Process | PATCH /status | Customer | order_status | "Pesananmu Sedang Diproses" | "Outlet {name} mulai memproses pesanan #{id}" |
| Status → Completed | PATCH /status | Customer | order_status | "Pesananmu Siap Diambil! 🎉" | "Cucian kamu sudah selesai, silakan ambil di {name}" |
| Status → Cancelled | PATCH /status | Customer | order_cancelled | "Pesananmu Dibatalkan" | "Pesanan #{id} di {name} telah dibatalkan" |
| Price Updated | PATCH /status (process) | Customer | price_updated | "Harga Final Pesananmu Sudah Diketahui" | "Total pembayaran pesanan #{id} adalah Rp {price}" |

---

## Lampiran C: Quick Reference Commands

### Backend Development
```bash
cd backend
go run cmd/api/main.go
# or
go build -o bin/api cmd/api/main.go && ./bin/api
```

### Frontend Development
```bash
cd frontend
npm install
npm run dev
```

### Database Migration
```bash
# Auto-migrate on backend start
go run cmd/api/main.go
```

### Test API
```bash
# Health check
curl http://localhost:8080/ping

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"+628123456789","password":"Password123"}'
```

---

**Dokumentasi ini dibuat berdasarkan codebase LaundryIn per Maret 2026.**
**Untuk pertanyaan atau update, hubungi tim development.**
