import { useAuthStore } from '~/stores/auth'

export const useApiFetch = (path: string | (() => string), options: any = {}) => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    // Debug logging for production
    if (import.meta.client) {
        console.log('[useApiFetch] Config:', {
            apiBase: config.public.apiBase,
            path,
            isProduction: config.public.apiBase?.startsWith('http')
        })
    }

    // Default to true if not explicitly false
    const needsAuth = options.authenticated !== false

    const p = typeof path === 'function' ? path() : path

    // Path mapping for production
    // If apiBase is a full URL (production), strip /api prefix
    let actualPath = p
    if (config.public.apiBase && config.public.apiBase !== '/api' && config.public.apiBase.startsWith('http')) {
        if (p.startsWith('/api/')) {
            actualPath = p.slice(5) // Remove '/api/'
        } else if (p === '/api') {
            actualPath = ''
        }
    }

    if (import.meta.client) {
        console.log('[useApiFetch] After mapping:', { actualPath, baseURL: config.public.apiBase })
    }

    const fetchOptions = {
        baseURL: config.public.apiBase,
        ...options,
        headers: {
            ...authStore.authHeader ? { Authorization: authStore.authHeader } : {},
            ...options.headers,
        },
    }

    // Guard: Don't fetch if authenticated request but no token exists
    if (needsAuth && !authStore.token) {
        fetchOptions.immediate = false
    }

    // Fix BUG-007: Handle 401 Unauthorized globally
    const { error: toastError } = useToast()
    const router = useRouter()

    const originalOnResponseError = fetchOptions.onResponseError
    fetchOptions.onResponseError = async (context: any) => {
        console.error(`[API FETCH ERROR] ${context.request} -> ${context.response.status}`, context.response._data)

        if (context.response.status === 401 && authStore.token) {
            authStore.logout()
            toastError('Sesi kamu telah kadaluarsa, silakan login ulang')
            router.push('/customer/login')
        }
        if (originalOnResponseError) {
            await originalOnResponseError(context)
        }
    }

    return useFetch(actualPath, fetchOptions)
}

export const useApiRaw = <T>(path: string, options: any = {}): Promise<T> => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    if (import.meta.client) {
        console.log('[useApiRaw] Called with:', { path, apiBase: config.public.apiBase })
    }

    const needsAuth = options.authenticated !== false

    // Path mapping for production
    let actualPath = path
    if (config.public.apiBase && config.public.apiBase !== '/api' && config.public.apiBase.startsWith('http')) {
        if (path.startsWith('/api/')) {
            actualPath = path.slice(5)
        } else if (path === '/api') {
            actualPath = ''
        }
    }

    if (import.meta.client) {
        console.log('[useApiRaw] After mapping:', { actualPath, baseURL: config.public.apiBase })
    }

    if (needsAuth && !authStore.token) {
        return Promise.reject(new Error('Authentication required'))
    }

    return $fetch<T>(actualPath, {
        baseURL: config.public.apiBase,
        ...options,
        headers: {
            ...authStore.authHeader ? { Authorization: authStore.authHeader } : {},
            ...options.headers,
        },
        onResponseError: async (context) => {
            console.error(`[API RAW ERROR] ${context.request} -> ${context.response.status}`, context.response._data)

            if (context.response.status === 401 && authStore.token) {
                authStore.logout()
                const { error: toastError } = useToast()
                toastError('Sesi kamu telah kadaluarsa, silakan login ulang')
                const router = useRouter()
                router.push('/customer/login')
            }
            if (options.onResponseError) {
                await options.onResponseError(context)
            }
        }
    })
}
