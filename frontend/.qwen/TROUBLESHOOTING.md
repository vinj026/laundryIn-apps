# 🔧 Troubleshooting: Frontend Tidak Muncul Outlet

## Masalah
- ✅ API langsung (curl) berhasil
- ❌ Frontend di Vercel tidak menampilkan outlet

## Penyebab Paling Mungkin

### 1. Environment Variable di Vercel Belum Di-Set ⚠️

Ini adalah penyebab 99% masalah. Vercel **TIDAK** otomatis membaca `.env.example` atau file `.env` apapun.

**Solusi:**

#### Langkah 1: Buka Vercel Dashboard
1. Buka https://vercel.com/dashboard
2. Klik project frontend lu (laundryIn-frontend atau apapun namanya)
3. Klik **Settings** di tab atas
4. Klik **Environment Variables** di sidebar kiri

#### Langkah 2: Tambahkan 3 Variables Ini

Klik **Add New Variable** untuk masing-masing:

| Variable Name | Value | Environments |
|--------------|-------|--------------|
| `NUXT_PUBLIC_API_BASE_URL` | `https://laundryin-apps-production-ebe2.up.railway.app/api/v1` | ✅ Production |
| `NUXT_PUBLIC_WS_BASE_URL` | `wss://laundryin-apps-production-ebe2.up.railway.app/api/v1/ws/connect` | ✅ Production |
| `BACKEND_URL` | `https://laundryin-apps-production-ebe2.up.railway.app` | ✅ Production |

**PENTING:** Pastikan centang **Production** saat add variable!

#### Langkah 3: Redeploy
Setelah add semua variables:
1. Buka **Deployments** tab
2. Klik tombol **⋯** di deployment terakhir
3. Klik **Redeploy**

ATAU push commit baru ke GitHub (akan auto trigger deploy).

---

### 2. Browser Console Error

Buka browser DevTools (F12) dan cek **Console** tab.

**Kemungkinan error:**

#### Error: "Failed to fetch" atau "Network Error"
```
TypeError: Failed to fetch
```
**Penyebab:** Environment variable belum di-set atau salah URL.

**Solusi:**
1. Cek di Vercel Dashboard > Environment Variables
2. Pastikan `NUXT_PUBLIC_API_BASE_URL` sudah di-set dengan benar
3. Redeploy

#### Error: CORS
```
Access to fetch at '...' from origin 'https://laundryin.vercel.app' has been blocked by CORS policy
```
**Penyebab:** Backend Railway tidak mengirim CORS header yang benar.

**Solusi:** Sudah ditambahkan di `vercel.json` rewrites. Pastikan sudah deploy perubahan terbaru.

#### Error: 404 Not Found
```
GET https://laundryin.vercel.app/api/v1/public/outlets 404
```
**Penyebab:** Vercel tidak forward request ke Railway.

**Solusi:**
1. Pastikan `vercel.json` sudah di-commit dan di-deploy
2. Atau set `NUXT_PUBLIC_API_BASE_URL` dengan full URL Railway

---

### 3. Network Tab Inspection

Buka **Network** tab di browser DevTools, lalu refresh halaman.

**Cek request ke `/api/v1/public/outlets`:**

#### Request URL yang BENAR:
```
Request URL: https://laundryin-apps-production-ebe2.up.railway.app/api/v1/public/outlets
```
✅ Ini artinya environment variable sudah benar!

#### Request URL yang SALAH:
```
Request URL: https://laundryin.vercel.app/api/v1/public/outlets
```
❌ Ini artinya Vercel rewrites mungkin tidak bekerja, atau environment variable belum di-set.

**Fix:** Set `NUXT_PUBLIC_API_BASE_URL` di Vercel dengan full URL Railway.

---

### 4. Test Manual dengan cURL

Test dari terminal untuk memastikan backend accessible:

```bash
# Test ping
curl https://laundryin-apps-production-ebe2.up.railway.app/ping

# Test get outlets
curl https://laundryin-apps-production-ebe2.up.railway.app/api/v1/public/outlets

# Test dengan header Vercel origin
curl -H "Origin: https://laundryin.vercel.app" \
     https://laundryin-apps-production-ebe2.up.railway.app/api/v1/public/outlets
```

Semua harus return JSON response.

---

## Checklist Lengkap

- [ ] Environment variables di-set di Vercel Dashboard
- [ ] Ketiga variables sudah di-add (API_BASE_URL, WS_BASE_URL, BACKEND_URL)
- [ ] Environment "Production" dicentang saat add variable
- [ ] Redeploy sudah dilakukan
- [ ] Browser cache sudah di-clear (Ctrl+Shift+R / Cmd+Shift+R)
- [ ] Browser DevTools Console tidak ada error
- [ ] Network tab menunjukkan request ke Railway URL

---

## Cara Cepat Debug

### 1. Tambah Log di Frontend

Tambahkan ini di `app/pages/customer/index.vue` untuk debug:

```vue
<script setup lang="ts">
const config = useRuntimeConfig()
console.log('🔧 API Base URL:', config.public.apiBase)
</script>
```

Deploy, lalu cek browser console. Harus muncul:
```
🔧 API Base URL: https://laundryin-apps-production-ebe2.up.railway.app/api/v1
```

Kalau muncul `/api`, berarti environment variable belum terbaca!

### 2. Test di Vercel Preview Deploy

Setiap push ke GitHub, Vercel buat preview deploy. Bisa test dulu sebelum production:
1. Buka link preview deploy di Vercel
2. Test apakah outlet muncul
3. Kalau berhasil, merge ke production

---

## Environment Variables yang Benar

**Di Vercel Dashboard:**

```
NUXT_PUBLIC_API_BASE_URL = https://laundryin-apps-production-ebe2.up.railway.app/api/v1
NUXT_PUBLIC_WS_BASE_URL = wss://laundryin-apps-production-ebe2.up.railway.app/api/v1/ws/connect
BACKEND_URL = https://laundryin-apps-production-ebe2.up.railway.app
```

**JANGAN gunakan `/api` di production!** Itu hanya untuk local development.

---

## Kalau Masih Tidak Muncul

1. **Screenshot error di browser console**
2. **Screenshot Network tab** (request URL dan response)
3. **Screenshot Environment Variables di Vercel**

Lalu share untuk debug lebih lanjut.
