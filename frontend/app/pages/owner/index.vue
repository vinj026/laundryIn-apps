<template>
  <div class="p-6 md:p-8 space-y-8">
    <div class="flex items-center justify-between mb-10">
      <div>
        <h2 class="text-2xl font-bold">Analytics Overview</h2>
        <p class="text-surface-onSurfaceVariant text-sm">Track your laundry business growth.</p>
      </div>
    </div>

    <div v-if="pending" class="flex justify-center py-20">
      <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <div v-else-if="analytics">
      <!-- Revenue Card -->
      <div class="bg-surface-container rounded-3xl p-8 border border-outline/10 shadow-[0_4px_30px_rgba(0,229,255,0.03)] relative overflow-hidden mb-8">
        <div class="absolute -right-4 -top-4 w-32 h-32 bg-primary/10 rounded-full blur-2xl"></div>
        <div class="flex items-center gap-2 text-primary font-bold mb-2 uppercase tracking-wide text-xs">
          <span class="material-symbols-outlined text-[18px]">account_balance_wallet</span>
          Total Omzet
        </div>
        <h3 class="text-4xl md:text-5xl font-mono font-bold tracking-tighter mt-2">Rp {{ Number(analytics.totalRevenue || 0).toLocaleString('id-ID') }}</h3>
      </div>

      <!-- Order Status Summary -->
      <div class="mb-8">
        <h3 class="text-lg font-medium mb-4">Pipeline Status</h3>
        <div class="grid grid-cols-2 gap-4">
          
          <div class="bg-surface-container rounded-2xl p-5 border border-outline/5 hover:border-primary/30 transition-colors">
            <div class="flex items-center justify-between mb-3 text-surface-onSurfaceVariant">
              <span class="material-symbols-outlined">inventory_2</span>
              <span class="text-xs font-bold bg-surface-containerHigh px-2 py-0.5 rounded-full">Pending</span>
            </div>
            <p class="text-3xl font-mono font-bold">{{ analytics.ordersPending || 0 }}</p>
          </div>

          <div class="bg-primary/10 rounded-2xl p-5 border border-primary/20">
            <div class="flex items-center justify-between mb-3 text-primary">
              <span class="material-symbols-outlined">local_laundry_service</span>
              <span class="text-xs font-bold bg-primary px-2 py-0.5 rounded-full text-primary-text">Process</span>
            </div>
            <p class="text-3xl font-mono font-bold text-primary">{{ analytics.ordersProcess || 0 }}</p>
          </div>

          <div class="bg-surface-container rounded-2xl p-5 border border-outline/5 hover:border-primary/30 transition-colors">
            <div class="flex items-center justify-between mb-3 text-surface-onSurfaceVariant">
              <span class="material-symbols-outlined">check_circle</span>
              <span class="text-xs font-bold bg-surface-containerHigh px-2 py-0.5 rounded-full">Completed</span>
            </div>
            <p class="text-3xl font-mono font-bold">{{ analytics.ordersCompleted || 0 }}</p>
          </div>

          <div class="bg-surface-container rounded-2xl p-5 border border-outline/5 hover:border-primary/30 transition-colors">
            <div class="flex items-center justify-between mb-3 text-surface-onSurfaceVariant">
              <span class="material-symbols-outlined">hail</span>
              <span class="text-xs font-bold bg-surface-containerHigh px-2 py-0.5 rounded-full">Picked Up</span>
            </div>
            <p class="text-3xl font-mono font-bold">{{ analytics.ordersPickedUp || 0 }}</p>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn - Owner Analytics'
})

interface AnalyticsWrapper {
  totalRevenue: number
  ordersPending: number
  ordersProcess: number
  ordersCompleted: number
  ordersPickedUp: number
}

// Mocks request or expects an implemented API.
const { data: wrapper, pending } = await useFetch<{data: AnalyticsWrapper}>('/api/analytics')
const analytics = computed(() => wrapper.value?.data)
</script>
