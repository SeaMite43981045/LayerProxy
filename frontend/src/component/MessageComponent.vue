<script setup lang="ts">
import { CircleAlertIcon, CircleXIcon, InfoIcon, CircleCheckIcon, type LucideIcon } from '@lucide/vue'
import { NCard, NFlex, NText, NIcon } from 'naive-ui'
import { h, ref, type Ref } from 'vue'

const props = defineProps<{
  type: 'info' | 'warn' | 'error' | 'success'
  title: string
  content: string
}>()

const cardClassd: Ref<string> = ref<string>(props.type + '-card')

const cardIcon: Ref<LucideIcon> = ref<LucideIcon>(InfoIcon)
if (props.type === 'info') {
  cardIcon.value = InfoIcon
} else if (props.type === 'warn') {
  cardIcon.value = CircleAlertIcon
} else if (props.type === 'error') {
  cardIcon.value = CircleXIcon
} else if (props.type === 'success') {
  cardIcon.value = CircleCheckIcon
}
</script>

<template>
  <n-card hoverable :class="cardClassd" :size="'small'">
    <template #header>
      <n-flex :justify="'start'" :align="'center'">
        <n-icon :component="() => h(cardIcon, { size: 18 })" :size="18" />
        <n-text style="font-size: 18px">{{ props.title }}</n-text>
      </n-flex>
    </template>
    <template #default>
      <n-flex :justify="'start'" :align="'center'">
        <n-text>{{ props.content }}</n-text>
      </n-flex>
    </template>
  </n-card>
</template>

<style scoped>
.info-card {
  background-color: #379dfd36;
  min-width: 320px;
}
.warn-card {
  background-color: #ffd90036;
  min-width: 320px;
}
.error-card {
  background-color: #fd375136;
  min-width: 320px;
}
.success-card {
  background-color: #37fd6236;
  min-width: 320px;
}
</style>
