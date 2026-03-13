export interface ApiResponse<T> {
  status: string
  message: string
  data: T
  unread_count?: number
}

export interface PaginatedResponse<T> {
  data: T
  page: number
  limit: number
  total: number
  total_pages: number
}

