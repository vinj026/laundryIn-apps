<template>
  <div class="px-4 md:px-8 py-6 space-y-6">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 mb-2">
      <div>
        <h2 class="text-xl font-bold mb-1">Analytics Overview</h2>
        <p class="text-surface-onSurfaceVariant text-sm">Track your laundry business growth.</p>
      </div>

      <!-- Date Filter -->
      <div class="flex bg-surface-containerHigh rounded-xl p-1 shrink-0 w-max">
        <button
          v-for="f in filterOptions"
          :key="f.value"
          @click="daysFilter = f.value"
          class="px-4 py-1.5 rounded-lg text-sm font-semibold transition-all duration-normal"
          :class="daysFilter === f.value ? 'bg-surface shadow text-primary' : 'text-surface-onSurfaceVariant hover:text-surface-onSurface'"
        >
          {{ f.label }}
        </button>
      </div>
    </div>

    <div v-if="pending" class="flex justify-center py-20">
      <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <div v-else-if="error" class="flex flex-col items-center justify-center text-center py-24 animate-fade-in card">
       <div class="w-16 h-16 rounded-2xl bg-danger-muted text-danger flex items-center justify-center mb-4 border border-danger/30">
        <span class="material-symbols-outlined text-3xl">error</span>
      </div>
      <p class="text-sm font-medium mb-1">Gagal memuat data dashboard</p>
      <button @click="refresh()" class="btn-primary py-2 px-6 rounded-xl text-sm font-semibold mt-4">
        Coba Lagi
      </button>
    </div>

    <div v-else-if="analytics" class="space-y-6 animate-fade-in">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Revenue Card (Kiri) -->
        <div class="card !p-6 md:!p-8 relative overflow-hidden flex flex-col justify-center min-h-[160px]">
          <div class="absolute -right-8 -top-8 w-40 h-40 bg-primary/5 rounded-full blur-2xl"></div>
          <div class="absolute -bottom-10 -left-10 w-32 h-32 bg-secondary/5 rounded-full blur-2xl"></div>
          
          <div class="flex items-center gap-2 text-primary font-semibold uppercase tracking-wider text-[10px] mb-4 relative z-10">
            <span class="material-symbols-outlined text-[18px]">account_balance_wallet</span>
            Total Omzet
          </div>
          <h3 class="text-4xl md:text-5xl font-mono font-bold tracking-tight relative z-10">Rp {{ Number(analytics.totalRevenue || 0).toLocaleString('id-ID') }}</h3>
          
          <div class="mt-4 pt-4 border-t border-border flex items-center justify-between text-sm relative z-10">
            <span class="text-surface-onSurfaceVariant">Berdasarkan periode terpilih</span>
          </div>
        </div>

        <!-- Pipeline Status (Kanan) -->
        <div>
          <h3 class="text-sm font-bold mb-3">Pipeline Status</h3>
          <div class="grid grid-cols-2 gap-3 h-[calc(100%-32px)]">

            <div class="card !p-4 flex flex-col justify-between hover:border-surface-onSurfaceVariant/20 transition-colors">
              <div class="flex items-center justify-between mb-3 text-surface-onSurfaceVariant">
                <span class="material-symbols-outlined text-[20px]">inventory_2</span>
                <span class="badge badge-pending">Pending</span>
              </div>
              <p class="text-3xl font-mono font-bold">{{ analytics.ordersPending || 0 }}</p>
            </div>

            <div class="card !p-4 !bg-primary/[0.05] !border-primary/20 flex flex-col justify-between">
              <div class="flex items-center justify-between mb-3 text-primary">
                <span class="material-symbols-outlined text-[20px]">local_laundry_service</span>
                <span class="badge badge-process shadow-sm">Process</span>
              </div>
              <p class="text-3xl font-mono font-bold text-primary">{{ analytics.ordersProcess || 0 }}</p>
            </div>

            <div class="card !p-4 flex flex-col justify-between hover:border-surface-onSurfaceVariant/20 transition-colors">
              <div class="flex items-center justify-between mb-3 text-surface-onSurfaceVariant">
                <span class="material-symbols-outlined text-[20px] text-secondary">check_circle</span>
                <span class="badge badge-completed">Completed</span>
              </div>
              <p class="text-3xl font-mono font-bold">{{ analytics.ordersCompleted || 0 }}</p>
            </div>

            <div class="card !p-4 flex flex-col justify-between hover:border-surface-onSurfaceVariant/20 transition-colors">
              <div class="flex items-center justify-between mb-3 text-surface-onSurfaceVariant">
                <span class="material-symbols-outlined text-[20px]">hail</span>
                <span class="badge badge-picked-up">Picked Up</span>
              </div>
              <p class="text-3xl font-mono font-bold">{{ analytics.ordersPickedUp || 0 }}</p>
            </div>

          </div>
        </div>
      </div>

      <!-- Top Services (Bawah Full Width) -->
      <div>
        <h3 class="text-sm font-bold mb-3">Top Services</h3>
        <div class="card !p-0 overflow-hidden">
          <div v-if="analytics.topServices.length === 0" class="p-8 text-center text-surface-onSurfaceVariant text-sm">
            <span class="material-symbols-outlined text-3xl mb-2 opacity-50">data_alert</span>
            <p>Belum ada data untuk periode ini.</p>
          </div>
          <div class="divide-y divide-border" v-else>
            <div
              v-for="(service, index) in analytics.topServices"
              :key="index"
              class="p-4 flex items-center justify-between hover:bg-surface-containerHigh/30 transition-colors"
            >
              <div class="flex items-center gap-3 md:gap-4 min-w-0">
                <div class="w-8 h-8 rounded-full bg-surface-container flex items-center justify-center font-bold text-xs shrink-0 text-surface-onSurfaceVariant">
                  {{ index + 1 }}
                </div>
                <div class="min-w-0">
                  <p class="font-bold text-sm truncate text-surface-onSurface">{{ service.service_name }}</p>
                  <div class="flex items-center gap-1.5 mt-0.5">
                    <span class="material-symbols-outlined text-[12px] text-surface-onSurfaceVariant">store</span>
                    <p class="text-xs text-surface-onSurfaceVariant truncate">{{ service.outlet_name }}</p>
                  </div>
                </div>
              </div>
              <div class="text-right shrink-0 ml-4 flex flex-col items-end">
                <p class="font-mono font-bold text-sm text-primary">Rp {{ Number(service.total_revenue || 0).toLocaleString('id-ID') }}</p>
                <div class="badge badge-live mt-1 px-1.5 bg-surface-container border border-border">
                  <span class="text-[10px] font-semibold text-surface-onSurface">{{ Number(service.total_qty || 0).toFixed(1).replace('.0', '') }} sales</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { ref, computed } from 'vue'
import type { ApiResponse } from '~/types/api'

definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn — Owner Analytics'
})

const authStore = useAuthStore()
const router = useRouter()
const { error: toastError } = useToast()

const filterOptions = [
  { label: '7 Hari', value: 7 },
  { label: '30 Hari', value: 30 },
  { label: '3 Bulan', value: 90 }
]
const daysFilter = ref(30)

const getDateRange = (days: number) => {
  const end = new Date()
  const start = new Date()
  start.setDate(end.getDate() - days)
  
  const endStr = end.toISOString().split('T')[0]
  const startStr = start.toISOString().split('T')[0]
  
  return { start_date: startStr, end_date: endStr }
}

watchEffect(() => {
  if (import.meta.client) {
    if (!authStore.isLoggedIn || authStore.user?.role !== 'owner') {
      router.push('/owner/login')
    }
  }
})

interface TopServiceResponse {
  service_name: string
  outlet_name: string
  total_qty: string
  total_revenue: string
}

interface AnalyticsWrapper {
  totalRevenue: number
  ordersPending: number
  ordersProcess: number
  ordersCompleted: number
  ordersPickedUp: number
  topServices: TopServiceResponse[]
}

interface OmzetResponse {
  total_omzet: string
}

interface OrderStatusSummaryResponse {
  pending: number
  process: number
  completed: number
  picked_up: number
  cancelled: number
}

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

watchEffect(() => {
  if (error.value) {
    const status = error.value?.statusCode || error.value?.status || (error.value?.data as any)?.statusCode
    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else {
      toastError('Gagal memuat data dashboard')
    }
  }
})
</script>
