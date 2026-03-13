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
            {{ authStore.user?.name ? authStore.user.name.charAt(0) : 'C' }}
          </button>

          <!-- Dropdown Profile -->
          <div v-if="showProfileDropdown" class="absolute bottom-0 left-[68px] w-48 bg-surface-raised border border-border rounded-xl shadow-xl py-1 animate-fade-in group">
            <div class="px-4 py-3 border-b border-border mb-1">
              <p class="text-sm font-bold truncate text-surface-onSurface leading-tight">{{ authStore.user?.name || 'Customer' }}</p>
              <p class="text-[10px] text-surface-onSurfaceVariant truncate mt-0.5">{{ authStore.user?.phone || '' }}</p>
            </div>
            <NuxtLink @click="showProfileDropdown = false" to="/customer/profile" class="w-full text-left px-4 py-2 mt-1 text-sm font-medium text-surface-onSurface hover:bg-surface-containerHigh transition-colors flex items-center gap-2">
              <span class="material-symbols-outlined text-[18px]">person</span>
              Lihat Profil
            </NuxtLink>
            <button @click="logout" class="w-full text-left px-4 py-2 mb-1 text-sm font-medium text-danger hover:bg-danger/10 transition-colors flex items-center gap-2">
              <span class="material-symbols-outlined text-[18px]">logout</span>
              Logout
            </button>
          </div>
        </div>
        <div v-else class="flex justify-center w-full">
          <NuxtLink to="/customer/login" class="h-9 w-9 rounded-xl bg-surface-containerHigh text-surface-onSurfaceVariant hover:text-primary hover:bg-primary/10 transition-colors flex items-center justify-center" title="Login">
            <span class="material-symbols-outlined text-[20px]">login</span>
          </NuxtLink>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 min-h-0 overflow-y-auto relative flex flex-col custom-scrollbar">
      <slot />

      <!-- Mobile Bottom Nav (hidden desktop) -->
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
  showProfileDropdown.value = false
  authStore.logout()
  router.push('/customer/login')
}

const navItems = [
  { to: '/customer', icon: 'explore', label: 'Explore' },
  { to: '/customer/orders', icon: 'receipt_long', label: 'Orders' },
  { to: '/customer/profile', icon: 'person', label: 'Profile' },
]

const isActive = (path: string) => {
  if (path === '/customer') return route.path === '/customer' || route.path.startsWith('/customer/outlet')
  return route.path.startsWith(path)
}
</script>
