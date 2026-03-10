<template>
  <div class="h-full w-full flex flex-col p-6 md:p-10 lg:p-12 overflow-y-auto custom-scrollbar">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-10">
      <div>
        <h1 class="text-3xl md:text-4xl font-bold mb-2">Find a Laundry</h1>
        <p class="text-surface-onSurfaceVariant">Discover the best laundry outlets around you.</p>
      </div>
      
      <!-- Search & Filter -->
      <div class="flex items-center gap-3">
        <div class="bg-surface-container rounded-full px-4 py-2.5 flex items-center gap-3 w-full md:w-64 border border-transparent focus-within:border-primary/50 transition-colors">
          <span class="material-symbols-outlined text-surface-onSurfaceVariant">search</span>
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Search outlet..." 
            class="bg-transparent border-none outline-none w-full text-sm placeholder-surface-onSurfaceVariant" 
          />
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="flex-1 flex items-center justify-center">
      <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <!-- Grid of Outlets -->
    <div v-else-if="filteredOutlets.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="outlet in filteredOutlets" 
        :key="outlet.id"
        class="bg-surface-container rounded-3xl p-6 group hover:ring-2 hover:ring-primary/50 transition-all cursor-pointer flex flex-col justify-between border border-transparent"
        @click="navigateTo(`/customer/outlet/${outlet.id}`)"
      >
        <div>
          <div class="flex items-start justify-between mb-4">
            <div class="h-14 w-14 rounded-2xl bg-primary/20 text-primary flex items-center justify-center">
              <span class="material-symbols-outlined text-3xl">storefront</span>
            </div>
            <div class="bg-green-500/20 text-green-500 font-bold text-xs px-2.5 py-1 rounded-full flex items-center gap-1">
              <div class="w-1.5 h-1.5 rounded-full bg-green-500"></div> Open
            </div>
          </div>
          
          <h2 class="text-xl font-bold mb-1 group-hover:text-primary transition-colors">{{ outlet.name }}</h2>
          <p class="text-sm text-surface-onSurfaceVariant line-clamp-2 mb-4">{{ outlet.address }}</p>
          <p class="text-xs text-surface-onSurfaceVariant opacity-80">{{ outlet.description }}</p>
        </div>

        <div class="mt-8 flex items-center justify-between pt-4 border-t border-outline/10">
          <div class="flex items-center gap-1 text-yellow-500">
            <span class="material-symbols-outlined text-[16px] text-yellow-400" style="font-variation-settings: 'FILL' 1;">star</span>
            <span class="text-sm font-bold">4.8</span>
          </div>
          <button class="text-sm font-bold text-primary flex items-center gap-1 group-hover:translate-x-1 transition-transform">
            Order Now <span class="material-symbols-outlined text-[18px]">arrow_forward</span>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Empty State -->
    <div v-else class="flex-1 flex flex-col items-center justify-center text-center opacity-70">
      <span class="material-symbols-outlined text-6xl mb-4">search_off</span>
      <p class="text-lg">No outlets found matching "{{ searchQuery }}"</p>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'customer'
})
useHead({
  title: 'LaundryIn - Find Outlet'
})

interface Outlet {
  id: string
  name: string
  address: string
  description: string
  is_active: boolean
}

const { data: outlets, pending } = await useFetch<{ data: Outlet[] }>('/api/outlets')
const searchQuery = ref('')

const filteredOutlets = computed(() => {
  if (!outlets.value?.data) return []
  const q = searchQuery.value.toLowerCase()
  return outlets.value.data.filter((o: Outlet) => 
    o.name.toLowerCase().includes(q) || o.address.toLowerCase().includes(q)
  )
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: 20px;
}
</style>
