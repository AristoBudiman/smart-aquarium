import { apiClient, useMockApi } from './client'
import { mockApi } from './mock'
import type { Me } from '../types'

export async function login(username: string, password: string): Promise<{ token: string }> {
  if (useMockApi) return mockApi.login(username, password)
  const response = await apiClient.post<{ token: string }>('/api/auth/login', {
    username,
    password,
  })
  return response.data
}

export async function register(username: string, password: string): Promise<{ token: string }> {
  if (useMockApi) return mockApi.register(username, password)
  const response = await apiClient.post<{ token: string }>('/api/auth/register', {
    username,
    password,
  })
  return response.data
}

export async function getMe(): Promise<Me> {
  if (useMockApi) return mockApi.getMe()
  const response = await apiClient.get<Me>('/api/me')
  return response.data
}
