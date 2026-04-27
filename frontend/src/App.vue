<script setup lang="ts">
import { ref, provide } from 'vue'
import { darkTheme, lightTheme, NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import NaiveProvider from '@/components/NaiveProvider.vue'

const isDark = ref(true)

function toggleTheme() {
  isDark.value = !isDark.value
}

provide('isDark', isDark)
provide('toggleTheme', toggleTheme)
</script>

<template>
  <n-config-provider :theme="isDark ? darkTheme : lightTheme">
    <n-message-provider>
      <n-dialog-provider>
        <NaiveProvider>
          <router-view />
        </NaiveProvider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>
