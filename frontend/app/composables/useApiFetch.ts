import { useAuthStore } from '~/stores/auth'

export const useApiFetch = (path: string | (() => string), options: any = {}) => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    // Default to true if not explicitly false
    const needsAuth = options.authenticated !== false

    // Map internal /api paths to actual base URL if provided
    // If apiBase is just "/api", we don't want to replace /api as it IS the prefix.
    // However, if we hit the proxy, we usually send the relative part.
    const p = typeof path === 'function' ? path() : path
    const actualPath = p.startsWith('/api') && config.public.apiBase !== '/api'
        ? p.replace('/api', '')
        : p

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

    const needsAuth = options.authenticated !== false

    // If apiBase is "/api", we keep the path as is if it starts with /api
    const actualPath = path.startsWith('/api') && config.public.apiBase !== '/api'
        ? path.replace('/api', '')
        : path

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
            if (context.response.status === 401 && authStore.token) {
                authStore.logout()
                const { error: toastError } = useToast()
                toastError('Sesi kamu telah kadaluarsa, silakan login ulang')
                useRouter().push('/customer/login')
            }
            if (options.onResponseError) {
                await options.onResponseError(context)
            }
        }
    })
}
