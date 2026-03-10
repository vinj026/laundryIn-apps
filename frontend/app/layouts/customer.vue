<template>
  <div class="flex h-[90vh] md:h-full w-full text-surface-onSurface bg-surface overflow-hidden">
    <!-- Sidebar Navigation -->
    <aside class="w-[80px] shrink-0 bg-surface border-r border-outline/10 flex flex-col justify-between py-6">
      <div class="flex flex-col items-center gap-8">
        <!-- Logo -->
        <div class="h-12 w-12 bg-primary rounded-xl flex items-center justify-center shadow-lg">
          <span class="material-symbols-outlined text-primary-text text-2xl">local_laundry_service</span>
        </div>

        <!-- Nav Items -->
        <nav class="flex flex-col gap-6 w-full items-center">
          <!-- Active Item -->
          <NuxtLink to="/customer" class="flex flex-col items-center gap-1 group">
            <div class="h-10 w-14 bg-surface-containerHigh text-primary rounded-full flex items-center justify-center mb-1 transition-all">
              <span class="material-symbols-outlined text-[24px]">add</span>
            </div>
            <span class="text-[10px] text-surface-onSurface font-medium">Booking</span>
          </NuxtLink>

          <!-- Inactive Items -->
          <NuxtLink to="/customer/orders" class="flex flex-col items-center gap-1 group opacity-60 hover:opacity-100 transition-opacity">
            <div class="h-10 w-14 rounded-full flex items-center justify-center mb-1 group-hover:bg-surface-onSurface/10 transition-all">
              <span class="material-symbols-outlined text-[24px]">receipt_long</span>
            </div>
            <span class="text-[10px] text-surface-onSurface">Orders</span>
          </NuxtLink>

          <NuxtLink to="/customer/profile" class="flex flex-col items-center gap-1 group opacity-60 hover:opacity-100 transition-opacity">
            <div class="h-10 w-14 rounded-full flex items-center justify-center mb-1 group-hover:bg-surface-onSurface/10 transition-all">
              <span class="material-symbols-outlined text-[24px]">person</span>
            </div>
            <span class="text-[10px] text-surface-onSurface">Profile</span>
          </NuxtLink>
        </nav>
      </div>

      <!-- Bottom Actions -->
      <div class="flex flex-col items-center gap-6">
        <button class="opacity-60 hover:opacity-100 transition-opacity text-surface-onSurface">
          <span class="material-symbols-outlined">settings</span>
        </button>
        <div class="h-10 w-10 rounded-full bg-surface-containerHigh overflow-hidden border border-outline/20">
          <!-- Mock Avatar -->
          <div class="w-full h-full bg-neutral-600 flex items-center justify-center">
            <span class="material-symbols-outlined text-white text-xl">face</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main Viewport -->
    <main class="flex-1 overflow-hidden relative">
      <slot />
    </main>
  </div>
</template>
