<script setup lang="ts">
import { darkTheme, NConfigProvider, NMessageProvider } from 'naive-ui'
</script>

<template>
  <n-config-provider :theme="darkTheme">
    <n-message-provider :placement="'top-right'">
      <router-view />
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped></style>
