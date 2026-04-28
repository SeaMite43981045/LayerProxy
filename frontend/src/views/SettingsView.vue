<template>
  <n-card style="max-width: 800px; margin: 0 auto">
    <n-tabs type="line" animated>
      <n-tab-pane :name="t('settings.systemConfig')" :tab="t('settings.systemConfig')">
        <n-form label-placement="left" label-width="160px" style="margin-top: 16px">
          <n-form-item :label="t('settings.webPort')">
            <n-input v-model:value="config.web_port" />
            <template #feedback>{{ t('settings.webPortDesc') }}</template>
          </n-form-item>
          <n-form-item :label="t('settings.portStartAt')">
            <n-input-number v-model:value="config.port_start_at" />
            <template #feedback>{{ t('settings.portStartAtDesc') }}</template>
          </n-form-item>
          <n-form-item :label="t('settings.wildcardDomain')">
            <n-input v-model:value="config.wildcard_domain" />
            <template #feedback>{{ t('settings.wildcardDomainDesc') }}</template>
          </n-form-item>
          <n-form-item :label="t('settings.wildcardMainPort')">
            <n-input v-model:value="config.wildcard_main_port" />
            <template #feedback>{{ t('settings.wildcardMainPortDesc') }}</template>
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="saveConfig">{{ t('settings.save') }}</n-button>
          </n-form-item>
        </n-form>
      </n-tab-pane>
      <n-tab-pane :name="t('settings.preferences')" :tab="t('settings.preferences')">
        <n-form label-placement="left" label-width="160px" style="margin-top: 16px">
          <n-form-item :label="t('settings.language')">
            <n-select v-model:value="prefs.language" :options="languageOptions" />
          </n-form-item>
          <n-form-item :label="t('settings.theme')">
            <n-select v-model:value="prefs.theme" :options="themeOptions" />
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="savePrefs">{{ t('settings.save') }}</n-button>
          </n-form-item>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NCard,
  NTabs,
  NTabPane,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSelect,
  NButton,
} from 'naive-ui'
import HttpRequest from '@/http/httpRequest'
import { createMessage } from '@/message/showMessage'

const httpRequest = new HttpRequest()
const { t, locale } = useI18n()

const config = ref({
  web_port: '',
  port_start_at: 25565,
  wildcard_domain: '',
  wildcard_main_port: '',
})

const prefs = ref({
  language: 'zh',
  theme: 'dark',
})

const languageOptions = [
  { label: '中文', value: 'zh' },
  { label: 'English', value: 'en' },
]

const themeOptions = [
  { label: t('settings.dark'), value: 'dark' },
  { label: t('settings.light'), value: 'light' },
]

async function loadConfig() {
  try {
    const res = await httpRequest.getConfig()
    config.value = res.data
  } catch {}
}

async function saveConfig() {
  try {
    await httpRequest.updateConfig(config.value)
    createMessage('success', t('settings.saved'))
  } catch {
    createMessage('error', t('common.error'))
  }
}

async function loadPrefs() {
  try {
    const res = await httpRequest.getPreferences()
    prefs.value = res.data
  } catch {}
}

async function savePrefs() {
  try {
    await httpRequest.updatePreferences(prefs.value)
    localStorage.setItem('lp_language', prefs.value.language)
    locale.value = prefs.value.language
    createMessage('success', t('settings.saved'))
  } catch {
    createMessage('error', t('common.error'))
  }
}

onMounted(() => {
  loadConfig()
  loadPrefs()
})
</script>
