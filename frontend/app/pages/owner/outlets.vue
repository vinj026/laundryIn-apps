<template>
  <div class="p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold">Outlet Management</h2>
        <p class="text-surface-onSurfaceVariant text-sm">Manage your laundry branch locations.</p>
      </div>
      <button class="bg-primary text-primary-text rounded-full h-12 w-12 flex items-center justify-center shadow-lg hover:brightness-110 active:scale-95 transition-all">
        <span class="material-symbols-outlined text-[24px]">add</span>
      </button>
    </div>

    <!-- Data Table Card -->
    <div class="bg-surface-container rounded-3xl border border-outline/10 overflow-hidden">
      <!-- Search Bar -->
      <div class="p-4 border-b border-outline/10">
        <div class="relative">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-surface-onSurfaceVariant text-[20px]">search</span>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Cari outlet..." 
            class="bg-surface-containerHigh w-full rounded-full py-2.5 pl-10 pr-4 text-sm outline-none focus:ring-1 focus:ring-primary border border-transparent transition-all">
        </div>
      </div>

      <div v-if="pending" class="flex justify-center py-10">
          <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
      </div>

      <!-- List -->
      <div v-else-if="filteredOutlets.length > 0" class="divide-y divide-outline/10">
        <div 
          class="p-4 flex flex-col gap-3 hover:bg-surface-onSurface/5 transition-colors" 
          v-for="outlet in filteredOutlets" 
          :key="outlet.id"
        >
          <div class="flex items-start justify-between">
            <div>
              <h3 class="font-bold text-primary">{{ outlet.name }}</h3>
              <p class="text-xs text-surface-onSurfaceVariant line-clamp-1 mt-0.5">{{ outlet.address }}</p>
            </div>
            <div class="flex gap-2">
              <button class="h-8 w-8 rounded-full bg-surface-containerHigh flex items-center justify-center text-surface-onSurface hover:text-primary transition-colors">
                <span class="material-symbols-outlined text-[18px]">edit</span>
              </button>
              <button class="h-8 w-8 rounded-full bg-red-500/10 flex items-center justify-center text-red-400 hover:bg-red-500/20 transition-colors">
                <span class="material-symbols-outlined text-[18px]">delete</span>
              </button>
            </div>
          </div>
          <div class="flex items-center gap-1 text-[11px] font-mono text-surface-onSurfaceVariant">
            <span class="material-symbols-outlined text-[14px] mr-1">{{ outlet.is_active ? 'check_circle' : 'cancel' }}</span>
            <span :class="outlet.is_active ? 'text-green-500' : 'text-red-500'">{{ outlet.is_active ? 'Active' : 'Inactive' }}</span>
          </div>
        </div>
      </div>
      
      <div v-else class="text-center py-10 opacity-70">
          No outlets found.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn - Outlets'
})

interface Outlet {
  id: string
  name: string
  address: string
  is_active: boolean
}

const { data: outletsWrapper, pending } = await useFetch<{data: Outlet[]}>('/api/outlets')
const searchQuery = ref('')

const filteredOutlets = computed(() => {
  if (!outletsWrapper.value?.data) return []
  const q = searchQuery.value.toLowerCase()
  return outletsWrapper.value.data.filter((o: Outlet) => 
    o.name.toLowerCase().includes(q) || o.address.toLowerCase().includes(q)
  )
})
</script>
