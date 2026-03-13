<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  modelValue: string // format YYYY-MM-DD
}>()
const emit = defineEmits(['update:modelValue'])

// Generate 14 hari: hari ini + 13 hari ke depan
const days = computed(() => {
  return Array.from({ length: 14 }, (_, i) => {
    const date = new Date()
    date.setDate(date.getDate() + i)
    return {
      label: date.toLocaleDateString('en-US', { weekday: 'short' }).toUpperCase(),
      date: date.getDate(),
      month: date.toLocaleDateString('en-US', { month: 'short' }),
      value: date.toISOString().substring(0, 10)
    }
  })
})

function select(value: string) {
  emit('update:modelValue', value)
}

const carouselRef = ref<HTMLElement | null>(null)
const canScrollLeft = ref(false)
const canScrollRight = ref(true) // assume there's always scrollable content first

const onScroll = () => {
  if (!carouselRef.value) return
  const { scrollLeft, scrollWidth, clientWidth } = carouselRef.value
  canScrollLeft.value = scrollLeft > 0
  // add a small 1px threshold for fractional rounding
  canScrollRight.value = Math.ceil(scrollLeft + clientWidth) < scrollWidth - 1
}

const scroll = (direction: 'left' | 'right') => {
  if (!carouselRef.value) return
  const elem = carouselRef.value
  // Scroll roughly 4 cards worth (80% of width)
  const scrollAmount = elem.clientWidth * 0.8 
  elem.scrollBy({
    left: direction === 'left' ? -scrollAmount : scrollAmount,
    behavior: 'smooth'
  })
}

onMounted(() => {
  setTimeout(onScroll, 100) // Ensure DOM layout settles
  if (carouselRef.value) {
    carouselRef.value.addEventListener('scroll', onScroll)
  }
  window.addEventListener('resize', onScroll)
})

onUnmounted(() => {
  if (carouselRef.value) {
    carouselRef.value.removeEventListener('scroll', onScroll)
  }
  window.removeEventListener('resize', onScroll)
})
</script>

<template>
  <div class="relative w-full group">
    <!-- Left Navigation Button -->
    <button
      v-show="canScrollLeft"
      @click.prevent="scroll('left')"
      class="absolute -left-3 top-1/2 -translate-y-1/2 z-10 w-9 h-9 rounded-full bg-surface-raised border border-border flex items-center justify-center text-surface-onSurface transition-all hover:bg-surface-containerHigh shadow-lg"
      aria-label="Scroll left"
      type="button"
    >
      <span class="material-symbols-outlined text-[20px]">chevron_left</span>
    </button>

    <div class="date-carousel relative z-0" ref="carouselRef">
      <div
        v-for="day in days"
        :key="day.value"
        class="date-card"
        :class="{ active: modelValue === day.value }"
        @click="select(day.value)"
      >
        <span class="day-label">{{ day.label }}</span>
        <span class="day-number">{{ day.date }}</span>
        <span class="day-month">{{ day.month }}</span>
      </div>
    </div>

    <!-- Right Navigation Button -->
    <button
      v-show="canScrollRight"
      @click.prevent="scroll('right')"
      class="absolute -right-3 top-1/2 -translate-y-1/2 z-10 w-9 h-9 rounded-full bg-surface-raised border border-border flex items-center justify-center text-surface-onSurface transition-all hover:bg-surface-containerHigh shadow-lg"
      aria-label="Scroll right"
      type="button"
    >
      <span class="material-symbols-outlined text-[20px]">chevron_right</span>
    </button>
  </div>
</template>

<style scoped>
.date-carousel {
  display: flex;
  gap: 8px;
  align-items: center;        /* card non-aktif lebih pendek, center secara vertikal */
  overflow-x: auto;
  padding: 8px 4px 12px;
  scroll-snap-type: x mandatory;

  /* sembunyikan scrollbar */
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.date-carousel::-webkit-scrollbar { display: none; }

/* ── Base card ── */
.date-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;

  /* Menampilkan tepat 5 card di layar dengan menghitung (100% viewport span - 4 gaps) / 5 */
  flex: 0 0 calc((100% - 32px) / 5);
  
  padding: 14px 8px;
  border-radius: 14px;
  border: 1.5px solid var(--border-border);
  background: var(--surface-container);
  cursor: pointer;
  scroll-snap-align: start;

  /* transisi smooth saat scale berubah */
  transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1),
              border-color 0.2s ease,
              background 0.2s ease,
              opacity 0.2s ease;

  /* default: mengecil & redup */
  transform: scale(0.85);
  opacity: 0.5;
}

.date-card:hover:not(.active) {
  opacity: 0.75;
  transform: scale(0.9);
}

/* ── Active card ── */
.date-card.active {
  transform: scale(1);           /* full size */
  opacity: 1;
  border-color: #2dd4bf;         /* accent */
  background: rgba(45, 212, 191, 0.1);
}

/* ── Text dalam card ── */
.day-label {
  font-size: 9px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--surface-onSurfaceVariant);
}

.day-number {
  font-size: 20px;
  font-weight: 800;
  font-family: 'Roboto Mono', monospace;
  color: var(--surface-onSurface);
  line-height: 1;
}

.day-month {
  font-size: 9px;
  font-weight: 500;
  color: var(--surface-onSurfaceVariant);
}

/* Active state: semua teks lebih terang */
.date-card.active .day-label,
.date-card.active .day-number,
.date-card.active .day-month {
  color: #2dd4bf; /* accent */
}

.date-card.active .day-number {
  color: #ffffff; /* angka utama tetap putih terang */
}
</style>
