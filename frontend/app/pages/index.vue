<template>
  <div class="flex flex-col h-full grow p-6 md:p-10 bg-surface relative">
    <header class="flex items-center justify-between mb-16">
      <div class="flex items-center gap-3">
        <div class="h-10 w-10 flex items-center justify-center bg-primary rounded-xl text-primary-text">
          <span class="material-symbols-outlined text-[24px]">local_laundry_service</span>
        </div>
        <span class="text-2xl font-medium tracking-tight text-surface-onSurface">LaundryIn</span>
      </div>
      <button class="material-symbols-outlined p-2 rounded-full hover:bg-surface-onSurface/10 transition-colors text-surface-onSurface">
        more_vert
      </button>
    </header>

    <div class="flex-grow flex flex-col justify-center items-center text-center w-full max-w-sm mx-auto gap-12">
      <div>
        <h1 class="text-5xl font-normal mb-4 text-surface-onSurface">Smart POS</h1>
        <p class="text-surface-onSurfaceVariant text-lg leading-relaxed px-4">
          Manage your laundry business seamlessly with accurate tracking.
        </p>
      </div>
      
      <div class="w-full flex flex-col gap-4 mt-4">
        <NuxtLink to="/owner/login" class="py-4 px-6 w-full rounded-full bg-primary text-primary-text font-medium flex items-center justify-center gap-2 hover:brightness-110 transition-all shadow-md">
          <span class="material-symbols-outlined">shield_person</span>
          Owner Access
        </NuxtLink>

        <NuxtLink to="/customer" class="py-4 px-6 w-full rounded-full bg-surface-containerHigh text-surface-onSurface font-medium flex items-center justify-center gap-2 hover:bg-surface-onSurface/10 transition-all border border-surface-outlineVariant">
          <span class="material-symbols-outlined">storefront</span>
          Find a Laundry
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
useHead({
  title: 'LaundryIn - Modern Kasir'
})
</script>
