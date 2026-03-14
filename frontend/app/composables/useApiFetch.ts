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
    })
}
