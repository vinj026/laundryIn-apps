<template>
  <div class="mx-auto max-w-[1200px] w-full px-4 md:px-8 py-6 space-y-5">
    <div>
      <h2 class="text-xl font-bold mb-0.5">Order Management</h2>
      <p class="text-surface-onSurfaceVariant text-sm">Pipeline & eksekusi status pesanan pelanggan.</p>
    </div>

    <!-- Filters Container -->
    <div class="space-y-4">
      <div v-if="outlets.length > 0" class="card !p-4">
        <label class="text-xs font-semibold text-surface-onSurfaceVariant block mb-2">Outlet</label>
        <select
          v-model="selectedOutletId"
          class="bg-surface-containerHigh w-full rounded-xl py-2.5 px-3 text-sm outline-none focus:ring-1 focus:ring-primary/30 border border-border transition-all"
        >
          <option v-for="o in outlets" :key="o.id" :value="o.id">
            {{ o.name }}
          </option>
        </select>
      </div>

      <!-- Status Tabs -->
      <div class="flex overflow-x-auto custom-scrollbar gap-2 pb-2 -mb-2">
        <button
          v-for="status in ['Semua', 'pending', 'process', 'completed', 'picked_up', 'cancelled']"
          :key="status"
          @click="statusFilter = status"
          class="whitespace-nowrap px-4 py-2 rounded-full text-xs font-semibold transition-colors border"
          :class="statusFilter === status ? 'bg-primary text-primary-text border-primary' : 'bg-surface border-border text-surface-onSurface hover:bg-surface-containerHigh'"
        >
          {{ status === 'Semua' ? 'Semua Status' : status.toUpperCase() }}
        </button>
      </div>
    </div>

    <div v-if="pending" class="flex justify-center py-10">
      <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex flex-col items-center justify-center text-center py-16 animate-fade-in card">
      <div class="w-14 h-14 rounded-2xl bg-danger-muted text-danger flex items-center justify-center mb-3 border border-danger/30">
        <span class="material-symbols-outlined text-2xl">error</span>
      </div>
      <p class="text-sm font-medium mb-1">Gagal memuat pesanan</p>
      <button @click="refresh()" class="btn-primary py-1.5 px-5 rounded-xl text-xs font-semibold mt-3">
        Coba Lagi
      </button>
    </div>

    <div v-else-if="filteredOrders.length > 0" class="space-y-3">
      <div
        v-for="order in filteredOrders"
        :key="order.id"
        class="card"
        :class="order.status === 'process' ? '!bg-primary/[0.04] !border-primary/20' : ''"
      >
        <div class="flex justify-between items-start mb-3">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="badge" :class="{
                'badge-pending': order.status === 'pending',
                'badge-process': order.status === 'process',
                'badge-completed': order.status === 'completed',
                'badge-picked-up': order.status === 'picked_up',
                'bg-surface-containerHigh text-surface-onSurfaceVariant border-border': order.status === 'cancelled'
              }">
              <span v-if="order.status === 'process'" class="w-1 h-1 bg-primary rounded-full animate-pulse"></span>
                {{ order.status === 'completed' ? 'READY' : order.status.toUpperCase() }}
              </span>
              <span class="text-xs text-surface-onSurfaceVariant font-mono">#{{ order.id.slice(0,8).toUpperCase() }}</span>
            </div>
            <h3 class="font-bold text-base">{{ order.customer_name || (order as any).user?.name || (order as any).Customer?.name || `Customer_${order.user_id.slice(0,4)}` }}</h3>
          </div>
          <div class="text-right">
            <span v-if="order.final_total_price" class="font-mono font-bold text-success text-sm block">Rp {{ Number(order.final_total_price).toLocaleString('id-ID') }}</span>
            <span v-else class="font-mono font-bold text-primary text-sm block">Rp {{ Number(order.total_price).toLocaleString('id-ID') }}</span>
            <span v-if="!order.final_total_price" class="text-[10px] text-surface-onSurfaceVariant block uppercase tracking-wide">Estimasi</span>
          </div>
        </div>

        <div class="text-xs text-surface-onSurfaceVariant flex flex-col gap-1 mb-4">
          <p class="font-mono line-clamp-1">
            <span v-for="(item, idx) in order.items" :key="item.id">
              {{ Number(item.actual_qty || item.qty).toFixed(1).replace('.0', '') }}{{ item.unit }} {{ item.service_name }}{{ idx < order.items.length - 1 ? ', ' : '' }}
            </span>
          </p>
          <p class="flex items-center gap-1 mt-0.5">
            <span class="material-symbols-outlined text-[13px]">location_on</span> {{ selectedOutlet?.name || order.outlet_id }}
          </p>
        </div>

        <!-- Input Berat Aktual -->
        <div v-if="order.status === 'pending' && order.items.some(i => i.unit === 'KG')" class="mb-4 bg-surface-containerHigh/50 p-3 rounded-xl border border-border">
          <p class="text-xs font-semibold mb-2">Input Berat Aktual (KG)</p>
          <div class="space-y-2">
            <div v-for="item in order.items.filter(i => i.unit === 'KG')" :key="item.id" class="flex items-center justify-between gap-2">
              <span class="text-xs text-surface-onSurfaceVariant truncate">{{ item.service_name }} (Est: {{ Number(item.qty).toFixed(1).replace('.0', '') }} {{ item.unit }})</span>
              <div class="relative w-24 shrink-0">
                <input 
                  v-if="actualQtyInputs[order.id]"
                  type="number" 
                  step="0.1" 
                  min="0.1"
                  v-model="actualQtyInputs[order.id][item.id]" 
                  class="w-full bg-surface text-xs rounded-lg py-1.5 pr-7 pl-2 outline-none focus:ring-1 focus:ring-primary/50 text-right border border-border"
                  placeholder="0.0"
                >
                <span class="absolute right-2 top-1/2 -translate-y-1/2 text-[10px] text-surface-onSurfaceVariant font-medium pointer-events-none">KG</span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex gap-2" v-if="order.status === 'pending' || order.status === 'process'">
           <button
            v-if="order.status === 'pending'"
            @click="updateStatus(order.id, 'process')"
            :disabled="updatingOrder === order.id || hasMissingKgInput(order)"
            class="btn-primary flex-1 py-2.5 rounded-xl text-sm disabled:opacity-50"
          >
            <span class="material-symbols-outlined text-[18px]" v-if="updatingOrder !== order.id">play_arrow</span>
            <span class="material-symbols-outlined text-[18px] animate-spin" v-else>progress_activity</span>
            Proses
          </button>

          <button
            v-else-if="order.status === 'process'"
            @click="updateStatus(order.id, 'completed')"
            :disabled="updatingOrder === order.id"
            class="btn-secondary flex-1 py-2.5 rounded-xl text-sm disabled:opacity-50"
          >
            <span class="material-symbols-outlined text-[18px]" v-if="updatingOrder !== order.id">check_circle</span>
            <span class="material-symbols-outlined text-[18px] animate-spin" v-else>progress_activity</span>
            Selesai
          </button>

          <button
             @click="updateStatus(order.id, 'cancelled')"
             :disabled="updatingOrder === order.id"
             class="btn-secondary py-2.5 px-3 rounded-xl text-sm text-danger hover:bg-danger-muted border-danger/30 disabled:opacity-50 flex items-center justify-center shrink-0"
          >
             <span class="material-symbols-outlined text-[18px]" v-if="updatingOrder !== order.id">cancel</span>
             <span class="material-symbols-outlined text-[18px] animate-spin" v-else>progress_activity</span>
          </button>
        </div>

        <button
          v-else-if="order.status === 'completed'"
          @click="updateStatus(order.id, 'picked_up')"
          :disabled="updatingOrder === order.id"
          class="btn-secondary w-full py-2.5 rounded-xl text-sm disabled:opacity-50"
        >
          <span class="material-symbols-outlined text-[18px]" v-if="updatingOrder !== order.id">hail</span>
          <span class="material-symbols-outlined text-[18px] animate-spin" v-else>progress_activity</span>
          Sudah Diambil
        </button>

        <!-- Final State Label -->
        <div v-else-if="order.status === 'picked_up' || order.status === 'cancelled'" class="py-2.5 px-3 rounded-xl bg-surface-container border border-border text-surface-onSurfaceVariant text-xs font-semibold flex items-center justify-center gap-1.5 mt-1 opacity-70">
          <span class="material-symbols-outlined text-[16px]">{{ order.status === 'cancelled' ? 'block' : 'done_all' }}</span>
          Selesai ({{ order.status === 'cancelled' ? 'Dibatalkan' : 'Diambil' }})
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col items-center text-center py-10 text-surface-onSurfaceVariant text-sm animate-fade-in">
      <div class="w-14 h-14 rounded-2xl bg-surface-container flex items-center justify-center mb-3 border border-border">
        <span class="material-symbols-outlined text-2xl text-surface-onSurfaceVariant">view_list</span>
      </div>
      Belum ada pesanan masuk
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useToast } from '~/composables/useToast'
import { ref, computed, watchEffect, watch } from 'vue'
import type { ApiResponse, PaginatedResponse } from '~/types/api'

definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn — Orders Pipeline'
})

const authStore = useAuthStore()
const router = useRouter()
const { success: toastSuccess, error: toastError } = useToast()

const statusFilter = ref('Semua')

watchEffect(() => {
  if (import.meta.client) {
    if (!authStore.isLoggedIn || authStore.user?.role !== 'owner') {
      router.push('/owner/login')
    }
  }
})

interface Outlet {
  id: string
  name: string
}

interface OrderItem {
  id: string
  service_name: string
  qty: string
  actual_qty?: string
  unit: string
}

interface Order {
  id: string
  status: string
  total_price: string
  final_total_price?: string
  user_id: string
  outlet_id: string
  customer_name?: string
  items: OrderItem[]
}

const { data: outletsWrapper } = await useApiFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/outlets', {
  server: false
})

const outlets = computed(() => outletsWrapper.value?.data?.data ?? [])
const selectedOutletId = ref('')

watch(
  outlets,
  (list) => {
    if (!selectedOutletId.value && list.length > 0) {
      selectedOutletId.value = list[0].id
    }
  },
  { immediate: true }
)

const selectedOutlet = computed(() => outlets.value.find((o) => o.id === selectedOutletId.value) || null)

const { data: ordersWrapper, pending, error, refresh } = await useAsyncData<ApiResponse<PaginatedResponse<Order[]>>>(
  'owner-orders',
  async () => {
    if (!selectedOutletId.value) {
      return { status: 'success', message: '', data: { data: [], page: 1, limit: 10, total: 0, total_pages: 0 } }
    }
    return useApiRaw<ApiResponse<PaginatedResponse<Order[]>>>(`/api/outlets/${selectedOutletId.value}/orders`, {
      params: { page: ordersPage.value, limit: 10 }
    })
  },
  { watch: [selectedOutletId], server: false }
)

const orders = computed(() => ordersWrapper.value?.data?.data ?? [])

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

const filteredOrders = computed(() => {
  if (statusFilter.value === 'Semua') return orders.value
  return orders.value.filter(o => o.status === statusFilter.value)
})

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

watchEffect(() => {
  if (error.value) {
    const status = error.value?.statusCode || error.value?.status || (error.value?.data as any)?.statusCode
    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else {
      toastError('Gagal memuat pesanan')
    }
  }
})

const updatingOrder = ref<string | null>(null)

const updateStatus = async (id: string, newStatus: string) => {
  if (updatingOrder.value) return // Prevent double clicks
  updatingOrder.value = id

  try {
    const body: any = { status: newStatus }
    
    if (newStatus === 'process') {
      const order = orders.value.find(o => o.id === id)
      if (order) {
        const kgItems = order.items.filter(i => i.unit === 'KG')
        if (kgItems.length > 0) {
          body.items = kgItems.map(i => ({
            id: i.id,
            actual_qty: String(actualQtyInputs.value[id]?.[i.id] || "0")
          }))
        }
      }
    }

    await useApiRaw(`/api/orders/${id}/status`, {
      method: 'PATCH',
      body
    })
    await refresh()
    toastSuccess(`Status pesanan diperbarui`)
  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    const apiMsg = err?.data?.message || ''

    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else if (status === 404) {
      toastError('Pesanan tidak ditemukan')
    } else if (status === 400) {
      toastError(apiMsg || 'Data tidak valid') 
    } else {
      toastError('Gagal memperbarui status, coba lagi')
    }
  } finally {
    updatingOrder.value = null
  }
}

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

/* Hilangkan spinner di Chrome, Safari, Edge */
input[type=number]::-webkit-inner-spin-button,
input[type=number]::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

/* Hilangkan spinner di Firefox */
input[type=number] {
  -moz-appearance: textfield;
}
</style>
