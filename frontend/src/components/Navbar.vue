<template>
  <n-layout-header bordered style="height: 48px; padding: 0 20px; display: flex; align-items: center; justify-content: space-between;">
    <div style="font-size: 16px; font-weight: bold;">
      LayerProxy
    </div>
    <n-space align="center">
      <n-button text tag="a" href="https://github.com/SeaMite43981045/LayerProxy" target="_blank">
        {{ t('navbar.docs') }}
      </n-button>
      <n-button text @click="handleCheckUpdate">
        {{ t('navbar.checkUpdate') }}
      </n-button>
      <n-button text @click="toggleTheme">
        <n-icon>
          <component :is="isDark ? Moon : Sun" />
        </n-icon>
      </n-button>
      <n-tag size="small">{{ t('navbar.user') }}</n-tag>
    </n-space>
  </n-layout-header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { NLayoutHeader, NSpace, NButton, NIcon, NTag } from 'naive-ui'
import { Moon, Sun } from 'lucide-vue-next'
import HttpRequest from '@/http/httpRequest'
import { createMessage } from '@/message/showMessage'

const httpRequest = new HttpRequest()
const { t } = useI18n()
const isDark = ref(true)

function toggleTheme() {
  isDark.value = !isDark.value
}

async function handleCheckUpdate() {
  try {
    const res = await httpRequest.checkUpdate()
    if (res.data.has_update) {
      createMessage('info', '发现新版本', res.data.latest_version)
    } else {
      createMessage('success', '当前已是最新版本')
    }
  } catch {
    createMessage('error', '检查更新失败')
  }
}
</script>
