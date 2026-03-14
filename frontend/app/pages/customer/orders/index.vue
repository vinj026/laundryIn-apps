<template>
  <div class="h-full w-full p-6 lg:p-10 overflow-y-auto custom-scrollbar">
    <div class="max-w-3xl mx-auto">
      <div class="mb-8">
        <h1 class="text-2xl font-bold mb-1">My Orders</h1>
        <p class="text-surface-onSurfaceVariant text-sm">Track and view your recent laundry lists.</p>
      </div>

      <!-- Loading State -->
      <div v-if="pending" class="flex justify-center py-20">
        <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex flex-col items-center justify-center text-center py-24 animate-fade-in">
        <div class="w-16 h-16 rounded-2xl bg-danger-muted text-danger flex items-center justify-center mb-4 border border-danger/30">
          <span class="material-symbols-outlined text-3xl">error</span>
        </div>
        <p class="text-sm font-medium mb-1">Gagal memuat pesanan</p>
        <p class="text-xs text-surface-onSurfaceVariant mb-4">Terjadi kesalahan jaringan atau server.</p>
        <button @click="refresh()" class="btn-primary py-2 px-6 rounded-xl text-sm font-semibold">
          Coba Lagi
        </button>
      </div>

      <!-- Order List -->
      <div v-else-if="ordersList.length > 0" class="space-y-3">
        <div
          v-for="order in ordersList"
          :key="order.id"
          @click="navigateTo(`/customer/orders/${order.id}`)"
          class="card-interactive flex flex-col sm:flex-row sm:items-center gap-4"
        >
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-1.5">
              <span class="text-xs font-mono font-bold text-surface-onSurfaceVariant">#{{ order.id.slice(0,8).toUpperCase() }}</span>
              <span class="badge" :class="{
                'badge-pending': order.status === 'pending',
                'badge-process': order.status === 'process',
                'badge-completed': order.status === 'completed',
                'badge-picked-up': order.status === 'picked_up',
              }">
                {{ order.status === 'completed' ? 'Ready Pickup' : order.status === 'picked_up' ? 'Finished' : order.status === 'process' ? 'In Process' : order.status }}
              </span>
            </div>
            <h3 class="text-base font-bold mb-0.5">Outlet #{{ order.outlet_id.slice(0, 8).toUpperCase() }}</h3>
            <p class="text-xs text-surface-onSurfaceVariant">{{ new Date(order.order_date).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) }}</p>
          </div>

          <div class="flex items-center justify-between sm:flex-col sm:items-end gap-1">
            <span class="text-lg font-bold text-primary font-mono" :class="order.final_total_price ? 'text-success' : 'text-primary'">
              Rp {{ Number(order.final_total_price || order.total_price).toLocaleString('id-ID') }}
            </span>
            <span class="text-xs text-surface-onSurfaceVariant flex items-center gap-1 group-hover:text-primary transition-colors">
              Details <span class="material-symbols-outlined text-[14px]">arrow_forward</span>
            </span>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="flex flex-col items-center justify-center text-center py-24 animate-fade-in">
        <div class="w-16 h-16 rounded-2xl bg-surface-container flex items-center justify-center mb-4 border border-border">
          <span class="material-symbols-outlined text-3xl text-surface-onSurfaceVariant">receipt_long</span>
        </div>
        <p class="text-sm font-medium mb-1">Belum ada pesanan</p>
        <button @click="navigateTo('/customer')" class="btn-primary py-2 px-6 rounded-xl text-sm font-semibold mt-4">
          Cari Laundry
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useToast } from '~/composables/useToast'
import type { ApiResponse, PaginatedResponse } from '~/types/api'

definePageMeta({
  layout: 'customer'
})
useHead({
  title: 'LaundryIn — My Orders'
})

const authStore = useAuthStore()
const router = useRouter()
const { error: toastError } = useToast()

watchEffect(() => {
  if (import.meta.client && !authStore.isLoggedIn) {
    router.push('/customer/login?redirect=/customer/orders')
  }
})

interface Order {
  id: string
  status: string
  total_price: string
  final_total_price?: string
  outlet_id: string
  order_date: string
}

const { data: ordersResponse, pending, error, refresh } = await useApiFetch<PaginatedResponse<Order[]>>('/api/orders')

watchEffect(() => {
  if (error.value) {
    const status = error.value?.statusCode || error.value?.status || (error.value?.data as any)?.statusCode
    if (status === 401) {
      if (authStore.token) {
        toastError('Sesi kamu habis, silakan login ulang')
        authStore.logout()
      }
      router.push('/customer/login')
    } else {
      toastError('Gagal memuat pesanan')
    }
  }
})

const ordersList = computed(() => ordersWrapper.value?.data?.data ?? [])

onActivated(() => {
  refresh()
})
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
