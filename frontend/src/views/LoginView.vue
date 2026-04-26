<script setup lang="ts">
import { KeyIcon, LogInIcon } from '@lucide/vue'
import { NLayout, NFlex, NCard, NInput, NText, NButton } from 'naive-ui'
import { h, onMounted, ref, type Ref } from 'vue'
import HttpRequest from '@/http/httpRequest'
import { authService } from '@/services/authService'
import { AxiosError } from 'axios'
import { useRouter, type Router } from 'vue-router'
import { createMessage } from '@/message/showMessage'
import type { ResponseError } from '@/interfaces/response'

const router: Router = useRouter()
const httpRequest: HttpRequest = new HttpRequest()

const isLogin: Ref<boolean> = ref<boolean>(false)
const inputValue: Ref<string> = ref<string>('')

const handleLogin = async () => {
  isLogin.value = true
  if (inputValue.value == '') {
    createMessage('error', '登录失败', '密钥不能为空！')
    isLogin.value = false
    return
  }

  httpRequest
    .login(inputValue.value)
    .then((data) => {
      authService.setToken(data.data.token)
      createMessage('success', '登录成功')
    })
    .catch((error: AxiosError<ResponseError>) => {
      if (error.response?.data.redirect) {
        createMessage('error', '登录失败', '请先对该系统进行初始化！')
        router.push({ name: 'Setup' })
      } else {
        createMessage(
          'error',
          '登录失败',
          error.response?.data?.error || `未知错误: ${error.message}`,
        )
      }
    })
    .finally(() => {
      isLogin.value = false
    })
}

onMounted(() => {
  httpRequest.status().then((data) => {
    if (!(data.data.hasKey as boolean)) {
      createMessage('info', '提示', '请先完成初始化')
      router.push({ name: 'Setup' })
    }
  })
})
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
