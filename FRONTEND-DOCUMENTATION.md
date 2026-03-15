# Dokumentasi Teknis Frontend LaundryIn

**Single Source of Truth** untuk semua developer dan AI agent yang mengerjakan project LaundryIn.

Dokumentasi ini dibuat berdasarkan pembacaan langsung seluruh kode frontend yang ada di repository.

---

## Section 1 — Konfigurasi & Environment

### 1.1 nuxt.config.ts

File konfigurasi utama Nuxt 4.

```typescript
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE_URL || '/api',
      wsBase: process.env.NUXT_PUBLIC_WS_BASE_URL || 'ws://localhost:8080/api/v1/ws/connect'
    }
  },
  devtools: { enabled: true },
  future: {
    compatibilityVersion: 4,
  },
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  routeRules: {
    '/api/**': {
      proxy: (process.env.BACKEND_URL || 'http://localhost:8080') + '/api/v1/**'
    }
  },
  css: ['~/assets/css/main.css'],
  app: {
    pageTransition: { name: 'page', mode: 'out-in' },
    head: {
      title: 'LaundryIn — Modern Laundry Platform',
      meta: [
        { name: 'description', content: 'Smart online laundry platform with real-time tracking and management dashboard.' },
        { name: 'theme-color', content: '#0a0a0a' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Roboto+Mono:wght@400;500;700&display=swap' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200&display=swap' }
      ]
    }
  },
  vite: {
    esbuild: {
      drop: []  // Console drop disabled untuk remote debugging
    }
  }
})
```

**Konfigurasi Penting:**

| Konfigurasi | Nilai | Keterangan |
|-------------|-------|------------|
| `compatibilityDate` | `'2024-11-01'` | Tanggal kompatibilitas Nuxt |
| `future.compatibilityVersion` | `4` | Menggunakan Nuxt 4 compatibility |
| `modules` | `@nuxtjs/tailwindcss`, `@pinia/nuxt` | Tailwind CSS dan Pinia state management |
| `pageTransition` | `{ name: 'page', mode: 'out-in' }` | Animasi transisi halaman |

---

### 1.2 Environment Variables

| Variable | Digunakan Di | Nilai Development | Nilai Production | Wajib? |
|----------|-------------|-------------------|------------------|--------|
| `NUXT_PUBLIC_API_BASE_URL` | `nuxt.config.ts` → `runtimeConfig.public.apiBase` | `/api` (fallback) | URL backend production (mis. `https://api.laundryin.com/api`) | Tidak (ada fallback) |
| `NUXT_PUBLIC_WS_BASE_URL` | `nuxt.config.ts` → `runtimeConfig.public.wsBase` | `ws://localhost:8080/api/v1/ws/connect` (fallback) | URL WebSocket production | Tidak (ada fallback) |
| `BACKEND_URL` | `nuxt.config.ts` → `routeRules` proxy | `http://localhost:8080` (fallback) | URL backend production untuk proxy | Tidak (ada fallback) |

---

### 1.3 Cara Kerja Proxy

#### Development (Local)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Browser   │────▶│  Nuxt Dev   │────▶│   Backend   │
│ localhost:  │     │  Server:3000│     │ localhost:  │
│    3000     │     │  (proxy /api│     │    8080     │
│             │     │   → :8080)  │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

**Step-by-step:**
1. Frontend melakukan fetch ke `/api/xxx`
2. Nuxt dev server intercept request ke `/api/**` berdasarkan `routeRules`
3. Request di-proxy ke `http://localhost:8080/api/v1/xxx` (dari `BACKEND_URL`)
4. Backend menerima request dan mengembalikan response

#### Production

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Browser   │────▶│   Vercel/   │────▶│   Backend   │
│             │     │   Hosting   │     │   Server    │
│             │     │  (direct to │     │             │
│             │     │  apiBase)   │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

**Step-by-step:**
1. Frontend melakukan fetch ke `/api/xxx`
2. Jika `NUXT_PUBLIC_API_BASE_URL` diset (mis. `https://api.laundryin.com/api`), request langsung ke URL tersebut
3. Jika tidak diset (default `/api`), request tetap ke path relatif `/api/xxx`
4. Backend menerima request dan mengembalikan response

**Catatan:** Di production, proxy Nuxt `routeRules` tidak digunakan karena hosting seperti Vercel tidak mendukung proxy rules yang sama dengan dev server.

---

## Section 2 — Stores (State Management)

### 2.1 auth.ts

**File:** `app/stores/auth.ts`

**STATE:**
- `token`: `string | null` — JWT token untuk autentikasi — default: `null`
- `user`: `User | null` — Data user yang login — default: `null`

**Interface User:**
```typescript
interface User {
  id: string
  name: string
  phone: string
  email?: string
  role: 'owner' | 'customer'
}
```

**GETTERS:**
- `isLoggedIn`: `boolean` — Returns `!!state.token`
- `isOwner`: `boolean` — Returns `state.user?.role === 'owner'`
- `isCustomer`: `boolean` — Returns `state.user?.role === 'customer'`
- `authHeader`: `string` — Returns `state.token ? \`Bearer ${state.token}\` : ''`

**ACTIONS:**
- `setAuth(token: string, user: User)`
  - Deskripsi: Menyimpan token dan user ke state dan localStorage
  - Side effects: `localStorage.setItem('laundryin_token', token)`, `localStorage.setItem('laundryin_user', JSON.stringify(user))`

- `logout()`
  - Deskripsi: Menghapus semua data autentikasi
  - Side effects: `localStorage.removeItem('laundryin_token')`, `localStorage.removeItem('laundryin_user')`

- `restoreSession()`
  - Deskripsi: Memulihkan session dari localStorage saat app diinisialisasi
  - Side effects: Membaca dari localStorage, bisa memanggil `logout()` jika data corrupt

**PERSISTED DATA (localStorage):**
- Key: `laundryin_token` — Value: JWT token string
- Key: `laundryin_user` — Value: JSON string dari object User

---

### 2.2 cart.ts

**File:** `app/stores/cart.ts`

**STATE:**
- `items`: `CartItem[]` — Array item yang dipilih untuk checkout — default: `[]`
- `outletId`: `string | null` — ID outlet yang sedang dipilih — default: `null`

**Interface CartItem:**
```typescript
interface CartItem {
  serviceId: string
  name: string
  price: number
  unit: string
  qty: string  // zero-trust frontend uses string
}
```

**GETTERS:**
- `totalPreview`: `number` — Menghitung total harga dari semua item: `items.reduce((total, item) => total + (item.price * parseFloat(item.qty || '0')), 0)`
- `itemCount`: `number` — Jumlah unique items di cart: `items.length`

**ACTIONS:**
- `setOutlet(id: string)`
  - Deskripsi: Set outlet ID, otomatis clear cart jika outlet berubah
  - Side effects: Clear `items` jika `outletId` berubah

- `addItem(item: CartItem)`
  - Deskripsi: Tambah item ke cart, merge qty jika serviceId sama
  - Side effects: None

- `updateQty(serviceId: string, qty: string)`
  - Deskripsi: Update qty item, hapus jika qty <= 0
  - Side effects: Bisa memanggil `removeItem()`

- `removeItem(serviceId: string)`
  - Deskripsi: Hapus item dari cart berdasarkan serviceId
  - Side effects: None

- `clearCart()`
  - Deskripsi: Kosongkan cart dan reset outletId
  - Side effects: None

**PERSISTED DATA:** Tidak ada — cart bersifat temporary dan tidak dipersist

---

### 2.3 notification.ts

**File:** `app/stores/notification.ts`

**STATE:**
- `notifications`: `Notification[]` — Array notifikasi — default: `[]`
- `unreadCount`: `number` — Jumlah notifikasi belum dibaca — default: `0`
- `isOpen`: `boolean` — Status dropdown notifikasi — default: `false`
- `loading`: `boolean` — Loading state saat fetch — default: `false`

**Interface Notification:**
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
```

**ACTIONS:**
- `addNotification(notif: Notification)`
  - Deskripsi: Tambah notifikasi baru (dari WebSocket)
  - Side effects: None
  - Logic: Hindari duplicate by ID, limit max 50 notifikasi

- `fetchNotifications(page = 1, limit = 20)`
  - Deskripsi: Fetch notifikasi dari API
  - Side effects: API call ke `/api/notifications`
  - Logic: Merge notifikasi dari API dengan yang dari WS, sort by date, limit 100

- `fetchUnreadCount()`
  - Deskripsi: Fetch jumlah notifikasi belum dibaca
  - Side effects: API call ke `/api/notifications/unread-count`

- `markAsRead(id: string)`
  - Deskripsi: Tandai notifikasi sebagai dibaca
  - Side effects: API call PATCH ke `/api/notifications/${id}/read`

- `markAllAsRead()`
  - Deskripsi: Tandai semua notifikasi sebagai dibaca
  - Side effects: API call PATCH ke `/api/notifications/read-all`

- `toggleDropdown(val?: boolean)`
  - Deskripsi: Toggle visibility dropdown notifikasi
  - Side effects: None

**PERSISTED DATA:** Tidak ada — notifikasi bersifat temporary

---

## Section 3 — Plugins & Composables

### 3.1 Plugins

#### auth.client.ts

**File:** `app/plugins/auth.client.ts`

**Tipe:** plugin (client-side only)

**Dijalankan:** Client-side only, saat aplikasi diinisialisasi

**Urutan eksekusi:** Setelah app mount, sebelum halaman pertama di-render

**Fungsi:**
```typescript
export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()
  if (import.meta.client) {
    authStore.restoreSession()
  }
})
```

**Dependencies:** `useAuthStore`

---

#### websocket.client.ts

**File:** `app/plugins/websocket.client.ts`

**Tipe:** plugin (client-side only)

**Dijalankan:** Client-side only, setelah auth.client.ts

**Urutan eksekusi:** Setelah session restored, setup WS connection

**Fungsi:**
```typescript
export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()
  const notifStore = useNotificationStore()
  const { connect, disconnect } = useWebSocket()

  if (import.meta.client && authStore.isLoggedIn) {
    notifStore.fetchUnreadCount()
  }

  watch(() => authStore.isLoggedIn, (loggedIn) => {
    if (import.meta.client) {
      if (loggedIn) connect()
      else disconnect()
    }
  }, { immediate: true })
})
```

**Dependencies:** `useAuthStore`, `useNotificationStore`, `useWebSocket`

---

### 3.2 Composables

#### useApiFetch.ts

**File:** `app/composables/useApiFetch.ts`

**Fungsi yang di-export:**

- `useApiFetch(path: string | (() => string), options: any = {})`
  - Return type: Nuxt `useFetch` result
  - Deskripsi: Wrapper untuk `useFetch` dengan auth header dan path transformation otomatis
  - Logic:
    - Transform path `/api/xxx` jika `apiBase !== '/api'`
    - Tambahkan `Authorization: Bearer ${token}` header jika ada token
    - Set `immediate: false` jika butuh auth tapi tidak ada token
    - Handle 401 error: logout dan redirect ke `/customer/login`

- `useApiRaw<T>(path: string, options: any = {}): Promise<T>`
  - Return type: `Promise<T>`
  - Deskripsi: Wrapper untuk `$fetch` dengan auth header dan error handling
  - Logic:
    - Transform path `/api/xxx` jika `apiBase !== '/api'`
    - Reject jika butuh auth tapi tidak ada token
    - Handle 401 error: logout dan redirect ke `/customer/login`

**Dependencies:** `useAuthStore`, `useRuntimeConfig`, `useToast`, `useRouter`

---

#### useToast.ts

**File:** `app/composables/useToast.ts`

**Fungsi yang di-export:**

- `useToast()`
  - Return type: `{ toasts, remove, success, error, info }`
  - Deskripsi: Composable untuk menampilkan toast notification

**Methods:**
- `success(msg: string)`: Tampilkan toast hijau (success)
- `error(msg: string)`: Tampilkan toast merah (error)
- `info(msg: string)`: Tampilkan toast abu-abu (info)
- `remove(id: number)`: Hapus toast tertentu
- `toasts`: Readonly array dari semua toast aktif

**Dependencies:** None (pure Vue reactivity)

---

#### useWebSocket.ts ⭐

**File:** `app/composables/useWebSocket.ts`

**Tipe:** composable

**Dijalankan:** Client-side only

**URL WebSocket:** Dari `runtimeConfig.public.wsBase` (env: `NUXT_PUBLIC_WS_BASE_URL`), fallback: `ws://localhost:8080/api/v1/ws/connect`

**Fungsi yang di-export:**

- `useWebSocket()`
  - Return type: `{ connect, disconnect }`
  - Deskripsi: Composable untuk manage WebSocket connection

**Cara koneksi dibuka:**
```typescript
const connect = () => {
  if (!import.meta.client || !authStore.isLoggedIn || ws) return

  const config = useRuntimeConfig()
  let wsUrl = config.public.wsBase

  if (!wsUrl.includes('token=')) {
    const separator = wsUrl.includes('?') ? '&' : '?'
    wsUrl = `${wsUrl}${separator}token=${authStore.token}`
  }

  ws = new WebSocket(wsUrl)
  // ... setup handlers
}
```

**Format pesan yang diterima:**
```typescript
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data)
  // Expected format:
  {
    id: string,        // atau akan digenerate fallback
    type: string,      // 'new_order', 'order_status', 'price_updated', 'order_cancelled'
    title: string,
    body: string,
    data: any,         // extra data (mis. order_id)
    is_read: false,    // selalu false dari WS
    created_at: string // ISO timestamp
  }
}
```

**Cara reconnect bekerja:**
```typescript
const scheduleReconnect = () => {
  if (reconnectTimer) clearTimeout(reconnectTimer)

  const delayWithJitter = reconnectDelay + Math.random() * RECONNECT_JITTER

  reconnectTimer = setTimeout(() => {
    reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY)
    connect()
  }, delayWithJitter)
}
```

- Initial delay: 2000ms
- Exponential backoff: delay × 2 setiap retry
- Max delay: 60000ms
- Jitter: 0-1000ms random untuk prevent thundering herd

**Kapan disconnect dipanggil:**
1. Saat user logout (via plugin watch effect)
2. Saat WebSocket onclose dan user tidak logged in
3. Manual call `disconnect()` dari komponen

**Dependencies:** `useAuthStore`, `useNotificationStore`, `useToast`, `useRuntimeConfig`

---

## Section 4 — Layouts

### 4.1 customer.vue

**File:** `app/layouts/customer.vue`

**Digunakan oleh halaman:**
- `/customer` (customer/index.vue)
- `/customer/outlet/[id]`
- `/customer/orders`
- `/customer/orders/[id]`
- `/customer/profile`
- `/customer/tracking`

**Struktur UI:**
- **Desktop Sidebar (72px, hidden mobile):**
  - Logo LaundryIn (icon laundry)
  - Nav items: Explore, Orders, Profile
  - Notification bell dengan badge unread count
  - Profile avatar button dengan dropdown (Logout)

- **Mobile Bottom Nav (16px height, hidden desktop):**
  - Nav items: Explore, Orders, Profile
  - Notification bell integrated

**Guard/Redirect:**
- Tidak ada guard di layout level
- Guard ada di masing-masing halaman via `watchEffect`

**Props/Data:**
- `navItems`: Array navigation items
  ```typescript
  [
    { to: '/customer', icon: 'explore', label: 'Explore' },
    { to: '/customer/orders', icon: 'receipt_long', label: 'Orders' },
    { to: '/customer/profile', icon: 'person', label: 'Profile' }
  ]
  ```

**Auth-dependent UI:**
- Notification bell: Hanya muncul jika `authStore.isLoggedIn`
- Profile dropdown:
  - Jika logged in: Avatar dengan initial nama + dropdown (Lihat Profil, Logout)
  - Jika tidak logged in: Login button

**Active State Logic:**
```typescript
const isActive = (path: string) => {
  if (path === '/customer') return route.path === '/customer' || route.path.startsWith('/customer/outlet')
  return route.path.startsWith(path)
}
```

---

### 4.2 owner.vue

**File:** `app/layouts/owner.vue`

**Digunakan oleh halaman:**
- `/owner` (owner/index.vue)
- `/owner/outlets`
- `/owner/services`
- `/owner/orders`

**Struktur UI:**
- **Desktop Sidebar (72px, hidden mobile):**
  - Logo LaundryIn (icon laundry)
  - Nav items: Analytics, Outlets, Services, Orders
  - Notification bell dengan badge unread count
  - Profile avatar button dengan dropdown (Logout)

- **Mobile Bottom Nav (16px height, hidden desktop):**
  - Nav items: Analytics, Outlets, Services, Orders
  - Profile button dengan dropdown (Logout)

**Guard/Redirect:**
```typescript
watchEffect(() => {
  if (!import.meta.client) return
  if (!authStore.isOwner) {
    navigateTo('/owner/login')
  }
})
```
- Redirect ke `/owner/login` jika user bukan owner

**Props/Data:**
- `navItems`: Array navigation items
  ```typescript
  [
    { to: '/owner', icon: 'bar_chart', label: 'Analytics' },
    { to: '/owner/outlets', icon: 'storefront', label: 'Outlets' },
    { to: '/owner/services', icon: 'dry_cleaning', label: 'Services' },
    { to: '/owner/orders', icon: 'view_list', label: 'Orders' }
  ]
  ```

**Auth-dependent UI:**
- Notification bell: Hanya muncul jika `authStore.isLoggedIn`
- Profile dropdown: Hanya avatar + Logout (tidak ada "Lihat Profil")

**Active State Logic:**
```typescript
const isActive = (path: string) => {
  if (path === '/owner') return route.path === '/owner'
  return route.path.startsWith(path)
}
```

---

## Section 5 — Halaman (Pages)

### 5.1 index.vue (Landing Page)

**URL:** `/`
**File:** `app/pages/index.vue`
**Layout:** none
**Auth Required:** tidak

**Struktur UI:**
- Header dengan logo LaundryIn
- Hero section dengan icon dan CTA buttons
- Footer dengan tagline

**User Actions:**
- Click "Owner Dashboard" → Navigate ke `/owner/login`
- Click "Find a Laundry" → Navigate ke `/customer`

---

### 5.2 customer/login.vue ⭐

**URL:** `/customer/login`
**File:** `app/pages/customer/login.vue`
**Layout:** false (standalone page)
**Auth Required:** tidak

**GUARD RULES:**
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
- Jika sudah login sebagai owner → redirect ke `/owner` dengan error toast
- Jika sudah login sebagai customer → redirect ke `/customer`

**FORM:**
- `phone`: text input (tel) — required
- `password`: password input — required, min 8 karakter, show/hide toggle

**Submit Function (SETIAP BARIS):**
```typescript
const login = async () => {
  console.log('🔑 Login attempt:', {
    phone: form.value.phone.trim(),
    formattedPhone: formattedPhone.value
  })

  // 1. Validasi phone tidak kosong
  if (!form.value.phone.trim()) {
    toastError('Nomor HP harus diisi')
    return
  }
  // 2. Validasi password min 8 karakter
  if (form.value.password.length < 8) {
    toastError('Password minimal 8 karakter')
    return
  }

  loading.value = true
  try {
    // 3. Call API login
    const res = await $fetch<{
      status: string
      message: string
      data: { token: string, user: any }
    }>('/api/auth/login', {
      method: 'POST',
      body: {
        phone: formattedPhone.value,  // Phone sudah diformat (+62)
        password: form.value.password
      }
    })

    console.log('✅ Login success:', res)
    
    // 4. Set auth di store (simpan token + user ke state & localStorage)
    authStore.setAuth(res.data.token, res.data.user)
    
    // 5. Redirect ke halaman tujuan atau default /customer
    router.push(redirectTo.value)

  } catch (err: any) {
    console.error('🔴 LOGIN_ERROR:', err)
    console.error('Status:', err?.statusCode || err?.status)
    console.error('Message:', err?.data?.message)
    console.error('Response:', err?.data)

    const status = err?.statusCode || err?.status || err?.response?.status
    const msg = err?.data?.message || ''

    // 6. Handle error berdasarkan status code
    if (status === 401) {
      toastError('Nomor HP atau password salah')
    } else if (status === 400) {
      toastError(msg || 'Data tidak valid')
    } else if (status === 429) {
      toastError('Terlalu banyak percobaan, coba lagi beberapa saat')
    } else if (!status || status === 0) {
      toastError('Tidak dapat terhubung ke server, coba lagi')
    } else {
      toastError(`Error ${status}: ${msg || 'Terjadi kesalahan'}`)
    }
  } finally {
    loading.value = false
  }
}
```

**Phone Formatting:**
```typescript
const formattedPhone = computed(() => {
  const p = form.value.phone.trim()
  if (p.startsWith('0')) return '+62' + p.slice(1)  // 0812 → +62812
  if (p.startsWith('62') && !p.startsWith('+')) return '+' + p  // 62812 → +62812
  return p  // sudah +62812 atau format lain
})
```

**ERROR STATES:**
- 401: "Nomor HP atau password salah"
- 400: "Data tidak valid" atau pesan dari API
- 429: "Terlalu banyak percobaan, coba lagi beberapa saat"
- 0/Network: "Tidak dapat terhubung ke server, coba lagi"
- Other: "Error {status}: {message}"

**LOADING STATES:**
- Button disabled dengan spinner icon saat `loading.value = true`

---

### 5.3 customer/register.vue ⭐

**URL:** `/customer/register`
**File:** `app/pages/customer/register.vue`
**Layout:** false (standalone page)
**Auth Required:** tidak

**GUARD RULES:**
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

**FORM:**
- `name`: text input — required
- `phone`: text input (tel) — required
- `password`: password input — required, min 8 karakter, show/hide toggle
- `confirmPassword`: password input — required, harus match dengan password

**Validasi password (hint di UI):**
> Min. 8 characters, must include uppercase, lowercase, and a number.

**Submit Function (SETIAP BARIS):**
```typescript
const register = async () => {
  console.log('📝 Register attempt:', {
    name: form.value.name.trim(),
    phone: form.value.phone.trim(),
    formattedPhone: formattedPhone.value
  })

  // 1. Validasi nama
  if (!form.value.name.trim()) { toastError('Nama harus diisi'); return }
  // 2. Validasi phone
  if (!form.value.phone.trim()) { toastError('Nomor HP harus diisi'); return }
  // 3. Validasi password min 8 karakter
  if (form.value.password.length < 8) { toastError('Password minimal 8 karakter'); return }
  // 4. Validasi confirm password match
  if (form.value.password !== form.value.confirmPassword) {
    toastError('Konfirmasi password tidak cocok'); return
  }

  loading.value = true
  try {
    // 5. Call API register dengan role: 'customer'
    const res = await $fetch<{
      status: string
      message: string
      data: { token: string, user: any }
    }>('/api/auth/register', {
      method: 'POST',
      body: {
        name: form.value.name.trim(),
        phone: formattedPhone.value,
        password: form.value.password,
        role: 'customer'
      }
    })

    console.log('✅ Register success:', res)
    
    // 6. Set auth di store (auto login setelah register)
    authStore.setAuth(res.data.token, res.data.user)
    
    // 7. Redirect ke /customer
    router.push('/customer')

  } catch (err: any) {
    console.error('🔴 REGISTER_ERROR:', err)
    console.error('Status:', err?.statusCode || err?.status)
    console.error('Message:', err?.data?.message)
    console.error('Response:', err?.data)

    const status = err?.statusCode || err?.status || err?.response?.status
    const apiMsg = err?.data?.message || ''

    // 8. Handle error berdasarkan status code
    if (status === 409 || apiMsg.includes('sudah terdaftar')) {
      toastError('Nomor HP sudah terdaftar, silakan login')
    } else if (status === 400) {
      if (apiMsg.includes('lemah') || apiMsg.includes('huruf besar') || apiMsg.includes('angka')) {
        toastError('Password harus mengandung huruf besar, huruf kecil, dan angka')
      } else {
        toastError(apiMsg || 'Data tidak valid')
      }
    } else if (status === 429) {
      toastError('Terlalu banyak percobaan, coba lagi beberapa saat')
    } else if (!status || status === 0) {
      toastError('Tidak dapat terhubung ke server, coba lagi')
    } else {
      toastError(`Error ${status}: ${apiMsg || 'Terjadi kesalahan'}`)
    }
  } finally {
    loading.value = false
  }
}
```

**ERROR STATES:**
- 409 / "sudah terdaftar": "Nomor HP sudah terdaftar, silakan login"
- 400 / password lemah: "Password harus mengandung huruf besar, huruf kecil, dan angka"
- 400 / other: Pesan dari API atau "Data tidak valid"
- 429: "Terlalu banyak percobaan, coba lagi beberapa saat"
- 0/Network: "Tidak dapat terhubung ke server, coba lagi"

---

### 5.4 customer/index.vue (Find Outlet)

**URL:** `/customer`
**File:** `app/pages/customer/index.vue`
**Layout:** customer
**Auth Required:** tidak

**DATA FETCHING:**
- `useApiFetch<{ data: { data: Outlet[] } }>('/api/public/outlets')`
  - Response diakses via: `outletsResponse.value?.data?.data`
  - Ditaruh di: `outletsResponse` (ref)

**INTERFACE Outlet:**
```typescript
interface Outlet {
  id: string
  name: string
  address: string
  description: string
  phone: string
  is_active: boolean
}
```

**USER ACTIONS:**
- Search: Filter outlet by name atau address (case-insensitive)
- Click outlet card: Navigate ke `/customer/outlet/[id]`
- Pull-to-refresh: `refresh()` dipanggil saat `onActivated`

**LOADING STATES:**
- Spinner besar di tengah saat `pending = true`

**EMPTY STATES:**
- "No outlets found" jika search tidak match

---

### 5.5 customer/outlet/[id].vue ⭐ (Checkout Flow)

**URL:** `/customer/outlet/[id]`
**File:** `app/pages/customer/outlet/[id].vue`
**Layout:** customer
**Auth Required:** tidak (tapi required saat checkout)

**DATA FETCHING:**
1. `useApiFetch<{ data: Outlet }>(\`/api/public/outlets/${outletId}\`)` — Detail outlet
2. `useApiFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/public/outlets')` — Semua outlet (untuk modal change)
3. `useApiFetch<{ data: Service[] }>(\`/api/public/outlets/${outletId}/services\`)` — Services di outlet

**FORM:**
- `name`: text input — required
- `phone`: text input (tel) — required
- `address`: textarea — required
- `date`: DateCarousel component — default hari ini, format YYYY-MM-DD
- `time`: radio chips (morning: 09:00-12:00, afternoon: 13:00-16:00) — default: afternoon
- `notes`: textarea — optional

**SERVICES PANEL:**
- Search: Filter by name atau description
- Category filter: All, Clothes (unit: kg), Shoes (unit: pcs + nama mengandung sepatu/tas/sandal), Others
- Add/Remove: Quantity control dengan + / - buttons

**CHECKOUT FLOW (Step by Step):**
```typescript
const checkout = async () => {
  // Guard 1: must be logged in
  if (!authStore.isLoggedIn) {
    toastInfo('Silakan login terlebih dahulu')
    router.push(`/customer/login?redirect=${route.fullPath}`)
    return
  }

  // Guard 2: cart must not be empty
  if (cartStore.items.length === 0) {
    toastError('Pilih minimal 1 layanan')
    return
  }

  // Guard 3: form validation
  if (!form.value.name.trim()) {
    toastError('Nama harus diisi')
    return
  }
  if (!form.value.phone.trim()) {
    toastError('Nomor WhatsApp harus diisi')
    return
  }
  if (!form.value.address.trim()) {
    toastError('Alamat pickup harus diisi')
    return
  }

  checkoutLoading.value = true

  try {
    const payload = {
      outlet_id: outletId,
      customer_name: form.value.name.trim(),
      customer_phone: form.value.phone.trim(),
      customer_address: form.value.address.trim(),
      pickup_date: form.value.date,
      pickup_time: form.value.time,
      notes: form.value.notes.trim(),
      items: cartStore.items.map(i => ({
        service_id: i.serviceId,
        qty: i.qty.toString()
      }))
    }

    await useApiRaw('/api/orders', {
      method: 'POST',
      body: payload
    })

    // Success: clear cart, toast, redirect
    cartStore.clearCart()
    toastSuccess('Pesanan berhasil dibuat!')
    router.push('/customer/orders')

  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    const apiMsg = err?.data?.message || ''

    if (status === 400) {
      toastError(apiMsg || 'Data tidak valid')
    } else if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push(`/customer/login?redirect=${route.fullPath}`)
    } else if (status === 404) {
      toastError('Layanan tidak tersedia, coba muat ulang halaman')
    } else {
      toastError('Gagal membuat pesanan, coba lagi')
    }
  } finally {
    checkoutLoading.value = false
  }
}
```

**CHANGE OUTLET:**
- Click "Change" button → Modal muncul dengan list semua outlet
- Click outlet lain → Clear cart, navigate ke `/customer/outlet/[newId]`

**ERROR STATES:**
- 400: "Data tidak valid" atau pesan dari API
- 401: Logout dan redirect ke login dengan redirect URL
- 404: "Layanan tidak tersedia, coba muat ulang halaman"

---

### 5.6 customer/orders/[id].vue ⭐ (Tracking)

**URL:** `/customer/orders/[id]`
**File:** `app/pages/customer/orders/[id].vue`
**Layout:** customer
**Auth Required:** ya (customer)

**DATA FETCHING:**
1. `useApiFetch<ApiResponse<PaginatedResponse<Order[]>>>('/api/orders', { server: false })` — Semua orders
2. `useAsyncData(\`/api/public/outlets/${order.outlet_id}\`)` — Detail outlet

**INTERFACE Order:**
```typescript
interface Order {
  id: string
  status: string  // 'pending' | 'process' | 'completed' | 'picked_up'
  total_price: string
  final_total_price?: string
  outlet_id: string
  order_date: string
  items: OrderItem[]
}

interface OrderItem {
  id: string
  service_name: string
  qty: string
  actual_qty?: string
  unit: string
  service_price: string
  subtotal: string
  final_price?: string
}
```

**STATUS CONFIG:**
```typescript
const statusConfig = {
  pending: {
    icon: 'inventory_2',
    label: 'Menunggu Konfirmasi',
    description: 'Pesanan kamu sudah masuk dan sedang menunggu konfirmasi dari outlet.'
  },
  process: {
    icon: 'local_laundry_service',
    label: 'Sedang Dicuci',
    description: 'Mesin sedang membersihkan bajumu. Proses pencucian dan pengeringan sedang berjalan.'
  },
  completed: {
    icon: 'check_circle',
    label: 'Siap Diambil',
    description: 'Cucian kamu sudah selesai! Silakan ambil di outlet.'
  },
  picked_up: {
    icon: 'hail',
    label: 'Selesai',
    description: 'Pesanan sudah diambil. Terima kasih!'
  }
}
```

**TIMELINE STEPS:**
```typescript
const statusOrder = ['pending', 'process', 'completed', 'picked_up']
const stepLabels = {
  pending: { label: 'Pesanan Dibuat', description: 'Order placed via App.' },
  process: { label: 'Sedang Diproses', description: 'Items being washed and dried.' },
  completed: { label: 'Siap Diambil', description: 'Ready for pickup at outlet.' },
  picked_up: { label: 'Selesai', description: 'Customer has picked up the order.' }
}
```

**Timeline Logic:**
- Step `active`: jika index <= current status index
- Step `current`: jika index === current status index (pulsing dot)
- Time ditampilkan hanya untuk step pertama (order placed time)

**RECEIPT DISPLAY:**
- List semua items dengan qty dan harga
- Jika `actual_qty` ada dan berbeda dari `qty`: tampilkan coret di qty lama
- Jika `final_price` ada dan berbeda dari `subtotal`: tampilkan coret di harga lama
- Total: tampilkan `final_total_price` jika ada, else `total_price`
- Badge "Harga Final" jika `final_total_price` ada, else "Estimasi"

**LOADING STATES:**
- Spinner besar saat `pending = true`

**ERROR STATES:**
- "Order not found" jika order ID tidak ada di list

---

### 5.7 customer/orders/index.vue

**URL:** `/customer/orders`
**File:** `app/pages/customer/orders/index.vue`
**Layout:** customer
**Auth Required:** ya (customer)

**GUARD RULES:**
```typescript
watchEffect(() => {
  if (import.meta.client && !authStore.isLoggedIn) {
    router.push('/customer/login?redirect=/customer/orders')
  }
})
```

**DATA FETCHING:**
- `useApiFetch<PaginatedResponse<Order[]>>('/api/orders')`
  - Response diakses via: `ordersResponse.value?.data?.data`

**INTERFACE Order:**
```typescript
interface Order {
  id: string
  status: string
  total_price: string
  final_total_price?: string
  outlet_id: string
  order_date: string
}
```

**BADGE MAPPING:**
```typescript
{
  'badge-pending': order.status === 'pending',
  'badge-process': order.status === 'process',
  'badge-completed': order.status === 'completed',  // Display: "Ready Pickup"
  'badge-picked-up': order.status === 'picked_up'   // Display: "Finished"
}
```

**USER ACTIONS:**
- Click order card: Navigate ke `/customer/orders/[id]`
- Click "Coba Lagi" saat error: `refresh()`
- Pull-to-refresh: `refresh()` dipanggil saat `onActivated`

**ERROR STATES:**
- 401: Logout dan redirect ke `/customer/login`
- Other: "Gagal memuat pesanan" dengan tombol retry

**EMPTY STATES:**
- "Belum ada pesanan" dengan button "Cari Laundry"

---

### 5.8 customer/profile.vue

**URL:** `/customer/profile`
**File:** `app/pages/customer/profile.vue`
**Layout:** customer
**Auth Required:** ya (customer)

**GUARD RULES:**
```typescript
watchEffect(() => {
  if (import.meta.client) {
    if (!authStore.isLoggedIn) {
      router.push('/customer/login')
    } else if (authStore.user?.role !== 'customer') {
      router.push('/owner/login')
    }
  }
})
```

**VIEW MODE:**
- Avatar circle dengan initial nama
- Nama, phone, role badge
- Detail rows: Nama Lengkap, Nomor HP, Role
- Buttons: Edit Profil, Keluar

**EDIT MODE:**
- Form: Nama Lengkap, Nomor HP
- Buttons: Batal, Simpan Perubahan

**SAVE EDIT LOGIC:**
```typescript
const saveEdit = async () => {
  editError.value = ''
  const nextName = editForm.value.name.trim()
  const nextPhone = editForm.value.phone.trim()

  if (!nextName) {
    editError.value = 'Nama tidak boleh kosong'
    return
  }

  editLoading.value = true
  try {
    // API call (fallback jika endpoint tidak ada)
    await useApiRaw('/api/users/me', {
      method: 'PUT',
      body: { name: nextName, phone: nextPhone }
    }).catch(err => {
      console.warn('API update failed, updating locally only:', err.message)
    })

    // Update store (dan localStorage)
    if (authStore.user && authStore.token) {
      authStore.setAuth(authStore.token, {
        ...authStore.user,
        name: nextName,
        phone: nextPhone
      })
    }

    editSuccess.value = true
    isEditing.value = false
    setTimeout(() => { editSuccess.value = false }, 3000)

  } catch (err: any) {
    editError.value = err?.data?.message || 'Gagal menyimpan perubahan'
  } finally {
    editLoading.value = false
  }
}
```

**⚠️ Inconsistency:** API call `/api/users/me` di-catch silent, jadi selalu success secara UI meskipun API endpoint tidak ada. Update hanya terjadi di localStorage.

---

### 5.9 customer/tracking.vue

**URL:** `/customer/tracking`
**File:** `app/pages/customer/tracking.vue`
**Layout:** customer
**Auth Required:** ya

**Behavior:**
- Redirect otomatis ke `/customer/orders` saat mounted
- Hanya menampilkan spinner loading selama redirect

**🚧 Not Implemented:** Halaman ini adalah placeholder, tidak ada tracking functionality di luar `/customer/orders/[id]`

---

### 5.10 owner/login.vue

**URL:** `/owner/login`
**File:** `app/pages/owner/login.vue`
**Layout:** false (standalone page)
**Auth Required:** tidak

**GUARD RULES:**
```typescript
watchEffect(() => {
  if (import.meta.client && authStore.isLoggedIn) {
    if (authStore.user?.role === 'customer') {
      toastError('Kamu sudah login sebagai customer, silakan logout terlebih dahulu')
      router.push('/customer')
    } else if (authStore.user?.role === 'owner') {
      router.push('/owner')
    }
  }
})
```

**FORM:**
- `phone`: text input (tel) — required
- `password`: password input — required, min 8 karakter, show/hide toggle
- "Forgot password?" button (🚧 not implemented, hanya UI)

**Submit Function:**
```typescript
const login = async () => {
  if (!form.value.phone.trim()) { toastError('Nomor HP harus diisi'); return }
  if (form.value.password.length < 8) { toastError('Password minimal 8 karakter'); return }

  loading.value = true
  try {
    const res = await $fetch<{
      status: string
      message: string
      data: { token: string, user: any }
    }>('/api/auth/login', {
      method: 'POST',
      body: {
        phone: formattedPhone.value,
        password: form.value.password
      }
    })

    if (res.data.user.role !== 'owner') {
      toastError('Akun ini bukan akun owner')
      return
    }

    authStore.setAuth(res.data.token, res.data.user)
    router.push('/owner')

  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    const msg = err?.data?.message || ''

    if (status === 401) {
      toastError('Nomor HP atau password salah')
    } else if (status === 403) {
      toastError('Kamu tidak memiliki akses ke dashboard owner')
    } else if (status === 429) {
      toastError('Terlalu banyak percobaan, coba lagi beberapa saat')
    } else if (!status || status === 0) {
      toastError('Tidak dapat terhubung ke server, coba lagi')
    } else {
      toastError(`Error ${status}: ${msg || 'Terjadi kesalahan'}`)
    }
  } finally {
    loading.value = false
  }
}
```

**⚠️ Inconsistency:** Ada validasi `role !== 'owner'` setelah login success, tapi seharusnya backend sudah mengembalikan 403 jika role bukan owner.

---

### 5.11 owner/register.vue

**URL:** `/owner/register`
**File:** `app/pages/owner/register.vue`
**Layout:** false (standalone page)
**Auth Required:** tidak

**GUARD RULES:**
```typescript
watchEffect(() => {
  if (import.meta.client && authStore.isLoggedIn) {
    if (authStore.user?.role === 'customer') {
      toastError('Kamu sudah login sebagai customer, silakan logout terlebih dahulu')
      router.push('/customer')
    } else if (authStore.user?.role === 'owner') {
      router.push('/owner')
    }
  }
})
```

**FORM:**
- `name`: text input — required
- `phone`: text input (tel) — required
- `password`: password input — required, min 8 karakter
- `confirmPassword`: password input — required

**Submit Function:**
```typescript
const register = async () => {
  if (!form.value.name.trim()) { toastError('Nama harus diisi'); return }
  if (!form.value.phone.trim()) { toastError('Nomor HP harus diisi'); return }
  if (form.value.password.length < 8) { toastError('Password minimal 8 karakter'); return }
  if (form.value.password !== form.value.confirmPassword) {
    toastError('Konfirmasi password tidak cocok'); return
  }

  loading.value = true
  try {
    const res = await useApiRaw<{
      status: string
      message: string
      data: { token: string, user: any }
    }>('/api/auth/register', {
      method: 'POST',
      body: {
        name: form.value.name.trim(),
        phone: formattedPhone.value,
        password: form.value.password,
        role: 'owner'
      }
    })

    authStore.setAuth(res.data.token, res.data.user)
    router.push('/owner')

  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    const apiMsg = err?.data?.message || ''

    if (status === 409 || apiMsg.includes('sudah terdaftar')) {
      toastError('Nomor HP sudah terdaftar, silakan login')
    } else if (status === 400 && (apiMsg.includes('lemah') || apiMsg.includes('huruf besar') || apiMsg.includes('angka'))) {
      toastError('Password harus mengandung huruf besar, huruf kecil, dan angka')
    } else if (status === 429) {
      toastError('Terlalu banyak percobaan, coba lagi beberapa saat')
    } else {
      toastError(err?.data?.message || 'Terjadi kesalahan, coba lagi')
    }
  } finally {
    loading.value = false
  }
}
```

---

### 5.12 owner/index.vue (Analytics)

**URL:** `/owner`
**File:** `app/pages/owner/index.vue`
**Layout:** owner
**Auth Required:** ya (owner)

**GUARD RULES:**
```typescript
watchEffect(() => {
  if (import.meta.client) {
    if (!authStore.isLoggedIn || authStore.user?.role !== 'owner') {
      router.push('/owner/login')
    }
  }
})
```

**DATE FILTER:**
- Options: 7 Hari, 30 Hari, 3 Bulan
- Default: 30 Hari

**DATA FETCHING:**
```typescript
const { data: analytics, pending, error, refresh } = await useAsyncData<AnalyticsWrapper | null>(
  'owner-analytics',
  async () => {
    const dates = getDateRange(daysFilter.value)
    const queryParams = `?start_date=${dates.start_date}&end_date=${dates.end_date}`

    const [omzetRes, summaryRes, servicesRes] = await Promise.all([
      useApiRaw<ApiResponse<OmzetResponse>>(`/api/reports/omzet${queryParams}`),
      useApiRaw<ApiResponse<OrderStatusSummaryResponse>>(`/api/reports/orders/summary${queryParams}`),
      useApiRaw<ApiResponse<TopServiceResponse[]>>(`/api/reports/services/top${queryParams}`)
    ])

    return {
      totalRevenue: Number(omzetRes.data?.total_omzet || 0),
      ordersPending: summaryRes.data?.pending || 0,
      ordersProcess: summaryRes.data?.process || 0,
      ordersCompleted: summaryRes.data?.completed || 0,
      ordersPickedUp: summaryRes.data?.picked_up || 0,
      topServices: servicesRes.data || []
    }
  },
  { watch: [daysFilter], server: false }
)
```

**UI COMPONENTS:**
1. **Total Omzet Card:** Display `totalRevenue` dengan format Rupiah
2. **Pipeline Status:** 4 cards (Pending, Process, Completed, Picked Up)
3. **Top Services:** List services dengan total revenue dan qty

**ERROR STATES:**
- 401: Logout dan redirect ke `/owner/login`
- Other: "Gagal memuat data dashboard" dengan tombol retry

---

### 5.13 owner/outlets.vue

**URL:** `/owner/outlets`
**File:** `app/pages/owner/outlets.vue`
**Layout:** owner
**Auth Required:** ya (owner)

**DATA FETCHING:**
- `useApiFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/outlets', { server: false })`

**USER ACTIONS:**
- Search: Filter outlet by name atau address
- Create: Button "+" → Modal form
- Edit: Button edit icon → Modal form dengan data existing
- Delete: Button delete icon → Confirmation modal

**FORM VALIDATION (Create/Edit):**
```typescript
const validateForm = () => {
  let valid = true
  formErrors.value = { name: '', address: '', phone: '' }

  if (!outletForm.value.name.trim()) {
    formErrors.value.name = 'Nama outlet wajib diisi'
    valid = false
  }
  if (!outletForm.value.address.trim()) {
    formErrors.value.address = 'Alamat wajib diisi'
    valid = false
  }
  if (!outletForm.value.phone.trim()) {
    formErrors.value.phone = 'Nomor telepon wajib diisi'
    valid = false
  } else if (!/^\+62\d{8,13}$/.test(outletForm.value.phone.trim().replace(/\s/g, ''))) {
    formErrors.value.phone = 'Format harus +62xxx (contoh: +628123456789)'
    valid = false
  }

  return valid
}
```

**API ENDPOINTS:**
- GET `/api/outlets` — List outlets
- POST `/api/outlets` — Create outlet
- PUT `/api/outlets/[id]` — Update outlet
- DELETE `/api/outlets/[id]` — Delete outlet

---

### 5.14 owner/services.vue

**URL:** `/owner/services`
**File:** `app/pages/owner/services.vue`
**Layout:** owner
**Auth Required:** ya (owner)

**DATA FETCHING:**
```typescript
const { data: outletsWrapper } = await useApiFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/outlets')

const { data: services, pending, refresh: refreshServices } = await useAsyncData<Service[]>(
  'owner-services',
  async () => {
    if (outlets.value.length === 0) return []

    if (selectedOutletId.value) {
      // Filter by specific outlet
      const res = await useApiRaw<ApiResponse<Service[]>>(`/api/outlets/${selectedOutletId.value}/services`)
      const outletName = outlets.value.find(o => o.id === selectedOutletId.value)?.name || 'Unknown'
      return (res.data || []).map(s => ({ ...s, outlet_name: outletName }))
    }

    // All outlets parallel query
    const promises = outlets.value.map(async (o) => {
      try {
        const res = await useApiRaw<ApiResponse<Service[]>>(`/api/outlets/${o.id}/services`)
        return (res.data || []).map(s => ({ ...s, outlet_name: o.name }))
      } catch {
        return []
      }
    })

    const results = await Promise.all(promises)
    return results.flat()
  },
  { watch: [selectedOutletId, outletsWrapper], server: false }
)
```

**FILTER:**
- Dropdown "Filter Outlet": All outlets atau specific outlet

**USER ACTIONS:**
- Create: Button "New Service" → Modal form
- Edit: Button edit icon → Modal form dengan data existing
- Delete: Button delete icon → Confirmation modal

**FORM VALIDATION:**
```typescript
const validateForm = () => {
  let valid = true
  formErrors.value = { outlet_id: '', name: '', price: '', unit: '' }

  if (!serviceForm.value.outlet_id) {
    formErrors.value.outlet_id = 'Pilih outlet terlebih dahulu'
    valid = false
  }
  if (!serviceForm.value.name.trim()) {
    formErrors.value.name = 'Nama layanan wajib diisi'
    valid = false
  }
  if (!serviceForm.value.price || serviceForm.value.price <= 0) {
    formErrors.value.price = 'Harga harus lebih dari 0'
    valid = false
  }

  return valid
}
```

**⚠️ Inconsistency:** Saat edit, `outlet_id` tidak bisa diubah (disabled), tapi field masih dikirim di payload (dihapus manual di submit).

**API ENDPOINTS:**
- GET `/api/outlets/[id]/services` — List services per outlet
- POST `/api/services` — Create service
- PUT `/api/services/[id]` — Update service
- DELETE `/api/services/[id]` — Delete service

---

### 5.15 owner/orders.vue

**URL:** `/owner/orders`
**File:** `app/pages/owner/orders.vue`
**Layout:** owner
**Auth Required:** ya (owner)

**DATA FETCHING:**
```typescript
const { data: outletsWrapper } = await useApiFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/outlets', { server: false })

const { data: ordersWrapper, pending, error, refresh } = await useAsyncData<ApiResponse<PaginatedResponse<Order[]>>>(
  'owner-orders',
  async () => {
    if (!selectedOutletId.value) {
      return { status: 'success', message: '', data: { data: [], page: 1, limit: 10, total: 0, total_pages: 0 } }
    }
    return useApiRaw<ApiResponse<PaginatedResponse<Order[]>>>(`/api/outlets/${selectedOutletId.value}/orders`, {
      params: { page: 1, limit: 10 }
    })
  },
  { watch: [selectedOutletId], server: false }
)
```

**FILTER:**
- Dropdown "Outlet": Pilih outlet (default: outlet pertama)
- Status tabs: Semua, pending, process, completed, picked_up, cancelled

**ORDER MANAGEMENT:**

**Status: pending**
- Input berat aktual untuk items dengan unit KG
- Button "Proses" → PATCH `/api/orders/[id]/status` dengan body `{ status: 'process', items: [...] }`
- Button "Cancel" → PATCH dengan `{ status: 'cancelled' }`

**Status: process**
- Button "Selesai" → PATCH dengan `{ status: 'completed' }`
- Button "Cancel" → PATCH dengan `{ status: 'cancelled' }`

**Status: completed**
- Button "Sudah Diambil" → PATCH dengan `{ status: 'picked_up' }`

**Status: picked_up / cancelled**
- Display label "Selesai (Diambil/Dibatalkan)" disabled

**INPUT BERAT AKTUAL:**
```typescript
const actualQtyInputs = ref<Record<string, Record<string, string>>>({})

watch(orders, (newOrders) => {
  newOrders.forEach(o => {
    if (!actualQtyInputs.value[o.id]) {
      actualQtyInputs.value[o.id] = {}
    }
    o.items.forEach(i => {
      if (i.unit === 'KG' && !actualQtyInputs.value[o.id][i.id]) {
        actualQtyInputs.value[o.id][i.id] = i.actual_qty || ''
      }
    })
  })
}, { immediate: true })
```

**VALIDASI SEBELUM PROSES:**
```typescript
const hasMissingKgInput = (order: Order) => {
  if (order.status !== 'pending') return false
  const kgItems = order.items.filter(i => i.unit === 'KG')
  if (kgItems.length === 0) return false

  for (const item of kgItems) {
    const val = actualQtyInputs.value[order.id]?.[item.id]
    if (!val || parseFloat(val) <= 0) return true
  }
  return false
}
```

**⚠️ Bug Potential:** Button "Proses" disabled jika ada KG item tanpa input, tapi tidak ada visual indication mana yang kosong.

---

## Section 6 — Components

### 6.1 DateCarousel.vue

**File:** `app/components/ui/DateCarousel.vue`

**Digunakan di:** `customer/outlet/[id].vue`

**Props:**
- `modelValue`: `string` (YYYY-MM-DD) — required — selected date

**Emits:**
- `update:modelValue`: `string` — saat date dipilih

**Behavior:**
- Generate 14 hari: hari ini + 13 hari ke depan
- Horizontal scrollable carousel
- Active card scale 1, inactive scale 0.85
- Left/Right navigation buttons muncul saat bisa scroll
- Scroll 80% width per click

**INTERFACE Day:**
```typescript
{
  label: 'MON',  // weekday short uppercase
  date: 16,      // day of month
  month: 'Mar',  // month short
  value: '2025-03-16'  // YYYY-MM-DD
}
```

---

### 6.2 NotificationDropdown.vue

**File:** `app/components/ui/NotificationDropdown.vue`

**Digunakan di:** `layouts/customer.vue`, `layouts/owner.vue`

**Props:** None

**Emits:** None

**Behavior:**
- Muncul saat `notifStore.isOpen = true`
- Fetch notifications saat dropdown dibuka
- Display list notifikasi dengan icon berdasarkan type
- Click notifikasi: mark as read + navigate berdasarkan type
- "Tandai semua baca" button: mark all as read

**NOTIFICATION TYPES:**
```typescript
const getIcon = (type: string) => {
  switch (type) {
    case 'new_order': return 'shopping_cart'
    case 'order_status': return 'local_laundry_service'
    case 'price_updated': return 'payments'
    case 'order_cancelled': return 'cancel'
    default: return 'notifications'
  }
}
```

**NAVIGATION LOGIC:**
```typescript
const handleNotifClick = async (notif: Notification) => {
  await notifStore.markAsRead(notif.id)
  notifStore.toggleDropdown(false)

  if (notif.type === 'new_order') {
    router.push('/owner/orders')
  } else if (notif.data?.order_id) {
    router.push(`/customer/orders/${notif.data.order_id}`)
  }
}
```

---

### 6.3 ToastContainer.vue

**File:** `app/components/ui/ToastContainer.vue`

**Digunakan di:** `app.vue` (global)

**Props:** None

**Emits:** None

**Behavior:**
- Teleport to body
- Display toasts di bottom-right
- Auto-dismiss tidak ada (manual close button)
- Slide-in dari kanan animation

**TOAST TYPES:**
- `success`: Green background, check_circle icon
- `error`: Red background, error icon
- `info`: Gray background, info icon

---

## Section 7 — CSS & Design System

### 7.1 Custom CSS Classes (main.css)

**Global Scrollbar:**
```css
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background-color: rgba(255, 255, 255, 0.06); border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background-color: rgba(255, 255, 255, 0.12); }
```

**Hide Scrollbar:**
```css
.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
```

**M3 Field (Material 3 Input):**
```css
.m3-field { position: relative; }
.m3-field input/textarea/select { bg-surface-container, border-border, rounded-xl, px-4 py-3.5, text-sm }
.m3-field label { absolute, floating label animation }
```

**Chip Radio:**
```css
.chip-radio:checked+label { bg-primary/10, border-primary, text-primary }
```

**Status Badges:**
```css
.badge { text-[10px], font-bold, px-2.5 py-1, rounded-md, uppercase }
.badge-pending { bg-surface-containerHigh, text-surface-onSurfaceVariant }
.badge-process { bg-primary/15, text-primary }
.badge-completed { bg-success-muted, text-success }
.badge-picked-up { bg-surface-containerHigh/50, text-surface-onSurfaceVariant, opacity-60 }
.badge-live { bg-danger-muted, text-danger, border-danger/30 }
```

**Buttons:**
```css
.btn-primary { bg-primary, text-primary-text, rounded-xl, px-5 py-3, hover:brightness-110 }
.btn-secondary { bg-surface-containerHigh, border-border, rounded-xl, px-5 py-3 }
.btn-ghost { bg-transparent, text-surface-onSurfaceVariant, hover:bg-surface-containerHigh }
.btn-danger { bg-danger, text-surface-onSurface, rounded-xl, px-5 py-3 }
```

**Cards:**
```css
.card { bg-surface-container, rounded-2xl, border-border, p-5 }
.card-hover { hover:border-border-hover, hover:shadow-lg }
.card-interactive { card + card-hover + cursor-pointer + hover:-translate-y-0.5 }
```

**Utility:**
```css
.text-gradient-primary { bg-gradient-to-r from-primary to-primary-light, bg-clip-text, text-transparent }
.glass { bg-surface-container/80, backdrop-blur-xl, border-border }
```

**Page Transitions:**
```css
.page-enter-active, .page-leave-active { transition: opacity 0.2s ease; }
.page-enter-from, .page-leave-to { opacity: 0; }
```

---

### 7.2 Tailwind Custom Config (tailwind.config.js)

**Font Family:**
```javascript
fontFamily: {
  sans: ['"Inter"', 'system-ui', 'sans-serif'],
  mono: ['"Roboto Mono"', 'ui-monospace', 'monospace'],
}
```

**Custom Colors:**

| Color | DEFAULT | light | dark | container | onContainer | text |
|-------|---------|-------|------|-----------|-------------|------|
| `primary` | `#2dd4bf` (teal-400) | `#5eead4` (teal-300) | `#0d9488` (teal-600) | `#042f2e` (teal-950) | `#99f6e4` (teal-200) | `#021a19` |
| `secondary` | `#a78bfa` (violet-400) | - | - | `#1e1b4b` (indigo-950) | `#c4b5fd` (violet-300) | - |
| `success` | `#4ade80` | - | - | - | - | - |
| `warning` | `#fbbf24` | - | - | - | - | - |
| `danger` | `#f87171` | - | - | - | - | - |
| `surface` | `#0a0a0a` | - | - | - | - | - |
| `surface.raised` | `#111111` | - | - | - | - | - |
| `surface.container` | `#161616` | - | - | - | - | - |
| `surface.containerHigh` | `#1e1e1e` | - | - | - | - | - |
| `surface.overlay` | `#262626` | - | - | - | - | - |
| `surface.onSurface` | `#f0f0f0` | - | - | - | - | - |
| `surface.onSurfaceVariant` | `#8a8a8a` | - | - | - | - | - |
| `border` | `#262626` | subtle: `#1e1e1e`, hover: `#404040` | - | - | - | - |
| `outline` | `#404040` | - | - | - | - | - |

**Border Radius:**
```javascript
borderRadius: {
  'xl': '0.875rem',    // 14px
  '2xl': '1rem',       // 16px
  '3xl': '1.25rem',    // 20px
  '4xl': '1.5rem',     // 24px
}
```

**Transition Duration:**
```javascript
transitionDuration: {
  fast: '150ms',
  normal: '250ms',
  slow: '400ms',
}
```

**Custom Animations:**
```javascript
animation: {
  'fade-in': 'fadeIn 0.3s ease-out',
  'slide-up': 'slideUp 0.35s ease-out',
  'pulse-soft': 'pulseSoft 2s ease-in-out infinite',
}
```

**Success/Warning/Danger Muted:**
```javascript
success: { muted: 'rgba(74, 222, 128, 0.15)' },
warning: { muted: 'rgba(251, 191, 36, 0.15)' },
danger: { muted: 'rgba(248, 113, 113, 0.15)' },
```

---

## Section 8 — API Integration Map

| Halaman/File | Method | Endpoint | Auth Header | Request Body | Response Diakses Via |
|-------------|--------|----------|-------------|--------------|---------------------|
| `customer/login.vue` | POST | `/api/auth/login` | Tidak | `{ phone, password }` | `res.data.token`, `res.data.user` |
| `customer/register.vue` | POST | `/api/auth/register` | Tidak | `{ name, phone, password, role: 'customer' }` | `res.data.token`, `res.data.user` |
| `owner/login.vue` | POST | `/api/auth/login` | Tidak | `{ phone, password }` | `res.data.token`, `res.data.user` |
| `owner/register.vue` | POST | `/api/auth/register` | Tidak | `{ name, phone, password, role: 'owner' }` | `res.data.token`, `res.data.user` |
| `customer/index.vue` | GET | `/api/public/outlets` | Tidak | - | `outletsResponse.value?.data?.data` |
| `customer/outlet/[id].vue` | GET | `/api/public/outlets/[id]` | Tidak | - | `outletWrapper.value?.data` |
| `customer/outlet/[id].vue` | GET | `/api/public/outlets` | Tidak | - | `allOutletsResponse.value?.data?.data` |
| `customer/outlet/[id].vue` | GET | `/api/public/outlets/[id]/services` | Tidak | - | `services.value?.data` |
| `customer/outlet/[id].vue` | POST | `/api/orders` | Ya | `{ outlet_id, customer_name, customer_phone, customer_address, pickup_date, pickup_time, notes, items }` | - |
| `customer/orders/index.vue` | GET | `/api/orders` | Ya | - | `ordersResponse.value?.data?.data` |
| `customer/orders/[id].vue` | GET | `/api/orders` | Ya | - | `ordersWrapper.value?.data?.data` |
| `customer/orders/[id].vue` | GET | `/api/public/outlets/[outlet_id]` | Tidak | - | `outletWrapper.value?.data` |
| `customer/profile.vue` | PUT | `/api/users/me` | Ya | `{ name, phone }` | - (silent fallback) |
| `owner/index.vue` | GET | `/api/reports/omzet?start_date=&end_date=` | Ya | - | `omzetRes.data?.total_omzet` |
| `owner/index.vue` | GET | `/api/reports/orders/summary?start_date=&end_date=` | Ya | - | `summaryRes.data?.{pending,process,completed,picked_up}` |
| `owner/index.vue` | GET | `/api/reports/services/top?start_date=&end_date=` | Ya | - | `servicesRes.data` |
| `owner/outlets.vue` | GET | `/api/outlets` | Ya | - | `outletsWrapper.value?.data?.data` |
| `owner/outlets.vue` | POST | `/api/outlets` | Ya | `{ name, address, phone }` | - |
| `owner/outlets.vue` | PUT | `/api/outlets/[id]` | Ya | `{ name, address, phone }` | - |
| `owner/outlets.vue` | DELETE | `/api/outlets/[id]` | Ya | - | - |
| `owner/services.vue` | GET | `/api/outlets/[id]/services` | Ya | - | `res.data` |
| `owner/services.vue` | POST | `/api/services` | Ya | `{ outlet_id, name, price, unit }` | - |
| `owner/services.vue` | PUT | `/api/services/[id]` | Ya | `{ name, price, unit }` | - |
| `owner/services.vue` | DELETE | `/api/services/[id]` | Ya | - | - |
| `owner/orders.vue` | GET | `/api/outlets` | Ya | - | `outletsWrapper.value?.data?.data` |
| `owner/orders.vue` | GET | `/api/outlets/[id]/orders?page=1&limit=10` | Ya | - | `ordersWrapper.value?.data?.data` |
| `owner/orders.vue` | PATCH | `/api/orders/[id]/status` | Ya | `{ status, items? }` | - |
| `notification.ts` | GET | `/api/notifications?page=1&limit=20` | Ya | - | `res.data?.data` |
| `notification.ts` | GET | `/api/notifications/unread-count` | Ya | - | `res.data?.count` |
| `notification.ts` | PATCH | `/api/notifications/[id]/read` | Ya | - | - |
| `notification.ts` | PATCH | `/api/notifications/read-all` | Ya | - | - |

---

## Section 9 — Known Issues & Inconsistencies

### 9.1 Hardcoded Values

| File | Line/Section | Issue | Rekomendasi |
|------|-------------|-------|-------------|
| `nuxt.config.ts` | `wsBase` fallback | `ws://localhost:8080/api/v1/ws/connect` hardcoded | Sebaiknya env variable wajib di production |
| `useWebSocket.ts` | `reconnectDelay` | `2000ms` hardcoded | Bisa dibuat configurable via env |
| `customer/outlet/[id].vue` | `statusConfig` | Label dan description hardcoded di frontend | Sebaiknya dari backend untuk konsistensi |

---

### 9.2 Unused Code

| File | Code | Keterangan |
|------|------|------------|
| `customer/tracking.vue` | Seluruh file | Hanya redirect ke `/customer/orders`, tidak ada tracking logic |
| `owner/login.vue` | "Forgot password?" button | Hanya UI, tidak ada functionality |
| `nuxt.config.ts` | `vite.esbuild.drop: []` | Comment "Temporarily disabled for remote debugging" — bisa dihapus jika sudah tidak perlu |

---

### 9.3 Inconsistent Patterns

| Issue | File 1 | File 2 | Rekomendasi |
|-------|--------|--------|-------------|
| Login API call | `customer/login.vue` menggunakan `$fetch` langsung | `owner/register.vue` menggunakan `useApiRaw` | Standardisasi ke `useApiRaw` untuk konsistensi error handling |
| Auth guard | Beberapa halaman menggunakan `watchEffect` dengan `import.meta.client` check | Beberapa tidak | Standardisasi pattern guard di semua halaman |
| Error handling | Beberapa halaman handle 401 dengan logout | Beberapa tidak | Pastikan semua halaman yang butuh auth handle 401 consistently |

---

### 9.4 Missing Error Handling

| File | Section | Issue | Impact |
|------|---------|-------|--------|
| `customer/profile.vue` | `saveEdit` | API call `/api/users/me` di-catch silent | User thinks save success padahal mungkin gagal |
| `useWebSocket.ts` | `ws.onerror` | Hanya `console.error`, tidak ada user notification | User tidak tahu kalau WS disconnected |
| `owner/services.vue` | Parallel fetch all outlets | Catch block return `[]` silent | User mungkin tidak lihat services dari outlet yang error |

---

### 9.5 Missing Loading States

| File | Action | Issue |
|------|--------|-------|
| `customer/outlet/[id].vue` | Change outlet modal | `pendingAllOutlets` ada, tapi tidak ada loading state saat fetch services refresh |
| `owner/orders.vue` | Update status | Ada `updatingOrder` ref, tapi tidak ada visual feedback selain disabled button |

---

### 9.6 Dead Code

| File | Code | Keterangan |
|------|------|------------|
| `customer/tracking.vue` | Seluruh halaman | Tidak pernah diakses, tracking ada di `/customer/orders/[id]` |
| `nuxt.config.ts` | `app.pageTransition` | Ada config tapi tidak ada CSS transition yang match dengan nama 'page' (ada, di main.css) — actually used ✅ |

---

### 9.7 TODO/FIXME Comments

**Tidak ada TODO/FIXME comments yang ditemukan di kode.**

---

### 9.8 Potential Bugs

| File | Issue | Severity | Rekomendasi |
|------|-------|----------|-------------|
| `customer/outlet/[id].vue` | `hasMissingKgInput` check tapi tidak ada visual indication mana KG item yang kosong | Medium | Tambah border merah atau helper text di input yang kosong |
| `owner/login.vue` | Validasi `role !== 'owner'` setelah login success | Low | Seharusnya backend return 403, frontend check redundant |
| `useApiFetch.ts` | `onResponseError` handler bisa conflict dengan custom `options.onResponseError` | Medium | Pastikan error handling tidak duplicate |
| `notification.ts` | `fetchNotifications` merge logic bisa duplicate jika WS dan API return same notification | Low | Sudah ada check `apiIds`, tapi bisa ada race condition |

---

## Appendix A — File Structure

```
frontend/
├── app/
│   ├── app.vue                    # Root component
│   ├── assets/
│   │   └── css/
│   │       └── main.css           # Global styles
│   ├── components/
│   │   └── ui/
│   │       ├── DateCarousel.vue
│   │       ├── NotificationDropdown.vue
│   │       └── ToastContainer.vue
│   ├── composables/
│   │   ├── useApiFetch.ts
│   │   ├── useToast.ts
│   │   └── useWebSocket.ts
│   ├── layouts/
│   │   ├── customer.vue
│   │   └── owner.vue
│   ├── pages/
│   │   ├── index.vue              # Landing page
│   │   ├── customer/
│   │   │   ├── index.vue          # Find outlet
│   │   │   ├── login.vue
│   │   │   ├── register.vue
│   │   │   ├── outlet/
│   │   │   │   └── [id].vue       # Checkout
│   │   │   ├── orders/
│   │   │   │   ├── index.vue
│   │   │   │   └── [id].vue       # Tracking
│   │   │   ├── profile.vue
│   │   │   └── tracking.vue       # 🚧 Placeholder
│   │   └── owner/
│   │       ├── index.vue          # Analytics
│   │       ├── login.vue
│   │       ├── register.vue
│   │       ├── outlets.vue
│   │       ├── services.vue
│   │       └── orders.vue
│   ├── plugins/
│   │   ├── auth.client.ts
│   │   └── websocket.client.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── cart.ts
│   │   └── notification.ts
│   └── types/
│       └── api.ts
├── nuxt.config.ts
├── package.json
├── tailwind.config.js
├── tsconfig.json
└── vercel.json
```

---

## Appendix B — Authentication Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        AUTHENTICATION FLOW                       │
└─────────────────────────────────────────────────────────────────┘

1. APP INITIALIZATION
   ┌──────────────┐
   │ app.vue      │
   └──────┬───────┘
          │
          ▼
   ┌──────────────┐
   │ auth.client. │──────▶ restoreSession() from localStorage
   │ plugin       │
   └──────────────┘

2. LOGIN (Customer/Owner)
   ┌──────────────┐     POST /api/auth/login     ┌──────────────┐
   │ Login Form   │─────────────────────────────▶│   Backend    │
   └──────────────┘                              └──────┬───────┘
          │                                             │
          │                                             │ { token, user }
          │◀────────────────────────────────────────────┘
          │
          │ authStore.setAuth(token, user)
          │ ├─ state.token = token
          │ ├─ state.user = user
          │ ├─ localStorage.laundryin_token = token
          │ └─ localStorage.laundryin_user = JSON(user)
          │
          ▼
   ┌──────────────┐
   │   Redirect   │
   └──────────────┘

3. AUTHENTICATED REQUEST
   ┌──────────────┐
   │ useApiFetch  │
   └──────┬───────┘
          │
          │ Get authHeader from authStore
          │ Authorization: Bearer {token}
          │
          ▼
   ┌──────────────┐
   │   API Call   │
   └──────────────┘

4. 401 HANDLING
   ┌──────────────┐
   │ API Response │
   │   401        │
   └──────┬───────┘
          │
          ▼
   ┌──────────────┐
   │ onResponse   │──────▶ authStore.logout()
   │ Error Hook   │        ├─ state.token = null
   └──────────────┘        ├─ state.user = null
          │                ├─ localStorage.removeItem('laundryin_token')
          │                └─ localStorage.removeItem('laundryin_user')
          │
          ▼
   ┌──────────────┐
   │   Redirect   │──────▶ /customer/login
   └──────────────┘

5. WEBSOCKET CONNECTION
   ┌──────────────┐
   │ websocket.   │──────▶ watch(authStore.isLoggedIn)
   │ client.plugin│
   └──────┬───────┘
          │
          │ If loggedIn: connect()
          │ └─ new WebSocket(wsBase + '?token=' + token)
          │
          │ If loggedOut: disconnect()
          │ └─ ws.close()
          │
          ▼
   ┌──────────────┐
   │   WebSocket  │◀─────▶ Realtime notifications
   └──────────────┘
```

---

## Appendix C — State Management Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                      PINIA STORES OVERVIEW                       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   authStore     │     │    cartStore    │     │  notifStore     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ STATE:          │     │ STATE:          │     │ STATE:          │
│ - token         │     │ - items[]       │     │ - notifications │
│ - user          │     │ - outletId      │     │ - unreadCount   │
│                 │     │                 │     │ - isOpen        │
│ GETTERS:        │     │ GETTERS:        │     │ - loading       │
│ - isLoggedIn    │     │ - totalPreview  │     │                 │
│ - isOwner       │     │ - itemCount     │     │ ACTIONS:        │
│ - isCustomer    │     │                 │     │ - addNotif      │
│ - authHeader    │     │ ACTIONS:        │     │ - fetchNotifs   │
│                 │     │ - setOutlet     │     │ - markAsRead    │
│ ACTIONS:        │     │ - addItem       │     │ - markAllRead   │
│ - setAuth       │     │ - updateQty     │     │ - toggleDrop    │
│ - logout        │     │ - removeItem    │     │                 │
│ - restoreSession│     │ - clearCart     │     │ PERSISTED:      │
│                 │     │                 │     │ - None          │
│ PERSISTED:      │     │ PERSISTED:      │     │                 │
│ - localStorage  │     │ - None (temp)   │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                       │                       │
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                    ▼                         ▼
           ┌────────────────┐      ┌────────────────┐
           │  Components    │      │   Composables  │
           │  - Pages       │      │  - useApiFetch │
           │  - Layouts     │      │  - useToast    │
           │  - UI          │      │  - useWebSocket│
           └────────────────┘      └────────────────┘
```

---

**Dokumentasi ini selesai dibuat pada:** March 16, 2026

**Versi:** 1.0.0

**Catatan:** Dokumentasi ini adalah living document. Update setiap ada perubahan signifikan di kode frontend.
