import { apiClient } from './client'
import type {
  ActuatorCommand,
  AdminNode,
  NodeDetail,
  NodeLogs,
  NodeMode,
  NodeSummary,
  TelemetryPoint,
} from '../types'

const isMock = import.meta.env.VITE_MOCK === '1'

async function getMock() {
  const { mockApi } = await import('./mock')
  return mockApi
}

export async function getNodes(): Promise<NodeSummary[]> {
  if (isMock) return (await getMock()).getNodes()
  const response = await apiClient.get<NodeSummary[]>('/api/nodes')
  return response.data
}

export async function getNode(nodeId: string): Promise<NodeDetail> {
  if (isMock) return (await getMock()).getNode(nodeId)
  const response = await apiClient.get<NodeDetail>(`/api/nodes/${encodeURIComponent(nodeId)}`)
  return response.data
}

export async function claimNode(nodeId: string, claimCode: string): Promise<NodeDetail> {
  if (isMock) return (await getMock()).claimNode(nodeId, claimCode)
  const response = await apiClient.post<NodeDetail>('/api/nodes/claim', {
    node_id: nodeId.trim(),
    claim_code: claimCode.trim().toUpperCase(),
  })
  return response.data
}

export async function unclaimNode(nodeId: string): Promise<void> {
  if (isMock) return (await getMock()).unclaimNode(nodeId)
  await apiClient.delete(`/api/nodes/${encodeURIComponent(nodeId)}`)
}

export async function renameNode(nodeId: string, name: string): Promise<NodeDetail> {
  if (isMock) return (await getMock()).renameNode(nodeId, name)
  const response = await apiClient.patch<NodeDetail>(`/api/nodes/${encodeURIComponent(nodeId)}`, {
    name,
  })
  return response.data
}

export async function setNodeMode(nodeId: string, mode: NodeMode): Promise<NodeDetail> {
  if (isMock) return (await getMock()).setMode(nodeId, mode)
  const response = await apiClient.put<NodeDetail>(
    `/api/nodes/${encodeURIComponent(nodeId)}/mode`,
    { mode },
  )
  return response.data
}

export async function setNodeActuators(nodeId: string, command: ActuatorCommand): Promise<void> {
  if (isMock) return (await getMock()).setActuators(nodeId, command)
  await apiClient.put(`/api/nodes/${encodeURIComponent(nodeId)}/actuators`, command)
}

export async function getNodeTelemetry(
  nodeId: string,
  from: string,
  to: string,
): Promise<TelemetryPoint[]> {
  if (isMock) return (await getMock()).getTelemetry(nodeId, from, to)
  const response = await apiClient.get<TelemetryPoint[]>(
    `/api/nodes/${encodeURIComponent(nodeId)}/telemetry`,
    { params: { from, to } },
  )
  return response.data
}

export async function getNodeLogs(nodeId: string): Promise<NodeLogs> {
  if (isMock) return (await getMock()).getLogs(nodeId)
  const response = await apiClient.get<NodeLogs>(`/api/nodes/${encodeURIComponent(nodeId)}/logs`)
  return response.data
}

export async function saveTelegramChatId(telegram_chat_id: string): Promise<void> {
  if (isMock) return (await getMock()).saveTelegramChatId(telegram_chat_id)
  await apiClient.put('/api/me/telegram', { telegram_chat_id })
}

export async function getAdminNodes(): Promise<AdminNode[]> {
  if (isMock) return (await getMock()).getAdminNodes()
  const response = await apiClient.get<AdminNode[]>('/api/admin/nodes')
  return response.data
}

export async function createAdminNode(nodeId: string): Promise<AdminNode> {
  if (isMock) return (await getMock()).createAdminNode(nodeId)
  const response = await apiClient.post<AdminNode>('/api/admin/nodes', {
    node_id: nodeId.trim().toLowerCase(),
  })
  return response.data
}

export async function regenerateClaimCode(nodeId: string): Promise<AdminNode> {
  if (isMock) return (await getMock()).regenerateClaimCode(nodeId)
  const response = await apiClient.post<AdminNode>(
    `/api/admin/nodes/${encodeURIComponent(nodeId)}/claim-code`,
  )
  return response.data
}

export async function adminUnclaimNode(nodeId: string): Promise<void> {
  if (isMock) return (await getMock()).adminUnclaimNode(nodeId)
  await apiClient.delete(`/api/admin/nodes/${encodeURIComponent(nodeId)}/owner`)
}

export async function deleteAdminNode(nodeId: string): Promise<void> {
  if (isMock) return (await getMock()).deleteAdminNode(nodeId)
  await apiClient.delete(`/api/admin/nodes/${encodeURIComponent(nodeId)}`)
}
