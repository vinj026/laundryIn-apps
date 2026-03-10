import { defineStore } from 'pinia'

export interface CartItem {
    serviceId: string
    name: string
    price: number
    unit: string
    qty: string // zero-trust frontend uses string
}

export const useCartStore = defineStore('cart', {
    state: () => ({
        items: [] as CartItem[],
        outletId: null as string | null
    }),
    getters: {
        totalPreview(state): number {
            // client-side preview only
            return state.items.reduce((total, item) => {
                return total + (item.price * parseFloat(item.qty || '0'))
            }, 0)
        },
        itemCount(state): number {
            return state.items.length
        }
    },
    actions: {
        setOutlet(id: string) {
            if (this.outletId !== id) {
                this.items = [] // clear cart if changing outlet
                this.outletId = id
            }
        },
        addItem(item: CartItem) {
            const existing = this.items.find(i => i.serviceId === item.serviceId)
            if (existing) {
                existing.qty = (parseFloat(existing.qty || '0') + parseFloat(item.qty || '0')).toString()
            } else {
                this.items.push({ ...item })
            }
        },
        updateQty(serviceId: string, qty: string) {
            const existing = this.items.find(i => i.serviceId === serviceId)
            if (existing) {
                if (parseFloat(qty || '0') <= 0) {
                    this.removeItem(serviceId)
                } else {
                    existing.qty = qty
                }
            }
        },
        removeItem(serviceId: string) {
            this.items = this.items.filter(i => i.serviceId !== serviceId)
        },
        clearCart() {
            this.items = []
            this.outletId = null
        }
    }
})
