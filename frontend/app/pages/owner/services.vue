<template>
  <div class="mx-auto max-w-[1200px] w-full px-4 md:px-8 py-6 space-y-5">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold mb-0.5">Services Catalog</h2>
        <p class="text-surface-onSurfaceVariant text-sm">Kelola produk & layanan laundry.</p>
      </div>
      <button @click="handleNewService" class="btn-primary rounded-xl px-4 h-10 text-sm gap-1.5 shrink-0">
        <span class="material-symbols-outlined text-[18px]">add</span>
        <span class="hidden sm:inline">New Service</span>
        <span class="sm:hidden">New</span>
      </button>
    </div>

    <!-- Filter Outlet -->
    <div v-if="outlets.length > 0" class="card !p-4">
      <label class="text-xs font-semibold text-surface-onSurfaceVariant block mb-2">Filter Outlet</label>
      <select
        v-model="selectedOutletId"
        class="bg-surface-containerHigh w-full rounded-xl py-2.5 px-3 text-sm outline-none focus:ring-1 focus:ring-primary/30 border border-border transition-all"
      >
        <option value="">Semua Outlet</option>
        <option v-for="o in outlets" :key="o.id" :value="o.id">
          {{ o.name }}
        </option>
      </select>
    </div>

    <div v-if="pending" class="flex justify-center py-10">
      <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
    </div>

    <!-- Service List -->
    <div v-else-if="services && services.length > 0" class="card !p-0 overflow-hidden divide-y divide-border">
      <div
        class="p-4 flex gap-3.5 items-center hover:bg-surface-containerHigh/40 transition-colors group"
        v-for="service in services"
        :key="service.id"
      >
        <div class="w-10 h-10 bg-primary/10 text-primary rounded-xl flex items-center justify-center shrink-0">
          <span class="material-symbols-outlined text-[20px]">styler</span>
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-sm truncate">{{ service.name }}</h3>
          </div>
          <div class="text-[11px] text-surface-onSurfaceVariant mt-0.5 flex items-center gap-1.5 truncate">
            <span class="material-symbols-outlined text-[12px]">store</span>
            {{ service.outlet_name }} &middot; Unit: <span class="font-semibold">{{ service.unit }}</span>
          </div>
          <div class="font-mono text-primary font-bold text-sm mt-1">Rp {{ Number(service.price).toLocaleString('id-ID') }}</div>
        </div>
        <div class="flex gap-1.5 shrink-0">
          <button @click="openEditModal(service)" class="h-8 w-8 rounded-lg bg-surface-containerHigh flex items-center justify-center text-surface-onSurface hover:text-primary transition-all">
            <span class="material-symbols-outlined text-[18px]">edit_square</span>
          </button>
          <button @click="openDeleteModal(service)" class="h-8 w-8 rounded-lg bg-danger-muted flex items-center justify-center text-danger hover:bg-danger/20 transition-all">
            <span class="material-symbols-outlined text-[18px]">delete</span>
          </button>
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col items-center text-center py-10 text-surface-onSurfaceVariant text-sm animate-fade-in">
      <div class="w-14 h-14 rounded-2xl bg-surface-container flex items-center justify-center mb-3 border border-border">
        <span class="material-symbols-outlined text-2xl text-surface-onSurfaceVariant">dry_cleaning</span>
      </div>
      Belum ada service yang terdaftar.
    </div>

    <!-- Modal Form (Tambah/Edit) -->
    <div v-if="showModal" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 animate-fade-in" @click.self="closeModal">
      <div class="bg-surface w-full max-w-lg rounded-2xl shadow-xl overflow-hidden flex flex-col">
        <div class="px-6 py-4 border-b border-border flex justify-between items-center bg-surface-containerHigh">
          <h3 class="font-bold text-lg">{{ isEditing ? 'Edit Layanan' : 'Tambah Layanan' }}</h3>
          <button @click="closeModal" class="text-surface-onSurfaceVariant hover:text-danger transition-colors">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="p-6 space-y-4">
          <!-- Outlet Selection (Only on Create) -->
          <div>
            <label class="block text-xs font-semibold mb-1 text-surface-onSurfaceVariant">Outlet <span class="text-danger">*</span></label>
            <div v-if="isEditing" class="w-full bg-surface-container p-2.5 rounded-xl text-sm border border-border opacity-70">
              {{ outlets.find(o => o.id === serviceForm.outlet_id)?.name || 'Outlet Selected' }}
            </div>
            <select 
              v-else
              v-model="serviceForm.outlet_id" 
              class="w-full bg-surface-containerHigh text-sm rounded-xl py-2.5 px-4 outline-none focus:ring-1 focus:ring-primary/50 border border-transparent focus:border-border"
            >
              <option value="">Pilih Outlet</option>
              <option v-for="o in outlets" :key="o.id" :value="o.id">{{ o.name }}</option>
            </select>
            <span v-if="formErrors.outlet_id" class="text-danger text-[10px] block mt-1">{{ formErrors.outlet_id }}</span>
          </div>

          <div>
            <label class="block text-xs font-semibold mb-1 text-surface-onSurfaceVariant">Nama Layanan <span class="text-danger">*</span></label>
            <input v-model="serviceForm.name" type="text" class="w-full bg-surface-containerHigh text-sm rounded-xl py-2.5 px-4 outline-none focus:ring-1 focus:ring-primary/50 border border-transparent focus:border-border" placeholder="Contoh: Cuci Kering Regular">
            <span v-if="formErrors.name" class="text-danger text-[10px] block mt-1">{{ formErrors.name }}</span>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold mb-1 text-surface-onSurfaceVariant">Harga (Rp) <span class="text-danger">*</span></label>
              <input v-model.number="serviceForm.price" type="number" class="w-full bg-surface-containerHigh text-sm rounded-xl py-2.5 px-4 outline-none focus:ring-1 focus:ring-primary/50 border border-transparent focus:border-border" placeholder="0">
              <span v-if="formErrors.price" class="text-danger text-[10px] block mt-1">{{ formErrors.price }}</span>
            </div>
            <div>
              <label class="block text-xs font-semibold mb-1 text-surface-onSurfaceVariant">Unit <span class="text-danger">*</span></label>
              <select v-model="serviceForm.unit" class="w-full bg-surface-containerHigh text-sm rounded-xl py-2.5 px-4 outline-none focus:ring-1 focus:ring-primary/50 border border-transparent focus:border-border">
                <option value="KG">Kilogram (KG)</option>
                <option value="PCS">Satuan (PCS)</option>
              </select>
            </div>
          </div>
        </div>
        <div class="px-6 py-4 border-t border-border flex justify-end gap-3 bg-surface-container/50">
          <button @click="closeModal" class="px-5 py-2.5 rounded-xl text-sm font-semibold text-surface-onSurfaceVariant hover:bg-surface-containerHigh transition-colors" :disabled="formLoading">
            Batal
          </button>
          <button @click="submitForm" class="btn-primary px-6 py-2.5 rounded-xl text-sm font-semibold flex items-center gap-2 disabled:opacity-50" :disabled="formLoading">
            <span v-if="formLoading" class="material-symbols-outlined text-[18px] animate-spin">progress_activity</span>
            Simpan
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Delete -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 animate-fade-in" @click.self="showDeleteModal = false">
      <div class="bg-surface w-full max-w-sm rounded-2xl shadow-xl overflow-hidden flex flex-col p-6 text-center">
        <div class="w-16 h-16 rounded-full bg-danger/10 text-danger flex items-center justify-center mx-auto mb-4 border border-danger/20">
          <span class="material-symbols-outlined text-3xl">delete</span>
        </div>
        <h3 class="font-bold text-xl mb-2">Hapus Layanan</h3>
        <p class="text-sm text-surface-onSurfaceVariant mb-6">
          Yakin ingin menghapus layanan <b>{{ serviceToDelete?.name }}</b>? Aksi ini tidak dapat dibatalkan.
        </p>
        <div class="flex flex-col gap-2">
          <button @click="confirmDelete" class="bg-danger text-white w-full py-3 rounded-xl text-sm font-semibold flex items-center justify-center gap-2 disabled:opacity-50 hover:bg-danger/90 transition-colors" :disabled="deleteLoading">
            <span v-if="deleteLoading" class="material-symbols-outlined text-[18px] animate-spin">progress_activity</span>
            Ya, Hapus
          </button>
          <button @click="showDeleteModal = false" class="w-full py-3 rounded-xl text-sm font-semibold text-surface-onSurface hover:bg-surface-container transition-colors" :disabled="deleteLoading">
            Batal
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useToast } from '~/composables/useToast'
import { ref, computed } from 'vue'
import type { ApiResponse, PaginatedResponse } from '~/types/api'

definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn — Services'
})

const authStore = useAuthStore()
const router = useRouter()
const { success: toastSuccess, error: toastError } = useToast()

interface Service {
  id: string
  outlet_id: string
  name: string
  price: string
  unit: string
  outlet_name?: string
}

interface Outlet {
  id: string
  name: string
}

const { data: outletsWrapper } = await useFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/outlets', {
  headers: { Authorization: authStore.authHeader },
  server: false
})

const outlets = computed(() => outletsWrapper.value?.data?.data ?? [])
const selectedOutletId = ref('')

const { data: services, pending, refresh: refreshServices } = await useAsyncData<Service[]>(
  'owner-services',
  async () => {
    if (outlets.value.length === 0) return []

    // Specific outlet query
    if (selectedOutletId.value) {
      const res = await $fetch<ApiResponse<Service[]>>(`/api/outlets/${selectedOutletId.value}/services`, {
        headers: { Authorization: authStore.authHeader }
      })
      const outletName = outlets.value.find(o => o.id === selectedOutletId.value)?.name || 'Unknown'
      return (res.data || []).map(s => ({ ...s, outlet_name: outletName }))
    }

    // All outlets parallel query
    const promises = outlets.value.map(async (o) => {
      try {
        const res = await $fetch<ApiResponse<Service[]>>(`/api/outlets/${o.id}/services`, {
          headers: { Authorization: authStore.authHeader }
        })
        return (res.data || []).map(s => ({ ...s, outlet_name: o.name }))
      } catch {
        return []
      }
    })

    const results = await Promise.all(promises)
    return results.flat()
  },
  { watch: [selectedOutletId, outletsWrapper], server: false }
)

// Modal States
const showModal = ref(false)
const isEditing = ref(false)
const formLoading = ref(false)

const serviceForm = ref({
  id: '',
  outlet_id: '',
  name: '',
  price: 0,
  unit: 'KG'
})

const formErrors = ref({
  outlet_id: '',
  name: '',
  price: '',
  unit: ''
})

// Delete Modal States
const showDeleteModal = ref(false)
const serviceToDelete = ref<Service | null>(null)
const deleteLoading = ref(false)

const closeModal = () => {
  showModal.value = false
  serviceForm.value = { id: '', outlet_id: '', name: '', price: 0, unit: 'KG' }
  formErrors.value = { outlet_id: '', name: '', price: '', unit: '' }
}

const handleNewService = () => {
  isEditing.value = false
  serviceForm.value = { 
    id: '', 
    outlet_id: selectedOutletId.value, 
    name: '', 
    price: 0, 
    unit: 'KG' 
  }
  formErrors.value = { outlet_id: '', name: '', price: '', unit: '' }
  showModal.value = true
}

const openEditModal = (s: Service) => {
  isEditing.value = true
  serviceForm.value = { 
    id: s.id, 
    outlet_id: s.outlet_id, 
    name: s.name, 
    price: Number(s.price), 
    unit: s.unit 
  }
  formErrors.value = { outlet_id: '', name: '', price: '', unit: '' }
  showModal.value = true
}

const validateForm = () => {
  let valid = true
  formErrors.value = { outlet_id: '', name: '', price: '', unit: '' }
  
  if (!serviceForm.value.outlet_id) {
    formErrors.value.outlet_id = 'Pilih outlet terlebih dahulu'
    valid = false
  }
  if (!serviceForm.value.name.trim()) {
    formErrors.value.name = 'Nama layanan wajib diisi'
    valid = false
  }
  if (!serviceForm.value.price || serviceForm.value.price <= 0) {
    formErrors.value.price = 'Harga harus lebih dari 0'
    valid = false
  }
  
  return valid
}

const submitForm = async () => {
  if (!validateForm()) return
  
  formLoading.value = true
  try {
    const payload = {
      outlet_id: serviceForm.value.outlet_id,
      name: serviceForm.value.name.trim(),
      price: String(serviceForm.value.price),
      unit: serviceForm.value.unit
    }

    if (isEditing.value) {
      // Body for PUT usually doesn't need outlet_id if it's not allowed to change
      const { outlet_id, ...editPayload } = payload
      await $fetch(`/api/services/${serviceForm.value.id}`, {
        method: 'PUT',
        headers: { Authorization: authStore.authHeader },
        body: editPayload
      })
      toastSuccess('Service berhasil diperbarui')
    } else {
      await $fetch('/api/services', {
        method: 'POST',
        headers: { Authorization: authStore.authHeader },
        body: payload
      })
      toastSuccess('Service berhasil ditambahkan')
    }
    
    closeModal()
    await refreshServices()
  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else {
      toastError(err?.data?.message || 'Gagal menyimpan service')
    }
  } finally {
    formLoading.value = false
  }
}

const openDeleteModal = (s: Service) => {
  serviceToDelete.value = s
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!serviceToDelete.value) return
  
  deleteLoading.value = true
  try {
    await $fetch(`/api/services/${serviceToDelete.value.id}`, {
      method: 'DELETE',
      headers: { Authorization: authStore.authHeader }
    })
    toastSuccess('Service berhasil dihapus')
    showDeleteModal.value = false
    await refreshServices()
  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else {
      toastError('Gagal menghapus service, coba lagi')
    }
  } finally {
    deleteLoading.value = false
    serviceToDelete.value = null
  }
}

onMounted(() => {
  if (import.meta.client) {
    window.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') {
        if (showModal.value) closeModal()
        if (showDeleteModal.value) showDeleteModal.value = false
      }
    })
  }
})

onActivated(() => {
  refreshServices()
})
</script>
