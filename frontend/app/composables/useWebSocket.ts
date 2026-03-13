import { useAuthStore } from '~/stores/auth'
import { useNotificationStore } from '~/stores/notification'
import { useToast } from '~/composables/useToast'

export const useWebSocket = () => {
    const authStore = useAuthStore()
    const notifStore = useNotificationStore()
    const { info } = useToast()

    let ws: WebSocket | null = null
    let reconnectTimer: any = null
    let reconnectDelay = 1000
    const MAX_RECONNECT_DELAY = 30000

    const connect = () => {
        if (!import.meta.client || !authStore.isLoggedIn || ws) return

        // Base URL from runtime config or hardcoded for now
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const host = window.location.hostname === 'localhost' ? 'localhost:8080' : window.location.host

        // We need to pass the token. Since standard WebSocket doesn't support headers easily, 
        // we use a ticket or just protocol sub-protocol? 
        // Actually, our backend middleware checks Authorization header. 
        // Browser WS API doesn't support custom headers. 
        // WORKAROUND: Pass token as query param, then middleware handles it.

        const wsUrl = `${protocol}//${host}/api/v1/ws/connect?token=${authStore.token}`

        ws = new WebSocket(wsUrl)

        ws.onopen = () => {
            reconnectDelay = 1000
        }

        ws.onmessage = (event) => {
            try {
                const msg = JSON.parse(event.data)
                notifStore.addNotification({
                    ...msg,
                    id: msg.id || Math.random().toString(36).substring(7), // Fallback if no ID from WS yet
                    is_read: false,
                    created_at: new Date().toISOString()
                })
                info(msg.title || 'Notifikasi Baru')
            } catch (e) {
                console.error('WS Message error', e)
            }
        }

        ws.onclose = () => {
            ws = null
            if (authStore.isLoggedIn) {
                scheduleReconnect()
            }
        }

        ws.onerror = (err) => {
            console.error('WebSocket Error', err)
            ws?.close()
        }
    }

    const scheduleReconnect = () => {
        if (reconnectTimer) clearTimeout(reconnectTimer)
        reconnectTimer = setTimeout(() => {
            reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY)
            connect()
        }, reconnectDelay)
    }

    const disconnect = () => {
        if (reconnectTimer) clearTimeout(reconnectTimer)
        if (ws) {
            ws.close()
            ws = null
        }
    }

    return { connect, disconnect }
}
