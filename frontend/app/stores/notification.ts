import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import type { ApiResponse } from '~/types/api'

export interface Notification {
    id: string
    type: string
    title: string
    body: string
    data: any
    is_read: boolean
    created_at: string
}

export const useNotificationStore = defineStore('notification', {
    state: () => ({
        notifications: [] as Notification[],
        unreadCount: 0,
        isOpen: false,
        loading: false
    }),

    actions: {
        addNotification(notif: Notification) {
            // Avoid duplicate by ID if any
            const exists = this.notifications.find(n => n.id === notif.id)
            if (exists) return

            this.notifications.unshift(notif)
            if (!notif.is_read) {
                this.unreadCount++
            }

            // Limit to 50 in memory
            if (this.notifications.length > 50) {
                this.notifications.pop()
            }
        },

        async fetchNotifications(page = 1, limit = 20) {
            const authStore = useAuthStore()
            if (!authStore.isLoggedIn) return

            this.loading = true
            try {
                const res = await useApiRaw<ApiResponse<any>>('/api/notifications', {
                    params: { page, limit }
                })

                const fetchedNotifs = res.data?.data || []

                if (page === 1) {
                    // Merge logic: Start with fetched, but keep anything that was added via WS and not in fetched
                    const apiIds = new Set(fetchedNotifs.map((n: Notification) => n.id))
                    const localOnly = this.notifications.filter(n => !apiIds.has(n.id))

                    // Combine and sort by date (descending)
                    this.notifications = [...fetchedNotifs, ...localOnly]
                        .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
                        .slice(0, 100) // Keep max 100
                } else {
                    // Pagination append
                    const existingIds = new Set(this.notifications.map(n => n.id))
                    const newItems = fetchedNotifs.filter((n: any) => !existingIds.has(n.id))
                    this.notifications.push(...newItems)
                }

                this.unreadCount = res.unread_count ?? 0
            } catch (err) {
                useToast().error('Gagal mengambil notifikasi')
            } finally {
                this.loading = false
            }
        },

        async fetchUnreadCount() {
            const authStore = useAuthStore()
            if (!authStore.isLoggedIn) return

            try {
                const res = await useApiRaw<ApiResponse<{ count: number }>>('/api/notifications/unread-count')
                this.unreadCount = res.data?.count || 0
            } catch (err) {
                // Silent for unread check
            }
        },

        async markAsRead(id: string) {
            const authStore = useAuthStore()
            const notif = this.notifications.find(n => n.id === id)
            if (!notif || notif.is_read) return

            try {
                await useApiRaw(`/api/notifications/${id}/read`, {
                    method: 'PATCH'
                })
                notif.is_read = true
                this.unreadCount = Math.max(0, this.unreadCount - 1)
            } catch (err) {
                useToast().error('Gagal memperbarui status notifikasi')
            }
        },

        async markAllAsRead() {
            const authStore = useAuthStore()
            try {
                await useApiRaw('/api/notifications/read-all', {
                    method: 'PATCH'
                })
                this.notifications.forEach(n => n.is_read = true)
                this.unreadCount = 0
            } catch (err) {
                useToast().error('Gagal memperbarui semua notifikasi')
            }
        },

        toggleDropdown(val?: boolean) {
            this.isOpen = val !== undefined ? val : !this.isOpen
        }
    }
})
