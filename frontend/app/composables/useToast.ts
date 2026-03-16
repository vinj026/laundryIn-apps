import { ref, readonly } from 'vue'

interface Toast {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

const toasts = ref<Toast[]>([])
let counter = 0
const timers = new Map<number, ReturnType<typeof setTimeout>>()

export const useToast = () => {
  const getDefaultDuration = (type: Toast['type']) => {
    if (type === 'success') return 2500
    if (type === 'error') return 4500
    return 3000
  }

  const show = (message: string, type: Toast['type'] = 'info', durationMs?: number) => {
    const id = ++counter
    toasts.value.push({ id, message, type })

    const duration = durationMs ?? getDefaultDuration(type)
    if (duration > 0) {
      const timer = setTimeout(() => {
        remove(id)
      }, duration)
      timers.set(id, timer)
    }
  }

  const remove = (id: number) => {
    const timer = timers.get(id)
    if (timer) {
      clearTimeout(timer)
      timers.delete(id)
    }
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  return {
    toasts: readonly(toasts),
    remove,
    success: (msg: string, durationMs?: number) => show(msg, 'success', durationMs),
    error: (msg: string, durationMs?: number) => show(msg, 'error', durationMs),
    info: (msg: string, durationMs?: number) => show(msg, 'info', durationMs)
  }
}
