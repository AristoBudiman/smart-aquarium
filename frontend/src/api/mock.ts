import { HttpError } from './client'
import type {
  ActuatorCommand,
  ActuatorLog,
  AdminNode,
  Alert,
  Me,
  NodeDetail,
  NodeLogs,
  NodeMode,
  NodeSummary,
  RealtimeEvent,
  StatusClass,
  TelemetryPoint,
} from '../types'

interface MockNodeRecord extends NodeDetail {
  created_at: string
  claim_code: string
  owner_username: string | null
}

const now = Date.now()
const iso = (offsetMs = 0) => new Date(now + offsetMs).toISOString()

const initialClientOwner = 'klien'

const details = new Map<string, MockNodeRecord>([
  [
    'aq-sim001',
    {
      node_id: 'aq-sim001',
      name: 'Akuarium Terumbu',
      online: true,
      mode: 'auto',
      last_seen: iso(-15_000),
      created_at: iso(-86_400_000 * 5),
      claim_code: 'A1B2C3D4',
      owner_username: initialClientOwner,
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
      created_at: iso(-86_400_000 * 3),
      claim_code: 'E5F6A7B8',
      owner_username: initialClientOwner,
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
      created_at: iso(-86_400_000),
      claim_code: 'KLIM2026',
      owner_username: null,
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

function buildInitialActuators(nodeId: string, count: number): ActuatorLog[] {
  const items: ActuatorLog[] = []
  for (let i = 0; i < count; i += 1) {
    items.push({
      id: (nodeId === 'aq-sim001' ? 100 : 200) + i + 1,
      ts: iso(-(i + 1) * 900_000),
      heater: i % 4 === 0 ? 1 : 0,
      fan: i % 3 === 0 ? 1 : 0,
      pump: i % 2 === 0 ? 1 : 0,
      source: i % 5 === 0 ? 'manual' : 'auto',
    })
  }
  return items
}

function buildInitialAlerts(nodeId: string, count: number): Alert[] {
  const levels: StatusClass[] = ['WARNING', 'CAUTION', 'SEVERE', 'WARNING']
  const messages = [
    'Suhu 32,4 °C melewati batas aman (WARNING)',
    'Kekeruhan air meningkat ke 78,2 NTU (CAUTION)',
    'Suhu air naik mendadak ke 33,8 °C (SEVERE)',
    'Kekeruhan air mencapai 94,0 NTU, pompa diaktifkan (WARNING)',
  ]
  const items: Alert[] = []
  for (let i = 0; i < count; i += 1) {
    items.push({
      id: (nodeId === 'aq-sim001' ? 300 : 400) + i + 1,
      ts: iso(-(i + 1) * 1_200_000),
      level: levels[i % levels.length]!,
      message: messages[i % messages.length]!,
    })
  }
  return items
}

const logData = new Map<string, NodeLogs>([
  [
    'aq-sim001',
    {
      actuators: buildInitialActuators('aq-sim001', 14),
      alerts: buildInitialAlerts('aq-sim001', 4),
    },
  ],
  [
    'aq-1a2b3c',
    {
      actuators: buildInitialActuators('aq-1a2b3c', 16),
      alerts: buildInitialAlerts('aq-1a2b3c', 15),
    },
  ],
  ['aq-sim003', { actuators: [], alerts: [] }],
])

const registeredUsers = new Set<string>(['admin'])
let activeClientUsername = initialClientOwner
let pendingActuator: { nodeId: string; command: ActuatorCommand } | null = null
let telegramChatId = '123456789'
let logId = 900

function copy<T>(value: T): T {
  return structuredClone(value)
}

function getCurrentUser(): { username: string; isAdmin: boolean } {
  const raw = localStorage.getItem('token') ?? ''
  if (raw.startsWith('mock-token-')) {
    const username = decodeURIComponent(raw.slice('mock-token-'.length)).toLowerCase()
    if (username) {
      return { username, isAdmin: username === 'admin' }
    }
  }
  return { username: activeClientUsername, isAdmin: false }
}

function bindDefaultClientNodes(username: string) {
  if (username === 'admin') return
  const prevOwner = activeClientUsername
  activeClientUsername = username
  for (const record of details.values()) {
    if (record.owner_username === prevOwner) {
      record.owner_username = username
    }
  }
}

function toNodeDetail(record: MockNodeRecord): NodeDetail {
  const { created_at: _c, claim_code: _cc, owner_username: _o, ...detail } = record
  return copy(detail)
}

function toAdminNode(record: MockNodeRecord): AdminNode {
  return {
    node_id: record.node_id,
    name: record.name,
    online: record.online,
    last_seen: record.last_seen,
    created_at: record.created_at,
    claim_code: record.claim_code,
    owner_username: record.owner_username,
  }
}

function requireClientNode(nodeId: string): MockNodeRecord {
  const { username, isAdmin } = getCurrentUser()
  if (isAdmin) throw new HttpError(404, 'Node tidak ditemukan')
  const node = details.get(nodeId)
  if (!node || node.owner_username !== username) {
    throw new HttpError(404, 'Node tidak ditemukan')
  }
  return node
}

function requireAdmin(): void {
  const { isAdmin } = getCurrentUser()
  if (!isAdmin) throw new HttpError(403, 'Khusus admin')
}

function randomClaimCode(): string {
  const chars = '0123456789ABCDEF'
  let out = ''
  for (let i = 0; i < 8; i += 1) {
    out += chars[Math.floor(Math.random() * chars.length)]
  }
  return out
}

const usernamePattern = /^[a-z0-9_.-]{3,32}$/
const nodeIdPattern = /^(aq-[0-9a-f]{6}|aq-sim[0-9]{3})$/

export const mockApi = {
  async login(username: string, password: string): Promise<{ token: string }> {
    const normalized = username.trim().toLowerCase()
    if (!normalized || !password) {
      throw new HttpError(401, 'Username atau kata sandi tidak valid')
    }
    bindDefaultClientNodes(normalized)
    return { token: `mock-token-${encodeURIComponent(normalized)}` }
  },

  async register(username: string, password: string): Promise<{ token: string }> {
    const normalized = username.trim().toLowerCase()
    if (!usernamePattern.test(normalized)) {
      throw new HttpError(
        400,
        'Username harus 3–32 karakter dan hanya boleh berisi huruf kecil a-z, angka 0-9, _, ., atau -',
      )
    }
    if (password.length < 8) {
      throw new HttpError(400, 'Password minimal 8 karakter')
    }
    if (registeredUsers.has(normalized)) {
      throw new HttpError(409, 'Username sudah dipakai')
    }
    registeredUsers.add(normalized)
    bindDefaultClientNodes(normalized)
    return { token: `mock-token-${encodeURIComponent(normalized)}` }
  },

  async getMe(): Promise<Me> {
    const { username, isAdmin } = getCurrentUser()
    return {
      id: isAdmin ? 1 : 2,
      username,
      is_admin: isAdmin,
      telegram_linked: Boolean(telegramChatId.trim()),
    }
  },

  async getNodes(): Promise<NodeSummary[]> {
    const { username, isAdmin } = getCurrentUser()
    if (isAdmin) return []
    return [...details.values()]
      .filter((node) => node.owner_username === username)
      .map(({ node_id, name, online, mode, last_seen, temp, ntu, status }) => ({
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
    return toNodeDetail(requireClientNode(nodeId))
  },

  async claimNode(nodeId: string, claimCode: string): Promise<NodeDetail> {
    const { username, isAdmin } = getCurrentUser()
    if (isAdmin) {
      throw new HttpError(403, 'Admin tidak dapat mengklaim akuarium')
    }
    const id = nodeId.trim().toLowerCase()
    const code = claimCode.trim().toUpperCase()
    const node = details.get(id)
    if (!node || node.claim_code.toUpperCase() !== code) {
      throw new HttpError(404, 'ID node atau kode klaim salah')
    }
    if (node.owner_username && node.owner_username !== username) {
      throw new HttpError(409, 'Akuarium sudah terdaftar di akun lain')
    }
    node.owner_username = username
    return toNodeDetail(node)
  },

  async unclaimNode(nodeId: string): Promise<void> {
    const node = requireClientNode(nodeId)
    node.owner_username = null
    node.mode = 'manual'
  },

  async renameNode(nodeId: string, name: string): Promise<NodeDetail> {
    const node = requireClientNode(nodeId)
    node.name = name.trim()
    return toNodeDetail(node)
  },

  async setMode(nodeId: string, mode: NodeMode): Promise<NodeDetail> {
    const node = requireClientNode(nodeId)
    node.mode = mode
    return toNodeDetail(node)
  },

  async setActuators(nodeId: string, command: ActuatorCommand): Promise<void> {
    const node = requireClientNode(nodeId)
    if (node.mode !== 'manual') {
      throw new HttpError(409, 'Node harus berada di mode Manual untuk mengubah aktuator.')
    }
    pendingActuator = { nodeId, command: copy(command) }
  },

  async getTelemetry(nodeId: string, from: string, to: string): Promise<TelemetryPoint[]> {
    const node = requireClientNode(nodeId)
    if (!node.online) return []
    const start = new Date(from).getTime()
    const end = new Date(to).getTime()
    const duration = Math.max(end - start, 60_000)
    const count =
      duration <= 60 * 60_000
        ? 12
        : duration <= 24 * 60 * 60_000
          ? 24
          : duration <= 7 * 86400_000
            ? 28
            : 30
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
    requireClientNode(nodeId)
    return copy(logData.get(nodeId) ?? { actuators: [], alerts: [] })
  },

  async saveTelegramChatId(chatId: string): Promise<void> {
    telegramChatId = chatId.trim()
  },

  async getAdminNodes(): Promise<AdminNode[]> {
    requireAdmin()
    return [...details.values()]
      .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
      .map(toAdminNode)
  },

  async createAdminNode(nodeId: string): Promise<AdminNode> {
    requireAdmin()
    const normalized = nodeId.trim().toLowerCase()
    if (!nodeIdPattern.test(normalized)) {
      throw new HttpError(400, 'Format node_id tidak valid (gunakan aq-xxxxxx atau aq-simxxx)')
    }
    if (details.has(normalized)) {
      throw new HttpError(409, 'Node sudah terdaftar')
    }
    const created: MockNodeRecord = {
      node_id: normalized,
      name: `Akuarium ${normalized}`,
      online: false,
      mode: 'manual',
      last_seen: null,
      created_at: new Date().toISOString(),
      claim_code: randomClaimCode(),
      owner_username: null,
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
    }
    details.set(normalized, created)
    logData.set(normalized, { actuators: [], alerts: [] })
    return toAdminNode(created)
  },

  async regenerateClaimCode(nodeId: string): Promise<AdminNode> {
    requireAdmin()
    const node = details.get(nodeId)
    if (!node) throw new HttpError(404, 'Node tidak ditemukan')
    node.claim_code = randomClaimCode()
    return toAdminNode(node)
  },

  async adminUnclaimNode(nodeId: string): Promise<void> {
    requireAdmin()
    const node = details.get(nodeId)
    if (!node) throw new HttpError(404, 'Node tidak ditemukan')
    node.owner_username = null
    node.mode = 'manual'
  },

  async deleteAdminNode(nodeId: string): Promise<void> {
    requireAdmin()
    if (!details.has(nodeId)) throw new HttpError(404, 'Node tidak ditemukan')
    details.delete(nodeId)
    logData.delete(nodeId)
  },
}

export function nextMockRealtimeEvent(): RealtimeEvent | null {
  const { username, isAdmin } = getCurrentUser()
  if (isAdmin) {
    const first = details.get('aq-sim001')
    if (!first) return null
    return {
      type: 'node_status',
      node_id: first.node_id,
      online: first.online,
    }
  }

  if (pendingActuator) {
    const { nodeId, command } = pendingActuator
    pendingActuator = null
    const node = details.get(nodeId)
    if (!node || node.owner_username !== username) return null
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

  const node = details.get('aq-sim001')
  if (!node || node.owner_username !== username) return null
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
