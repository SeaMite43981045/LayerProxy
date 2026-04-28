<template>
  <div style="display: flex; gap: 16px; height: 100%">
    <n-card
      style="flex: 1; display: flex; flex-direction: column"
      content-style="flex: 1; display: flex; flex-direction: column; padding: 0;"
    >
      <template #header>
        <n-space align="center" justify="space-between" style="width: 100%">
          <span>{{ t('logging.liveStream') }}</span>
          <n-space>
            <n-tag :type="connected ? 'success' : 'error'" size="small">
              {{ connected ? t('logging.connected') : t('logging.disconnected') }}
            </n-tag>
            <n-button size="small" @click="clearLogs">{{ t('logging.clear') }}</n-button>
          </n-space>
        </n-space>
      </template>
      <div
        ref="logContainer"
        style="
          flex: 1;
          background: var(--n-color);
          font-family: monospace;
          padding: 12px;
          overflow-y: auto;
          white-space: pre-wrap;
          font-size: 13px;
        "
      >
        <div v-if="logs.length === 0" style="color: var(--n-text-color-3)">
          {{ t('logging.noLogs') }}
        </div>
        <div v-for="(log, index) in logs" :key="index">
          <span style="color: var(--n-text-color-3)">[{{ log.time }}]</span>
          <n-tag :type="logLevelColor(log.level)" size="tiny" style="margin: 0 4px">{{
            log.level
          }}</n-tag>
          <span>{{ log.message }}</span>
        </div>
      </div>
    </n-card>
    <n-card style="width: 240px" content-style="padding: 0;">
      <template #header>{{ t('logging.logFiles') }}</template>
      <n-list hoverable clickable>
        <n-list-item v-for="file in logFiles" :key="file.name">
          <n-thing :title="file.name" :description="formatBytes(file.size)" />
          <template #suffix>
            <n-space>
              <n-button text @click="downloadFile(file.name)">
                <n-icon><Download /></n-icon>
              </n-button>
              <n-button text type="error" @click="deleteFile(file.name)">
                <n-icon><Trash2 /></n-icon>
              </n-button>
            </n-space>
          </template>
        </n-list-item>
      </n-list>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { NCard, NSpace, NTag, NButton, NList, NListItem, NThing, NIcon } from 'naive-ui'
import { Download, Trash2 } from 'lucide-vue-next'
import HttpRequest from '@/http/httpRequest'
import { createMessage } from '@/message/showMessage'

const httpRequest = new HttpRequest()

const { t } = useI18n()

const logs = ref<{ time: string; level: string; message: string }[]>([])
const logFiles = ref<{ name: string; size: number }[]>([])
const connected = ref(false)
const logContainer = ref<HTMLDivElement>()
let eventSource: EventSource | null = null

function logLevelColor(level: string) {
  switch (level) {
    case 'ERROR':
      return 'error'
    case 'WARN':
      return 'warning'
    default:
      return 'info'
  }
}

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function clearLogs() {
  logs.value = []
}

function scrollToBottom() {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

let reconnectTimer: ReturnType<typeof setTimeout> | null = null

function connectSSE() {
  const token = localStorage.getItem('lp_token')
  if (!token) {
    connected.value = false
    return
  }

  if (eventSource) {
    eventSource.close()
  }

  eventSource = new EventSource(`/api/v1/logs/stream?token=${token}`)

  eventSource.addEventListener('connected', () => {
    connected.value = true
  })

  eventSource.addEventListener('log', (e) => {
    try {
      const data = JSON.parse(e.data)
      logs.value.push(data)
      if (logs.value.length > 500) {
        logs.value.shift()
      }
      scrollToBottom()
    } catch {}
  })

  eventSource.onopen = () => {
    connected.value = true
  }

  eventSource.onerror = () => {
    connected.value = false
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
    // 避免无限快速重连，延迟 3 秒后重试
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
    }
    reconnectTimer = setTimeout(() => {
      if (localStorage.getItem('lp_token')) {
        connectSSE()
      }
    }, 3000)
  }
}

async function loadLogFiles() {
  try {
    const res = await httpRequest.listLogFiles()
    logFiles.value = res.data
  } catch {}
}

async function downloadFile(name: string) {
  try {
    const res = await httpRequest.downloadLogFile(name)
    const blob = new Blob([res.data])
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = name
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    createMessage('error', '下载失败')
  }
}

async function deleteFile(name: string) {
  try {
    await httpRequest.deleteLogFile(name)
    createMessage('success', '删除成功')
    loadLogFiles()
  } catch {
    createMessage('error', '删除失败')
  }
}

onMounted(() => {
  connectSSE()
  loadLogFiles()
})

onUnmounted(() => {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
  }
  eventSource?.close()
})
</script>
