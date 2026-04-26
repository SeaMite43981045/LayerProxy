<script setup lang="ts">
import { darkTheme, NConfigProvider, NMessageProvider } from 'naive-ui'
import NaiveProvider from '@/components/NaiveProvider.vue'
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider :placement="'top-right'">
      <naive-provider>
        <router-view />
      </naive-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped></style>
