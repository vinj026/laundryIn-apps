<template>
  <div class="min-h-screen bg-black text-surface-onSurface font-sans app-background flex justify-center sm:p-4 md:p-8">
    <div class="w-full max-w-[768px] min-h-[90vh] bg-surface relative shadow-2xl md:rounded-[2.5rem] border border-outline/20 flex flex-col overflow-x-hidden md:ring-1 md:ring-white/5">
      <NuxtLayout>
        <NuxtPage />
      </NuxtLayout>
    </div>
  </div>
</template>

<style>
.app-background {
  background-image: 
    linear-gradient(to right, rgba(255,255,255,0.04) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(255,255,255,0.04) 1px, transparent 1px);
  background-size: 40px 40px;
}
.material-symbols-outlined {
  font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24;
}
</style>
