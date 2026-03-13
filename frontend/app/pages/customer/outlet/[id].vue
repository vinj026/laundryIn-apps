<template>
  <div class="h-full w-full flex-1 overflow-y-auto md:overflow-hidden bg-surface relative p-0 sm:p-4 lg:p-8">
    <!-- Outer container: scrolls on mobile, fixed hidden on desktop -->

    <!-- Card: Auto height on mobile, full height on desktop -->
    <div class="booking-card w-full max-w-6xl mx-auto flex flex-col md:grid md:grid-cols-[2fr_1px_3fr] lg:grid-cols-[1fr_1px_1fr] bg-surface-raised border-x-0 sm:border md:border-border sm:rounded-2xl md:rounded-[24px] overflow-hidden h-auto md:h-full sm:shadow-2xl">

      <!-- LEFT COLUMN: form-panel -->
      <div class="form-panel flex flex-col relative h-auto md:h-full md:overflow-hidden">
        <!-- Header / Back -->
        <div class="sticky top-0 z-10 bg-surface-raised/95 backdrop-blur px-5 py-3 border-b border-border/50">
          <button @click="navigateTo('/customer')" class="flex items-center gap-2 text-surface-onSurfaceVariant hover:text-primary transition-colors text-sm w-fit mb-4">
            <span class="material-symbols-outlined text-[18px]">arrow_back</span>
            Back
          </button>
          <h1 class="text-2xl font-bold mb-1">Order Details</h1>
          <p class="text-surface-onSurfaceVariant text-[13px]">Fill in your information to schedule a pickup.</p>
        </div>

        <div class="p-3 lg:p-4 flex flex-col gap-3 flex-1 min-h-0">
          <!-- Loading -->
          <div v-if="pendingOutlet" class="py-10 flex justify-center">
            <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
          </div>

          <form v-else-if="outlet" class="flex flex-col gap-6 h-full">
            <!-- Outlet Info -->
            <div class="space-y-2 shrink-0">
              <h2 class="text-[10px] font-bold text-surface-onSurfaceVariant uppercase tracking-wider">Selected Branch</h2>
              <div class="flex items-center justify-between p-3 rounded-xl border border-border bg-surface-container">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-surface-containerHigh flex items-center justify-center shrink-0">
                    <span class="material-symbols-outlined text-primary">storefront</span>
                  </div>
                  <div>
                    <h3 class="font-bold text-sm text-surface-onSurface">{{ outlet.name }}</h3>
                    <p class="text-xs text-surface-onSurfaceVariant mt-0.5 line-clamp-1">{{ outlet.address }}</p>
                  </div>
                </div>
                <button @click="showOutletModal = true" type="button" class="text-[13px] font-bold text-primary hover:text-primary-light transition-colors px-2 py-1">Change</button>
              </div>
            </div>

            <!-- Contact Info -->
            <div class="flex flex-col gap-4 shrink-0">
              <div class="space-y-2">
              <h2 class="text-[10px] font-bold text-surface-onSurfaceVariant uppercase tracking-wider">Contact Info</h2>
              <div class="grid grid-cols-2 gap-4">
                <div class="m3-field min-w-0">
                  <input id="name" name="name" placeholder=" " required type="text" v-model="form.name" class="!py-3.5">
                  <label for="name" class="whitespace-nowrap">Name</label>
                </div>
                <div class="m3-field min-w-0">
                  <input id="phone" name="phone" placeholder=" " required type="tel" v-model="form.phone" class="!py-3.5">
                  <label for="phone" class="whitespace-nowrap overflow-hidden text-ellipsis max-w-full">Phone</label>
                </div>
              </div>
              </div>
              <div class="m3-field">
                <textarea class="resize-none !py-3.5" id="address" name="address" placeholder=" " rows="2" v-model="form.address"></textarea>
                <label for="address">Pickup Address</label>
              </div>
            </div>

            <!-- Pickup Date -->
          <div class="space-y-2 w-full shrink-0 relative">
            <h2 class="text-[10px] font-bold text-surface-onSurfaceVariant uppercase tracking-wider">Pickup Date</h2>
            <ClientOnly>
              <DateCarousel v-model="form.date" />
              <template #fallback>
                <div class="h-[96px] flex items-center justify-center text-surface-onSurfaceVariant/50 text-sm">
                  Loading dates...
                </div>
              </template>
            </ClientOnly>
          </div>

            <!-- Pickup Time & Notes -->
            <div class="flex flex-col gap-4 shrink-0">
              <div class="space-y-2">
                <h2 class="text-[10px] font-bold text-surface-onSurfaceVariant uppercase tracking-wider">Pickup Time</h2>
                <div class="grid grid-cols-2 gap-4">
                  <div class="w-full">
                    <input class="chip-radio hidden" id="time-morning" name="time" type="radio" value="morning" v-model="form.time">
                    <label class="flex items-center gap-2 justify-center w-full px-2 py-3.5 rounded-xl border border-border bg-surface-container cursor-pointer hover:bg-surface-containerHigh transition-all text-surface-onSurfaceVariant text-[13px]" for="time-morning">
                      <span class="material-symbols-outlined text-[18px] icon">wb_sunny</span>
                      <span class="font-medium text-surface-onSurface">09:00 - 12:00</span>
                    </label>
                  </div>
                  <div class="w-full">
                    <input class="chip-radio hidden" id="time-afternoon" name="time" type="radio" value="afternoon" v-model="form.time">
                    <label class="flex items-center gap-2 justify-center w-full px-2 py-3.5 rounded-xl border border-border bg-surface-container cursor-pointer hover:bg-surface-containerHigh transition-all text-surface-onSurfaceVariant text-[13px]" for="time-afternoon">
                      <span class="material-symbols-outlined text-[18px] icon">partly_cloudy_day</span>
                      <span class="font-medium text-surface-onSurface">13:00 - 16:00</span>
                    </label>
                  </div>
                </div>
              </div>

              <div class="m3-field shrink-0">
                <textarea class="resize-none !py-3.5" id="notes" name="notes" placeholder=" " rows="2" v-model="form.notes"></textarea>
                <label for="notes">Special Instructions (Optional)</label>
              </div>
            </div>
          </form>
        </div>
      </div>

      <!-- divider-vertical -->
      <div class="hidden md:block w-px bg-border h-full"></div>

      <!-- border horizontal for mobile -->
      <div class="block md:hidden h-2 bg-surface-base w-full"></div>

      <!-- RIGHT COLUMN: services-panel -->
      <div class="services-panel flex flex-col h-[700px] md:h-full bg-surface-container relative overflow-hidden">
        
        <!-- services-header -->
        <header class="flex-shrink-0 px-6 py-5 border-b border-border/50 bg-surface-containerHigh/50 backdrop-blur z-10 w-full">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-[17px] font-bold text-surface-onSurface flex items-center gap-2">
              Select Services
              <span class="bg-surface-overlay border border-border text-[11px] font-semibold px-2 py-0.5 rounded text-surface-onSurfaceVariant">{{ filteredServices?.length || 0 }} Available</span>
            </h2>
          </div>

          <div class="flex flex-col sm:flex-row items-center gap-3">
            <!-- Search -->
            <div class="relative w-full sm:flex-1">
              <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-surface-onSurfaceVariant text-[18px]">search</span>
              <input 
                type="text" 
                v-model="searchQuery" 
                placeholder="Search services..." 
                class="w-full bg-surface-raised border border-border rounded-lg pl-9 pr-4 py-2.5 text-[13px] text-surface-onSurface focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20 transition-all"
              >
            </div>
            <!-- Category Filter -->
            <div class="relative w-full sm:w-auto shrink-0">
              <select v-model="categoryFilter" class="w-full sm:w-[130px] appearance-none bg-surface-raised border border-border rounded-lg pl-4 pr-8 py-2.5 text-[13px] text-surface-onSurface font-medium focus:outline-none focus:border-primary transition-all">
                <option value="all">All Categories</option>
                <option value="clothes">Clothes</option>
                <option value="shoes">Shoes</option>
                <option value="others">Others</option>
              </select>
              <span class="material-symbols-outlined absolute right-2.5 top-1/2 -translate-y-1/2 text-surface-onSurfaceVariant pointer-events-none text-[18px]">expand_more</span>
            </div>
          </div>
        </header>

        <!-- services-list -->
        <div class="services-list flex-1 overflow-y-auto custom-scrollbar p-6 bg-surface-raised/20">
          <div v-if="pendingServices" class="py-10 flex justify-center">
            <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
          </div>

          <div v-else-if="filteredServices.length > 0" class="flex flex-col gap-3">
            <div
              v-for="service in filteredServices"
              :key="service.id"
              class="service-item flex items-center gap-4 p-4 rounded-2xl border transition-all cursor-pointer"
              :class="getQty(service.id) > 0 ? 'border-primary bg-primary/5' : 'border-border bg-surface-container hover:border-border-hover hover:bg-surface-containerHigh'"
              @click="getQty(service.id) === 0 ? updateQuantity(service, 1) : null"
            >
              <!-- icon box -->
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl shrink-0 transition-colors"
                   :class="getQty(service.id) > 0 ? 'bg-primary/20' : 'bg-surface-raised'">
                {{ getServiceIcon(service) }}
              </div>
              
              <!-- nama + harga -->
              <div class="flex flex-col flex-1 min-w-0">
                <h3 class="text-[14px] font-bold leading-tight truncate" :class="getQty(service.id) > 0 ? 'text-primary' : 'text-surface-onSurface'">{{ service.name }}</h3>
                <p v-if="service.description" class="text-[11px] text-surface-onSurfaceVariant leading-snug line-clamp-1 mt-0.5">{{ service.description }}</p>
                <p v-else class="text-[11px] text-surface-onSurfaceVariant leading-snug line-clamp-1 mt-0.5">Pencucian dan perawatan standar.</p>
                <span class="font-mono font-bold text-[13px] mt-1" :class="getQty(service.id) > 0 ? 'text-primary' : 'text-surface-onSurface'">
                  Rp {{ Number(service.price).toLocaleString('id-ID') }} <span class="text-[10px] font-sans font-normal opacity-70">/ {{ service.unit }}</span>
                </span>
              </div>
              
              <!-- State Controls -->
              <div class="shrink-0 w-[96px] flex justify-end">
                <!-- State 2: Selected -->
                <div class="qty-control flex items-center justify-between w-full h-[36px] bg-surface-raised rounded-lg p-1 border border-border border-primary/20" v-if="getQty(service.id) > 0">
                  <button
                    @click.stop="updateQuantity(service, getQty(service.id) - 1)"
                    class="w-7 h-7 flex items-center justify-center rounded bg-surface-container hover:bg-surface-overlay text-surface-onSurface transition-colors"
                  >
                    <span class="material-symbols-outlined text-[15px]">remove</span>
                  </button>
                  <span class="text-[13px] font-bold w-5 text-center font-mono">{{ getQty(service.id) }}</span>
                  <button
                    @click.stop="updateQuantity(service, getQty(service.id) + 1)"
                    class="w-7 h-7 flex items-center justify-center rounded bg-primary text-primary-text hover:brightness-110 transition-colors"
                  >
                    <span class="material-symbols-outlined text-[15px]">add</span>
                  </button>
                </div>
                <!-- State 1: Not Selected -->
                <div v-else class="w-full">
                  <button class="w-full h-[34px] flex items-center justify-center gap-1 rounded-lg border border-border bg-surface-raised text-[12px] font-bold text-surface-onSurface hover:text-primary hover:border-primary/50 transition-colors">
                    <span class="material-symbols-outlined text-[14px]">add</span>
                    Add
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <div v-else class="py-12 flex flex-col items-center justify-center text-surface-onSurfaceVariant">
            <span class="material-symbols-outlined text-4xl mb-2 opacity-50">search_off</span>
            <p class="text-[13px]">No services found matching your criteria.</p>
          </div>
        </div>

        <!-- services-footer -->
        <footer class="flex-shrink-0 p-5 lg:p-6 border-t border-border/50 bg-surface-raised z-10 w-full relative before:absolute before:-top-6 before:left-0 before:right-0 before:h-6 before:bg-gradient-to-t before:from-surface-raised before:to-transparent before:pointer-events-none">
          <div class="flex flex-col gap-2">
            <div class="flex items-center justify-between">
              <div class="flex flex-col">
                <span class="text-xs text-surface-onSurfaceVariant font-medium mb-0.5">
                  Estimated Total 
                  <span v-if="cartStore.itemCount > 0" class="ml-1 inline-flex items-center justify-center px-1.5 py-0.5 rounded bg-surface-containerHigh text-[10px] text-surface-onSurface">{{ cartStore.itemCount }} items</span>
                </span>
                <span class="text-[22px] font-extrabold font-mono leading-none tracking-tight">Rp {{ cartStore.totalPreview.toLocaleString('id-ID') }}</span>
              </div>
              <button
                @click="checkout"
                :disabled="checkoutLoading || cartStore.itemCount === 0"
                class="btn-primary py-3.5 px-6 rounded-xl text-sm min-w-[160px]"
              >
                <span v-if="checkoutLoading">Memproses...</span>
                <span v-else class="flex items-center gap-1">Confirm Booking <span class="material-symbols-outlined text-[18px]">arrow_forward</span></span>
              </button>
            </div>
          </div>
        </footer>

      </div>

    </div>

    <!-- Change Outlet Modal -->
    <div v-if="showOutletModal" class="fixed inset-0 z-[100] flex items-end md:items-center justify-center animate-fade-in p-0 md:p-4">
      <div class="absolute inset-0 bg-surface-overlay/80 backdrop-blur-sm" @click="showOutletModal = false"></div>
      
      <div class="relative w-full md:w-[480px] bg-surface-raised rounded-t-[24px] md:rounded-[24px] max-h-[80vh] flex flex-col shadow-2xl border border-border animate-slide-up md:animate-scale-in">
        <div class="flex items-center justify-between p-5 border-b border-border/50 shrink-0">
          <h2 class="text-lg font-bold">Pilih Outlet Lain</h2>
          <button @click="showOutletModal = false" class="w-8 h-8 flex items-center justify-center rounded-full bg-surface-container hover:bg-surface-containerHigh transition-colors">
            <span class="material-symbols-outlined text-[20px]">close</span>
          </button>
        </div>
        
        <div class="p-5 overflow-y-auto custom-scrollbar flex-1 relative flex flex-col gap-3 min-h-[300px]">
          <div v-if="pendingAllOutlets" class="flex justify-center py-12">
            <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
          </div>
          <div v-else-if="allOutlets.length > 0" class="flex flex-col gap-3 pb-safe">
            <div
              v-for="o in allOutlets"
              :key="o.id"
              @click="changeOutlet(o.id)"
              class="flex items-center gap-3 p-4 rounded-xl border transition-all cursor-pointer"
              :class="o.id === outletId ? 'border-primary bg-primary/5' : 'border-border bg-surface-container hover:border-surface-onSurfaceVariant/20'"
            >
              <div class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0" :class="o.id === outletId ? 'bg-primary/20 text-primary' : 'bg-surface-containerHigh text-surface-onSurfaceVariant'">
                <span class="material-symbols-outlined text-[22px]">{{ o.id === outletId ? 'check_circle' : 'storefront' }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-sm text-surface-onSurface truncate">{{ o.name }}</h3>
                <p class="text-xs text-surface-onSurfaceVariant mt-0.5 line-clamp-2 leading-relaxed">{{ o.address }}</p>
              </div>
            </div>
          </div>
          <div v-else class="text-center text-surface-onSurfaceVariant py-10 flex flex-col items-center">
             <span class="material-symbols-outlined text-4xl mb-2 opacity-50">store_off</span>
             <p class="text-sm">Tidak ada outlet lain tersedia.</p>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useCartStore } from '~/stores/cart'
import { useAuthStore } from '~/stores/auth'
import { ref, computed, onActivated } from 'vue'
import DateCarousel from '~/components/ui/DateCarousel.vue'
import type { ApiResponse, PaginatedResponse } from '~/types/api'

definePageMeta({
  layout: 'customer'
})
useHead({
  title: 'LaundryIn — Order & Booking'
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
  price: string
  unit: string
}

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const authStore = useAuthStore()
const outletId = route.params.id as string

const form = ref({
  name: '',
  phone: '',
  address: '',
  date: new Date().toISOString().substring(0, 10),
  time: 'afternoon',
  notes: ''
})

const showOutletModal = ref(false)

const { data: outletWrapper, pending: pendingOutlet, refresh: refreshOutlet } = await useFetch<{ data: Outlet }>(`/api/public/outlets/${outletId}`)
const outlet = computed(() => outletWrapper.value?.data)

const { data: allOutletsResponse, pending: pendingAllOutlets } = await useFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/public/outlets', {
  lazy: true,
  server: false
})
const allOutlets = computed(() => allOutletsResponse.value?.data?.data || [])

const changeOutlet = (newId: string) => {
  if (newId === outletId) {
    showOutletModal.value = false
    return
  }
  cartStore.clearCart()
  router.push(`/customer/outlet/${newId}`)
}

const { data: services, pending: pendingServices, refresh: refreshServices } = await useFetch<{ data: Service[] }>(`/api/public/outlets/${outletId}/services`)

onActivated(() => {
  if (outletId) {
    refreshOutlet()
    refreshServices()
  }
})

// --- Filtering & Search ---
const searchQuery = ref('')
const categoryFilter = ref('all')

const filteredServices = computed(() => {
  if (!services.value?.data) return []
  let filtered = services.value.data

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    filtered = filtered.filter(s => s.name.toLowerCase().includes(q) || (s.description && s.description.toLowerCase().includes(q)))
  }

  if (categoryFilter.value !== 'all') {
    filtered = filtered.filter(s => {
      const lname = s.name.toLowerCase()
      const lunit = (s.unit || '').toLowerCase().trim()
      
      if (categoryFilter.value === 'clothes') {
        return lunit === 'kg'
      } else if (categoryFilter.value === 'shoes') {
        return lunit === 'pcs' && (lname.includes('sepatu') || lname.includes('tas') || lname.includes('sandal'))
      } else {
        return lunit !== 'kg' && !(lunit === 'pcs' && (lname.includes('sepatu') || lname.includes('tas') || lname.includes('sandal')))
      }
    })
  }

  return filtered
})

const getServiceIcon = (service: Service) => {
  const lname = service.name.toLowerCase()
  const lunit = (service.unit || '').toLowerCase().trim()
  if (lunit === 'kg') return '👕'
  if (lunit === 'pcs' && (lname.includes('sepatu') || lname.includes('tas') || lname.includes('sandal'))) return '👟'
  return '🧹'
}
// ----------------------------

if (outletId) {
  cartStore.setOutlet(outletId)
}

const getQty = (serviceId: string) => {
  const item = cartStore.items.find(i => i.serviceId === serviceId)
  return item ? parseFloat(item.qty) : 0
}

const updateQuantity = (service: Service, value: number) => {
  if (value <= 0) {
    cartStore.removeItem(service.id)
    return
  }

  cartStore.updateQty(service.id, value.toString())

  const existing = cartStore.items.find(i => i.serviceId === service.id)
  if (!existing && value > 0) {
    cartStore.addItem({
      serviceId: service.id,
      name: service.name,
      price: Number(service.price),
      unit: service.unit,
      qty: value.toString()
    })
  }
}

// --- Checkout ---
const checkoutLoading = ref(false)
const { info: toastInfo, success: toastSuccess, error: toastError } = useToast()

const checkout = async () => {
  // Guard 1: must be logged in
  if (!authStore.isLoggedIn) {
    toastInfo('Silakan login terlebih dahulu')
    router.push(`/customer/login?redirect=${route.fullPath}`)
    return
  }

  // Guard 2: cart must not be empty
  if (cartStore.items.length === 0) {
    toastError('Pilih minimal 1 layanan')
    return
  }

  // Guard 3: form validation checks individually
  if (!form.value.name.trim()) {
    toastError('Nama harus diisi')
    return
  }
  if (!form.value.phone.trim()) {
    toastError('Nomor WhatsApp harus diisi')
    return
  }
  if (!form.value.address.trim()) {
    toastError('Alamat pickup harus diisi')
    return
  }

  checkoutLoading.value = true

  try {
    await $fetch('/api/orders', {
      method: 'POST',
      headers: {
        Authorization: authStore.authHeader
      },
      body: {
        outlet_id: outletId,
        customer_name: form.value.name.trim(),
        customer_phone: form.value.phone.trim(),
        customer_address: form.value.address.trim(),
        pickup_date: form.value.date,
        pickup_time: form.value.time,
        notes: form.value.notes.trim(),
        items: cartStore.items.map(i => ({
          service_id: i.serviceId,
          qty: i.qty.toString()
        }))
      }
    })

    // Success: clear cart, toast, and redirect
    cartStore.clearCart()
    toastSuccess('Pesanan berhasil dibuat!')
    router.push('/customer/orders')

  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    const apiMsg = err?.data?.message || ''

    if (status === 400) {
      // Show exact API message for validation failures if provided
      toastError(apiMsg || 'Data tidak valid')
    } else if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push(`/customer/login?redirect=${route.fullPath}`)
    } else if (status === 404) {
      toastError('Layanan tidak tersedia, coba muat ulang halaman')
    } else {
      toastError('Gagal membuat pesanan, coba lagi')
    }
  } finally {
    checkoutLoading.value = false
  }
}
</script>

<style scoped>
/* No specific transitions needed for basic scroll panel */
</style>
