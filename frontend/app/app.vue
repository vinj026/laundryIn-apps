<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <div class="h-screen w-screen bg-surface text-surface-onSurface font-sans overflow-hidden flex flex-col">
    <NuxtLoadingIndicator color="#2dd4bf" :height="2" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UiToastContainer />
  </div>
</template>
