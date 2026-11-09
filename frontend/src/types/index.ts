export type StatusClass = 'NORMAL' | 'CAUTION' | 'WARNING' | 'SEVERE'
export type NodeMode = 'auto' | 'manual'
export type ActuatorValue = 0 | 1 | null

export interface Me {
  id: number
  username: string
  is_admin: boolean
  telegram_linked: boolean
}

export interface AdminNode {
  node_id: string
  name: string
  online: boolean
  last_seen: string | null
  created_at: string
  claim_code: string
  owner_username: string | null
}

export interface NodeSummary {
  node_id: string
  name: string
  online: boolean
  mode: NodeMode
  last_seen: string | null
  temp: number | null
  ntu: number | null
  status: StatusClass | null
}

export interface NodeDetail extends NodeSummary {
  class_temp: StatusClass | null
  class_turbidity: StatusClass | null
  class_overall: StatusClass | null
  heater: ActuatorValue
  fan: ActuatorValue
  pump: ActuatorValue
  failsafe: boolean
}

export interface TelemetryPoint {
  ts: string
  temp: number
  ntu: number
  status: StatusClass
}

export interface ActuatorLog {
  id: number
  ts: string
  heater: 0 | 1
  fan: 0 | 1
  pump: 0 | 1
  source: 'auto' | 'manual' | 'failsafe'
}

export interface Alert {
  id: number
  ts: string
  level: StatusClass
  message: string
}

export interface NodeLogs {
  actuators: ActuatorLog[]
  alerts: Alert[]
}

export interface ActuatorCommand {
  heater: 0 | 1
  fan: 0 | 1
  pump: 0 | 1
}

export type RealtimeEvent =
  | {
      type: 'telemetry'
      node_id: string
      ts: string
      temp: number
      ntu: number
      class_temp: StatusClass
      class_turbidity: StatusClass
      class_overall: StatusClass
      status: StatusClass
    }
  | {
      type: 'state'
      node_id: string
      heater: ActuatorValue
      fan: ActuatorValue
      pump: ActuatorValue
      failsafe: boolean
    }
  | { type: 'node_status'; node_id: string; online: boolean }
  | { type: 'mode'; node_id: string; mode: NodeMode }
  | { type: 'alert'; node_id: string; level: StatusClass; message: string }
