<template>
  <div v-if="notifStore.isOpen" class="absolute right-0 top-full mt-2 w-80 sm:w-96 bg-surface rounded-2xl shadow-2xl border border-border overflow-hidden z-[200] animate-fade-in">
    <div class="px-5 py-4 border-b border-border flex items-center justify-between bg-surface-container/30">
      <h3 class="font-bold text-sm">Notifikasi</h3>
      <button 
        v-if="notifStore.unreadCount > 0"
        @click="notifStore.markAllAsRead" 
        class="text-[11px] font-semibold text-primary hover:underline transition-all"
      >
        Tandai semua baca
      </button>
    </div>

    <div class="max-h-[400px] overflow-y-auto custom-scrollbar overflow-x-hidden">
      <div v-if="notifStore.loading && notifStore.notifications.length === 0" class="py-10 flex justify-center">
        <span class="material-symbols-outlined animate-spin text-primary">progress_activity</span>
      </div>

      <div v-else-if="notifStore.notifications.length === 0" class="py-12 px-6 text-center text-surface-onSurfaceVariant">
        <div class="w-12 h-12 rounded-2xl bg-surface-containerHigh flex items-center justify-center mx-auto mb-3">
          <span class="material-symbols-outlined text-2xl text-surface-onSurfaceVariant/50">notifications_off</span>
        </div>
        <p class="text-xs font-medium">Belum ada notifikasi.</p>
      </div>

      <div v-else class="divide-y divide-border/50">
        <div 
          v-for="notif in notifStore.notifications" 
          :key="notif.id"
          class="p-4 hover:bg-surface-containerHigh/50 transition-all cursor-pointer relative group"
          :class="{ 'bg-primary/5': !notif.is_read }"
          @click="handleNotifClick(notif)"
        >
          <div class="flex gap-3">
            <div 
              class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0"
              :class="getIconClass(notif.type)"
            >
              <span class="material-symbols-outlined text-[20px]">{{ getIcon(notif.type) }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-2 mb-0.5">
                <h4 class="font-bold text-[13px] truncate" :class="{ 'text-primary': !notif.is_read }">{{ notif.title }}</h4>
                <span class="text-[10px] text-surface-onSurfaceVariant shrink-0 whitespace-nowrap">{{ formatTime(notif.created_at) }}</span>
              </div>
              <p class="text-[11px] text-surface-onSurfaceVariant/90 leading-relaxed line-clamp-2">{{ notif.body }}</p>
            </div>
          </div>
          <!-- Unread Dot -->
          <div v-if="!notif.is_read" class="absolute right-4 bottom-4 w-2 h-2 rounded-full bg-primary shadow-[0_0_8px_rgba(var(--color-primary),0.5)]"></div>
        </div>
      </div>
    </div>

    <!-- View All Footer (Optional but good UX) -->
    <div class="px-5 py-3 border-t border-border bg-surface-container/30 text-center">
      <button class="text-[11px] font-bold text-surface-onSurfaceVariant hover:text-primary transition-colors">
        Tutup
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore, type Notification } from '~/stores/notification'
import { watch } from 'vue'

const notifStore = useNotificationStore()
const router = useRouter()

// Fetch when opened
watch(() => notifStore.isOpen, (newVal) => {
  if (newVal) {
    notifStore.fetchNotifications()
  }
})

const getIcon = (type: string) => {
  switch (type) {
    case 'new_order': return 'shopping_cart'
    case 'order_status': return 'local_laundry_service'
    case 'price_updated': return 'payments'
    case 'order_cancelled': return 'cancel'
    default: return 'notifications'
  }
}

const getIconClass = (type: string) => {
  switch (type) {
    case 'new_order': return 'bg-success/10 text-success'
    case 'order_status': return 'bg-primary/10 text-primary'
    case 'price_updated': return 'bg-warning/10 text-warning'
    case 'order_cancelled': return 'bg-danger/10 text-danger'
    default: return 'bg-surface-containerHigh text-surface-onSurface'
  }
}

const formatTime = (ts: string) => {
  const d = new Date(ts)
  const now = new Date()
  const diff = (now.getTime() - d.getTime()) / 1000 // seconds

  if (diff < 60) return 'Sekarang'
  if (diff < 3600) return `${Math.floor(diff / 60)} mnt`
  if (diff < 86400) return `${Math.floor(diff / 3600)} jam`
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

const handleNotifClick = async (notif: Notification) => {
  await notifStore.markAsRead(notif.id)
  notifStore.toggleDropdown(false)

  // Navigate based on type
  if (notif.type === 'new_order') {
    router.push('/owner/orders')
  } else if (notif.data?.order_id) {
    router.push(`/customer/orders/${notif.data.order_id}`)
  }
}
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.2s cubic-bezier(0, 0, 0.2, 1);
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
