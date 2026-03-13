<template>
  <div class="h-full w-full overflow-y-auto custom-scrollbar px-6 py-6 lg:px-10 flex flex-col xl:flex-row gap-8 max-w-6xl mx-auto">

    <!-- Loading -->
    <div v-if="pending" class="flex-1 flex justify-center py-20">
      <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <template v-else-if="order">
      <!-- Left: Tracking -->
      <div class="flex-1 space-y-6 max-w-3xl">
        <!-- Breadcrumbs -->
        <div>
          <div class="flex items-center gap-1.5 mb-2">
            <NuxtLink to="/customer/orders" class="text-xs font-medium text-surface-onSurfaceVariant hover:text-primary transition-colors">My Orders</NuxtLink>
            <span class="material-symbols-outlined text-xs text-surface-onSurfaceVariant">chevron_right</span>
            <span class="text-xs font-medium">Tracking</span>
          </div>
          <h1 class="text-3xl font-bold mb-1">Order #{{ order.id.slice(0, 8).toUpperCase() }}</h1>
          <p class="text-sm text-surface-onSurfaceVariant">
            Placed {{ new Date(order.order_date).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
          </p>
        </div>

        <!-- Live Status Banner -->
        <div class="card !bg-primary/10 !border-primary/20 p-6 flex flex-col sm:flex-row items-center gap-6 relative overflow-hidden">
          <div class="flex items-center justify-center h-16 w-16 rounded-2xl bg-primary/20 shrink-0">
            <span class="material-symbols-outlined text-primary text-[36px]" :class="order.status === 'process' ? 'animate-spin' : ''" style="animation-duration: 4s;">
              {{ statusConfig[order.status]?.icon || 'local_laundry_service' }}
            </span>
          </div>
          <div class="flex-1 text-center sm:text-left">
            <h2 class="text-xl font-bold text-primary mb-1">{{ statusConfig[order.status]?.label || order.status }}</h2>
            <p class="text-surface-onSurfaceVariant text-sm leading-relaxed max-w-md">
              {{ statusConfig[order.status]?.description || '' }}
            </p>
          </div>
        </div>

        <!-- Timeline -->
        <div>
          <h3 class="text-base font-bold mb-5">Tracking History</h3>
          <div class="relative ml-[68px] sm:ml-[116px] space-y-0">
            <div v-for="(step, idx) in timelineSteps" :key="idx" class="relative group" :class="{ 'opacity-50 pb-6': !step.active, 'pb-8': step.active, 'pb-2': idx === timelineSteps.length - 1 }">
              <!-- Connecting Line -->
              <div v-if="idx !== timelineSteps.length - 1" class="absolute left-0 top-6 bottom-[-24px] w-0.5 -ml-px"
                :class="step.active && !step.current ? 'bg-primary' : 'bg-border'">
              </div>

              <!-- Timeline Dot -->
              <div class="absolute left-0 top-1 -ml-[7px] flex items-center justify-center bg-surface">
                <div v-if="step.current" class="h-3.5 w-3.5 rounded-full bg-primary shadow-[0_0_8px_rgba(45,212,191,0.5)] animate-pulse relative z-10 flex items-center justify-center">
                  <div class="absolute -inset-1.5 rounded-full bg-primary/20 animate-ping"></div>
                </div>
                <div v-else-if="step.active" class="h-3.5 w-3.5 rounded-full bg-primary relative z-10 flex items-center justify-center">
                  <span class="material-symbols-outlined text-[10px] text-primary-text font-bold" style="font-variation-settings: 'wght' 700;">check</span>
                </div>
                <div v-else class="h-3.5 w-3.5 rounded-full border-2 border-border bg-surface relative z-10"></div>
              </div>

              <div class="flex flex-col sm:flex-row sm:items-baseline gap-2 sm:gap-6 w-full pl-6 relative">
                 <!-- Time (Moves left on sm screens layout shift) -->
                <div class="absolute left-[-68px] sm:left-[-116px] top-0 w-12 sm:w-[90px] text-right shrink-0 pr-2">
                  <span v-if="step.active" class="text-[11px] font-bold font-mono" :class="step.current ? 'text-primary' : 'text-surface-onSurfaceVariant'">
                    {{ step.time }}
                  </span>
                </div>

                <div class="w-full transition-all" :class="step.current ? 'card !border-primary/30 !bg-primary/5 -mt-1 shadow-md' : ''">
                  <div class="flex justify-between items-start" :class="step.current ? 'mb-1.5' : 'mb-0.5'">
                    <h4 class="text-sm transition-colors" :class="step.current ? 'font-bold text-primary' : (step.active ? 'font-bold text-surface-onSurface' : 'font-medium text-surface-onSurfaceVariant')">{{ step.label }}</h4>
                    <span v-if="step.current" class="badge badge-live ml-2">
                      <span class="w-1.5 h-1.5 bg-danger rounded-full animate-pulse mr-1"></span> LIVE
                    </span>
                  </div>
                  <p class="text-xs transition-colors" :class="step.current ? 'text-surface-onSurface' : 'text-surface-onSurfaceVariant'">{{ step.description }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Receipt -->
      <div class="w-full xl:w-[320px] flex flex-col gap-4 shrink-0 pb-8 mt-2 md:mt-0">
        <div class="card !p-5">
          <div class="flex items-center justify-between mb-4 pb-3 border-b border-border">
            <h3 class="text-sm font-bold">Order Summary</h3>
            <span class="material-symbols-outlined text-surface-onSurfaceVariant text-[18px]">receipt_long</span>
          </div>

          <div class="space-y-4 mb-5">
            <div v-for="item in order.items" :key="item.id" class="flex justify-between items-start text-sm">
              <div class="flex flex-col gap-0.5 mr-3">
                <span class="font-medium text-surface-onSurface">
                  <template v-if="item.actual_qty">
                    <span class="line-through text-surface-onSurfaceVariant text-xs mr-1" v-if="Number(item.actual_qty) !== Number(item.qty)">{{ Number(item.qty).toFixed(1).replace('.0', '') }}</span>
                    <span>{{ Number(item.actual_qty).toFixed(1).replace('.0', '') }}</span>
                  </template>
                  <template v-else>
                    {{ Number(item.qty).toFixed(1).replace('.0', '') }}
                  </template>
                  {{ item.unit }} &middot; {{ item.service_name }}
                </span>
              </div>
              <div class="flex flex-col items-end">
                <span v-if="item.final_price && Number(item.final_price) !== Number(item.subtotal || (Number(item.qty) * Number(item.service_price)))" class="text-[10px] text-surface-onSurfaceVariant line-through">
                  Rp {{ Number(item.subtotal || (Number(item.qty) * Number(item.service_price))).toLocaleString('id-ID') }}
                </span>
                <span class="font-medium font-mono text-surface-onSurface shrink-0" :class="item.final_price ? 'text-success' : ''">
                  Rp {{ Number(item.final_price || item.subtotal || (Number(item.qty) * Number(item.service_price))).toLocaleString('id-ID') }}
                </span>
              </div>
            </div>
          </div>

          <div class="pt-4 border-t border-border border-dashed">
            <div class="flex justify-between items-end">
              <div class="flex flex-col gap-1">
                <span class="text-sm font-medium text-surface-onSurfaceVariant">Total Payment</span>
                <span v-if="order.final_total_price" class="badge badge-completed px-1.5 py-0.5 rounded-md text-[10px] w-max font-bold">Harga Final</span>
                <span v-else class="badge bg-surface-containerHigh text-surface-onSurfaceVariant border-border border px-1.5 py-0.5 rounded-md text-[10px] w-max font-bold">Estimasi</span>
              </div>
              <div class="flex flex-col items-end">
                <span v-if="order.final_total_price && Number(order.final_total_price) !== Number(order.total_price)" 
                      class="text-xs text-surface-onSurfaceVariant line-through mb-0.5">
                  Rp {{ Number(order.total_price).toLocaleString('id-ID') }}
                </span>
                <span class="text-lg font-bold font-mono tracking-tight" :class="order.final_total_price ? 'text-success' : 'text-primary'">
                  Rp {{ Number(order.final_total_price || order.total_price).toLocaleString('id-ID') }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Outlet info -->
        <div class="card !p-5">
          <div class="flex items-start gap-3 text-sm">
            <div class="w-8 h-8 rounded-full bg-primary/10 text-primary flex items-center justify-center shrink-0 mt-0.5">
              <span class="material-symbols-outlined text-[16px]">storefront</span>
            </div>
            <div class="flex flex-col gap-1 min-w-0">
              <span class="font-bold text-surface-onSurface line-clamp-1">{{ outlet?.name || `Outlet #${order.outlet_id.slice(0, 8).toUpperCase()}` }}</span>
              <span class="text-xs text-surface-onSurfaceVariant leading-relaxed">{{ outlet?.address || 'Address not available' }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Not Found -->
    <div v-else class="flex-1 flex flex-col items-center justify-center text-center py-24">
      <span class="material-symbols-outlined text-4xl text-surface-onSurfaceVariant mb-3">search_off</span>
      <p class="text-sm font-medium mb-1">Order not found</p>
      <NuxtLink to="/customer/orders" class="mt-4 text-primary text-sm font-semibold hover:underline">Back to Orders</NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { computed } from 'vue'
import type { ApiResponse, PaginatedResponse } from '~/types/api'

definePageMeta({
  layout: 'customer'
})

const route = useRoute()
const authStore = useAuthStore()
const orderId = route.params.id as string

useHead({
  title: `LaundryIn — Order #${orderId.slice(0, 8).toUpperCase()}`
})

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

interface Order {
  id: string
  status: string
  total_price: string
  final_total_price?: string
  outlet_id: string
  order_date: string
  items: OrderItem[]
}

interface Outlet {
  id: string
  name: string
  address: string
  phone: string
}

const { data: ordersWrapper, pending } = await useFetch<ApiResponse<PaginatedResponse<Order[]>>>(
  '/api/orders?limit=100&page=1',
  { headers: { Authorization: authStore.authHeader }, server: false }
)

const order = computed(() => (ordersWrapper.value?.data?.data ?? []).find((o) => o.id === orderId) || null)

const { data: outletWrapper } = await useAsyncData<ApiResponse<Outlet> | null>(
  'customer-order-outlet',
  async () => {
    if (!order.value) return null
    return $fetch(`/api/public/outlets/${order.value.outlet_id}`)
  },
  { watch: [order], server: false }
)

const outlet = computed(() => outletWrapper.value?.data || null)

// Status config for the live banner
const statusConfig: Record<string, { icon: string; label: string; description: string }> = {
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

// Build timeline steps based on order status
const statusOrder = ['pending', 'process', 'completed', 'picked_up']
const stepLabels: Record<string, { label: string; description: string }> = {
  pending: { label: 'Pesanan Dibuat', description: 'Order placed via App.' },
  process: { label: 'Sedang Diproses', description: 'Items being washed and dried.' },
  completed: { label: 'Siap Diambil', description: 'Ready for pickup at outlet.' },
  picked_up: { label: 'Selesai', description: 'Customer has picked up the order.' }
}

const timelineSteps = computed(() => {
  if (!order.value) return []
  const currentIdx = statusOrder.indexOf(order.value.status)
  const createdDate = new Date(order.value.order_date)

  return statusOrder.map((status, idx) => ({
    label: stepLabels[status].label,
    description: stepLabels[status].description,
    active: idx <= currentIdx,
    current: idx === currentIdx,
    time: idx === 0
      ? createdDate.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
      : null
  }))
})
</script>
