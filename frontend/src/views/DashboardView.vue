<template>
  <div style="display: flex; flex-direction: column; gap: 16px; height: 100%">
    <div style="display: flex; gap: 16px">
      <n-card style="flex: 1">
        <n-statistic :label="t('dashboard.cpu')" :value="systemInfo.cpu_model">
          <template #prefix>
            <n-icon><Cpu /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px">
          <n-text depth="3"
            >{{ systemInfo.cpu_cores }} {{ t('dashboard.cores') }} / {{ systemInfo.cpu_threads }}
            {{ t('dashboard.threads') }}</n-text
          >
          <n-progress
            type="line"
            :percentage="Math.round(systemInfo.cpu_usage)"
            indicator-placement="inside"
          />
        </n-space>
      </n-card>
      <n-card style="flex: 1">
        <n-statistic :label="t('dashboard.memory')" :value="formatBytes(systemInfo.memory_total)">
          <template #prefix>
            <n-icon><MemoryStick /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px">
          <n-text depth="3"
            >{{ t('dashboard.used') }}: {{ formatBytes(systemInfo.memory_used) }} /
            {{ t('dashboard.free') }}: {{ formatBytes(systemInfo.memory_free) }}</n-text
          >
          <n-progress type="line" :percentage="memoryUsagePercent" indicator-placement="inside" />
        </n-space>
      </n-card>
      <n-card style="flex: 1">
        <n-statistic :label="t('dashboard.system')" :value="systemInfo.os_name">
          <template #prefix>
            <n-icon><Monitor /></n-icon>
          </template>
        </n-statistic>
        <n-space vertical style="margin-top: 8px">
          <n-text depth="3">{{ t('dashboard.version') }}: {{ systemInfo.os_version }}</n-text>
          <n-text depth="3"
            >{{ t('dashboard.uptime') }}: {{ formatUptime(systemInfo.uptime) }}</n-text
          >
        </n-space>
      </n-card>
    </div>
    <n-card style="flex: 1; overflow: auto">
      <template #header>
        <n-space align="center">
          <span>{{ t('dashboard.instances') }}</span>
          <n-button size="small" @click="openAddModal">{{ t('dashboard.addInstance') }}</n-button>
        </n-space>
      </template>
      <n-data-table :columns="columns" :data="instances" :bordered="false" size="small" />
    </n-card>
  </div>

  <n-modal
    v-model:show="showModal"
    :title="isEditing ? t('dashboard.editInstance') : t('dashboard.addInstance')"
    preset="dialog"
    :show-icon="false"
    @positive-click="handleSave"
    @negative-click="
      () => {
        showModal = false
      }
    "
    :positive-text="t('common.confirm')"
    :negative-text="t('common.cancel')"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="auto"
    >
      <n-form-item :label="t('dashboard.name')" path="name">
        <n-input v-model:value="formData.name" :disabled="isEditing" />
      </n-form-item>
      <n-form-item :label="t('dashboard.backendIP')" path="backend_ip">
        <n-input v-model:value="formData.backend_ip" placeholder="127.0.0.1:25565" />
      </n-form-item>
      <n-form-item :label="t('dashboard.subdomain')" path="subdomain">
        <n-input v-model:value="formData.subdomain" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import {
  NCard,
  NStatistic,
  NIcon,
  NProgress,
  NSpace,
  NText,
  NDataTable,
  NButton,
  NModal,
  NForm,
  NFormItem,
  NInput,
} from 'naive-ui'
import { Cpu, MemoryStick, Monitor, Pencil, Trash2 } from 'lucide-vue-next'
import HttpRequest from '@/http/httpRequest'
import { createMessage } from '@/message/showMessage'

const httpRequest = new HttpRequest()
const { t } = useI18n()

const systemInfo = ref({
  cpu_model: '',
  cpu_cores: 0,
  cpu_threads: 0,
  cpu_usage: 0,
  memory_total: 0,
  memory_used: 0,
  memory_free: 0,
  os_name: '',
  os_version: '',
  uptime: 0,
})

interface Instance {
  name: string
  backend_ip: string
  subdomain: string
}

const instances = ref<Instance[]>([])
const showModal = ref(false)
const isEditing = ref(false)
const formRef = ref<FormInst | null>(null)
const formData = ref({ name: '', backend_ip: '', subdomain: '' })

const formRules: FormRules = {
  name: [{ required: true, message: '名称不能为空', trigger: 'blur' }],
  backend_ip: [{ required: true, message: '后端 IP 不能为空', trigger: 'blur' }],
  subdomain: [{ required: true, message: '子域名不能为空', trigger: 'blur' }],
}

const memoryUsagePercent = computed(() => {
  if (!systemInfo.value.memory_total) return 0
  return Math.round((systemInfo.value.memory_used / systemInfo.value.memory_total) * 100)
})

const columns = computed<DataTableColumns<Instance>>(() => [
  { title: t('dashboard.name'), key: 'name' },
  { title: t('dashboard.backendIP'), key: 'backend_ip' },
  { title: t('dashboard.subdomain'), key: 'subdomain' },
  {
    title: t('dashboard.status'),
    key: 'status',
    render() {
      return h('span', { style: 'color: var(--n-success-color);' }, '运行中')
    },
  },
  {
    title: t('dashboard.actions'),
    key: 'actions',
    render(row) {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              text: true,
              onClick: () => openEditModal(row),
            },
            {
              icon: () => h(NIcon, null, { default: () => h(Pencil) }),
            },
          ),
          h(
            NButton,
            {
              size: 'small',
              text: true,
              type: 'error',
              onClick: () => handleDelete(row),
            },
            {
              icon: () => h(NIcon, null, { default: () => h(Trash2) }),
            },
          ),
        ],
      })
    },
  },
])

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

function openAddModal() {
  isEditing.value = false
  formData.value = { name: '', backend_ip: '', subdomain: '' }
  showModal.value = true
}

function openEditModal(row: Instance) {
  isEditing.value = true
  formData.value = { name: row.name, backend_ip: row.backend_ip, subdomain: row.subdomain }
  showModal.value = true
}

async function handleSave() {
  try {
    await formRef.value?.validate()
  } catch {
    return false
  }

  try {
    await httpRequest.addServer(formData.value)
    createMessage('success', isEditing.value ? '修改成功' : '添加成功')
    showModal.value = false
    loadData()
  } catch {
    createMessage('error', isEditing.value ? '修改失败' : '添加失败')
    return false
  }
}

async function handleDelete(row: Instance) {
  if (!confirm(t('common.deleteConfirm'))) return
  try {
    await httpRequest.deleteServer(row.name)
    createMessage('success', '删除成功')
    loadData()
  } catch {
    createMessage('error', '删除失败')
  }
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
