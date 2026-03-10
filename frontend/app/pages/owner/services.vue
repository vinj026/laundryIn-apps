<template>
  <div class="p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold">Services Catalog</h2>
        <p class="text-surface-onSurfaceVariant text-sm">Kelola produk & layanan laundry.</p>
      </div>
      <button class="bg-primary text-primary-text rounded-[1rem] px-4 h-11 flex items-center justify-center shadow-lg hover:brightness-110 active:scale-95 transition-all gap-2 font-bold text-sm">
        <span class="material-symbols-outlined text-[20px]">add</span>
        New
      </button>
    </div>

    <div v-if="pending" class="flex justify-center py-10">
        <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <!-- Data List -->
    <div v-else-if="services && services.data && services.data.length > 0" class="bg-surface-container rounded-3xl border border-outline/10 overflow-hidden divide-y divide-outline/10">
      
      <div 
        class="p-4 flex gap-4 items-center hover:bg-surface-onSurface/5 transition-colors" 
        v-for="service in services.data" 
        :key="service.id"
      >
        <div class="w-12 h-12 bg-primary-container text-primary-onContainer rounded-full flex items-center justify-center shrink-0">
          <span class="material-symbols-outlined">styler</span>
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
              <h3 class="font-bold text-sm truncate">{{ service.name }}</h3>
              <div v-if="!service.is_active" class="bg-red-500/20 text-red-500 text-[10px] font-bold px-1.5 py-0.5 rounded-sm">Inactive</div>
          </div>
          <div class="text-[11px] text-surface-onSurfaceVariant mt-0.5">
            Unit: <span class="font-bold">{{ service.unit }}</span>
          </div>
          <div class="font-mono text-primary font-bold text-sm mt-1">Rp {{ service.price }}</div>
        </div>
        <div class="flex flex-col gap-2 shrink-0">
          <button class="flex items-center justify-center h-7 w-7 rounded-sm bg-surface-containerHigh text-surface-onSurface hover:text-primary transition-colors">
            <span class="material-symbols-outlined text-[16px]">edit_square</span>
          </button>
          <button class="flex items-center justify-center h-7 w-7 rounded-sm bg-red-500/10 text-red-500 hover:bg-red-500/20 transition-colors">
            <span :class="service.is_active ? 'text-red-500' : 'text-green-500'" class="material-symbols-outlined text-[16px]">
                {{ service.is_active ? 'block' : 'check_circle' }}
            </span>
          </button>
        </div>
      </div>
      
    </div>
    
    <div v-else class="text-center py-10 opacity-70">
        You have no services listed.
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn - Services'
})

interface Service {
  id: string
  name: string
  description: string
  price: number
  unit: string
  is_active: boolean
}

// In real app we fetch owner services specifically or via outlet ID param / context
const { data: services, pending } = await useFetch<{ data: Service[] }>('/api/services')
</script>
