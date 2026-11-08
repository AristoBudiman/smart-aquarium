import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import {
  getNode,
  getNodes,
  renameNode as renameNodeRequest,
  setNodeActuators,
  setNodeMode,
} from '../api/nodes'
import type { ActuatorCommand, Alert, NodeDetail, NodeSummary, RealtimeEvent } from '../types'

export const useNodesStore = defineStore('nodes', () => {
  const nodes = ref<NodeSummary[]>([])
  const details = ref<Record<string, NodeDetail>>({})
  const realtimeAlerts = ref<Record<string, Alert[]>>({})
  const pendingActuators = ref<Record<string, boolean>>({})
  const activeNodeId = ref('')
  const loading = ref(false)
  const error = ref('')
  let nextRealtimeAlertId = -1
  const activeNode = computed(() => details.value[activeNodeId.value] ?? null)

  async function loadNodes() {
    loading.value = true
    error.value = ''
    try {
      nodes.value = await getNodes()
      if (!activeNodeId.value && nodes.value[0]) activeNodeId.value = nodes.value[0].node_id
    } catch (cause) {
      error.value = cause instanceof Error ? cause.message : 'Node tidak dapat dimuat.'
      throw cause
    } finally {
      loading.value = false
    }
  }

  async function loadNode(nodeId: string) {
    activeNodeId.value = nodeId
    const detail = await getNode(nodeId)
    details.value[nodeId] = detail
    const index = nodes.value.findIndex((node) => node.node_id === nodeId)
    const { class_temp: _classTemp, class_turbidity: _classTurbidity, class_overall: _classOverall, heater: _heater, fan: _fan, pump: _pump, failsafe: _failsafe, ...summary } = detail
    if (index < 0) nodes.value.push(summary)
    else nodes.value[index] = summary
    return detail
  }

  function updateDetail(detail: NodeDetail) {
    details.value[detail.node_id] = detail
    const { class_temp: _classTemp, class_turbidity: _classTurbidity, class_overall: _classOverall, heater: _heater, fan: _fan, pump: _pump, failsafe: _failsafe, ...summary } = detail
    const index = nodes.value.findIndex((node) => node.node_id === detail.node_id)
    if (index < 0) nodes.value.push(summary)
    else nodes.value[index] = summary
  }

  async function renameNode(nodeId: string, name: string) {
    updateDetail(await renameNodeRequest(nodeId, name))
  }

  async function changeMode(nodeId: string, mode: NodeDetail['mode']) {
    updateDetail(await setNodeMode(nodeId, mode))
  }

  async function changeActuators(nodeId: string, command: ActuatorCommand) {
    await setNodeActuators(nodeId, command)
    pendingActuators.value[nodeId] = true
  }

  function applyRealtimeEvent(event: RealtimeEvent) {
    const detail = details.value[event.node_id]
    const node = nodes.value.find((item) => item.node_id === event.node_id)
    if (event.type === 'telemetry') {
      if (detail) {
        Object.assign(detail, {
          last_seen: event.ts,
          temp: event.temp,
          ntu: event.ntu,
          class_temp: event.class_temp,
          class_turbidity: event.class_turbidity,
          class_overall: event.class_overall,
          status: event.status,
        })
      }
      if (node) {
        Object.assign(node, {
          last_seen: event.ts,
          temp: event.temp,
          ntu: event.ntu,
          status: event.status,
        })
      }
    } else if (event.type === 'state') {
      if (detail) {
        detail.heater = event.heater
        detail.fan = event.fan
        detail.pump = event.pump
        detail.failsafe = event.failsafe
      }
      pendingActuators.value[event.node_id] = false
    } else if (event.type === 'node_status') {
      if (detail) detail.online = event.online
      if (node) node.online = event.online
    } else if (event.type === 'mode') {
      if (detail) detail.mode = event.mode
      if (node) node.mode = event.mode
    } else if (event.type === 'alert') {
      const alerts = realtimeAlerts.value[event.node_id] ?? []
      alerts.unshift({
        id: nextRealtimeAlertId,
        ts: new Date().toISOString(),
        level: event.level,
        message: event.message,
      })
      nextRealtimeAlertId -= 1
      realtimeAlerts.value[event.node_id] = alerts
    }
  }

  return {
    nodes,
    details,
    realtimeAlerts,
    pendingActuators,
    activeNodeId,
    activeNode,
    loading,
    error,
    loadNodes,
    loadNode,
    updateDetail,
    renameNode,
    changeMode,
    changeActuators,
    applyRealtimeEvent,
  }
})
