import { useAuthStore } from '~/stores/auth'

export const useApiFetch: typeof useFetch = (path, options = {}) => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    // Map internal /api paths to actual base URL if provided
    const actualPath = String(path).startsWith('/api')
        ? String(path).replace('/api', '')
        : path

    return useFetch(actualPath, {
        baseURL: config.public.apiBase,
        ...options,
        headers: {
            ...authStore.authHeader ? { Authorization: authStore.authHeader } : {},
            ...options.headers,
        },
    })
}

export const useApiRaw = <T>(path: string, options: any = {}): Promise<T> => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    const actualPath = path.startsWith('/api') ? path.replace('/api', '') : path

    return $fetch<T>(actualPath, {
        baseURL: config.public.apiBase,
        ...options,
        headers: {
            ...authStore.authHeader ? { Authorization: authStore.authHeader } : {},
            ...options.headers,
        },
    })
}
