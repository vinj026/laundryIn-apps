<template>
  <div class="p-6 space-y-6">
    <div>
      <h2 class="text-2xl font-bold">Order Management</h2>
      <p class="text-surface-onSurfaceVariant text-sm">Pipeline & eksekusi status pesanan pelanggan.</p>
    </div>

    <!-- Active Orders List -->
    <div v-if="pending" class="flex justify-center py-10">
        <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <div v-else-if="orders && orders.data && orders.data.length > 0" class="space-y-4">
      
      <div 
        v-for="order in orders.data" 
        :key="order.id"
        :class="[
          'rounded-3xl border p-5 transition-all',
          order.status === 'process' ? 'bg-primary/5 border-primary/20' : 'bg-surface-container border-outline/10'
        ]"
      >
        <div class="flex justify-between items-start mb-4">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span v-if="order.status === 'process'" class="bg-primary text-primary-text text-xs font-bold px-2 py-0.5 rounded-sm flex items-center gap-1">
                <span class="w-1.5 h-1.5 bg-primary-text rounded-full animate-pulse"></span> PROCESS
              </span>
              <span v-else-if="order.status === 'pending'" class="bg-surface-containerHigh text-xs font-bold px-2 py-0.5 rounded-sm uppercase">{{ order.status }}</span>
              <span v-else-if="order.status === 'completed'" class="bg-green-500/20 text-green-500 text-xs font-bold px-2 py-0.5 rounded-sm uppercase">READY</span>
              <span v-else class="bg-surface-containerHigh/50 text-surface-onSurfaceVariant text-xs font-bold px-2 py-0.5 rounded-sm uppercase">{{ order.status }}</span>

              <span class="text-xs text-surface-onSurfaceVariant font-mono">#{{ order.id.slice(0,8).toUpperCase() }}</span>
            </div>
            <!-- Normally fetched from order.customer_id but dummy for UI -->
            <h3 class="font-bold text-lg">Customer_{{ order.customer_id.slice(0,4) }}</h3>
          </div>
          <span class="font-mono font-bold text-primary">Rp {{ Number(order.total_price).toLocaleString('id-ID') }}</span>
        </div>
        
        <div class="text-xs text-surface-onSurfaceVariant flex flex-col gap-1 mb-5">
          <p class="font-mono line-clamp-1">
             <span v-for="(item, idx) in order.Items" :key="item.id">
                {{ item.qty }}{{ item.unit }} {{ item.service_name }}{{ idx < order.Items.length - 1 ? ', ' : '' }}
             </span>
          </p>
          <p class="flex items-center gap-1 mt-1">
            <span class="material-symbols-outlined text-[14px]">location_on</span> {{ order.Outlet?.name }}
          </p>
        </div>

        <button 
           v-if="order.status === 'pending'"
           @click="updateStatus(order.id, 'process')"
           class="w-full py-3 bg-primary text-primary-text font-bold rounded-xl flex items-center justify-center gap-2 hover:brightness-110 transition-all"
        >
          <span class="material-symbols-outlined text-[18px]">play_arrow</span>
          Set to Process
        </button>

        <button 
           v-else-if="order.status === 'process'"
           @click="updateStatus(order.id, 'completed')"
           class="w-full py-3 bg-surface-containerHigh text-surface-onSurface font-bold border border-outline/20 rounded-xl flex items-center justify-center gap-2 hover:bg-surface-onSurface/10 transition-all"
        >
          <span class="material-symbols-outlined text-[18px]">check_circle</span>
          Complete Task
        </button>
        
        <button 
           v-else-if="order.status === 'completed'"
           @click="updateStatus(order.id, 'picked_up')"
           class="w-full py-3 bg-surface-containerHigh text-surface-onSurface font-bold border border-outline/20 rounded-xl flex items-center justify-center gap-2 hover:bg-surface-onSurface/10 transition-all"
        >
          <span class="material-symbols-outlined text-[18px]">hail</span>
          Mark as Picked Up
        </button>

      </div>
      
    </div>
    
    <div v-else class="text-center py-10 opacity-70">
        No active orders right now.
    </div>

  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn - Orders Pipeline'
})

interface OrderItem {
  id: string
  service_name: string
  qty: string
  unit: string
}

interface Order {
  id: string
  status: string
  total_price: number
  customer_id: string
  Items: OrderItem[]
  Outlet: {
    name: string
  }
}

// In real app, fetch /api/orders?outlet_id=xx for owners
const { data: orders, pending, refresh } = await useFetch<{ data: Order[] }>('/api/orders')

const updateStatus = async (id: string, newStatus: string) => {
   console.log(`FSM Transitioning Order ${id} -> ${newStatus}`)
   /*
   await $fetch(`/api/orders/${id}`, {
      method: 'PATCH',
      body: { status: newStatus }
   })
   refresh()
   */
}
</script>
