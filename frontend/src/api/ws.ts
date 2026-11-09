import type { ActuatorValue, NodeMode, RealtimeEvent, StatusClass } from '../types'

const reconnectLimitMs = 30_000
const isMock = import.meta.env.VITE_MOCK === '1'

function isStatusClass(value: unknown): value is StatusClass {
  return value === 'NORMAL' || value === 'CAUTION' || value === 'WARNING' || value === 'SEVERE'
}

function isActuatorValue(value: unknown): value is ActuatorValue {
  return value === 0 || value === 1 || value === null
}

function parseEvent(data: string): RealtimeEvent | null {
  try {
    const value: unknown = JSON.parse(data)
    if (!value || typeof value !== 'object' || !('type' in value) || !('node_id' in value)) return null
    const event = value as Record<string, unknown>
    if (typeof event.node_id !== 'string') return null
    if (event.type === 'node_status') {
      return typeof event.online === 'boolean'
        ? { type: 'node_status', node_id: event.node_id, online: event.online }
        : null
    }
    if (event.type === 'mode') {
      const mode: NodeMode | null = event.mode === 'auto' || event.mode === 'manual' ? event.mode : null
      return mode ? { type: 'mode', node_id: event.node_id, mode } : null
    }
    if (event.type === 'alert') {
      return typeof event.message === 'string' && isStatusClass(event.level)
        ? { type: 'alert', node_id: event.node_id, level: event.level, message: event.message }
        : null
    }
    if (event.type === 'state') {
      return typeof event.failsafe === 'boolean' &&
        isActuatorValue(event.heater) &&
        isActuatorValue(event.fan) &&
        isActuatorValue(event.pump)
        ? {
            type: 'state',
            node_id: event.node_id,
            heater: event.heater,
            fan: event.fan,
            pump: event.pump,
            failsafe: event.failsafe,
          }
        : null
    }
    if (event.type === 'telemetry') {
      return typeof event.ts === 'string' &&
        typeof event.temp === 'number' &&
        typeof event.ntu === 'number' &&
        isStatusClass(event.status) &&
        isStatusClass(event.class_temp) &&
        isStatusClass(event.class_turbidity) &&
        isStatusClass(event.class_overall)
        ? {
            type: 'telemetry',
            node_id: event.node_id,
            ts: event.ts,
            temp: event.temp,
            ntu: event.ntu,
            status: event.status,
            class_temp: event.class_temp,
            class_turbidity: event.class_turbidity,
            class_overall: event.class_overall,
          }
        : null
    }
    return null
  } catch (error) {
    console.error('Pesan WebSocket tidak dapat dibaca.', error)
    return null
  }
}

export function connectRealtime(onEvent: (event: RealtimeEvent) => void): () => void {
  if (isMock) {
    let timer: number | undefined
    let cancelled = false
    void import('./mock').then(({ nextMockRealtimeEvent }) => {
      if (cancelled) return
      timer = window.setInterval(() => {
        const event = nextMockRealtimeEvent()
        if (event) onEvent(event)
      }, 2_500)
    })
    return () => {
      cancelled = true
      if (timer !== undefined) window.clearInterval(timer)
    }
  }

  let socket: WebSocket | null = null
  let reconnectTimer: number | undefined
  let retryDelay = 1_000
  let closed = false

  const connect = () => {
    if (closed) return
    const token = localStorage.getItem('token')
    if (!token) return
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = new URL('/ws', `${protocol}//${window.location.host}`)
    url.searchParams.set('token', token)
    socket = new WebSocket(url)
    socket.onopen = () => {
      retryDelay = 1_000
    }
    socket.onmessage = (message) => {
      if (typeof message.data !== 'string') return
      const event = parseEvent(message.data)
      if (event) onEvent(event)
    }
    socket.onerror = () => socket?.close()
    socket.onclose = () => {
      if (closed) return
      reconnectTimer = window.setTimeout(connect, retryDelay)
      retryDelay = Math.min(retryDelay * 2, reconnectLimitMs)
    }
  }

  connect()
  return () => {
    closed = true
    if (reconnectTimer !== undefined) window.clearTimeout(reconnectTimer)
    socket?.close()
  }
}
