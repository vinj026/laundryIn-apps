<template>
  <div class="px-4 md:px-8 py-6 space-y-5">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold mb-0.5">Outlet Management</h2>
        <p class="text-surface-onSurfaceVariant text-sm">Manage your laundry branch locations.</p>
      </div>
      <button @click="createOutlet" class="btn-primary h-10 w-10 !p-0 rounded-xl shadow-lg">
        <span class="material-symbols-outlined text-[22px]">add</span>
      </button>
    </div>

    <!-- Data Card -->
    <div class="card !p-0 overflow-hidden">
      <!-- Search -->
      <div class="p-4 border-b border-border">
        <div class="relative">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-surface-onSurfaceVariant text-[18px]">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari outlet..."
            class="bg-surface-containerHigh w-full rounded-xl py-2.5 pl-10 pr-4 text-sm outline-none focus:ring-1 focus:ring-primary/30 border border-transparent transition-all"
          >
        </div>
      </div>

      <div v-if="pending" class="flex justify-center py-10">
        <span class="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
      </div>

      <div v-else-if="error" class="flex flex-col items-center justify-center text-center py-16 animate-fade-in">
        <div class="w-14 h-14 rounded-2xl bg-danger-muted text-danger flex items-center justify-center mb-3 border border-danger/30">
          <span class="material-symbols-outlined text-2xl">error</span>
        </div>
        <p class="text-sm font-medium mb-1">Gagal memuat data outlet</p>
        <button @click="refresh()" class="btn-primary py-1.5 px-5 rounded-xl text-xs font-semibold mt-3">
          Coba Lagi
        </button>
      </div>

      <!-- List -->
      <div v-else-if="filteredOutlets.length > 0" class="divide-y divide-border">
        <div
          class="p-4 flex flex-col gap-2 hover:bg-surface-containerHigh/40 transition-colors"
          v-for="outlet in filteredOutlets"
          :key="outlet.id"
        >
          <div class="flex items-start justify-between">
            <div>
              <h3 class="font-bold text-sm text-primary">{{ outlet.name }}</h3>
              <p class="text-xs text-surface-onSurfaceVariant line-clamp-1 mt-0.5">{{ outlet.address }}</p>
            </div>
            <div class="flex gap-1.5">
              <button @click="editOutlet(outlet)" class="h-7 w-7 rounded-lg bg-surface-containerHigh flex items-center justify-center text-surface-onSurface hover:text-primary transition-colors">
                <span class="material-symbols-outlined text-[16px]">edit</span>
              </button>
              <button @click="deleteOutlet(outlet)" class="h-7 w-7 rounded-lg bg-danger-muted flex items-center justify-center text-danger hover:bg-danger/20 transition-colors">
                <span class="material-symbols-outlined text-[16px]">delete</span>
              </button>
            </div>
          </div>
          <div class="flex items-center gap-1 text-[11px] font-mono text-surface-onSurfaceVariant">
            <span class="material-symbols-outlined text-[13px] text-success">check_circle</span>
            <span class="text-success">Active</span>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-16 text-surface-onSurfaceVariant text-sm animate-fade-in">
        <div class="w-14 h-14 mx-auto rounded-2xl bg-surface-container flex items-center justify-center mb-3 border border-border">
          <span class="material-symbols-outlined text-2xl">storefront</span>
        </div>
        <p class="font-medium mb-3">Belum ada outlet</p>
        <button @click="createOutlet" class="btn-primary py-2 px-5 rounded-xl text-sm font-semibold">
          Tambah Outlet Pertama
        </button>
      </div>
    </div>

    <!-- Modal Form (Tambah/Edit) -->
    <div v-if="showModal" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 animate-fade-in" @click.self="closeModal">
      <div class="bg-surface w-full max-w-lg rounded-2xl shadow-xl overflow-hidden flex flex-col">
        <div class="px-6 py-4 border-b border-border flex justify-between items-center bg-surface-containerHigh">
          <h3 class="font-bold text-lg">{{ isEditing ? 'Edit Outlet' : 'Tambah Outlet' }}</h3>
          <button @click="closeModal" class="text-surface-onSurfaceVariant hover:text-danger transition-colors">
            <span class="material-symbols-outlined">close</span>
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-semibold mb-1 text-surface-onSurfaceVariant">Nama Outlet <span class="text-danger">*</span></label>
            <input v-model="outletForm.name" type="text" class="w-full bg-surface-containerHigh text-sm rounded-xl py-2.5 px-4 outline-none focus:ring-1 focus:ring-primary/50 border border-transparent focus:border-border" placeholder="Contoh: Laundry Express Cabang 1">
            <span v-if="formErrors.name" class="text-danger text-[10px] block mt-1">{{ formErrors.name }}</span>
          </div>
          <div>
            <label class="block text-xs font-semibold mb-1 text-surface-onSurfaceVariant">Alamat <span class="text-danger">*</span></label>
            <textarea v-model="outletForm.address" rows="3" class="w-full bg-surface-containerHigh text-sm rounded-xl py-2.5 px-4 outline-none focus:ring-1 focus:ring-primary/50 border border-transparent focus:border-border resize-none" placeholder="Alamat lengkap outlet..."></textarea>
            <span v-if="formErrors.address" class="text-danger text-[10px] block mt-1">{{ formErrors.address }}</span>
          </div>
          <div>
            <label class="block text-xs font-semibold mb-1 text-surface-onSurfaceVariant">Nomor Telepon <span class="text-danger">*</span></label>
            <input v-model="outletForm.phone" type="text" class="w-full bg-surface-containerHigh text-sm rounded-xl py-2.5 px-4 outline-none focus:ring-1 focus:ring-primary/50 border border-transparent focus:border-border" placeholder="Contoh: +628123456789">
            <span v-if="formErrors.phone" class="text-danger text-[10px] block mt-1">{{ formErrors.phone }}</span>
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
        <h3 class="font-bold text-xl mb-2">Hapus Outlet</h3>
        <p class="text-sm text-surface-onSurfaceVariant mb-6">
          Yakin ingin menghapus <b>{{ outletToDelete?.name }}</b>? Aksi ini tidak dapat dibatalkan.
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
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useToast } from '~/composables/useToast'
import type { ApiResponse, PaginatedResponse } from '~/types/api'

definePageMeta({
  layout: 'owner'
})
useHead({
  title: 'LaundryIn — Outlets'
})

const authStore = useAuthStore()
const router = useRouter()
const { success: toastSuccess, error: toastError } = useToast()

watchEffect(() => {
  if (import.meta.client) {
    if (!authStore.isLoggedIn || authStore.user?.role !== 'owner') {
      router.push('/owner/login')
    }
  }
})

interface Outlet {
  id: string
  name: string
  address: string
  phone: string
}

const { data: outletsWrapper, pending, error, refresh } = await useFetch<ApiResponse<PaginatedResponse<Outlet[]>>>('/api/outlets', {
  headers: { Authorization: authStore.authHeader },
  server: false
})
const searchQuery = ref('')

watchEffect(() => {
  if (error.value) {
    const status = error.value?.statusCode || error.value?.status || (error.value?.data as any)?.statusCode
    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else {
      toastError('Gagal memuat data outlet')
    }
  }
})

const outlets = computed(() => outletsWrapper.value?.data?.data ?? [])

const filteredOutlets = computed(() => {
  if (!outlets.value.length) return []
  const q = searchQuery.value.toLowerCase()
  return outlets.value.filter((o: Outlet) =>
    o.name.toLowerCase().includes(q) || o.address.toLowerCase().includes(q)
  )
})

// Modal States
const showModal = ref(false)
const isEditing = ref(false)
const formLoading = ref(false)

const outletForm = ref({
  id: '',
  name: '',
  address: '',
  phone: ''
})

const formErrors = ref({
  name: '',
  address: '',
  phone: ''
})

// Delete Modal States
const showDeleteModal = ref(false)
const outletToDelete = ref<Outlet | null>(null)
const deleteLoading = ref(false)

// Reset form
const closeModal = () => {
  showModal.value = false
  outletForm.value = { id: '', name: '', address: '', phone: '' }
  formErrors.value = { name: '', address: '', phone: '' }
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

const createOutlet = () => {
  isEditing.value = false
  outletForm.value = { id: '', name: '', address: '', phone: '' }
  formErrors.value = { name: '', address: '', phone: '' }
  showModal.value = true
}

const editOutlet = (o: Outlet) => {
  isEditing.value = true
  outletForm.value = { id: o.id, name: o.name, address: o.address, phone: o.phone }
  formErrors.value = { name: '', address: '', phone: '' }
  showModal.value = true
}

const validateForm = () => {
  let valid = true
  formErrors.value = { name: '', address: '', phone: '' }
  
  if (!outletForm.value.name.trim()) {
    formErrors.value.name = 'Nama outlet wajib diisi'
    valid = false
  }
  if (!outletForm.value.address.trim()) {
    formErrors.value.address = 'Alamat wajib diisi'
    valid = false
  }
  if (!outletForm.value.phone.trim()) {
    formErrors.value.phone = 'Nomor telepon wajib diisi'
    valid = false
  } else if (!/^\+62\d{8,13}$/.test(outletForm.value.phone.trim().replace(/\s/g, ''))) {
    formErrors.value.phone = 'Format harus +62xxx (contoh: +628123456789)'
    valid = false
  }
  
  return valid
}

const submitForm = async () => {
  if (!validateForm()) return
  
  formLoading.value = true
  try {
    const payload = {
      name: outletForm.value.name.trim(),
      address: outletForm.value.address.trim(),
      phone: outletForm.value.phone.trim().replace(/\s/g, '')
    }

    if (isEditing.value) {
      await $fetch(`/api/outlets/${outletForm.value.id}`, {
        method: 'PUT',
        headers: { Authorization: authStore.authHeader },
        body: payload
      })
      toastSuccess('Outlet berhasil diperbarui')
    } else {
      await $fetch('/api/outlets', {
        method: 'POST',
        headers: { Authorization: authStore.authHeader },
        body: payload
      })
      toastSuccess('Outlet berhasil ditambahkan')
    }
    
    closeModal()
    await refresh()
  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else {
      toastError(err?.data?.message || 'Gagal menyimpan outlet')
    }
  } finally {
    formLoading.value = false
  }
}

const deleteOutlet = (o: Outlet) => {
  outletToDelete.value = o
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!outletToDelete.value) return
  
  deleteLoading.value = true
  try {
    await $fetch(`/api/outlets/${outletToDelete.value.id}`, {
      method: 'DELETE',
      headers: { Authorization: authStore.authHeader }
    })
    toastSuccess('Outlet berhasil dihapus')
    showDeleteModal.value = false
    await refresh()
  } catch (err: any) {
    const status = err?.statusCode || err?.status || err?.response?.status
    if (status === 401) {
      toastError('Sesi kamu habis, silakan login ulang')
      authStore.logout()
      router.push('/owner/login')
    } else {
      toastError('Gagal menghapus outlet, coba lagi')
    }
  } finally {
    deleteLoading.value = false
    outletToDelete.value = null
  }
}

onActivated(() => {
  refresh()
})
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
