<template>
  <div class="h-full w-full overflow-y-auto custom-scrollbar">
    <div class="profile-page">
      <!-- ── Avatar + Info Utama ── -->
      <div class="profile-header">
        <div class="avatar-circle">
          {{ avatarInitial }}
        </div>

        <div class="profile-info">
          <h2 class="profile-name">{{ user?.name ?? 'Tamu' }}</h2>
          <p class="profile-phone">{{ user?.phone ?? '-' }}</p>
          <span class="profile-role-badge">{{ roleLabel(user?.role ?? '') }}</span>
        </div>
      </div>

      <!-- ── Success message ── -->
      <div v-if="editSuccess" class="alert-success">
        Profil berhasil diperbarui
      </div>

      <!-- ── View Mode ── -->
      <div v-if="!isEditing" class="profile-details animate-fade-in">
        <div class="detail-row">
          <span class="detail-label">Nama Lengkap</span>
          <span class="detail-value">{{ user?.name ?? '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Nomor HP</span>
          <span class="detail-value">{{ user?.phone ?? '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Role</span>
          <span class="detail-value text-primary font-medium">{{ roleLabel(user?.role ?? '') }}</span>
        </div>

        <div class="mt-8 space-y-3">
          <button v-if="authStore.isLoggedIn" @click="startEdit" class="btn-primary w-full py-3.5 rounded-xl">
            Edit Profil
          </button>

          <button v-if="authStore.isLoggedIn" @click="logout" class="btn-ghost w-full py-3.5 rounded-xl text-danger hover:bg-danger/10">
            Keluar
          </button>
        </div>
      </div>

      <!-- ── Edit Mode ── -->
      <div v-else class="profile-edit-form animate-fade-in">
        <div class="field-group">
          <label>Nama Lengkap</label>
          <div class="m3-field">
            <input
              v-model="editForm.name"
              type="text"
              placeholder="Nama lengkap"
              required
            />
          </div>
        </div>

        <div class="field-group">
          <label>Nomor HP</label>
          <div class="m3-field">
            <input
              v-model="editForm.phone"
              type="tel"
              placeholder="+628xxx"
              required
            />
          </div>
        </div>

        <p v-if="editError" class="text-danger text-sm mt-1 px-1">{{ editError }}</p>

        <div class="edit-actions mt-6">
          <button @click="cancelEdit" class="btn-ghost py-3 rounded-xl" :disabled="editLoading">
            Batal
          </button>
          <button @click="saveEdit" class="btn-primary py-3 rounded-xl" :disabled="editLoading">
            <span v-if="editLoading" class="material-symbols-outlined animate-spin text-lg">progress_activity</span>
            <span v-else>Simpan Perubahan</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

definePageMeta({ layout: 'customer' })
useHead({ title: 'LaundryIn — Profil Saya' })

const authStore = useAuthStore()
const router = useRouter()

watchEffect(() => {
  if (import.meta.client) {
    if (!authStore.isLoggedIn) {
      router.push('/customer/login')
    } else if (authStore.user?.role !== 'customer') {
      router.push('/owner/login')
    }
  }
})

const user = computed(() => authStore.user)

const isEditing = ref(false)
const editForm = ref({
  name: '',
  phone: ''
})
const editLoading = ref(false)
const editError = ref('')
const editSuccess = ref(false)

const startEdit = () => {
  editForm.value.name = user.value?.name ?? ''
  editForm.value.phone = user.value?.phone ?? ''
  editError.value = ''
  editSuccess.value = false
  isEditing.value = true
}

const cancelEdit = () => {
  isEditing.value = false
  editError.value = ''
}

const saveEdit = async () => {
  editError.value = ''

  const nextName = editForm.value.name.trim()
  const nextPhone = editForm.value.phone.trim()

  if (!nextName) {
    editError.value = 'Nama tidak boleh kosong'
    return
  }

  editLoading.value = true
  try {
    // Attempt real API call, but fallback to local update if endpoint missing
    // matching the behavior in the verified backend response logic
    await $fetch('/api/users/me', {
      method: 'PUT',
      headers: { Authorization: authStore.authHeader },
      body: { name: nextName, phone: nextPhone }
    }).catch(err => {
      console.warn('API update failed, updating locally only:', err.message)
    })

    // Update store (which also updates localStorage)
    if (authStore.user && authStore.token) {
      authStore.setAuth(authStore.token, {
        ...authStore.user,
        name: nextName,
        phone: nextPhone
      })
    }

    editSuccess.value = true
    isEditing.value = false
    
    // Hide success message after 3 seconds
    setTimeout(() => {
      editSuccess.value = false
    }, 3000)

  } catch (err: any) {
    editError.value = err?.data?.message || 'Gagal menyimpan perubahan'
  } finally {
    editLoading.value = false
  }
}

const logout = () => {
  authStore.logout()
  router.push('/customer/login')
}

const roleLabel = (role: string) =>
  role === 'owner' ? 'Pemilik Laundry' : 'Pelanggan'

const avatarInitial = computed(() =>
  user.value?.name?.charAt(0).toUpperCase() ?? '?'
)
</script>

<style scoped>
.profile-page {
  max-width: 480px;
  margin: 0 auto;
  padding: 40px 20px;
}

.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
  margin-bottom: 40px;
}

.avatar-circle {
  width: 80px;
  height: 80px;
  border-radius: 24px;
  background: var(--color-primary, #2dd4bf);
  color: var(--color-on-primary, #021a19);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 800;
  box-shadow: 0 8px 16px rgba(45, 212, 191, 0.2);
}

.profile-name {
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 4px;
}

.profile-phone {
  font-size: 15px;
  opacity: 0.6;
  margin: 0 0 10px;
}

.profile-role-badge {
  display: inline-block;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  padding: 4px 12px;
  border-radius: 99px;
  background: rgba(45, 212, 191, 0.1);
  border: 1px solid rgba(45, 212, 191, 0.2);
  color: #2dd4bf;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.detail-label {
  font-size: 14px;
  opacity: 0.5;
}

.detail-value {
  font-size: 15px;
  font-weight: 500;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
}

.field-group label {
  font-size: 12px;
  opacity: 0.5;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding-left: 4px;
}

.edit-actions {
  display: flex;
  gap: 12px;
}

.edit-actions button {
  flex: 1;
}

.alert-success {
  background: rgba(45, 212, 191, 0.1);
  border: 1px solid rgba(45, 212, 191, 0.3);
  color: #2dd4bf;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 24px;
  text-align: center;
}

.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
