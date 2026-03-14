import { useAuthStore } from '~/stores/auth'

export const useApiFetch = (path: string | (() => string), options: any = {}) => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    // Default to true if not explicitly false
    const needsAuth = options.authenticated !== false

    // Map internal /api paths to actual base URL if provided
    const actualPath = String(path).startsWith('/api')
        ? String(path).replace('/api', '')
        : path

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
    const actualPath = path.startsWith('/api') ? path.replace('/api', '') : path

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
