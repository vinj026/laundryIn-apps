<template>
  <div class="h-full w-full overflow-y-auto px-6 py-8 md:py-10 lg:px-12 flex flex-col xl:flex-row gap-10">
    <!-- Left Section: Tracking Logic -->
    <div class="flex-1 space-y-8 max-w-3xl">
      <!-- Breadcrumbs & Headline -->
      <div>
        <div class="flex items-center gap-2 mb-2">
          <NuxtLink to="/customer/orders" class="text-sm font-medium text-surface-onSurfaceVariant hover:text-primary transition-colors">My Orders</NuxtLink>
          <span class="material-symbols-outlined text-sm text-surface-onSurfaceVariant">chevron_right</span>
          <span class="text-sm font-medium text-surface-onSurface">Tracking</span>
        </div>
        <h1 class="text-4xl md:text-5xl font-normal text-surface-onSurface mb-2">Order #LND-8821</h1>
        <p class="text-lg text-surface-onSurfaceVariant">
          Estimated Completion: <span class="font-bold text-primary">Today, 5:00 PM</span>
        </p>
      </div>

      <!-- Live Status Banner -->
      <div class="bg-primary-container text-primary-onContainer rounded-3xl p-6 md:p-8 flex flex-col sm:flex-row items-center gap-6 md:gap-8 relative overflow-hidden">
        <div class="flex items-center justify-center h-20 w-20 md:h-24 md:w-24 rounded-full bg-white/20 backdrop-blur-sm shadow-sm shrink-0">
          <span class="material-symbols-outlined text-[40px] md:text-[48px] animate-spin-slow" style="animation-duration: 4s;">local_laundry_service</span>
        </div>
        <div class="flex-1 text-center sm:text-left z-10 w-full">
          <h2 class="text-2xl md:text-3xl font-medium mb-2">Sedang Dicuci</h2>
          <p class="text-primary-onContainer/80 mb-6 leading-relaxed text-sm md:text-base max-w-lg">
              Mesin sedang membersihkan bajumu dengan deterjen eco-friendly. Menunggu siklus pengeringan.
          </p>
          <div class="flex flex-col gap-2">
            <div class="flex justify-between text-xs font-bold uppercase tracking-wider opacity-70">
              <span>Washing Cycle</span>
              <span>75%</span>
            </div>
            <div class="h-1 w-full bg-primary-onContainer/20 rounded-full overflow-hidden">
              <div class="h-full bg-primary-onContainer rounded-full w-[75%] transition-all duration-1000"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Timeline History -->
      <div class="py-4">
        <h3 class="text-xl font-normal text-surface-onSurface mb-6">Tracking History</h3>
        <div class="relative pl-4 border-l-2 border-surface-containerHigh ml-3 space-y-10">
          
          <div class="relative pl-8">
            <div class="absolute -left-[9px] top-1 h-4 w-4 rounded-full bg-primary border-2 border-surface"></div>
            <div class="flex flex-col sm:flex-row sm:items-baseline gap-1 sm:gap-4">
              <span class="text-sm font-bold text-surface-onSurfaceVariant w-32 shrink-0">10:30 AM</span>
              <div>
                <h4 class="text-lg font-medium text-surface-onSurface">Pesanan Dibuat</h4>
                <p class="text-sm text-surface-onSurfaceVariant">Order successfully placed via App.</p>
              </div>
            </div>
          </div>
          
          <div class="relative pl-8">
            <div class="absolute -left-[9px] top-1 h-4 w-4 rounded-full bg-primary border-2 border-surface"></div>
            <div class="flex flex-col sm:flex-row sm:items-baseline gap-1 sm:gap-4">
              <span class="text-sm font-bold text-surface-onSurfaceVariant w-32 shrink-0">11:15 AM</span>
              <div>
                <h4 class="text-lg font-medium text-surface-onSurface">Diterima di Cabang</h4>
                <p class="text-sm text-surface-onSurfaceVariant">Items received by front desk.</p>
              </div>
            </div>
          </div>
          
          <div class="relative pl-8">
            <div class="absolute -left-[11px] top-0 h-6 w-6 rounded-full bg-primary border-[3px] border-surface flex items-center justify-center shadow-[0_0_10px_rgba(0,229,255,0.4)]">
              <div class="h-1.5 w-1.5 bg-[#003641] rounded-full"></div>
            </div>
            <div class="flex flex-col sm:flex-row sm:items-baseline gap-1 sm:gap-4">
              <span class="text-sm font-bold text-primary w-32 shrink-0">12:45 PM</span>
              <div class="bg-surface-container rounded-2xl p-4 -mt-2 w-full border border-primary/20">
                <div class="flex justify-between items-start mb-2">
                  <h4 class="text-lg font-bold text-primary">Sedang Dicuci</h4>
                  <span class="bg-red-500/20 text-red-500 text-[10px] font-bold px-2 py-0.5 rounded-md flex items-center gap-1 border border-red-500/30">
                    <span class="w-1.5 h-1.5 bg-red-500 rounded-full animate-pulse"></span> LIVE
                  </span>
                </div>
                <p class="text-sm text-surface-onSurfaceVariant">Machine #04 - Premium Wash Cycle</p>
                <div class="mt-4 bg-black/50 rounded-xl h-24 flex items-center justify-center border border-outline/10 text-surface-onSurfaceVariant text-xs font-mono">
                  [ Live Camera Feed Blocked ]
                </div>
              </div>
            </div>
          </div>
          
          <div class="relative pl-8 opacity-40">
            <div class="absolute -left-[9px] top-1 h-4 w-4 rounded-full bg-surface-containerHigh border-2 border-surface"></div>
            <div class="flex flex-col sm:flex-row sm:items-baseline gap-1 sm:gap-4">
              <span class="text-sm font-medium text-surface-onSurfaceVariant w-32 shrink-0">--:--</span>
              <div>
                <h4 class="text-lg font-medium text-surface-onSurface">Ready to Pickup</h4>
                <p class="text-sm text-surface-onSurfaceVariant">Pending drying and iron completion.</p>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>

    <!-- Right Section: Receipt & Drivers -->
    <div class="w-full xl:w-[320px] flex flex-col gap-6 shrink-0 mt-4 xl:mt-0 pb-10">
      
      <!-- Receipt Structure -->
      <div class="bg-surface-container rounded-3xl p-6 shadow-sm border border-outline/10">
        <div class="flex items-center justify-between mb-6 pb-4 border-b border-outline/10">
          <h3 class="text-lg font-medium">Order Details</h3>
          <span class="material-symbols-outlined text-surface-onSurfaceVariant">receipt_long</span>
        </div>
        
        <div class="space-y-4 mb-6">
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-3">
              <div class="bg-surface-containerHigh w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold font-mono">2x</div>
              <div class="flex flex-col">
                <span class="text-sm font-medium">Cuci Kering</span>
              </div>
            </div>
            <span class="text-sm font-medium font-mono text-primary">Rp 16k</span>
          </div>
          
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-3">
              <div class="bg-surface-containerHigh w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold font-mono">1x</div>
              <div class="flex flex-col">
                <span class="text-sm font-medium">Cuci Karpet</span>
              </div>
            </div>
            <span class="text-sm font-medium font-mono text-primary">Rp 15k</span>
          </div>
        </div>
        
        <div class="pt-4 border-t border-outline/10 border-dashed space-y-2 font-mono text-sm">
          <div class="flex justify-between text-surface-onSurfaceVariant">
            <span>Subtotal</span>
            <span>Rp 31.000</span>
          </div>
          <div class="flex justify-between text-indigo-400">
            <span>Voucher Diskon</span>
            <span>-Rp 1.000</span>
          </div>
          <div class="flex justify-between text-xl font-bold mt-2 pt-2 text-primary border-t border-outline/10">
            <span>Total</span>
            <span>Rp 30.000</span>
          </div>
        </div>
      </div>

      <!-- Action Required Card -->
      <button class="w-full py-4 rounded-2xl bg-secondary-container text-secondary-onContainer font-medium flex items-center justify-center gap-2 hover:brightness-110 active:scale-95 transition-all shadow-md">
        <span class="material-symbols-outlined">support_agent</span>
        Hubungi Kasir
      </button>

    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'customer'
})
useHead({
  title: 'LaundryIn - Order Tracking'
})
</script>
