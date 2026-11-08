import { apiClient, useMockApi } from './client'
import { mockApi } from './mock'

export async function login(username: string, password: string): Promise<{ token: string }> {
  if (useMockApi) return mockApi.login(username, password)
  const response = await apiClient.post<{ token: string }>('/api/auth/login', { username, password })
  return response.data
}
