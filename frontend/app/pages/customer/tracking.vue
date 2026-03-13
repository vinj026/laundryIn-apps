<template>
  <div class="flex items-center justify-center h-full">
    <span class="material-symbols-outlined animate-spin text-primary text-4xl">progress_activity</span>
  </div>
</template>

<script setup lang="ts">
onMounted(() => {
  navigateTo('/customer/orders')
})

definePageMeta({
  layout: 'customer'
})
</script>
