<template>
  <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="180" show-trigger>
    <n-menu
      :value="activeKey"
      :collapsed-width="64"
      :collapsed-icon-size="22"
      :options="menuOptions"
      @update:value="handleMenuSelect"
    />
  </n-layout-sider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NLayoutSider, NMenu } from 'naive-ui'
import { LayoutDashboard, ScrollText, Settings } from 'lucide-vue-next'
import { h } from 'vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const activeKey = computed(() => route.name as string)

const menuOptions = [
  {
    label: t('sidebar.dashboard'),
    key: 'Dashboard',
    icon: () => h(LayoutDashboard),
  },
  {
    label: t('sidebar.logging'),
    key: 'Logging',
    icon: () => h(ScrollText),
  },
  {
    label: t('sidebar.settings'),
    key: 'Settings',
    icon: () => h(Settings),
  },
]

function handleMenuSelect(key: string) {
  router.push({ name: key })
}
</script>
