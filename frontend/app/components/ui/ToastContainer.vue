<template>
  <Teleport to="body">
    <div class="fixed bottom-20 right-4 z-[9999] flex flex-col gap-2 pointer-events-none">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="toast-item pointer-events-auto"
          :class="{
            'toast-success': toast.type === 'success',
            'toast-error': toast.type === 'error',
            'toast-info': toast.type === 'info'
          }"
        >
          <span class="material-symbols-outlined text-base">
            {{ toast.type === 'success' ? 'check_circle' : toast.type === 'error' ? 'error' : 'info' }}
          </span>
          <div class="flex-1 min-w-0 pr-1">
            <span class="text-sm font-medium">{{ toast.message }}</span>
          </div>
          <button 
            @click="remove(toast.id)"
            class="h-6 w-6 rounded-lg hover:bg-black/5 flex items-center justify-center transition-colors shrink-0"
          >
            <span class="material-symbols-outlined text-[18px]">close</span>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const { toasts, remove } = useToast()
</script>

<style scoped>
.toast-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 12px;
  min-width: 200px;
  max-width: 320px;
  backdrop-filter: blur(8px);
}

.toast-success {
  background: rgba(45, 212, 191, 0.15);
  border: 1px solid rgba(45, 212, 191, 0.3);
  color: #2dd4bf;
}

.toast-error {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.toast-info {
  background: rgba(148, 163, 184, 0.15);
  border: 1px solid rgba(148, 163, 184, 0.3);
  color: #94a3b8;
}

/* Slide from right animation */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
