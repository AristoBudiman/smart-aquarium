import { HttpError } from './client'
import type {
  ActuatorCommand,
  Alert,
  NodeDetail,
  NodeLogs,
  NodeMode,
  NodeSummary,
  RealtimeEvent,
  StatusClass,
  TelemetryPoint,
} from '../types'

const now = Date.now()
const iso = (offsetMs = 0) => new Date(now + offsetMs).toISOString()

const details = new Map<string, NodeDetail>([
  [
    'aq-sim001',
    {
      node_id: 'aq-sim001',
      name: 'Akuarium Terumbu',
      online: true,
      mode: 'auto',
      last_seen: iso(),
      temp: 27.3,
      ntu: 12.4,
      status: 'NORMAL',
      class_temp: 'NORMAL',
      class_turbidity: 'NORMAL',
      class_overall: 'NORMAL',
      heater: 0,
      fan: 0,
      pump: 0,
      failsafe: false,
    },
  ],
  [
    'aq-1a2b3c',
    {
      node_id: 'aq-1a2b3c',
      name: 'Akuarium Laguna',
      online: true,
      mode: 'manual',
      last_seen: iso(-45_000),
      temp: 32.4,
      ntu: 78.2,
      status: 'WARNING',
      class_temp: 'WARNING',
      class_turbidity: 'CAUTION',
      class_overall: 'WARNING',
      heater: 0,
      fan: 0,
      pump: 1,
      failsafe: false,
    },
  ],
  [
    'aq-sim003',
    {
      node_id: 'aq-sim003',
      name: 'Akuarium Karantina',
      online: false,
      mode: 'manual',
      last_seen: iso(-3_600_000),
      temp: null,
      ntu: null,
      status: null,
      class_temp: null,
      class_turbidity: null,
      class_overall: null,
      heater: null,
      fan: null,
      pump: null,
      failsafe: false,
    },
  ],
])

const logData = new Map<string, NodeLogs>([
  [
    'aq-sim001',
    {
      actuators: [
        { id: 1, ts: iso(-1_200_000), heater: 0, fan: 0, pump: 0, source: 'auto' },
        { id: 2, ts: iso(-3_600_000), heater: 0, fan: 0, pump: 1, source: 'auto' },
      ],
      alerts: [],
    },
  ],
  [
    'aq-1a2b3c',
    {
      actuators: [
        { id: 3, ts: iso(-45_000), heater: 0, fan: 0, pump: 1, source: 'manual' },
        { id: 4, ts: iso(-5_400_000), heater: 0, fan: 1, pump: 1, source: 'auto' },
      ],
      alerts: [
        { id: 1, ts: iso(-45_000), level: 'WARNING', message: 'Suhu 32,4 °C (WARNING)' },
        { id: 2, ts: iso(-7_200_000), level: 'CAUTION', message: 'Kekeruhan meningkat' },
      ],
    },
  ],
  ['aq-sim003', { actuators: [], alerts: [] }],
])

let pendingActuator: { nodeId: string; command: ActuatorCommand } | null = null
let telegramChatId = ''
let logId = 5

function copy<T>(value: T): T {
  return structuredClone(value)
}

function requireNode(nodeId: string): NodeDetail {
  const node = details.get(nodeId)
  if (!node) throw new HttpError(404, 'Node tidak ditemukan.')
  return node
}

export const mockApi = {
  async login(username: string, password: string): Promise<{ token: string }> {
    if (!username.trim() || !password) {
      throw new HttpError(401, 'Username atau kata sandi tidak valid.')
    }
    return { token: `mock-token-${encodeURIComponent(username.trim())}` }
  },

  async getNodes(): Promise<NodeSummary[]> {
    return [...details.values()].map(({ node_id, name, online, mode, last_seen, temp, ntu, status }) => ({
      node_id,
      name,
      online,
      mode,
      last_seen,
      temp,
      ntu,
      status,
    }))
  },

  async getNode(nodeId: string): Promise<NodeDetail> {
    return copy(requireNode(nodeId))
  },

  async renameNode(nodeId: string, name: string): Promise<NodeDetail> {
    const node = requireNode(nodeId)
    node.name = name.trim()
    return copy(node)
  },

  async setMode(nodeId: string, mode: NodeMode): Promise<NodeDetail> {
    const node = requireNode(nodeId)
    node.mode = mode
    return copy(node)
  },

  async setActuators(nodeId: string, command: ActuatorCommand): Promise<void> {
    const node = requireNode(nodeId)
    if (node.mode !== 'manual') {
      throw new HttpError(409, 'Node harus berada di mode Manual untuk mengubah aktuator.')
    }
    pendingActuator = { nodeId, command: copy(command) }
  },

  async getTelemetry(nodeId: string, from: string, to: string): Promise<TelemetryPoint[]> {
    if (!requireNode(nodeId).online) return []
    const start = new Date(from).getTime()
    const end = new Date(to).getTime()
    const duration = Math.max(end - start, 60_000)
    const count = duration <= 60 * 60_000 ? 12 : duration <= 24 * 60 * 60_000 ? 24 : duration <= 7 * 86400_000 ? 28 : 30
    const points: TelemetryPoint[] = []
    for (let index = 0; index < count; index += 1) {
      const progress = index / (count - 1)
      const timestamp = start + duration * progress
      const temp = 27.2 + Math.sin(index / 2.8) * 1.1 + (nodeId === 'aq-1a2b3c' ? 3.2 : 0)
      const ntu = 13 + Math.cos(index / 3.2) * 5 + (nodeId === 'aq-1a2b3c' ? 58 : 0)
      points.push({
        ts: new Date(timestamp).toISOString(),
        temp: Number(temp.toFixed(1)),
        ntu: Number(ntu.toFixed(1)),
        status: nodeId === 'aq-1a2b3c' ? 'WARNING' : 'NORMAL',
      })
    }
    return points
  },

  async getLogs(nodeId: string): Promise<NodeLogs> {
    requireNode(nodeId)
    return copy(logData.get(nodeId) ?? { actuators: [], alerts: [] })
  },

  async saveTelegramChatId(chatId: string): Promise<void> {
    telegramChatId = chatId.trim()
  },
}

export function nextMockRealtimeEvent(): RealtimeEvent | null {
  if (pendingActuator) {
    const { nodeId, command } = pendingActuator
    pendingActuator = null
    const node = requireNode(nodeId)
    Object.assign(node, command)
    node.last_seen = new Date().toISOString()
    const logs = logData.get(nodeId) ?? { actuators: [], alerts: [] }
    logs.actuators.unshift({
      id: logId,
      ts: node.last_seen,
      ...command,
      source: 'manual',
    })
    logId += 1
    logData.set(nodeId, logs)
    return {
      type: 'state',
      node_id: nodeId,
      heater: command.heater,
      fan: command.fan,
      pump: command.pump,
      failsafe: false,
    }
  }

  const node = requireNode('aq-sim001')
  const temp = Number((27.3 + Math.sin(Date.now() / 18_000) * 0.5).toFixed(1))
  const ntu = Number((12.4 + Math.cos(Date.now() / 22_000) * 2.2).toFixed(1))
  const status: StatusClass = temp > 28 ? 'CAUTION' : 'NORMAL'
  node.temp = temp
  node.ntu = ntu
  node.last_seen = new Date().toISOString()
  node.status = status
  node.class_temp = status
  node.class_turbidity = 'NORMAL'
  node.class_overall = status
  return {
    type: 'telemetry',
    node_id: node.node_id,
    ts: node.last_seen,
    temp,
    ntu,
    class_temp: status,
    class_turbidity: 'NORMAL',
    class_overall: status,
    status,
  }
}

export function getMockTelegramChatId(): string {
  return telegramChatId
}
