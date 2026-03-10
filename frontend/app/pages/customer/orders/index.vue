<template>
  <div class="h-full w-full p-6 md:p-10 lg:p-12 overflow-y-auto custom-scrollbar">
    <div class="max-w-4xl mx-auto">
      <div class="flex items-center justify-between mb-10">
        <div>
          <h1 class="text-3xl md:text-4xl font-bold mb-2">My Orders</h1>
          <p class="text-surface-onSurfaceVariant">Track and view your recent laundry lists.</p>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="pending" class="flex justify-center py-20">
        <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
      </div>

      <!-- Order Lists -->
      <div v-else-if="orders && orders.data && orders.data.length > 0" class="space-y-4">
        <div 
          v-for="order in orders.data" 
          :key="order.id"
          @click="navigateTo(`/customer/orders/${order.id}`)"
          class="bg-surface-container rounded-3xl p-5 md:p-6 group hover:ring-1 hover:ring-primary/50 transition-all cursor-pointer border border-transparent flex flex-col sm:flex-row sm:items-center gap-6"
        >
          <div class="flex-1">
             <div class="flex items-center gap-3 mb-2">
               <span class="text-sm font-mono font-bold text-surface-onSurfaceVariant">#{{ order.id.slice(0,8).toUpperCase() }}</span>
               <span v-if="order.status === 'pending'" class="bg-surface-containerHigh text-surface-onSurfaceVariant text-[10px] font-bold px-2 py-1 rounded-md uppercase tracking-wider">Pending</span>
               <span v-if="order.status === 'process'" class="bg-primary/20 text-primary text-[10px] font-bold px-2 py-1 rounded-md uppercase tracking-wider">In Process</span>
               <span v-if="order.status === 'completed'" class="bg-green-500/20 text-green-500 text-[10px] font-bold px-2 py-1 rounded-md uppercase tracking-wider">Ready Pickup</span>
               <span v-if="order.status === 'picked_up'" class="bg-surface-containerHigh/50 text-surface-onSurfaceVariant text-[10px] font-bold px-2 py-1 rounded-md uppercase tracking-wider opacity-60">Finished</span>
             </div>
             
             <h3 class="text-lg font-bold mb-1">{{ order.Outlet.name }}</h3>
             <p class="text-sm text-surface-onSurfaceVariant">{{ new Date(order.created_at).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) }}</p>
          </div>
          
          <div class="flex items-center justify-between sm:flex-col sm:items-end gap-1">
              <span class="text-xl font-bold text-primary">Rp {{ order.total_price }}</span>
              <span class="text-xs text-surface-onSurfaceVariant flex items-center gap-1 group-hover:text-primary transition-colors">
                View Tracking <span class="material-symbols-outlined text-[16px]">arrow_forward</span>
              </span>
          </div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-else class="text-center py-24 opacity-60">
        <span class="material-symbols-outlined text-6xl mb-4">receipt_long</span>
        <p class="text-lg">You have no active orders.</p>
        <button @click="navigateTo('/customer')" class="mt-6 text-primary font-bold hover:underline">Find a laundry nearby</button>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'customer'
})
useHead({
  title: 'LaundryIn - My Orders'
})

interface Order {
  id: string;
  status: string;
  total_price: number;
  created_at: string;
  Outlet: {
    name: string;
  }
}

// Assuming the API returns a 'data' array
const { data: orders, pending } = await useFetch<{data: Order[]}>('/api/orders/my')
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
