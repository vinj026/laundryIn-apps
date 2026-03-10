# Implementation Plan: API Security & Logic Validation (Phase 4)

**Status:** Phase 4 (Service Management) Integration Testing
**Security Principle:** Zero-Trust & Anti-IDOR (Insecure Direct Object Reference)
**Standard Response:** All unauthorized/cross-tenant access MUST return `404 Not Found` (to prevent resource enumeration).

## 1. Lingkungan & State Awal (Prerequisites)
Pastikan dua user berikut tersedia di database (atau buat via Register):
- **Actor A (Owner A):** phone: `+6281111111111` (Punya Outlet A)
- **Actor B (Owner B):** phone: `+6282222222222` (Punya Outlet B)
- **Token:** Dapatkan JWT Token untuk masing-masing melalui endpoint `/api/v1/auth/login`.

## 2. Matriks Pengujian Unit (Validation Check)
Gunakan Token Owner A untuk pengujian berikut:

| ID | Skenario Uji | Payload | Expected Result | Target Validasi |
|----|--------------|---------|-----------------|-----------------|
| VAL-1 | Negative Price | `price: -1000` | 400 Bad Request | DTO Tag `gt=0` |
| VAL-2 | Invalid Enum Unit | `unit: "LITER"` | 400 Bad Request | DTO Tag `oneof=KG PCS METER` |
| VAL-3 | Empty Name | `name: ""` | 400 Bad Request | DTO Tag `required` |
| VAL-4 | Invalid Outlet UUID | `outlet_id: "123-abc"` | 400 Bad Request | DTO Tag `uuid` |

## 3. Matriks Pengujian Security (Anti-IDOR Check)
Gunakan skenario Cross-Tenant Attack untuk memastikan SQL Join/Ownership logic bekerja:

| ID | Aktor | Action | Target Resource | Expected |
|----|-------|--------|-----------------|----------|
| SEC-1 | Owner B | POST `/services` | outlet_id milik Owner A | 404 Not Found |
| SEC-2 | Owner B | PUT `/services/:id` | service_id milik Owner A | 404 Not Found |
| SEC-3 | Owner B | DELETE `/services/:id` | service_id milik Owner A | 404 Not Found |
| SEC-4 | Owner B | GET `/outlets/:id/services` | :id outlet milik Owner A | 404 Not Found |

## 4. Requirement Implementasi untuk AI Lain
Berikan instruksi ini kepada model pelaksana (Claude):
- **Repository Pattern:** Pencarian Service wajib menggunakan SQL JOIN dengan tabel outlets untuk memverifikasi `outlets.user_id = current_user_id`.
- **Context Management:** Setiap query wajib dibungkus `WithContext(ctx)` dan menangani `ctx.Err()` untuk mengembalikan `408 Request Timeout` jika eksekusi melebihi limit (5s).
- **No Leaky Errors:** Jangan pernah mengembalikan internal error database (Postgres error) ke user. Gunakan Generic Error Message.
- **Transaction Safety:** Untuk Phase 5 nanti, pastikan operasi mutasi (Order) menggunakan `db.Transaction`.

## 5. Script Otomasi (Checklist)
Jika menggunakan Bash/CURL:
- [ ] Verify `Content-Type: application/json` is enforced.
- [ ] Verify `Authorization: Bearer <token>` is not empty.
- [ ] Verify `X-Content-Type-Options: nosniff` header is present in response.
