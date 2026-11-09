import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import {
  claimNode as claimNodeRequest,
  getAdminNodes,
  getNode,
  getNodes,
  renameNode as renameNodeRequest,
  setNodeActuators,
  setNodeMode,
  unclaimNode as unclaimNodeRequest,
} from '../api/nodes'
import type {
  ActuatorCommand,
  AdminNode,
  NodeDetail,
  NodeSummary,
  RealtimeEvent,
} from '../types'

export const useNodesStore = defineStore('nodes', () => {
  const nodes = ref<NodeSummary[]>([])
  const details = ref<Record<string, NodeDetail>>({})
  const adminNodes = ref<AdminNode[]>([])
  const alertVersion = ref<Record<string, number>>({})
  const pendingActuators = ref<Record<string, boolean>>({})
  const activeNodeId = ref('')
  const loading = ref(false)
  const error = ref('')
  const activeNode = computed(() => details.value[activeNodeId.value] ?? null)

  function reset() {
    nodes.value = []
    details.value = {}
    adminNodes.value = []
    alertVersion.value = {}
    pendingActuators.value = {}
    activeNodeId.value = ''
    loading.value = false
    error.value = ''
  }

  async function loadNodes() {
    loading.value = true
    error.value = ''
    try {
      nodes.value = await getNodes()
      if (nodes.value.length === 0) {
        activeNodeId.value = ''
      } else if (!activeNodeId.value || !nodes.value.some((n) => n.node_id === activeNodeId.value)) {
        activeNodeId.value = nodes.value[0]!.node_id
      }
    } catch (cause) {
      error.value = cause instanceof Error ? cause.message : 'Node tidak dapat dimuat.'
      throw cause
    } finally {
      loading.value = false
    }
  }

  async function loadAdminNodes() {
    loading.value = true
    error.value = ''
    try {
      adminNodes.value = await getAdminNodes()
    } catch (cause) {
      error.value = cause instanceof Error ? cause.message : 'Daftar node admin tidak dapat dimuat.'
      throw cause
    } finally {
      loading.value = false
    }
  }

  async function loadNode(nodeId: string) {
    activeNodeId.value = nodeId
    const detail = await getNode(nodeId)
    updateDetail(detail)
    return detail
  }

  function updateDetail(detail: NodeDetail) {
    details.value[detail.node_id] = detail
    const {
      class_temp: _classTemp,
      class_turbidity: _classTurbidity,
      class_overall: _classOverall,
      heater: _heater,
      fan: _fan,
      pump: _pump,
      failsafe: _failsafe,
      ...summary
    } = detail
    const index = nodes.value.findIndex((node) => node.node_id === detail.node_id)
    if (index < 0) nodes.value.push(summary)
    else nodes.value[index] = summary
  }

  async function claimNode(nodeId: string, claimCode: string) {
    const detail = await claimNodeRequest(nodeId, claimCode)
    updateDetail(detail)
    activeNodeId.value = detail.node_id
    return detail
  }

  async function unclaimNode(nodeId: string) {
    await unclaimNodeRequest(nodeId)
    delete details.value[nodeId]
    delete pendingActuators.value[nodeId]
    nodes.value = nodes.value.filter((item) => item.node_id !== nodeId)
    if (activeNodeId.value === nodeId) {
      activeNodeId.value = nodes.value[0]?.node_id ?? ''
    }
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
    const adminNode = adminNodes.value.find((item) => item.node_id === event.node_id)
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
      if (adminNode) {
        adminNode.online = event.online
        if (event.online) {
          adminNode.last_seen = new Date().toISOString()
        }
      }
    } else if (event.type === 'mode') {
      if (detail) detail.mode = event.mode
      if (node) node.mode = event.mode
    } else if (event.type === 'alert') {
      alertVersion.value[event.node_id] = (alertVersion.value[event.node_id] ?? 0) + 1
    }
  }

  return {
    nodes,
    details,
    adminNodes,
    alertVersion,
    pendingActuators,
    activeNodeId,
    activeNode,
    loading,
    error,
    reset,
    loadNodes,
    loadAdminNodes,
    loadNode,
    updateDetail,
    claimNode,
    unclaimNode,
    renameNode,
    changeMode,
    changeActuators,
    applyRealtimeEvent,
  }
})
