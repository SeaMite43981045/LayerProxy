<script setup lang="ts">
import MessageComponent from '@/component/MessageComponent.vue'
import { KeyIcon, LogInIcon } from '@lucide/vue'
import { NLayout, NFlex, NCard, NInput, NText, NButton, useMessage } from 'naive-ui'
import type { MessageApiInjection } from 'naive-ui/es/message/src/MessageProvider'
import { h, ref, type Ref } from 'vue'
import HttpRequest from '@/http/httpRequest'
import { type UnauthorizedResponse } from '@/services/authService'
import { AxiosError } from 'axios'

const message: MessageApiInjection = useMessage()
const httpRequest: HttpRequest = new HttpRequest()

const isLogin: Ref<boolean> = ref<boolean>(false)
const inputValue: Ref<string> = ref<string>('')

const handleLogin = async () => {
  isLogin.value = true
  if (inputValue.value == '') {
    message.create("", {
      keepAliveOnHover: true,
      render: () =>
        h(MessageComponent, {
          type: 'error',
          title: '登录失败',
          content: '密钥不能为空！',
        }),
    })
    isLogin.value = false
    return
  }

  httpRequest.login(inputValue.value)
  .then((data) => {
    message.create("", {
      keepAliveOnHover: true,
      render: () => h(MessageComponent, {
        type: 'success',
        title: "登陆成功",
        content: data?.data?.message || "未知错误"
      })
    })
  })
  .catch((error: AxiosError<UnauthorizedResponse>) => {
    message.create("", {
      keepAliveOnHover: true,
      render: () => h(MessageComponent, {
        type: 'error',
        title: "登录失败",
        content: error.response?.data?.error || `未知错误: ${error.message}`
      })
    })
  })
  .finally(() => {
    isLogin.value = false
  })
}
</script>

<template>
  <n-layout style="height: 100vh">
    <n-flex :justify="'center'" :align="'center'" style="height: 100%">
      <n-card style="width: 70%; max-width: 640px">
        <template #header>
          <n-flex :justify="'start'" :align="'center'">
            <KeyIcon />
            <n-text>登录</n-text>
          </n-flex>
        </template>
        <template #header-extra>
          <n-text style="font-weight: bold; font-style: italic; font-size: 16px; color: #aaaaaa">
            Layer-Proxy
          </n-text>
        </template>
        <template #default>
          <n-flex>
            <n-input
              :type="'password'"
              placeholder="请输入访问密钥"
              :disabled="isLogin"
              style="flex: 9"
              v-model:value="inputValue"
            />
            <n-button
              :loading="isLogin"
              style="flex: 1"
              :render-icon="() => h(LogInIcon)"
              @click="handleLogin"
            >
              登录
            </n-button>
          </n-flex>
        </template>
      </n-card>
    </n-flex>
  </n-layout>
</template>

<style scoped></style>
