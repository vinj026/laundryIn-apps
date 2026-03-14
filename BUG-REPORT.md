# Bug Report — LaundryIn Production

**Tanggal audit:** 14 Maret 2026  
**Environment:** Production (Vercel + Railway)  
**Auditor:** QA Engineer

---

## Daftar Bug

### BUG-001 — WebSocket URL Hardcoded ke `ws://` (Bukan `wss://`)

**Severity:** Critical  
**File:** `frontend/nuxt.config.ts`  
**Line:** 8

**Deskripsi:**
WebSocket base URL di `nuxt.config.ts` menggunakan `ws://` (non-secure) sebagai default value:
```typescript
wsBase: process.env.NUXT_PUBLIC_WS_BASE_URL || 'ws://localhost:8080/api/v1/ws/connect'
```

Di production (Vercel), browser akan memblokir koneksi WebSocket dari HTTPS ke `ws://` karena **mixed content policy**. Browser modern tidak mengizinkan koneksi insecure dari halaman secure.

**Root Cause:**
- Default value hardcoded `ws://localhost:8080`
- Tidak ada fallback ke `wss://` untuk production
- Environment variable `NUXT_PUBLIC_WS_BASE_URL` mungkin tidak di-set di Vercel

**Dampak:**
- **WebSocket tidak akan connect sama sekali di production**
- Notifikasi real-time tidak akan masuk ke user
- User tidak mendapat notifikasi order baru, status change, atau price update
- Owner tidak tahu ada order masuk
- Customer tidak tahu status orderannya berubah

**Fix yang Dibutuhkan:**

1. **Update `nuxt.config.ts`:**
```typescript
runtimeConfig: {
  public: {
    // Gunakan https/wss untuk production, http/ws untuk local
    apiBase: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1',
    wsBase: process.env.NUXT_PUBLIC_WS_BASE_URL || 
            (process.env.VERCEL ? 'wss://laundryin-backend-production.up.railway.app/api/v1/ws/connect' : 'ws://localhost:8080/api/v1/ws/connect')
  }
}
```

2. **Atau lebih baik, buat composable yang dinamis:**
```typescript
// composables/useRuntimeConfig.ts
export const getWebSocketUrl = () => {
  if (process.server) return null // Skip WS on server
  if (process.env.NUXT_PUBLIC_WS_BASE_URL) return process.env.NUXT_PUBLIC_WS_BASE_URL
  
  // Production fallback
  if (window.location.hostname !== 'localhost') {
    const backendHost = 'laundryin-backend-production.up.railway.app'
    return `wss://${backendHost}/api/v1/ws/connect`
  }
  
  // Local fallback
  return 'ws://localhost:8080/api/v1/ws/connect'
}
```

3. **Set environment variable di Vercel:**
```
NUXT_PUBLIC_WS_BASE_URL=wss://laundryin-backend-production.up.railway.app/api/v1/ws/connect
```

---

### BUG-002 — Proxy URL Hardcoded ke localhost

**Severity:** Critical  
**File:** `frontend/nuxt.config.ts`  
**Line:** 13

**Deskripsi:**
Proxy route rule hardcoded ke localhost:
```typescript
routeRules: {
  '/api/**': { proxy: 'http://localhost:8080/api/v1/**' }
}
```

Di production (Vercel), proxy ini akan mencoba connect ke `localhost:8080` yang tidak ada, menyebabkan **semua API call gagal**.

**Root Cause:**
- Hardcoded localhost tanpa conditional check untuk production
- Tidak menggunakan environment variable untuk backend URL

**Dampak:**
- **Semua API call akan gagal di production**
- Login/register tidak akan bekerja
- Tidak bisa fetch outlets, services, orders
- App tidak bisa digunakan sama sekali

**Fix yang Dibutuhkan:**

**Option 1 — Gunakan runtimeConfig (Recommended):**
```typescript
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1'
    }
  },
  // Hapus routeRules proxy, gunakan useApiFetch composable instead
})
```

**Option 2 — Conditional proxy berdasarkan environment:**
```typescript
const isProduction = process.env.VERCEL || process.env.NODE_ENV === 'production'
const backendUrl = isProduction 
  ? 'https://laundryin-backend-production.up.railway.app/api/v1'
  : 'http://localhost:8080/api/v1'

export default defineNuxtConfig({
  routeRules: {
    '/api/**': { proxy: `${backendUrl}/**` }
  }
})
```

**Option 3 — Set di Vercel environment variables:**
```
NUXT_PUBLIC_API_BASE_URL=https://laundryin-backend-production.up.railway.app/api/v1
```

---

### BUG-003 — Database SSL Mode Disable untuk Production

**Severity:** Critical  
**File:** `backend/.env`  
**Line:** 6

**Deskripsi:**
```bash
DB_SSLMODE=disable
```

Railway PostgreSQL **require SSL connection** untuk production. Dengan `sslmode=disable`, koneksi ke database akan ditolak.

**Root Cause:**
- `.env` file menggunakan setting development untuk production
- Tidak ada conditional SSL mode berdasarkan environment
- Railway menolak non-SSL connection untuk security

**Dampak:**
- **Backend tidak bisa connect ke database di Railway**
- Semua API endpoint akan return 500 error
- App crash saat startup
- Error log: "no pg_hba.conf entry for host" atau "SSL connection required"

**Fix yang Dibutuhkan:**

1. **Update `backend/internal/database/postgres.go` untuk auto-detect SSL:**
```go
func ConnectDB() *gorm.DB {
    // ... existing code ...
    
    // Auto-detect Railway environment
    isRailway := os.Getenv("RAILWAY_ENVIRONMENT") != "" || 
                 os.Getenv("RAILWAY_STATIC_URL") != "" ||
                 os.Getenv("DATABASE_URL") != ""
    
    if sslMode == "" {
        if isRailway {
            sslMode = "require" // Railway requires SSL
        } else {
            sslMode = "disable" // Local development
        }
    }
    
    // ... rest of code ...
}
```

2. **Atau set environment variable di Railway:**
```
DB_SSLMODE=require
```

3. **Better: Gunakan DATABASE_URL langsung (Railway menyediakan):**
```
DATABASE_URL=postgresql://user:pass@host:port/dbname?sslmode=require
```

---

### BUG-004 — JWT_SECRET Hardcoded di .env

**Severity:** Critical  
**File:** `backend/.env`  
**Line:** 8

**Deskripsi:**
```bash
JWT_SECRET=rahasia_negara_lu
```

Secret JWT hardcoded di file `.env` yang mungkin ter-commit ke Git atau digunakan sama untuk semua environment (development, staging, production).

**Root Cause:**
- Secret yang sama untuk semua environment
- Potensi secret terekspos jika repo public atau .env ter-upload
- Tidak ada rotation mechanism

**Dampak:**
- **Security vulnerability** — jika secret bocor, attacker bisa forge JWT token
- Bisa impersonate user manapun (termasuk owner)
- Bisa akses semua data sensitif
- Token tidak bisa di-invalidate tanpa ganti secret (breaking change)

**Fix yang Dibutuhkan:**

1. **Generate unique secret untuk production:**
```bash
# Generate random 32-byte secret
openssl rand -base64 32
# Output contoh: x7dK9pL2mN4qR6sT8vW0yZ1aB3cD5eF7gH9iJ0kL2mN4=
```

2. **Set di Railway environment variables (JANGAN commit ke Git):**
```
JWT_SECRET=<random-generated-secret>
```

3. **Update .env.example (untuk dokumentasi):**
```bash
# .env.example
JWT_SECRET=change-this-in-production-use-openssl-rand-base64-32
```

4. **Add .env ke .gitignore:**
```bash
# .gitignore
.env
.env.local
.env.production
```

---

### BUG-005 — CORS Allow Origin Wildcard (`*`)

**Severity:** High  
**File:** `backend/internal/delivery/http/middleware.go`  
**Line:** 145

**Deskripsi:**
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        // ...
    }
}
```

Wildcard CORS (`*`) mengizinkan **semua website** untuk make API kita. Ini security risk karena website malicious bisa make API kita dengan credentials user.

**Root Cause:**
- Development convenience (allow semua origin untuk testing)
- Tidak ada validasi Origin header
- Comment di code bilang "In production, you should ideally restrict" tapi tidak di-implement

**Dampak:**
- **Security vulnerability** — website lain bisa akses API kita
- CSRF attack potential (walaupun ada JWT, tetap risk)
- API abuse dari domain yang tidak authorized
- Potential data leak jika ada bug di auth

**Fix yang Dibutuhkan:**

```go
func CORSMiddleware() gin.HandlerFunc {
    allowedOrigins := map[string]bool{
        "https://laundryin.vercel.app":     true,
        "https://www.laundryin.vercel.app": true,
        "http://localhost:3000":            true, // Development
        "http://localhost:3001":            true, // Development
    }
    
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        
        // Check if origin is allowed
        if !allowedOrigins[origin] {
            // For production, reject unknown origins
            if os.Getenv("GIN_MODE") == "release" {
                c.AbortWithStatus(http.StatusForbidden)
                return
            }
            // Development fallback
            origin = "http://localhost:3000"
        }
        
        c.Header("Access-Control-Allow-Origin", origin)
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}
```

---

### BUG-006 — WebSocket CheckOrigin Return True (No Validation)

**Severity:** High  
**File:** `backend/internal/delivery/http/notification_handler.go`  
**Line:** 27-33

**Deskripsi:**
```go
var upgrader = gorilla.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        if origin == "" {
            return true
        }
        return true // Always allow!
    },
}
```

WebSocket upgrader menerima koneksi dari **origin manapun**. Ini bisa menyebabkan CSRF WebSocket attack.

**Root Cause:**
- Development convenience
- Comment bilang "For now, allow any origin" tapi tidak pernah di-fix
- Tidak ada whitelist untuk production

**Dampak:**
- **Security vulnerability** — website malicious bisa establish WebSocket connection
- Real-time notification bisa di-intercept atau di-inject
- Potential data leak

**Fix yang Dibutuhkan:**

```go
var allowedOrigins = map[string]bool{
    "https://laundryin.vercel.app":     true,
    "https://www.laundryin.vercel.app": true,
    "http://localhost:3000":            true,
    "http://localhost:3001":            true,
}

var upgrader = gorilla.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        if origin == "" {
            return true // Allow requests with no origin (mobile apps, etc.)
        }
        
        // Check against whitelist
        if allowedOrigins[origin] {
            return true
        }
        
        // Log for debugging
        log.Printf("WebSocket connection rejected from origin: %s", origin)
        return false
    },
}
```

---

### BUG-007 — Token Expired Tidak Ditangani (401 Handling)

**Severity:** High  
**File:** `frontend/app/stores/auth.ts`, `frontend/app/composables/useApiFetch.ts`  
**Line:** N/A (multiple files)

**Deskripsi:**
Tidak ada logic untuk handle JWT token expired (401 response). Ketika token expired (setelah 24 jam), user akan stuck dalam state "logged in" tapi semua API call gagal.

**Root Cause:**
- Tidak ada interceptor untuk catch 401 response
- Tidak ada automatic token refresh mechanism
- Tidak ada redirect ke login page saat session expired

**Dampak:**
- **User experience broken** — user klik tombol tapi tidak ada response
- User tidak tahu session sudah expired
- Data tidak ter-load tapi tidak ada error message
- User harus manual clear localStorage dan login ulang

**Fix yang Dibutuhkan:**

1. **Create Nuxt plugin untuk intercept 401:**
```typescript
// frontend/plugins/api-error-handler.client.ts
export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()
  const router = useRouter()
  const { error: toastError } = useToast()
  
  // Intercept $fetch errors
  const originalFetch = globalThis.$fetch
  globalThis.$fetch = async (request, options = {}) => {
    try {
      return await originalFetch(request, options)
    } catch (error) {
      if (error.status === 401 || error.statusCode === 401) {
        // Token expired or invalid
        authStore.logout()
        toastError('Sesi kamu telah kadaluarsa, silakan login ulang')
        
        // Redirect to login with current page as redirect param
        const currentPath = window.location.pathname
        router.push(`/customer/login?redirect=${encodeURIComponent(currentPath)}`)
      }
      throw error
    }
  }
})
```

2. **Add 401 handling di useApiFetch:**
```typescript
export const useApiFetch = (path: string | (() => string), options: any = {}) => {
  const authStore = useAuthStore()
  const router = useRouter()
  const { error: toastError } = useToast()
  
  const fetchOptions = {
    ...options,
    onResponseError: async ({ response }) => {
      if (response.status === 401) {
        authStore.logout()
        toastError('Sesi kamu telah kadaluarsa')
        router.push('/customer/login')
      }
    }
  }
  
  return useFetch(path, fetchOptions)
}
```

---

### BUG-008 — Missing Error Handling di $fetch Calls

**Severity:** High  
**File:** Multiple files (pages/*.vue)

**Deskripsi:**
Banyak API call menggunakan `useFetch` atau `$fetch` tanpa proper error handling:

```typescript
// frontend/app/pages/customer/index.vue
const { data: outletsResponse, pending, refresh } = await useFetch<{ data: { data: Outlet[] } }>('/api/public/outlets')
// No error handling!
```

**Root Cause:**
- Relying on Nuxt's default error handling
- Tidak ada try/catch atau onError callback
- Error state tidak ditampilkan ke user dengan proper message

**Dampak:**
- User tidak tahu kenapa data tidak muncul
- Silent failures — app looks broken tapi no error message
- Hard to debug production issues

**Fix yang Dibutuhkan:**

```typescript
// Example fix for customer/index.vue
const { data: outletsResponse, pending, error, refresh } = await useFetch<{ data: { data: Outlet[] } }>('/api/public/outlets', {
  onError: (err) => {
    console.error('Failed to fetch outlets:', err)
    toastError('Gagal memuat daftar outlet, silakan coba lagi')
  }
})

// Or with useApiFetch composable
const { data, pending, error } = await useApiFetch('/api/public/outlets')

watch(error, (err) => {
  if (err) {
    toastError(err.message || 'Gagal memuat data')
  }
})
```

---

### BUG-009 — Console.log di Production Code

**Severity:** Low  
**File:** `frontend/app/stores/notification.ts`  
**Line:** 24, 27, 31, 32, 47, 50, 55, 56, 57

**Deskripsi:**
Banyak `console.log` statements yang masih aktif di production:
```typescript
console.log('addNotification received:', notif)
console.log('addNotification rejected: duplicate ID', notif.id)
console.log('addNotification successful. New count:', this.notifications.length)
console.log('Notification API response:', res)
console.log('Extracted fetchedNotifs:', fetchedNotifs)
```

**Root Cause:**
- Debugging statements tidak di-remove sebelum deploy
- Tidak ada build-time stripping untuk console.log

**Dampak:**
- **Performance impact** — unnecessary logging di production
- **Security risk** — sensitive data (notifications, user info) ter-expose di console
- **Noisy logs** — hard to find real issues

**Fix yang Dibutuhkan:**

1. **Remove console.log di production:**
```typescript
// Replace with conditional logging
if (process.env.NODE_ENV === 'development') {
  console.log('addNotification received:', notif)
}
```

2. **Or use a logger utility:**
```typescript
// composables/useLogger.ts
export const useLogger = () => {
  const isDev = process.env.NODE_ENV === 'development'
  
  const log = (...args: any[]) => {
    if (isDev) console.log(...args)
  }
  
  const error = (...args: any[]) => {
    if (isDev) console.error(...args)
  }
  
  return { log, error }
}

// Usage in store
const { log, error } = useLogger()
log('addNotification received:', notif)
```

3. **Build-time stripping (Vite config):**
```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  vite: {
    define: {
      'console.log': '() => {}'
    }
  }
})
```

---

### BUG-010 — process.client Tidak Didefinisikan di Beberapa Tempat

**Severity:** Medium  
**File:** `frontend/app/stores/auth.ts`  
**Line:** 30, 38, 46

**Deskripsi:**
```typescript
if (process.client) {
  localStorage.setItem('laundryin_token', token)
}
```

Di Nuxt 3, `process.client` mungkin tidak selalu terdefinisi dengan benar, terutama di edge cases atau certain build configurations.

**Root Cause:**
- Nuxt 2 vs Nuxt 3 API difference
- `process.client` deprecated di Nuxt 3
- Should use `import.meta.client` instead

**Dampak:**
- **Potential SSR/SSG issues** — localStorage access during server-side rendering
- Build warnings atau errors
- Runtime errors di certain deployment configurations

**Fix yang Dibutuhkan:**

```typescript
// Replace all process.client with import.meta.client
export const useAuthStore = defineStore('auth', {
  actions: {
    setAuth(token: string, user: User) {
      this.token = token
      this.user = user
      if (import.meta.client) {
        localStorage.setItem('laundryin_token', token)
        localStorage.setItem('laundryin_user', JSON.stringify(user))
      }
    },
    // ... rest of actions
  }
})
```

---

### BUG-011 — Missing Owner Register Page Implementation

**Severity:** Medium  
**File:** `frontend/app/pages/owner/register.vue`  
**Line:** N/A (file exists but may not be linked)

**Deskripsi:**
File `owner/register.vue` ada di codebase tapi tidak ada link ke halaman ini dari manapun. Owner registration flow tidak jelas apakah sudah implemented atau belum.

**Root Cause:**
- Incomplete feature implementation
- Missing navigation links
- Possibly placeholder file

**Dampak:**
- **User confusion** — owner tidak bisa register sendiri
- Owner harus dibuat manual di database
- Inconsistent UX (customer bisa register, owner tidak)

**Fix yang Dibutuhkan:**

1. **Verify owner/register.vue implementation:**
```vue
<!-- Check if file has proper implementation -->
<script setup lang="ts">
// Should have similar logic to customer/register.vue
// But with role: 'owner' hardcoded
</script>
```

2. **Add link from owner/login.vue:**
```vue
<!-- frontend/app/pages/owner/login.vue -->
<NuxtLink to="/owner/register" class="...">
  Don't have an account? <span class="text-primary font-semibold">Register</span>
</NuxtLink>
```

3. **Or remove file jika tidak diperlukan:**
```bash
# Jika owner registration memang harus manual
rm frontend/app/pages/owner/register.vue
```

---

### BUG-012 — Tracking Page Menggunakan Static/Hardcoded Data

**Severity:** Medium  
**File:** `frontend/app/pages/customer/tracking.vue`  
**Line:** N/A (entire file)

**Deskripsi:**
Halaman `/customer/tracking` menggunakan static data yang di-hardcode, bukan fetch dari API:
```vue
<h1 class="text-3xl font-bold mb-1">Order #LND-8821</h1>
<p class="text-sm text-surface-onSurfaceVariant">
  Estimated Completion: <span class="font-semibold text-primary">Today, 5:00 PM</span>
</p>
```

**Root Cause:**
- Demo/placeholder page belum di-implement dengan real API
- Missing route parameter untuk order ID
- Tidak ada API call untuk fetch order details

**Dampak:**
- **Broken feature** — user tidak bisa tracking order real
- Confusing UX — user expect real tracking data
- Wasted development effort (page exists but useless)

**Fix yang Dibutuhkan:**

1. **Convert to dynamic route:**
```
frontend/app/pages/customer/tracking.vue
→ frontend/app/pages/customer/tracking/[id].vue
```

2. **Fetch real order data:**
```typescript
const route = useRoute()
const orderId = route.params.id

const { data: order } = await useApiFetch(`/api/orders/${orderId}`)
```

3. **Or redirect to existing order detail page:**
```typescript
// If /customer/orders/[id] already has tracking info
navigateTo(`/customer/orders/${orderId}`)
```

4. **Or remove page jika tidak diperlukan:**
```bash
rm frontend/app/pages/customer/tracking.vue
```

---

### BUG-013 — Missing Loading State di Beberapa Halaman

**Severity:** Medium  
**File:** Multiple files

**Deskripsi:**
Beberapa halaman tidak memiliki proper loading state atau loading state tidak konsisten:

```typescript
// Some pages have:
<div v-if="pending" class="...">Loading...</div>

// But others might not handle loading properly
```

**Root Cause:**
- Inconsistent implementation across pages
- Missing loading state for secondary data fetches

**Dampak:**
- **Poor UX** — user tidak tahu apakah data sedang loading
- Perceived as broken/slow
- User might click multiple times thinking nothing happened

**Fix yang Dibutuhkan:**

Add consistent loading state di semua halaman yang fetch data:

```vue
<template>
  <div v-if="pending" class="flex justify-center py-20">
    <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
  </div>
  
  <div v-else-if="error" class="...">
    <!-- Error state -->
  </div>
  
  <div v-else-if="data" class="...">
    <!-- Data content -->
  </div>
  
  <div v-else class="...">
    <!-- Empty state -->
  </div>
</template>
```

---

### BUG-014 — Notification Store Console.error Tidak Ditampilkan ke User

**Severity:** Low  
**File:** `frontend/app/stores/notification.ts`  
**Line:** 59, 74, 87, 100

**Deskripsi:**
```typescript
catch (err) {
  console.error('Failed to fetch notifications', err)
}
```

Error hanya di-log ke console, tidak ditampilkan ke user dengan toast atau UI feedback.

**Root Cause:**
- Silent error handling
- Assuming notifications are "nice to have" not critical

**Dampak:**
- User tidak tahu notifikasi gagal di-fetch
- Silent failures
- Hard to debug

**Fix yang Dibutuhkan:**

```typescript
async fetchNotifications(page = 1, limit = 20) {
  const authStore = useAuthStore()
  const { error: toastError } = useToast()
  
  if (!authStore.isLoggedIn) return

  this.loading = true
  try {
    const res = await useApiRaw<ApiResponse<any>>('/api/notifications', {
      params: { page, limit }
    })
    // ... existing logic
  } catch (err) {
    console.error('Failed to fetch notifications', err)
    // Don't toast every time, but maybe show indicator
    // toastError('Gagal memuat notifikasi')
  } finally {
    this.loading = false
  }
}
```

---

### BUG-015 — WebSocket Reconnect Delay Too Aggressive

**Severity:** Low  
**File:** `frontend/app/composables/useWebSocket.ts`  
**Line:** 13-14

**Deskripsi:**
```typescript
let reconnectDelay = 1000
const MAX_RECONNECT_DELAY = 30000
```

Initial reconnect delay 1 second dengan exponential backoff mungkin terlalu agresif untuk mobile networks atau unstable connections.

**Root Cause:**
- Default values tanpa consideration untuk production conditions
- No jitter/randomization

**Dampak:**
- **Battery drain** on mobile devices
- **Network congestion** with many simultaneous reconnects
- Server might rate-limit aggressive reconnects

**Fix yang Dibutuhkan:**

```typescript
let reconnectDelay = 2000 // Start with 2s
const MAX_RECONNECT_DELAY = 60000 // Max 1 minute
const RECONNECT_JITTER = 1000 // Add randomness

const scheduleReconnect = () => {
  if (reconnectTimer) clearTimeout(reconnectTimer)
  
  // Add jitter to prevent thundering herd
  const delayWithJitter = reconnectDelay + Math.random() * RECONNECT_JITTER
  
  reconnectTimer = setTimeout(() => {
    reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY)
    connect()
  }, delayWithJitter)
}
```

---

### BUG-016 — Missing Build/Deploy Configuration Files

**Severity:** Medium  
**File:** N/A (missing files)

**Deskripsi:**
Tidak ada configuration files untuk deployment:
- No `railway.toml` atau `Dockerfile` untuk backend
- No `vercel.json` untuk frontend (selain nuxt.config.ts)
- No Procfile

**Root Cause:**
- Relying on auto-detection
- Development-first approach

**Dampak:**
- **Deployment uncertainty** — relying on platform auto-detection
- **Build failures** if platform changes detection logic
- **No control** over build process

**Fix yang Dibutuhkan:**

1. **Create `backend/Dockerfile`:**
```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/api .

EXPOSE 8080
CMD ["./api"]
```

2. **Create `backend/railway.toml`:**
```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "Dockerfile"

[deploy]
startCommand = "./api"
healthcheckPath = "/ping"
healthcheckTimeout = 100
restartPolicyType = "ON_FAILURE"
```

3. **Create `frontend/vercel.json`:**
```json
{
  "buildCommand": "npm run build",
  "devCommand": "npm run dev",
  "installCommand": "npm install",
  "framework": "nuxtjs",
  "outputDirectory": ".output"
}
```

---

### BUG-017 — Database Migration Error Handling

**Severity:** Medium  
**File:** `backend/cmd/api/main.go`  
**Line:** 27-31

**Deskripsi:**
```go
err := db.AutoMigrate(&models.User{}, &models.Outlet{}, ...)
if err != nil {
  log.Printf("⚠️  Gagal migrasi database: %v", err)
} else {
  fmt.Println("🚀 Database Migration Successful!")
}
```

Migration error hanya di-log, app tetap continue. Ini bisa menyebabkan app running dengan schema yang tidak lengkap.

**Root Cause:**
- Lenient error handling
- Wanting app to start even if migration fails

**Dampak:**
- **Data corruption** — app running with incomplete schema
- **Silent failures** — tables might not exist
- **Hard to debug** — error only in logs, not obvious

**Fix yang Dibutuhkan:**

```go
err := db.AutoMigrate(&models.User{}, &models.Outlet{}, ...)
if err != nil {
  log.Fatalf("❌ CRITICAL: Database migration failed: %v", err)
  // App should NOT start if migration fails
}
fmt.Println("🚀 Database Migration Successful!")
```

---

### BUG-018 — Rate Limiter Memory Leak Potential

**Severity:** Low  
**File:** `backend/internal/delivery/http/middleware.go`  
**Line:** 93-103

**Deskripsi:**
```go
go func() {
  for {
    time.Sleep(1 * time.Minute)
    mu.Lock()
    for ip, v := range visitors {
      if time.Since(v.lastSeen) > 3*time.Minute {
        delete(visitors, ip)
      }
    }
    mu.Unlock()
  }
}()
```

Background goroutine untuk cleanup visitors ada, tapi:
- Tidak ada context untuk graceful shutdown
- Goroutine akan run forever (potential leak on hot reload)

**Root Cause:**
- Simple implementation without shutdown handling
- No context propagation

**Dampak:**
- **Memory leak** on hot reload/restart
- Goroutine accumulation
- Minor in production (single instance) tapi problematic di development

**Fix yang Dibutuhkan:**

```go
func RateLimiter() gin.HandlerFunc {
  var (
    mu       sync.Mutex
    visitors = make(map[string]*visitor)
    ctx, cancel = context.WithCancel(context.Background())
  )

  // Background worker with context
  go func() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for {
      select {
      case <-ticker.C:
        mu.Lock()
        for ip, v := range visitors {
          if time.Since(v.lastSeen) > 3*time.Minute {
            delete(visitors, ip)
          }
        }
        mu.Unlock()
      case <-ctx.Done():
        return
      }
    }
  }()

  return func(c *gin.Context) {
    // ... existing logic
  }
}
```

---

### BUG-019 — Password Validation Inconsistency

**Severity:** Medium  
**File:** `backend/internal/usecase/auth_usecase.go`  
**Line:** 37-43

**Deskripsi:**
```go
hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(req.Password)
hasLower := regexp.MustCompile(`[a-z]`).MatchString(req.Password)
hasDigit := regexp.MustCompile(`[0-9]`).MatchString(req.Password)

if !hasUpper || !hasLower || !hasDigit {
  return nil, ErrWeakPassword
}
```

Password validation di backend, tapi error message di frontend tidak selalu match:
- Frontend: "Password harus mengandung huruf besar, huruf kecil, dan angka"
- Backend: "password harus mengandung setidaknya satu huruf besar, satu huruf kecil, dan satu angka"

**Root Cause:**
- Duplicate validation logic
- Error messages not synchronized

**Dampak:**
- **User confusion** — different error messages
- Inconsistent UX

**Fix yang Dibutuhkan:**

1. **Standardize error messages:**
```go
// Backend
var ErrWeakPassword = errors.New("Password harus mengandung huruf besar, huruf kecil, dan angka")
```

2. **Or better, have single source of truth:**
- Move validation to shared package
- Or validate only at backend, frontend just shows generic message

---

### BUG-020 — Missing Input Sanitization di Beberapa Tempat

**Severity:** Medium  
**File:** Multiple backend files

**Deskripsi:**
Beberapa handler melakukan sanitization, tapi tidak semua:
- `auth_handler.go`: Sanitizes name, phone, email ✅
- `outlet_handler.go`: Sanitizes name, address, phone ✅
- `service_handler.go`: Sanitizes name, unit ✅
- `order_handler.go`: **Does NOT sanitize** ❌

**Root Cause:**
- Inconsistent implementation
- No global sanitization middleware

**Dampak:**
- **XSS potential** — unsanitized input stored to database
- **Data integrity issues** — whitespace, null bytes in data

**Fix yang Dibutuhkan:**

```go
// order_handler.go - CreateOrder
func (h *OrderHandler) CreateOrder(c *gin.Context) {
  // ...
  
  var req dto.OrderRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    // ...
  }
  
  // Add sanitization
  req.OutletID = utils.Sanitize(req.OutletID)
  // Note: Items should also be sanitized if they contain user input
  
  // ... rest of logic
}
```

---

## Summary

| Severity | Jumlah |
|----------|--------|
| Critical | 4 |
| High     | 5 |
| Medium   | 8 |
| Low      | 3 |
| **Total**| **20** |

---

## Critical & High Priority (HARUS FIX SEKARANG)

### Critical (App tidak bisa digunakan):
- **BUG-001** — WebSocket URL hardcoded `ws://` → **Fix: Set wss:// URL di Vercel env**
- **BUG-002** — Proxy URL hardcoded localhost → **Fix: Update nuxt.config.ts atau set env di Vercel**
- **BUG-003** — Database SSL mode disable → **Fix: Set DB_SSLMODE=require di Railway**
- **BUG-004** — JWT_SECRET hardcoded → **Fix: Generate new secret, set di Railway env**

### High (Fitur utama tidak berfungsi):
- **BUG-005** — CORS wildcard `*` → **Fix: Whitelist specific origins**
- **BUG-006** — WebSocket CheckOrigin always true → **Fix: Add origin whitelist**
- **BUG-007** — Token expired tidak ditangani → **Fix: Add 401 interceptor**
- **BUG-008** — Missing error handling di API calls → **Fix: Add onError callbacks**
- **BUG-017** — Migration error handling lenient → **Fix: Fatal on migration error**

---

## Medium & Low Priority (Bisa Fix Kemudian)

### Medium (Fitur sekunder bermasalah):
- **BUG-010** — process.client deprecated → **Fix: Use import.meta.client**
- **BUG-011** — Owner register page not linked → **Fix: Add link or remove file**
- **BUG-012** — Tracking page static data → **Fix: Implement real API or remove**
- **BUG-013** — Missing loading states → **Fix: Add consistent loading UI**
- **BUG-016** — Missing deploy config files → **Fix: Create Dockerfile, railway.toml, vercel.json**
- **BUG-019** — Password validation inconsistency → **Fix: Standardize messages**
- **BUG-020** — Missing input sanitization → **Fix: Sanitize all user input**

### Low (Cosmetic/minor):
- **BUG-009** — Console.log in production → **Fix: Remove or conditionally log**
- **BUG-014** — Notification errors not shown → **Fix: Add toast or indicator**
- **BUG-015** — WebSocket reconnect aggressive → **Fix: Add jitter, increase delay**
- **BUG-018** — Rate limiter memory leak potential → **Fix: Add context for shutdown**

---

## Checklist Deployment Production

### Frontend (Vercel):
- [ ] Set `NUXT_PUBLIC_API_BASE_URL=https://laundryin-backend-production.up.railway.app/api/v1`
- [ ] Set `NUXT_PUBLIC_WS_BASE_URL=wss://laundryin-backend-production.up.railway.app/api/v1/ws/connect`
- [ ] Update `nuxt.config.ts` untuk conditional proxy
- [ ] Remove console.log statements
- [ ] Replace `process.client` dengan `import.meta.client`

### Backend (Railway):
- [ ] Set `DATABASE_URL` atau `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
- [ ] Set `DB_SSLMODE=require`
- [ ] Generate new `JWT_SECRET` (openssl rand -base64 32)
- [ ] Set `GIN_MODE=release`
- [ ] Set `PORT=8080` (atau sesuai Railway)
- [ ] Create `Dockerfile` untuk consistent build
- [ ] Update CORS whitelist dengan domain production
- [ ] Update WebSocket CheckOrigin whitelist

### Database:
- [ ] Ensure PostgreSQL di Railway sudah running
- [ ] Connect database service ke backend service di Railway
- [ ] Verify SSL connection enabled

### Testing:
- [ ] Test login/register flow
- [ ] Test WebSocket connection (notifikasi real-time)
- [ ] Test CRUD outlets
- [ ] Test CRUD services
- [ ] Test create order
- [ ] Test order status update
- [ ] Test notification receive
- [ ] Test token expiration (wait 24h atau force expire)

---

**Dokumentasi ini harus di-review ulang setelah semua Critical & High bugs di-fix.**
