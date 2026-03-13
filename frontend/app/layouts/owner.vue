<template>
  <div class="flex flex-col h-full w-full flex-1 min-h-0 relative">

    <!-- Top App Bar -->
    <header class="h-14 shrink-0 glass flex justify-center z-40 w-full border-b border-border">
      <div class="w-full max-w-[1200px] px-4 md:px-8 flex items-center justify-between h-full">
        <div class="flex items-center gap-3">
          <div class="h-8 w-8 bg-primary/15 rounded-lg flex items-center justify-center">
            <span class="material-symbols-outlined text-primary text-lg">local_laundry_service</span>
          </div>
          <h1 class="text-base font-semibold tracking-tight">LaundryIn</h1>
        </div>
        <div class="flex items-center gap-2 profile-menu">
          <!-- Notification Bell -->
          <div class="relative notif-menu">
            <button 
              @click="notifStore.toggleDropdown()"
              class="material-symbols-outlined text-surface-onSurfaceVariant hover:text-surface-onSurface p-2 rounded-xl hover:bg-surface-containerHigh transition-all duration-normal text-[22px] relative"
            >
              notifications
              <span v-if="notifStore.unreadCount > 0" class="absolute top-1.5 right-1.5 w-4 h-4 bg-danger text-white text-[9px] font-bold rounded-full flex items-center justify-center border-2 border-surface">
                {{ notifStore.unreadCount > 99 ? '99+' : notifStore.unreadCount }}
              </span>
            </button>
            <UiNotificationDropdown />
          </div>
          
          <!-- Dropdown Profile -->
          <div class="relative">
            <button @click="showProfileDropdown = !showProfileDropdown" class="h-8 w-8 bg-primary/15 text-primary hover:bg-primary/25 transition-colors rounded-full flex items-center justify-center font-bold text-xs uppercase relative">
              {{ authStore.user?.name ? authStore.user.name.charAt(0) : 'O' }}
            </button>
            <div v-if="showProfileDropdown" class="absolute right-0 top-full mt-2 w-48 bg-surface-raised border border-border rounded-xl shadow-lg py-1 z-50 animate-fade-in group">
              <div class="px-4 py-3 border-b border-border mb-1">
                <p class="text-sm font-bold truncate text-surface-onSurface leading-tight">{{ authStore.user?.name || 'Owner' }}</p>
                <p class="text-[10px] text-surface-onSurfaceVariant truncate mt-0.5">{{ authStore.user?.phone || '' }}</p>
              </div>
              <button @click="logout" class="w-full text-left px-4 py-2 mt-1 mb-1 text-sm font-medium text-danger hover:bg-danger/10 transition-colors flex items-center gap-2">
                <span class="material-symbols-outlined text-[18px]">logout</span>
                Logout
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto w-full flex flex-col items-center pb-20 custom-scrollbar">
      <div class="w-full max-w-[1200px] flex-1">
        <slot />
      </div>
    </main>

    <!-- Bottom Navigation Bar -->
    <nav class="absolute bottom-0 left-0 right-0 h-16 bg-surface-raised/95 backdrop-blur-lg border-t border-border flex justify-center z-50">
      <div class="w-full max-w-[1200px] flex items-center justify-around px-4 h-full">
        <NuxtLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex flex-col items-center justify-center gap-0.5 py-1 min-w-[64px]"
        >
          <div
            class="h-8 w-14 rounded-full flex items-center justify-center transition-all duration-normal"
            :class="isActive(item.to) ? 'bg-primary/15' : ''"
          >
            <span
              class="material-symbols-outlined text-[22px] transition-colors"
              :class="isActive(item.to) ? 'text-primary' : 'text-surface-onSurfaceVariant'"
              :style="isActive(item.to) ? 'font-variation-settings: \'FILL\' 1' : ''"
            >{{ item.icon }}</span>
          </div>
          <span class="text-[10px] font-medium" :class="isActive(item.to) ? 'text-primary' : 'text-surface-onSurfaceVariant'">{{ item.label }}</span>
        </NuxtLink>
      </div>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useNotificationStore } from '~/stores/notification'
import { ref, onMounted, onUnmounted } from 'vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const notifStore = useNotificationStore()

const showProfileDropdown = ref(false)

const closeDropdown = (e: Event) => {
  const target = e.target as HTMLElement
  if (!target.closest('.profile-menu') && !target.closest('.notif-menu')) {
    showProfileDropdown.value = false
    notifStore.toggleDropdown(false)
  }
}

onMounted(() => document.addEventListener('click', closeDropdown))
onUnmounted(() => document.removeEventListener('click', closeDropdown))

const logout = () => {
  authStore.logout()
  router.push('/owner/login')
}

watchEffect(() => {
  if (!import.meta.client) return
  if (!authStore.isOwner) {
    navigateTo('/owner/login')
  }
})

const navItems = [
  { to: '/owner', icon: 'bar_chart', label: 'Analytics' },
  { to: '/owner/outlets', icon: 'storefront', label: 'Outlets' },
  { to: '/owner/services', icon: 'dry_cleaning', label: 'Services' },
  { to: '/owner/orders', icon: 'view_list', label: 'Orders' },
]

const isActive = (path: string) => {
  if (path === '/owner') return route.path === '/owner'
  return route.path.startsWith(path)
}
</script>
