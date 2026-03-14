<template>
  <div class="min-h-screen flex items-center justify-center p-6 bg-surface relative overflow-hidden">

    <!-- Background -->
    <div class="absolute inset-0 pointer-events-none">
      <div class="absolute top-1/3 left-1/2 -translate-x-1/2 w-[500px] h-[500px] bg-primary/[0.03] rounded-full blur-[120px]"></div>
    </div>

    <!-- Register Card -->
    <div class="relative z-10 w-full max-w-sm card p-8 flex flex-col items-center animate-fade-in">

      <!-- Brand -->
      <div class="mb-8 text-center">
        <div class="h-16 w-16 bg-primary/10 text-primary rounded-2xl flex items-center justify-center mx-auto border border-primary/20 mb-4">
          <span class="material-symbols-outlined text-[36px]">business_center</span>
        </div>
        <h1 class="text-2xl font-bold tracking-tight mb-1">Owner Account</h1>
        <p class="text-surface-onSurfaceVariant text-sm">Register your laundry business today.</p>
      </div>

      <!-- Form -->
      <form class="w-full space-y-4 mb-6" @submit.prevent="register">
        <div class="m3-field">
          <input id="name" v-model="form.name" name="name" type="text" placeholder=" " required />
          <label for="name">Full Name</label>
        </div>
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
        <div class="m3-field relative">
          <input
            id="confirm"
            v-model="form.confirmPassword"
            name="confirm"
            :type="showPassword ? 'text' : 'password'"
            placeholder=" "
            required
          />
          <label for="confirm">Confirm Password</label>
        </div>

        <!-- Password hint -->
        <p class="text-[11px] text-surface-onSurfaceVariant leading-relaxed px-1">
          Min. 8 characters, must include uppercase, lowercase, and a number.
        </p>
      </form>

      <!-- Actions -->
      <div class="w-full flex flex-col gap-3">
        <button
          class="btn-primary py-3.5 rounded-xl text-sm w-full"
          :disabled="loading"
          @click="register"
        >
          <span v-if="loading" class="material-symbols-outlined text-lg animate-spin">progress_activity</span>
          <template v-else>
            Create Owner Account
            <span class="material-symbols-outlined text-lg">arrow_forward</span>
          </template>
        </button>
        <NuxtLink
          to="/owner/login"
          class="py-2.5 w-full text-surface-onSurfaceVariant text-xs font-medium hover:text-primary transition-colors text-center"
        >
          Already have an account? <span class="text-primary font-semibold">Sign In</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useToast } from '~/composables/useToast'

useHead({ title: 'LaundryIn — Register Owner' })
definePageMeta({ layout: false })

const authStore = useAuthStore()
const { error: toastError } = useToast()

const router = useRouter()

// Guard: Redirect if already logged in
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

const form = ref({
  name: '',
  phone: '',
  password: '',
  confirmPassword: ''
})
const showPassword = ref(false)
const loading = ref(false)

const formattedPhone = computed(() => {
  const p = form.value.phone.trim()
  if (p.startsWith('0')) return '+62' + p.slice(1)
  if (p.startsWith('62') && !p.startsWith('+')) return '+' + p
  return p
})

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
    } else if (status === 400 && (apiMsg.includes('lemah') || apiMsg.includes('huruf besar'))) {
      toastError('Password harus mengandung huruf besar, huruf kecil, dan angka')
    } else if (status === 429) {
      toastError('Terlalu banyak percobaan, coba lagi beberapa saat')
    } else {
      toastError('Terjadi kesalahan, coba lagi')
    }
  } finally {
    loading.value = false
  }
}
</script>
