import axios, { type AxiosError } from 'axios'

export const apiClient = axios.create({
  baseURL: '',
  timeout: 10_000,
})

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError<{ error?: string }>) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.dispatchEvent(new Event('auth:unauthorized'))
    }
    return Promise.reject(error)
  },
)

export class HttpError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message)
    this.name = 'HttpError'
  }
}

export function getApiErrorStatus(error: unknown): number | undefined {
  if (error instanceof HttpError) return error.status
  if (axios.isAxiosError(error)) return error.response?.status
  return undefined
}

export function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof HttpError) return error.message
  if (axios.isAxiosError<{ error?: string }>(error)) {
    return error.response?.data?.error ?? fallback
  }
  if (error instanceof Error) return error.message
  return fallback
}

export const useMockApi = import.meta.env.VITE_MOCK === '1'
