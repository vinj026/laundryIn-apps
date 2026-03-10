<template>
  <div class="flex flex-col h-full md:min-h-screen w-full text-surface-onSurface bg-surface overflow-hidden relative">
    
    <!-- Top App Bar -->
    <header class="h-16 shrink-0 bg-surface-container flex items-center px-4 justify-between border-b border-outline/10 z-40 relative">
      <div class="flex items-center gap-3">
        <span class="material-symbols-outlined text-primary">admin_panel_settings</span>
        <h1 class="text-xl font-medium tracking-tight">Owner Dashboard</h1>
      </div>
      <div class="flex items-center gap-4">
        <button class="material-symbols-outlined hover:bg-surface-onSurface/10 p-2 rounded-full transition-colors">notifications</button>
        <button class="h-8 w-8 bg-primary-container text-primary-onContainer rounded-full flex items-center justify-center font-bold text-sm shadow-md hover:brightness-110">
          A
        </button>
      </div>
    </header>

    <!-- Main Content Viewport -->
    <main class="flex-1 overflow-y-auto w-full flex flex-col pb-24 relative custom-scrollbar">
      <slot />
    </main>

    <!-- Bottom Navigation Bar -->
    <nav class="absolute bottom-0 left-0 right-0 h-20 bg-surface-container border-t border-outline/20 flex items-center justify-around px-2 z-50">
      
      <NuxtLink to="/owner" class="flex flex-col items-center justify-center gap-1 group w-[72px]" active-class="text-primary-onContainer">
        <div class="w-16 h-8 rounded-full flex items-center justify-center transition-colors group-[.router-link-exact-active]:bg-primary group-[.router-link-exact-active]:shadow-md">
          <span class="material-symbols-outlined transition-colors group-[.router-link-exact-active]:text-primary-text">bar_chart</span>
        </div>
        <span class="text-[11px] font-medium opacity-80 group-[.router-link-exact-active]:opacity-100 group-[.router-link-exact-active]:font-bold">Analytics</span>
      </NuxtLink>

      <NuxtLink to="/owner/outlets" class="flex flex-col items-center justify-center gap-1 group w-[72px]" active-class="text-primary-onContainer">
        <div class="w-16 h-8 rounded-full flex items-center justify-center transition-colors group-[.router-link-exact-active]:bg-primary group-[.router-link-exact-active]:shadow-md">
          <span class="material-symbols-outlined transition-colors group-[.router-link-exact-active]:text-primary-text">storefront</span>
        </div>
        <span class="text-[11px] font-medium opacity-80 group-[.router-link-exact-active]:opacity-100 group-[.router-link-exact-active]:font-bold">Outlets</span>
      </NuxtLink>

      <NuxtLink to="/owner/services" class="flex flex-col items-center justify-center gap-1 group w-[72px]" active-class="text-primary-onContainer">
        <div class="w-16 h-8 rounded-full flex items-center justify-center transition-colors group-[.router-link-exact-active]:bg-primary group-[.router-link-exact-active]:shadow-md">
          <span class="material-symbols-outlined transition-colors group-[.router-link-exact-active]:text-primary-text">dry_cleaning</span>
        </div>
        <span class="text-[11px] font-medium opacity-80 group-[.router-link-exact-active]:opacity-100 group-[.router-link-exact-active]:font-bold">Services</span>
      </NuxtLink>

      <NuxtLink to="/owner/orders" class="flex flex-col items-center justify-center gap-1 group w-[72px]" active-class="text-primary-onContainer">
        <div class="w-16 h-8 rounded-full flex items-center justify-center transition-colors group-[.router-link-exact-active]:bg-primary group-[.router-link-exact-active]:shadow-md">
          <span class="material-symbols-outlined transition-colors group-[.router-link-exact-active]:text-primary-text">view_list</span>
        </div>
        <span class="text-[11px] font-medium opacity-80 group-[.router-link-exact-active]:opacity-100 group-[.router-link-exact-active]:font-bold">Orders</span>
      </NuxtLink>

    </nav>
  </div>
</template>

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
