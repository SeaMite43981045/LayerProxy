<template>
  <div style="display: flex; flex-direction: column; gap: 16px; height: 100%;">
    <div style="display: flex; gap: 16px;">
      <n-card style="flex: 1;">
        <n-statistic :label="t('dashboard.cpu')" :value="systemInfo.cpu_model">
          <template #prefix>
            <n-icon><Cpu /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px;">
          <n-text depth="3">{{ systemInfo.cpu_cores }} {{ t('dashboard.cores') }} / {{ systemInfo.cpu_threads }} {{ t('dashboard.threads') }}</n-text>
          <n-progress type="line" :percentage="Math.round(systemInfo.cpu_usage)" indicator-placement="inside" />
        </n-space>
      </n-card>
      <n-card style="flex: 1;">
        <n-statistic :label="t('dashboard.memory')" :value="formatBytes(systemInfo.memory_total)">
          <template #prefix>
            <n-icon><MemoryStick /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px;">
          <n-text depth="3">{{ t('dashboard.used') }}: {{ formatBytes(systemInfo.memory_used) }} / {{ t('dashboard.free') }}: {{ formatBytes(systemInfo.memory_free) }}</n-text>
          <n-progress type="line" :percentage="memoryUsagePercent" indicator-placement="inside" />
        </n-space>
      </n-card>
      <n-card style="flex: 1;">
        <n-statistic :label="t('dashboard.system')" :value="systemInfo.os_name">
          <template #prefix>
            <n-icon><Monitor /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px;">
          <n-text depth="3">{{ t('dashboard.version') }}: {{ systemInfo.os_version }}</n-text>
          <n-text depth="3">{{ t('dashboard.uptime') }}: {{ formatUptime(systemInfo.uptime) }}</n-text>
        </n-space>
      </n-card>
    </div>
    <n-card style="flex: 1; overflow: auto;">
      <template #header>
        <n-space align="center">
          <span>{{ t('dashboard.instances') }}</span>
          <n-button size="small" @click="showModal = true">{{ t('dashboard.addInstance') }}</n-button>
        </n-space>
      </template>
      <n-data-table :columns="columns" :data="instances" :bordered="false" size="small" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard, NStatistic, NIcon, NProgress, NSpace, NText,
  NDataTable, NButton,
} from 'naive-ui'
import { Cpu, MemoryStick, Monitor } from 'lucide-vue-next'
import HttpRequest from '@/http/httpRequest'
import { createMessage } from '@/message/showMessage'

const httpRequest = new HttpRequest()
const { t } = useI18n()

const systemInfo = ref({
  cpu_model: '', cpu_cores: 0, cpu_threads: 0, cpu_usage: 0,
  memory_total: 0, memory_used: 0, memory_free: 0,
  os_name: '', os_version: '', uptime: 0,
})

const instances = ref([])
const showModal = ref(false)

const memoryUsagePercent = computed(() => {
  if (!systemInfo.value.memory_total) return 0
  return Math.round((systemInfo.value.memory_used / systemInfo.value.memory_total) * 100)
})

const columns = [
  { title: t('dashboard.name'), key: 'name' },
  { title: t('dashboard.backendIP'), key: 'backend_ip' },
  { title: t('dashboard.subdomain'), key: 'subdomain' },
  { title: t('dashboard.status'), key: 'status' },
  { title: t('dashboard.actions'), key: 'actions' },
]

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatUptime(seconds: number) {
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  return `${d}d ${h}h ${m}m`
}

async function loadData() {
  try {
    const sysRes = await httpRequest.systemInfo()
    systemInfo.value = sysRes.data
    const instRes = await httpRequest.getServers()
    instances.value = instRes.data
  } catch {
    createMessage('error', '加载数据失败')
  }
}

onMounted(loadData)
</script>
