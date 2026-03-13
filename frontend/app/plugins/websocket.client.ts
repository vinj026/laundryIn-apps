export default defineNuxtPlugin(() => {
    const authStore = useAuthStore()
    const notifStore = useNotificationStore()
    const { connect, disconnect } = useWebSocket()

    // Initial fetch unread count if already logged in
    if (import.meta.client && authStore.isLoggedIn) {
        notifStore.fetchUnreadCount()
    }

    // Watch login state for WS connection
    watch(() => authStore.isLoggedIn, (loggedIn) => {
        if (import.meta.client) {
            if (loggedIn) connect()
            else disconnect()
        }
    }, { immediate: true })
})
