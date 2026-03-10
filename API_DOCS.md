# LaundryIn API Documentation

Dokumentasi ini menjelaskan endpoints API untuk aplikasi LaundryIn. Aplikasi ini dirancang sebagai sistem Kasir (Point of Sales) untuk pengelolaan bisnis laundry. Semua data finansial (seperti harga dan subtotal) memiliki ketelitian desimal logis `float` dalam pengiriman payload namun disimpan/diproses via sistem desimal/string untuk menghindari hilangnya presisi (loss of precision/floating point anomaly) saat masuk ke database/JSON.

## Base URL
Kecuali ditentukan lain, base URL untuk sistem ini adalah `/api/v1`.

### Standard Response Format
Semua output API distandarkan mengikuti format berikut:

**Success Response (2xx)**
```json
{
  "status": "success",
  "message": "Pesan sukses generik",
  "data": { ... } // atau array [...]
}
```

**Error Response (4xx, 5xx)**
```json
{
  "status": "error",
  "message": "Pesan error spesifik/generik",
  "errors": "Detail error jika ada (opsional)" // Dihilangkan jika error keamanan/internal
}
```

---

## Endpoint Overview Table

| Method | Endpoint | Group | Role | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/ping` | Base | Public | Health Check API |
| `POST` | `/api/v1/auth/register` | Auth | Public | Daftar akun baru (bawaan: role `owner`) |
| `POST` | `/api/v1/auth/login` | Auth | Public | Autentikasi dan dapatkan JWT Token |
| `POST` | `/api/v1/outlets` | Outlets | Owner | Membuat data cabang (outlet) baru |
| `GET` | `/api/v1/outlets` | Outlets | Owner | Menampilkan seluruh outlet milik kasir (mendukung pagination) |
| `GET` | `/api/v1/outlets/:id` | Outlets | Owner | Menampilkan 1 data spesifik outlet |
| `PUT` | `/api/v1/outlets/:id` | Outlets | Owner | Menyunting data outlet |
| `DELETE` | `/api/v1/outlets/:id` | Outlets | Owner | Menghapus (Soft-Delete) outlet |
| `POST` | `/api/v1/services` | Services | Owner | Membuat layanan kebersihan di outlet |
| `GET` | `/api/v1/outlets/:id/services` | Services | Owner | Seluruh list layanan kepunyaan satu outlet |
| `PUT` | `/api/v1/services/:id` | Services | Owner | Update harga/nama layanan (Kirim Price sebagai String) |
| `DELETE` | `/api/v1/services/:id` | Services | Owner | Menghapus (Soft-Delete) layanan |
| `POST` | `/api/v1/orders` | Orders | Owner | Membuat pesanan baru / Checkout (Otomatis generate status 'pending') |
| `GET` | `/api/v1/orders` | Orders | Owner | Seluruh history pesanan (semua outlet milik owner ini) |
| `GET` | `/api/v1/outlets/:id/orders` | Orders | Owner | History pesanan difilter untuk lokasi outlet tertentu |
| `PATCH` | `/api/v1/orders/:id/status` | Orders | Owner | Mengubah status state pesanan: `pending`, `process`, `completed`, `picked_up` |
| `GET` | `/api/v1/reports/omzet` | Analytics | Owner | Menghitung total omzet (pendapatan kotor) |
| `GET` | `/api/v1/reports/orders/summary` | Analytics | Owner | Menghitung ringkasan jumlah pesanan per status |
| `GET` | `/api/v1/reports/services/top` | Analytics | Owner | Menampilkan top 5 layanan yang paling laris dipesan |

---

## 1. Authentication (Public)

Group rute ini dapat diakses secara publik (Tanpa token).

### Register User
Digunakan untuk mendaftarkan akun `owner` (pemilik outlet/kasir).
- **URL**: `POST /auth/register`
- **Body**:
  ```json
  {
    "name": "Budi Laundry",
    "phone": "+6281234567890",
    "email": "budi@email.com", 
    "password": "PasswordKuat123!",
    "role": "owner" // role standar untuk kasir adalah 'owner'
  }
  ```
- **Response**: Mengembalikan objek `User` (tanpa password) beserta `Token` JWT.

### Login User
- **URL**: `POST /auth/login`
- **Body**:
  ```json
  {
    "phone": "+6281234567890",
    "password": "PasswordKuat123!"
  }
  ```
- **Response**: Mengembalikan Data User + Token JWT. Bearer Token ini harus disertakan di `Authorization` header pada endpoint terproteksi (`Authorization: Bearer <token>`).

---

## 2. Outlet Management (Protected — Role: Owner)

Semua rute di bawah, mewajibkan JWT dengan Role `owner`. Segala aksi memiliki filter _Tenant Isolation_, di mana Owner hanya dapat memodifikasi/melihat outlet miliknya sendiri (Anti-IDOR).

### Bikin Outlet Baru
- **URL**: `POST /outlets`
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "name": "LaundryIn Cabang A",
    "address": "Jl. Kemerdekaan No. 10",
    "phone": "+628999999999"
  }
  ```

### Ambil Daftar Outlet
Mendukung Pagination dengan default page `1` limit `10`.
- **URL**: `GET /outlets?page=1&limit=10`
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Daftar outlet berhasil diambil",
    "data": {
      "data": [
        { "id": "...", "name": "LaundryIn Cabang A", "address": "...", "phone": "..." }
      ],
      "page": 1,
      "limit": 10,
      "total_count": 1
    }
  }
  ```

### Ambil Data Outlet Spesifik
- **URL**: `GET /outlets/:outletId`
- **Header**: `Authorization: Bearer <token>`

### Update Outlet
- **URL**: `PUT /outlets/:outletId`
- **Header**: `Authorization: Bearer <token>`
- **Body**: (Sama dengan registrasi/Create Outlet)

### Hapus Outlet
Menerapkan Soft-Delete. Data akan ditandai `deleted_at`, referensi FK masih akan tersimpan dengan aman (tidak langsung cascade hilang paksa di-runtime Gorm soft delete).
- **URL**: `DELETE /outlets/:outletId`
- **Header**: `Authorization: Bearer <token>`

---

## 3. Service Management (Protected — Role: Owner)

Manajemen layanan laundry/katalog berdasarkan Outlet. Di-filter dengan Zero-Trust (hanya bisa insert layanan ke outlet miliknya). Harga diparsing secara desimal.

### Bikin Layanan
- **URL**: `POST /services`
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "outlet_id": "<uuid_outlet>",
    "name": "Cuci Kering Ekstra",
    "price": "8000.50", // Kirim sebagai string agar tidak loss decimal
    "unit": "kg"
  }
  ```

### Ambil Semua Layanan per Outlet
- **URL**: `GET /outlets/:outletId/services`
- **Header**: `Authorization: Bearer <token>`

### Update Layanan
- **URL**: `PUT /services/:serviceId`
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "outlet_id": "<uuid_outlet>",
    "name": "Cuci Kering Saja",
    "price": "7500.00",
    "unit": "kg"
  }
  ```

### Hapus Layanan
- **URL**: `DELETE /services/:serviceId`
- **Header**: `Authorization: Bearer <token>`

---

## 4. Order Management (Protected — Role: Owner)

Sistem pembuatan struk pesanan & tracking status audit (OrderLog). Dilengkapi deep Anti-IDOR: Order takkan diproses bila outlet tak valid, atau product/service ID tidak berada pada outlet yang sama. Total harga dihitung melalui ACID Transaction & Desimal murni di backend (backend adalah Single Source of Truth).

### Buat Orderan Baru (Kasir)
- **URL**: `POST /orders`
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "outlet_id": "<uuid_outlet>",
    "items": [
      {
        "service_id": "<uuid_service>",
        "qty": "2.5" // String untuk akurasi presisi (2 setengah kg)
      }
    ]
  }
  ```
> **Catatan**: Frontend/Client tidak perlu mengirimkan subtotal/totalPrice. Server secara otomatis akan merujuk harga database untuk menjamin *Zero-Trust Pricing*. Saat Order berhasil dibuat, 1 log awal (state `""` menjadi `"pending"`) otomatis dibuat.

### Lihat History Semua Transaksi Order Kasir (Berdasarkan User)
Menarik semua transaksi milik owner (global dari seluruh outlet). Mendukung Pagination (`?page=1&limit=10`).
- **URL**: `GET /orders`
- **Header**: `Authorization: Bearer <token>`

### Lihat History Order per Outlet
Mendukung Pagination (`?page=1&limit=10`).
- **URL**: `GET /outlets/:outletId/orders`
- **Header**: `Authorization: Bearer <token>`

### Update Status Order
Digunakan secara _Finite State Machine_ linear: `pending` -> `process` -> `completed` -> `picked_up`
Atau `pending` -> `cancelled`. Transisi status otomatis diverifikasi dan direkam ke dalam `order_log`.
- **URL**: `PATCH /orders/:orderId/status`
- **Header**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "status": "process"
  }
  ```

---

## 5. Reports & Analytics (Protected — Role: Owner)

Menampilkan data analitik dan reporting untuk kebutuhan dashboard bisnis aplikasi. Semua rute ini dilengkapi **Deep Anti-IDOR** berdasarkan parameter otomatis dari `JWT Token` Backend, dan tabel *Relationship Contextual*. Kasir Owner sama sekali tidak bisa melihat omzet atau laporan dari lapak kompetitornya (Owner silang).

Setiap endpoint di bawah mendukung parameter Filter Query API (opsional) berikut:
- `start_date` (format *YYYY-MM-DD*): Digunakan untuk melakukan filter perhitungan mulai dari batas bawah tanggal tertentu.
- `end_date` (format *YYYY-MM-DD*): Digunakan untuk melakukan filter perhitungan maksimal hingga hari/tanggal tersebut.
- `outlet_id` (format UUID Valid): Digunakan untuk melihat performa tersekat murni spesifik di 1 cabang bangunan outlet saja.

### Dapatkan Total Omzet
Menghitung kalkulasi bersih dari total pendapatan yang mana status transaksinya **BUKAN** `cancelled`.
- **URL**: `GET /reports/omzet` (contoh filter query lengkap: `GET /reports/omzet?start_date=2023-01-01&end_date=2025-12-31&outlet_id=xxxxx`)
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Berhasil mengambil data omzet",
    "data": {
      "total_omzet": "150000.50"
    }
  }
  ```

### Dapatkan Ringkasan Status Pesanan
Menampilkan grafik jumlah agregasi FSM (*Finite State Machine*) berjenjang berdasarkan status map pesanan.
- **URL**: `GET /reports/orders/summary`
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Berhasil mengambil ringkasan status pesanan",
    "data": {
      "pending": 5,
      "process": 2,
      "completed": 1,
      "picked_up": 10,
      "cancelled": 0
    }
  }
  ```

### Dapatkan Top 5 Layanan Terlaris
Menghitung agregasi dari kuantitas volume penjualan dan hasil total pendapatan yang diseleksi dari histori pesanan. \n*(Keamanan Arsitektur: Data di Query murni menggunakan logika Snapshot item order (`oi.service_name`), sehingga apabila data master layanan Services di-Soft-Delete oleh Owner sekalipun, rekap kalkulasinya akan selalu tetap utuh terarsip murni di grafik Analytics historis).*
- **URL**: `GET /reports/services/top`
- **Header**: `Authorization: Bearer <token>`
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Berhasil mengambil data layanan terlaris",
    "data": [
      {
        "service_name": "Cuci Kering Ekstra",
        "outlet_name": "LaundryIn Cabang A",
        "total_qty": "100.50",
        "total_revenue": "80000.00"
      },
      {
        // 2nd top...
      }
    ]
  }
  ```
