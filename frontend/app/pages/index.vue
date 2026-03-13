<template>
  <div class="flex flex-col min-h-screen bg-surface relative overflow-hidden">

    <!-- Subtle gradient background -->
    <div class="absolute inset-0 pointer-events-none">
      <div class="absolute top-1/4 left-1/2 -translate-x-1/2 w-[600px] h-[600px] bg-primary/[0.04] rounded-full blur-[150px]"></div>
    </div>

    <!-- Header -->
    <header class="relative z-10 flex items-center justify-between px-6 py-5">
      <div class="flex items-center gap-2.5">
        <div class="h-9 w-9 bg-primary/15 rounded-lg flex items-center justify-center">
          <span class="material-symbols-outlined text-primary text-lg">local_laundry_service</span>
        </div>
        <span class="text-lg font-semibold tracking-tight">LaundryIn</span>
      </div>
    </header>

    <!-- Hero -->
    <div class="relative z-10 flex-1 flex flex-col justify-center items-center text-center px-6 pb-12">
      <div class="max-w-md mx-auto w-full">

        <!-- Icon -->
        <div class="mb-8 flex justify-center">
          <div class="h-20 w-20 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center">
            <span class="material-symbols-outlined text-primary text-[40px]">local_laundry_service</span>
          </div>
        </div>

        <h1 class="text-4xl md:text-5xl font-bold mb-4 tracking-tight leading-tight">
          Smart Laundry<br/>
          <span class="text-gradient-primary">Platform</span>
        </h1>
        <p class="text-surface-onSurfaceVariant text-base leading-relaxed mb-10 max-w-sm mx-auto">
          Manage your laundry business or order pickup with real-time tracking — all in one place.
        </p>

        <div class="flex flex-col gap-3 w-full">
          <NuxtLink to="/owner/login" class="btn-primary py-4 rounded-2xl text-sm shadow-lg shadow-primary/10">
            <span class="material-symbols-outlined text-[20px]">shield_person</span>
            Owner Dashboard
          </NuxtLink>

          <NuxtLink to="/customer" class="btn-secondary py-4 rounded-2xl text-sm">
            <span class="material-symbols-outlined text-[20px]">storefront</span>
            Find a Laundry
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="relative z-10 text-center py-4 text-xs text-surface-onSurfaceVariant">
      Built with precision. Zero-tolerance pricing.
    </footer>
  </div>
</template>

<script setup>
useHead({
  title: 'LaundryIn — Smart Laundry Platform'
})
</script>
