<template>
  <div class="h-full w-full flex flex-col overflow-hidden relative">

    <!-- Header -->
    <div class="px-6 lg:px-10 pt-6 pb-4">
      <div class="max-w-5xl mx-auto flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold mb-1">Find a Laundry</h1>
          <p class="text-surface-onSurfaceVariant text-sm">Discover the best laundry outlets around you.</p>
        </div>

        <!-- Search Bar -->
        <div class="flex items-center gap-3">
          <div class="bg-surface-container rounded-xl px-4 py-2.5 flex items-center gap-3 w-full md:w-64 border border-border focus-within:border-primary/50 transition-all duration-normal">
            <span class="material-symbols-outlined text-surface-onSurfaceVariant text-[20px]">search</span>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search outlet..."
              class="bg-transparent border-none outline-none w-full text-sm text-surface-onSurface placeholder-surface-onSurfaceVariant/60"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Content Area -->
    <div class="flex-1 overflow-y-auto custom-scrollbar px-6 lg:px-10 pb-6">
      <div class="max-w-5xl mx-auto">

        <!-- Loading State -->
        <div v-if="pending" class="flex items-center justify-center py-20">
          <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
        </div>

        <!-- Outlet Count -->
        <p v-if="!pending && filteredOutlets.length > 0" class="text-xs text-surface-onSurfaceVariant mb-4">
          Showing <span class="text-surface-onSurface font-semibold">{{ filteredOutlets.length }}</span> outlets
        </p>

        <!-- Grid of Outlets -->
        <div v-if="!pending && filteredOutlets.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          <div
            v-for="outlet in filteredOutlets"
            :key="outlet.id"
            class="card-interactive group"
            @click="navigateTo(`/customer/outlet/${outlet.id}`)"
          >
            <div class="flex items-start gap-3.5 mb-3">
              <div class="h-11 w-11 rounded-xl bg-primary/10 text-primary flex items-center justify-center shrink-0 group-hover:bg-primary/15 transition-all">
                <span class="material-symbols-outlined text-[22px]">storefront</span>
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-0.5">
                  <h2 class="text-sm font-bold group-hover:text-primary transition-colors truncate">{{ outlet.name }}</h2>
                  <div class="bg-success-muted text-success font-bold text-[9px] px-2 py-0.5 rounded-full flex items-center gap-1 shrink-0">
                    <div class="w-1 h-1 rounded-full bg-success animate-pulse-soft"></div>
                    Open
                  </div>
                </div>
                <div class="flex items-center gap-1 text-surface-onSurfaceVariant">
                  <span class="material-symbols-outlined text-[13px]">location_on</span>
                  <p class="text-xs line-clamp-1">{{ outlet.address }}</p>
                </div>
              </div>
            </div>

            <!-- Footer -->
            <div class="flex items-center justify-between pt-3 border-t border-border">
              <div class="flex items-center gap-2.5">
                <div class="flex items-center gap-1">
                  <span class="material-symbols-outlined text-[13px] text-warning" style="font-variation-settings: 'FILL' 1;">star</span>
                  <span class="text-xs font-bold">4.8</span>
                </div>
                <div class="w-0.5 h-0.5 rounded-full bg-border-hover"></div>
                <div class="flex items-center gap-1 text-surface-onSurfaceVariant">
                  <span class="material-symbols-outlined text-[12px]">call</span>
                  <span class="text-[11px]">{{ outlet.phone || 'Available' }}</span>
                </div>
              </div>
              <span class="text-primary text-xs font-semibold flex items-center gap-1 group-hover:gap-2 transition-all">
                Order
                <span class="material-symbols-outlined text-[14px]">arrow_forward</span>
              </span>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="!pending" class="flex flex-col items-center justify-center text-center py-20">
          <div class="w-16 h-16 rounded-2xl bg-surface-container flex items-center justify-center mb-4 border border-border">
            <span class="material-symbols-outlined text-3xl text-surface-onSurfaceVariant">search_off</span>
          </div>
          <p class="text-sm font-medium mb-1">No outlets found</p>
          <p class="text-xs text-surface-onSurfaceVariant">Try adjusting your search for "{{ searchQuery }}"</p>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onActivated } from 'vue'

definePageMeta({
  layout: 'customer'
})
useHead({
  title: 'LaundryIn — Find Outlet'
})

interface Outlet {
  id: string
  name: string
  address: string
  description: string
  phone: string
  is_active: boolean
}

const { data: outletsResponse, pending, refresh } = await useFetch<{ data: { data: Outlet[] } }>('/api/public/outlets')
const searchQuery = ref('')

onActivated(() => {
  refresh()
})

const filteredOutlets = computed(() => {
  if (!outletsResponse.value?.data?.data) return []
  const q = searchQuery.value.toLowerCase()
  return outletsResponse.value.data.data.filter((o: Outlet) =>
    o.name.toLowerCase().includes(q) || o.address.toLowerCase().includes(q)
  )
})
</script>
