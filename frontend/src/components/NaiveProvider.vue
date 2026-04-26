<script setup lang="ts">
import { useMessage } from 'naive-ui'
import { registerGlobalMessage } from '@/message/showMessage'

const message = useMessage()

registerGlobalMessage(message)
</script>

<template>
  <slot />
</template>
