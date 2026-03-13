# 🔑 PHASE 1 — Token Fix

Masalah utama: token tidak tersimpan dengan benar setelah login,
menyebabkan semua request authenticated gagal (401/403).

**Selesaikan phase ini sampai acceptance criteria terpenuhi sebelum lanjut ke Phase 2.**

---

## Konteks

Backend membungkus semua response dalam format:
```json
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "token": "eyJ...",
    "user": { "id": "...", "name": "...", "role": "customer" }
  }
}
```

Frontend saat ini kemungkinan mengakses struktur yang salah,
menyebabkan token `undefined` dan tidak tersimpan ke localStorage.

---

## Langkah 1 — Verifikasi Struktur Response Asli

Jalankan curl ini dan perhatikan strukturnya:

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6281200001111","password":"Budi@12345"}' | jq .
```

Catat apakah token ada di `data.token` atau struktur lain. Hasil ini jadi acuan fix di bawah.

---

## Langkah 2 — Fix Auth Store

**File:** `app/stores/auth.ts`

Pastikan implementasinya **persis** seperti berikut:

```typescript
import { defineStore } from 'pinia'

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

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: null,
    user: null
  }),

  getters: {
    isLoggedIn: (state): boolean => !!state.token,
    isOwner: (state): boolean => state.user?.role === 'owner',
    isCustomer: (state): boolean => state.user?.role === 'customer',
    authHeader: (state): string =>
      state.token ? `Bearer ${state.token}` : ''
  },

  actions: {
    setAuth(token: string, user: User) {
      this.token = token
      this.user = user
      if (process.client) {
        localStorage.setItem('laundryin_token', token)
        localStorage.setItem('laundryin_user', JSON.stringify(user))
      }
    },

    logout() {
      this.token = null
      this.user = null
      if (process.client) {
        localStorage.removeItem('laundryin_token')
        localStorage.removeItem('laundryin_user')
      }
    },

    restoreSession() {
      if (process.client) {
        const token = localStorage.getItem('laundryin_token')
        const userStr = localStorage.getItem('laundryin_user')
        if (token && userStr) {
          try {
            this.token = token
            this.user = JSON.parse(userStr)
          } catch {
            this.logout()
          }
        }
      }
    }
  }
})
```

---

## Langkah 3 — Fix Login Customer

**File:** `app/pages/customer/login.vue`

Ganti seluruh fungsi `login()` dengan ini:

```typescript
const login = async () => {
  errorMsg.value = ''

  if (!form.value.phone.trim()) {
    errorMsg.value = 'Nomor HP harus diisi'
    return
  }
  if (form.value.password.length < 8) {
    errorMsg.value = 'Password minimal 8 karakter'
    return
  }

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

    // Debug — hapus setelah konfirmasi token tersimpan
    console.log('[LOGIN] token:', res.data?.token)

    authStore.setAuth(res.data.token, res.data.user)
    router.push(redirectTo.value)

  } catch (err: any) {
    errorMsg.value = err?.data?.message || 'Login gagal, coba lagi'
  } finally {
    loading.value = false
  }
}
```

---

## Langkah 4 — Fix Register Customer

**File:** `app/pages/customer/register.vue`

Ganti seluruh fungsi `register()` dengan ini:

```typescript
const register = async () => {
  errorMsg.value = ''

  if (!form.value.name.trim()) { errorMsg.value = 'Nama harus diisi'; return }
  if (!form.value.phone.trim()) { errorMsg.value = 'Nomor HP harus diisi'; return }
  if (form.value.password.length < 8) { errorMsg.value = 'Password minimal 8 karakter'; return }
  if (form.value.password !== form.value.confirmPassword) {
    errorMsg.value = 'Password tidak cocok'; return
  }

  loading.value = true
  try {
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

    authStore.setAuth(res.data.token, res.data.user)
    router.push('/customer')

  } catch (err: any) {
    const apiMsg = err?.data?.message || ''
    if (apiMsg.includes('sudah terdaftar')) {
      errorMsg.value = 'Nomor HP sudah terdaftar, silakan login'
    } else if (apiMsg.includes('huruf besar')) {
      errorMsg.value = apiMsg
    } else {
      errorMsg.value = apiMsg || 'Registrasi gagal, coba lagi'
    }
  } finally {
    loading.value = false
  }
}
```

---

## Langkah 5 — Fix Login Owner

**File:** `app/pages/owner/login.vue`

Ganti seluruh fungsi `login()` dengan ini:

```typescript
const login = async () => {
  errorMsg.value = ''
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
      errorMsg.value = 'Akun ini bukan akun owner'
      return
    }

    authStore.setAuth(res.data.token, res.data.user)
    router.push('/owner')

  } catch (err: any) {
    errorMsg.value = err?.data?.message || 'Login gagal'
  } finally {
    loading.value = false
  }
}
```

---

## Langkah 6 — Pastikan restoreSession Dipanggil

**File:** `app/app.vue`

```vue
<script setup>
const authStore = useAuthStore()
onMounted(() => {
  authStore.restoreSession()
})
</script>
```

---

## ✅ Acceptance Criteria Phase 1

- [ ] `curl` login → response menampilkan token dengan jelas di field `data.token`
- [ ] Login sebagai customer → DevTools → Application → Local Storage → `laundryin_token` terisi
- [ ] Login sebagai owner → `laundryin_token` dan `laundryin_user` terisi
- [ ] Refresh halaman setelah login → token masih ada, tidak hilang
- [ ] `console.log('[LOGIN] token:', res.data?.token)` menampilkan string JWT, bukan `undefined`
- [ ] `authStore.authHeader` menghasilkan `"Bearer eyJ..."` bukan string kosong

**Semua centang? → Lanjut ke Phase 2**