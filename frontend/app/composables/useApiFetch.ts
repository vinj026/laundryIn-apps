import { useAuthStore } from '~/stores/auth'

export const useApiFetch = <T>(path: any, options: any = {}) => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    // Map internal /api paths to actual base URL if provided
    let actualPath = String(path)
    if (actualPath.startsWith('/api/v1')) {
        actualPath = actualPath.replace('/api/v1', '')
    } else if (actualPath.startsWith('/api')) {
        actualPath = actualPath.replace('/api', '')
    }

    return useFetch<T>(actualPath, {
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
