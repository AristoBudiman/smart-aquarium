import { apiClient, useMockApi } from './client'
import { mockApi } from './mock'
import type {
  ActuatorCommand,
  NodeDetail,
  NodeLogs,
  NodeMode,
  NodeSummary,
  TelemetryPoint,
} from '../types'

export async function getNodes(): Promise<NodeSummary[]> {
  if (useMockApi) return mockApi.getNodes()
  const response = await apiClient.get<NodeSummary[]>('/api/nodes')
  return response.data
}

export async function getNode(nodeId: string): Promise<NodeDetail> {
  if (useMockApi) return mockApi.getNode(nodeId)
  const response = await apiClient.get<NodeDetail>(`/api/nodes/${encodeURIComponent(nodeId)}`)
  return response.data
}

export async function renameNode(nodeId: string, name: string): Promise<NodeDetail> {
  if (useMockApi) return mockApi.renameNode(nodeId, name)
  const response = await apiClient.patch<NodeDetail>(`/api/nodes/${encodeURIComponent(nodeId)}`, { name })
  return response.data
}

export async function setNodeMode(nodeId: string, mode: NodeMode): Promise<NodeDetail> {
  if (useMockApi) return mockApi.setMode(nodeId, mode)
  const response = await apiClient.put<NodeDetail>(`/api/nodes/${encodeURIComponent(nodeId)}/mode`, { mode })
  return response.data
}

export async function setNodeActuators(nodeId: string, command: ActuatorCommand): Promise<void> {
  if (useMockApi) return mockApi.setActuators(nodeId, command)
  await apiClient.put(`/api/nodes/${encodeURIComponent(nodeId)}/actuators`, command)
}

export async function getNodeTelemetry(
  nodeId: string,
  from: string,
  to: string,
): Promise<TelemetryPoint[]> {
  if (useMockApi) return mockApi.getTelemetry(nodeId, from, to)
  const response = await apiClient.get<TelemetryPoint[]>(
    `/api/nodes/${encodeURIComponent(nodeId)}/telemetry`,
    { params: { from, to } },
  )
  return response.data
}

export async function getNodeLogs(nodeId: string): Promise<NodeLogs> {
  if (useMockApi) return mockApi.getLogs(nodeId)
  const response = await apiClient.get<NodeLogs>(`/api/nodes/${encodeURIComponent(nodeId)}/logs`)
  return response.data
}

export async function saveTelegramChatId(telegram_chat_id: string): Promise<void> {
  if (useMockApi) return mockApi.saveTelegramChatId(telegram_chat_id)
  await apiClient.put('/api/me/telegram', { telegram_chat_id })
}
