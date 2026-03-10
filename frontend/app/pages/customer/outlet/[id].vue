<template>
  <div class="h-full w-full flex flex-col md:flex-row overflow-hidden bg-[#131b1f] text-surface-onSurface">
    
    <!-- Left Column: Outlet Info & Order Details -->
    <div class="flex-1 overflow-y-auto p-6 md:p-10 border-r border-outline/10 custom-scrollbar">
      <div class="max-w-md">
        <button @click="navigateTo('/customer')" class="flex items-center gap-2 text-surface-onSurfaceVariant hover:text-primary transition-colors mb-6">
          <span class="material-symbols-outlined">arrow_back</span> Back to Catalog
        </button>

        <!-- Loading State for Outlet -->
        <div v-if="pendingOutlet" class="py-10">
          <span class="material-symbols-outlined animate-spin text-4xl text-primary block">progress_activity</span>
        </div>

        <div v-else-if="outlet">
          <h1 class="text-3xl font-bold mb-2">{{ outlet.name }}</h1>
          <p class="text-surface-onSurfaceVariant mb-10">{{ outlet.address }}</p>
          
          <!-- Contact Info Form -->
          <section class="mb-10">
            <h2 class="text-sm font-bold text-primary tracking-wider mb-4 uppercase">Contact Info</h2>
            <div class="space-y-4">
              <div class="bg-surface-container rounded-2xl px-4 py-3 border border-transparent focus-within:border-primary/50 transition-colors">
                <input type="text" placeholder="Full Name" class="w-full bg-transparent border-none outline-none text-surface-onSurface placeholder-surface-onSurfaceVariant" />
              </div>
              <div class="bg-surface-container rounded-2xl px-4 py-3 border border-transparent focus-within:border-primary/50 transition-colors">
                <input type="text" placeholder="WhatsApp Number" class="w-full bg-transparent border-none outline-none text-surface-onSurface placeholder-surface-onSurfaceVariant" />
              </div>
              <div class="bg-surface-container rounded-2xl px-4 py-3 border border-transparent focus-within:border-primary/50 transition-colors">
                <textarea placeholder="Pickup Address" rows="2" class="w-full bg-transparent border-none outline-none text-surface-onSurface placeholder-surface-onSurfaceVariant resize-none mt-1"></textarea>
              </div>
            </div>
          </section>

          <!-- Notes -->
          <section>
            <div class="bg-surface-container rounded-2xl px-4 py-3 border border-transparent focus-within:border-primary/50 transition-colors">
              <input type="text" placeholder="Special Instructions (Optional)" class="w-full bg-transparent border-none outline-none text-surface-onSurface placeholder-surface-onSurfaceVariant text-sm" />
            </div>
          </section>
        </div>
      </div>
    </div>

    <!-- Right Column: Select Services -->
    <div class="w-full md:w-[45%] flex flex-col bg-[#161f24] relative">
      <div class="p-6 md:p-8 flex-1 overflow-y-auto custom-scrollbar pb-32">
        <div class="flex items-center justify-between mb-8">
          <h2 class="text-2xl font-bold">Select Services</h2>
        </div>

        <div v-if="pendingServices" class="py-10 flex justify-center">
          <span class="material-symbols-outlined animate-spin text-4xl text-primary block">progress_activity</span>
        </div>

        <div v-else-if="services?.data" class="space-y-4">
          <div 
            v-for="service in services.data" 
            :key="service.id" 
            class="bg-surface-container rounded-3xl p-3 flex gap-4 items-center"
            :class="{ 'border border-primary/50': getQty(service.id) > 0 }"
          >
            <div class="w-16 h-16 bg-surface-containerHigh rounded-2xl flex items-center justify-center shrink-0">
               <span class="material-symbols-outlined text-[32px] text-surface-onSurfaceVariant">local_laundry_service</span>
            </div>
            <div class="flex-1 min-w-0 pr-2">
              <h3 class="font-bold text-[15px] truncate">{{ service.name }}</h3>
              <p class="text-[11px] text-surface-onSurfaceVariant mb-1 line-clamp-2">{{ service.description }}</p>
              <div class="text-sm mt-1">
                <span class="font-bold text-primary">Rp {{ service.price }}</span> 
                <span class="text-xs text-surface-onSurfaceVariant">/ {{ service.unit }}</span>
              </div>
            </div>
            
            <div class="flex flex-col gap-1 w-20">
              <input 
                 type="number" 
                 :value="getQty(service.id)"
                 @input="updateQuantity(service, $event.target.value)"
                 class="w-full bg-surface-containerHigh text-center text-sm font-bold border-none outline-none py-1.5 rounded-lg text-surface-onSurface placeholder-surface-onSurfaceVariant" 
                 pattern="[0-9]*" 
                 inputmode="decimal"
                 min="0"
                 step="any"
              />
            </div>
          </div>
          
          <div v-if="services.data.length === 0" class="text-center opacity-60 py-10">
            <p>No services available here currently.</p>
          </div>
        </div>
      </div>

      <!-- Sticky Footer / Cart -->
      <div class="absolute bottom-0 left-0 right-0 py-5 px-6 bg-[#161f24] border-t border-outline/10">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-[10px] text-surface-onSurfaceVariant mb-0.5">Estimated Total <span class="bg-surface-containerHigh px-1.5 py-0.5 rounded text-white ml-2">{{ cartStore.itemCount }} services</span></p>
            <div class="flex items-baseline gap-2">
              <span class="text-2xl font-bold">Rp {{ cartStore.totalPreview }}</span>
              <span class="text-[10px] text-primary font-medium">Free Delivery</span>
            </div>
          </div>
          <button 
            @click="checkout" 
            :disabled="cartStore.itemCount === 0"
            class="bg-primary text-primary-text px-6 py-3.5 rounded-full font-bold flex items-center gap-2 hover:brightness-110 active:scale-95 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Checkout
            <span class="material-symbols-outlined text-[20px]">arrow_forward</span>
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useCartStore } from '~/stores/cart'

definePageMeta({
  layout: 'customer'
})
useHead({
  title: 'LaundryIn - Outlet Services'
})

interface Outlet {
  id: string
  name: string
  address: string
  description: string
}

interface Service {
  id: string
  name: string
  description: string
  price: number
  unit: string
}

const route = useRoute()
const cartStore = useCartStore()
const outletId = route.params.id as string

const { data: outletWrapper, pending: pendingOutlet } = await useFetch<{ data: Outlet }>(`/api/outlets/${outletId}`)
const outlet = computed(() => outletWrapper.value?.data)

const { data: services, pending: pendingServices } = await useFetch<{ data: Service[] }>(`/api/outlets/${outletId}/services`)

// Initialize this outlet to cart (clears other outlet items if there is any mismatch in cart state)
if (outletId) {
  cartStore.setOutlet(outletId)
}

const getQty = (serviceId: string) => {
  const item = cartStore.items.find(i => i.serviceId === serviceId)
  return item ? parseFloat(item.qty) : 0
}

const updateQuantity = (service: Service, value: string) => {
  if (!value || value === '') {
    cartStore.removeItem(service.id)
    return
  }
  
  cartStore.updateQty(service.id, value)
  
  // If it didn't exist and they typed something, we need to add it manually instead of just updating qty
  const existing = cartStore.items.find(i => i.serviceId === service.id)
  if (!existing && parseFloat(value) > 0) {
    cartStore.addItem({
      serviceId: service.id,
      name: service.name,
      price: service.price,
      unit: service.unit,
      qty: value
    })
  }
}

const checkout = () => {
    // Navigate to actual checkout/auth flow later
    console.log('Sending this format back to the API zero-trust:')
    const payload = {
        outlet_id: outletId,
        items: cartStore.items.map(i => ({
            service_id: i.serviceId,
            qty: i.qty // sent as string
        }))
    }
    console.log(JSON.stringify(payload, null, 2))
    alert('Mock Checkout! Payload output to console.')
}
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
