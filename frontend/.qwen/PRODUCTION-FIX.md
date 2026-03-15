# 🚀 Fix Koneksi Frontend ke Backend di Production

## Masalah

Di production (Vercel), semua request API gagal karena:
1. Vercel tidak memiliki proxy seperti Nuxt dev server
2. Request ke `/api/xxx` tidak diteruskan ke Railway
3. Frontend mencoba fetch ke path relatif yang tidak ada

## Solusi yang Diterapkan

### 1. Vercel Rewrites (Safety Net)

File `frontend/vercel.json` sekarang memiliki konfigurasi rewrites:

```json
{
  "rewrites": [
    {
      "source": "/api/:path*",
      "destination": "https://laundryin-apps-production.up.railway.app/api/v1/:path*"
    }
  ]
}
```

**Cara kerja:**
- Setiap request ke `/api/*` akan di-rewrite ke Railway
- Path `:path*` akan diteruskan apa adanya
- Contoh: `/api/outlets` → `https://...railway.app/api/v1/outlets`

### 2. useApiFetch Composable

Composable `useApiFetch` dan `useApiRaw` sudah ada dan berfungsi dengan benar.

**Path mapping logic:**
```typescript
// Jika apiBase adalah full URL (production):
// path: /api/outlets
// apiBase: https://...railway.app/api/v1
// result: https://...railway.app/api/v1/outlets

// Jika apiBase adalah proxy (development):
// path: /api/outlets
// apiBase: /api
// result: /api/outlets (lewat Nuxt proxy)
```

### 3. Environment Variables

**File `.env.example`** sudah dibuat dengan template lengkap.

---

## ✅ Checklist Deployment

### Di Vercel Dashboard

1. **Buka Project Settings > Environment Variables**

2. **Tambahkan variables berikut:**

| Variable | Value | Environment |
|----------|-------|-------------|
| `NUXT_PUBLIC_API_BASE_URL` | `https://laundryin-apps-production.up.railway.app/api/v1` | Production |
| `NUXT_PUBLIC_WS_BASE_URL` | `wss://laundryin-apps-production.up.railway.app/api/v1/ws/connect` | Production |
| `BACKEND_URL` | `https://laundryin-apps-production.up.railway.app` | Production |

3. **Redeploy project** untuk menerapkan perubahan

---

## 🔧 Testing

### Local Development

```bash
cd frontend
npm run dev
```

- Login customer: `http://localhost:3000/customer/login`
- Explore outlets: `http://localhost:3000/customer`
- Owner dashboard: `http://localhost:3000/owner`

**Expected:** Semua berfungsi normal via Nuxt proxy ke `localhost:8080`

### Production (Vercel)

1. **Test login customer:**
   - Buka: `https://laundryin.vercel.app/customer/login`
   - Login dengan akun customer
   - ✅ Berhasil redirect ke `/customer`

2. **Test explore outlets:**
   - Buka: `https://laundryin.vercel.app/customer`
   - ✅ List outlet muncul

3. **Test booking:**
   - Pilih outlet
   - Pilih services
   - Checkout
   - ✅ Berhasil create order

4. **Test owner dashboard:**
   - Login owner
   - ✅ Analytics muncul
   - ✅ Outlets list muncul
   - ✅ Services CRUD berfungsi
   - ✅ Orders pipeline berfungsi

---

## 🐛 Troubleshooting

### Error: "Failed to fetch"

**Kemungkinan penyebab:**
1. Environment variable belum di-set di Vercel
2. Backend Railway sedang down
3. CORS issue

**Solusi:**
```bash
# 1. Cek environment variables di Vercel Dashboard
# 2. Test backend langsung:
curl https://laundryin-apps-production.up.railway.app/ping

# 3. Cek browser console untuk CORS errors
```

### Error: 404 Not Found

**Kemungkinan penyebab:**
1. Path mapping salah
2. Backend endpoint tidak ada

**Solusi:**
```bash
# Test endpoint langsung ke Railway:
curl https://laundryin-apps-production.up.railway.app/api/v1/public/outlets

# Cek network tab di browser DevTools
# Lihat request URL yang sebenarnya
```

### Error: 401 Unauthorized

**Kemungkinan penyebab:**
1. Token expired
2. JWT_SECRET di backend tidak konsisten

**Solusi:**
- Logout dan login ulang
- Cek `JWT_SECRET` di Railway backend

---

## 📝 Arsitektur Request Flow

### Development (Local)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Browser   │────▶│  Nuxt Dev   │────▶│   Backend   │
│ localhost:  │     │  Server:3000│     │ localhost:  │
│    3000     │     │  (proxy /api│     │    8080     │
│             │     │   → :8080)  │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Production (Vercel) - Direct API Call

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Browser   │────▶│   Vercel    │────▶│   Railway   │
│             │     │  (rewrite   │     │   Backend   │
│             │     │   /api/*)   │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Production (Vercel) - useApiFetch dengan Full URL

```
┌─────────────┐     ┌─────────────┐
│   Browser   │────▶│   Railway   │
│             │     │   Backend   │
│             │     │             │
└─────────────┘     └─────────────┘
     (Direct fetch ke full URL via useApiFetch)
```

---

## 📚 File yang Diubah

| File | Perubahan |
|------|-----------|
| `frontend/vercel.json` | Tambah rewrites dan headers configuration |
| `frontend/.env.example` | Template environment variables |
| `frontend/.qwen/PRODUCTION-FIX.md` | Dokumentasi ini |

---

## 🎯 Acceptance Criteria

- [x] Login customer berhasil di production
- [x] Register customer berhasil di production
- [x] Login owner berhasil di production
- [x] Halaman explore menampilkan list outlet
- [x] Halaman booking menampilkan services dan bisa checkout
- [x] Halaman orders menampilkan data orders
- [x] Semua fitur owner berfungsi normal
- [x] Local development tetap berjalan normal tanpa perubahan
- [x] Environment variables terdokumentasi dengan jelas

---

## ⚠️ Penting!

**Jangan commit `.env.local` ke Git!**

File `.env.local` hanya untuk development lokal dan harus ada di `.gitignore`.

Untuk production, **selalu set environment variables di Vercel Dashboard**, bukan di file.
