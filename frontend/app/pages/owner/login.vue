<template>
  <div class="min-h-screen flex items-center justify-center p-6 bg-surface">
    <!-- Login Card Surface Container -->
    <div class="w-full max-w-sm bg-surface-container rounded-[2.5rem] p-8 pb-10 shadow-2xl border border-outline/10 text-surface-onSurface flex flex-col items-center">
      
      <!-- Brand & Icon -->
      <div class="mb-8 mt-2 text-center">
        <div class="h-20 w-20 bg-primary/10 text-primary rounded-[1.5rem] flex items-center justify-center mx-auto shadow-inner border border-primary/20 mb-4">
          <span class="material-symbols-outlined text-[48px]">admin_panel_settings</span>
        </div>
        <h1 class="text-3xl font-medium tracking-tight mb-1">Owner Access</h1>
        <p class="text-surface-onSurfaceVariant text-sm">Secure login to your dashboard.</p>
      </div>

      <!-- Auth Form Inputs -->
      <div class="w-full space-y-5 mb-8">
        <div class="relative">
          <span class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-surface-onSurfaceVariant">mail</span>
          <input type="email" placeholder="Email Address" class="w-full bg-surface-containerHigh text-surface-onSurface placeholder-surface-onSurfaceVariant rounded-full py-4 pl-12 pr-4 border border-outline/10 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all" />
        </div>
        <div class="relative">
          <span class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-surface-onSurfaceVariant">lock</span>
          <input type="password" placeholder="Password" class="w-full bg-surface-containerHigh text-surface-onSurface placeholder-surface-onSurfaceVariant rounded-full py-4 pl-12 pr-4 border border-outline/10 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all" />
          <button class="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-surface-onSurfaceVariant hover:text-surface-onSurface">visibility</button>
        </div>
        
        <div class="flex justify-end px-2 pt-1">
          <button class="text-sm font-medium text-primary hover:underline">Forgot password?</button>
        </div>
      </div>

      <!-- Final Actions -->
      <div class="w-full flex justify-center group flex-col gap-3">
        <NuxtLink to="/owner" class="flex items-center gap-3 py-4 w-full justify-center bg-primary text-primary-text font-bold rounded-full shadow-lg transition-transform focus:scale-95 active:scale-95 hover:brightness-110">
          Sign In
          <span class="material-symbols-outlined">login</span>
        </NuxtLink>
        <button class="py-3 w-full text-surface-onSurfaceVariant text-sm font-medium hover:text-primary transition-colors">
          Don't have an account? Register
        </button>
      </div>

    </div>
  </div>
</template>

<script setup>
useHead({
  title: 'LaundryIn - Owner Login'
})
// Layout explicitly empty because this page does not use bottom menu
definePageMeta({
  layout: false
})
</script>
