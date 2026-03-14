<template>
  <div class="flex h-full w-full relative">
    <!-- Desktop Sidebar (hidden mobile) -->
    <aside class="hidden lg:flex w-[72px] shrink-0 bg-surface-raised border-r border-border flex-col justify-between py-5 z-40 relative">
      <div class="flex flex-col items-center gap-7">
        <!-- Logo -->
        <div class="h-10 w-10 bg-primary/15 rounded-xl flex items-center justify-center">
          <span class="material-symbols-outlined text-primary text-xl">local_laundry_service</span>
        </div>

        <!-- Nav Items -->
        <nav class="flex flex-col gap-1.5 w-full items-center">
           <!-- Notification Bell (Shared) -->
           <div v-if="authStore.isLoggedIn" class="relative notif-menu w-full flex flex-col items-center mb-2">
            <button 
              @click="notifStore.toggleDropdown()"
              class="nav-item flex flex-col items-center gap-0.5 w-full py-1.5 group"
              :class="{ 'text-primary': notifStore.isOpen, 'text-surface-onSurfaceVariant': !notifStore.isOpen }"
            >
              <div
                class="h-9 w-12 rounded-xl flex items-center justify-center transition-all duration-normal relative"
                :class="notifStore.isOpen ? 'bg-primary/15 text-primary' : 'text-surface-onSurfaceVariant group-hover:bg-surface-containerHigh group-hover:text-surface-onSurface'"
              >
                <span class="material-symbols-outlined text-[22px]" :style="notifStore.isOpen ? 'font-variation-settings: \'FILL\' 1' : ''">notifications</span>
                <span v-if="notifStore.unreadCount > 0" class="absolute -top-1 -right-1 w-4 h-4 bg-danger text-white text-[9px] font-bold rounded-full flex items-center justify-center border-2 border-surface">
                  {{ notifStore.unreadCount > 9 ? '9+' : notifStore.unreadCount }}
                </span>
              </div>
              <span class="text-[10px] font-medium">Notif</span>
            </button>
            <UiNotificationDropdown class="!bottom-auto !top-0 !left-[68px] !right-auto" />
          </div>

          <NuxtLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="nav-item flex flex-col items-center gap-0.5 w-full py-1.5 group"
            :class="{ 'is-active': isActive(item.to) }"
          >
            <div
              class="h-9 w-12 rounded-xl flex items-center justify-center transition-all duration-normal"
              :class="isActive(item.to) ? 'bg-primary/15 text-primary' : 'text-surface-onSurfaceVariant group-hover:bg-surface-containerHigh group-hover:text-surface-onSurface'"
            >
              <span class="material-symbols-outlined text-[22px]" :style="isActive(item.to) ? 'font-variation-settings: \'FILL\' 1' : ''">{{ item.icon }}</span>
            </div>
            <span class="text-[10px] font-medium" :class="isActive(item.to) ? 'text-primary' : 'text-surface-onSurfaceVariant'">{{ item.label }}</span>
          </NuxtLink>
        </nav>
      </div>

      <!-- Bottom -->
      <div class="flex flex-col items-center justify-end flex-1 profile-menu w-full pb-4">
        <div v-if="authStore.isLoggedIn" class="relative flex justify-center w-full">
          <button @click="showProfileDropdown = !showProfileDropdown" class="h-9 w-9 rounded-full bg-primary/15 text-primary hover:bg-primary/25 transition-colors flex items-center justify-center font-bold text-sm uppercase">
            {{ authStore.user?.name ? authStore.user.name.charAt(0) : 'O' }}
          </button>

          <!-- Dropdown Profile -->
          <div v-if="showProfileDropdown" class="absolute bottom-0 left-[68px] w-48 bg-surface-raised border border-border rounded-xl shadow-xl py-1 animate-fade-in group">
            <div class="px-4 py-3 border-b border-border mb-1">
              <p class="text-sm font-bold truncate text-surface-onSurface leading-tight">{{ authStore.user?.name || 'Owner' }}</p>
              <p class="text-[10px] text-surface-onSurfaceVariant truncate mt-0.5">{{ authStore.user?.phone || '' }}</p>
            </div>
            <button @click="logout" class="w-full text-left px-4 py-2 mb-1 text-sm font-medium text-danger hover:bg-danger/10 transition-colors flex items-center gap-2">
              <span class="material-symbols-outlined text-[18px]">logout</span>
              Logout
            </button>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 min-h-0 overflow-y-auto relative flex flex-col custom-scrollbar">
      <slot />

      <!-- Mobile Bottom Nav -->
      <nav class="lg:hidden shrink-0 h-16 bg-surface-raised/95 backdrop-blur-lg border-t border-border flex items-center justify-around px-4 z-50">
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
        <button
          @click="showProfileDropdown = !showProfileDropdown"
          class="flex flex-col items-center justify-center gap-0.5 py-1 min-w-[64px] relative"
        >
          <div class="h-8 w-8 rounded-full bg-primary/15 text-primary flex items-center justify-center font-bold text-xs uppercase overflow-hidden">
            {{ authStore.user?.name ? authStore.user.name.charAt(0) : 'O' }}
          </div>
          <span class="text-[10px] font-medium text-surface-onSurfaceVariant">Profile</span>
          
           <!-- Mobile Profile Dropdown -->
           <div v-if="showProfileDropdown" class="absolute bottom-full right-4 mb-2 w-48 bg-surface-raised border border-border rounded-xl shadow-xl py-1 animate-fade-in z-[60]">
            <div class="px-4 py-3 border-b border-border mb-1">
              <p class="text-sm font-bold truncate text-surface-onSurface leading-tight">{{ authStore.user?.name || 'Owner' }}</p>
            </div>
            <button @click="logout" class="w-full text-left px-4 py-2 text-sm font-medium text-danger flex items-center gap-2">
              <span class="material-symbols-outlined text-[18px]">logout</span>
              Logout
            </button>
          </div>
        </button>
      </nav>
    </main>
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
