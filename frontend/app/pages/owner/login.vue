<template>
  <div class="min-h-screen flex items-center justify-center p-6 bg-surface relative overflow-hidden">

    <!-- Background -->
    <div class="absolute inset-0 pointer-events-none">
      <div class="absolute top-1/3 left-1/2 -translate-x-1/2 w-[500px] h-[500px] bg-primary/[0.03] rounded-full blur-[120px]"></div>
    </div>

    <!-- Login Card -->
    <div class="relative z-10 w-full max-w-sm card p-8 flex flex-col items-center animate-fade-in">

      <!-- Brand -->
      <div class="mb-8 text-center">
        <div class="h-16 w-16 bg-primary/10 text-primary rounded-2xl flex items-center justify-center mx-auto border border-primary/20 mb-4">
          <span class="material-symbols-outlined text-[36px]">admin_panel_settings</span>
        </div>
        <h1 class="text-2xl font-bold tracking-tight mb-1">Welcome back</h1>
        <p class="text-surface-onSurfaceVariant text-sm">Sign in to your owner dashboard.</p>
      </div>

      <!-- Form -->
      <form class="w-full space-y-4 mb-6" @submit.prevent="login">
        <div class="m3-field">
          <input id="phone" v-model="form.phone" name="phone" type="tel" placeholder=" " required />
          <label for="phone">Phone Number</label>
        </div>
        <div class="m3-field relative">
          <input
            id="password"
            v-model="form.password"
            name="password"
            :type="showPassword ? 'text' : 'password'"
            placeholder=" "
            required
          />
          <label for="password">Password</label>
          <button
            type="button"
            class="absolute right-3 top-3 text-surface-onSurfaceVariant hover:text-surface-onSurface transition-colors"
            @click="showPassword = !showPassword"
          >
            <span class="material-symbols-outlined text-xl">
              {{ showPassword ? 'visibility_off' : 'visibility' }}
            </span>
          </button>
        </div>

        <div class="flex justify-end pt-1">
          <button type="button" class="text-xs font-medium text-primary hover:underline">Forgot password?</button>
        </div>
      </form>

      <!-- Actions -->
      <div class="w-full flex flex-col gap-3">
        <button
          class="btn-primary py-3.5 rounded-xl text-sm w-full"
          :disabled="loading"
          @click="login"
        >
          <span v-if="loading" class="material-symbols-outlined text-lg animate-spin">progress_activity</span>
          <template v-else>
            Sign In
            <span class="material-symbols-outlined text-lg">arrow_forward</span>
          </template>
        </button>
        <NuxtLink to="/owner/register" class="py-2.5 w-full text-surface-onSurfaceVariant text-xs font-medium hover:text-primary transition-colors text-center">
          Don't have an account? <span class="text-primary font-semibold">Register</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useToast } from '~/composables/useToast'

useHead({ title: 'LaundryIn — Owner Login' })
definePageMeta({ layout: false })

const authStore = useAuthStore()
const { error: toastError } = useToast()

const route = useRoute()
const router = useRouter()

// Guard: Redirect if already logged in
watchEffect(() => {
  if (import.meta.client && authStore.isLoggedIn) {
    if (authStore.user?.role === 'customer') {
      if (authStore.token) {
        toastError('Kamu sudah login sebagai customer, silakan logout terlebih dahulu')
      }
      router.push('/customer')
    } else if (authStore.user?.role === 'owner') {
      router.push('/owner')
    }
  }
})

const form = ref({
  phone: '',
  password: ''
})
const showPassword = ref(false)
const loading = ref(false)

const formattedPhone = computed(() => {
  const p = form.value.phone.trim()
  if (p.startsWith('0')) return '+62' + p.slice(1)
  if (p.startsWith('62') && !p.startsWith('+')) return '+' + p
  return p
})

const login = async () => {
  console.log('👑 Owner login attempt:', {
    phone: form.value.phone.trim(),
    formattedPhone: formattedPhone.value
  })

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

    console.log('✅ Owner login success:', res)
    if (res.data.user.role !== 'owner') {
      toastError('Akun ini bukan akun owner')
      return
    }

    authStore.setAuth(res.data.token, res.data.user)
    router.push('/owner')

  } catch (err: any) {
    console.error('🔴 OWNER_LOGIN_ERROR:', err)
    console.error('Status:', err?.statusCode || err?.status)
    console.error('Message:', err?.data?.message)
    console.error('Response:', err?.data)

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
</script>
