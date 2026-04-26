<script setup lang="ts">
import {
  NLayout,
  NFlex,
  NCard,
  NText,
  NTabPane,
  NTabs,
  NInput,
  NInputNumber,
  NButton,
  type FormValidationStatus,
} from 'naive-ui'
import {
  ArrowBigLeft,
  ArrowBigRightIcon,
  EyeClosedIcon,
  EyeIcon,
  TextSearchIcon,
  CheckIcon,
} from '@lucide/vue'
import { h, onMounted, ref, type Ref } from 'vue'
import { createMessage } from '@/message/showMessage'
import HttpRequest from '@/http/httpRequest'
import type { AxiosError } from 'axios'
import type { ResponseError } from '@/interfaces/response'
import { useRouter, type Router } from 'vue-router'

const httpRequest: HttpRequest = new HttpRequest()
const router: Router = useRouter()

const currentTab: Ref<string> = ref<string>('key-setup')

const strongPasswordRegex: RegExp = /^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[^\sa-zA-Z0-9]).{8,}$/

const docsLink: Ref<string> = ref<string>('https://proxy.laystudio.top/docs/setup#key')

const isLoading: Ref<boolean> = ref<boolean>(false)

const inputKeyValue: Ref<string | undefined> = ref<string | undefined>('')
const inputKeyStatus: Ref<FormValidationStatus | undefined> = ref<
  FormValidationStatus | undefined
>()

const portStartPort: Ref<number | undefined> = ref<number | undefined>(23755)

const wildPort: Ref<number | undefined> = ref<number | undefined>(23755)
const wildDomain: Ref<string | undefined> = ref<string | undefined>('')

const tabsBeforeLeave = () => {
  return false
}

const handleKeyNext = () => {
  if (inputKeyValue.value == '') {
    inputKeyStatus.value = 'error'
    createMessage('error', '密钥不能为空！')
    return
  }
  if (!strongPasswordRegex.test(inputKeyValue.value || '')) {
    inputKeyStatus.value = 'error'
    createMessage('error', '密钥强度太低！', '密码必须包含大小写字母、数字及特殊字符，且不少于8位')
    return
  }
  currentTab.value = 'config-setup'
  docsLink.value = 'https://proxy.laystudio.top/docs/setup#config'
}

const handleConfigBack = () => {
  currentTab.value = 'key-setup'
  docsLink.value = 'https://proxy.laystudio.top/docs/setup#key'
}

const handleSetupCommit = () => {
  isLoading.value = true

  if (inputKeyValue.value == undefined) {
    createMessage('error', '提交失败', '密钥的值为空！')
    isLoading.value = false
    handleConfigBack()
    return
  }

  if (portStartPort.value == undefined) {
    createMessage('error', '提交失败', 'Port 起始端口的值为空！')
    isLoading.value = false
    currentTab.value = 'config-setup'
    docsLink.value = 'https://proxy.laystudio.top/docs/setup#config'
    return
  }

  if (wildPort.value == undefined) {
    createMessage('error', '提交失败', 'Wildcard 端口的值为空！')
    isLoading.value = false
    currentTab.value = 'config-setup'
    docsLink.value = 'https://proxy.laystudio.top/docs/setup#config'
    return
  }

  httpRequest
    .setup(
      inputKeyValue.value,
      portStartPort.value,
      wildDomain.value || "",
      wildPort.value?.toString(),
    )
    .then((data) => {
      if (data.data.status == "ok") {
        createMessage("success", "Success", data.data.message)
        router.push({ name: "Login" })
      }
    })
    .catch((error: AxiosError<ResponseError>) => {
      createMessage("error", "Failed", error.response?.data.error)
    })
    .finally(() => {
      isLoading.value = false
    })
}

onMounted(() => {
  httpRequest.status().then((data) => {
    if (data.data.hasKey as boolean) {
      createMessage("error", "提示", "你已经完成了初始化")
      router.push({ name: "Login" })
    }
  })
})
</script>

<template>
  <n-layout style="height: 100vh">
    <n-flex :justify="'center'" :align="'center'" style="height: 100%">
      <n-card style="width: 70%; max-width: 640px">
        <template #header>
          <n-text style="font-size: 22px"> 初始化配置 </n-text>
        </template>
        <template #header-extra>
          <n-flex :align="'center'" style="gap: 16px">
            <n-button
              :tag="'a'"
              :href="docsLink"
              target="_blank"
              :render-icon="() => h(TextSearchIcon)"
              text
              >查看文档</n-button
            >
            <n-text style="font-weight: bold; font-style: italic; font-size: 16px; color: #aaaaaa">
              Layer-Proxy
            </n-text>
          </n-flex>
        </template>
        <template #default>
          <n-flex :justify="'center'" style="padding: 8px 0">
            <n-tabs
              :type="'line'"
              v-model:value="currentTab"
              @before-leave="tabsBeforeLeave"
              style="width: 80%"
              animated
            >
              <n-tab-pane name="key-setup" tab="密钥设置" :disabled="isLoading">
                <n-flex vertical>
                  <n-input
                    :type="'password'"
                    :status="inputKeyStatus"
                    :show-password-on="'mousedown'"
                    placeholder="请输入密钥"
                    v-model:value="inputKeyValue"
                    @input="inputKeyStatus = undefined"
                  >
                    <template #password-invisible-icon>
                      <eye-closed-icon :size="16" />
                    </template>
                    <template #password-visible-icon>
                      <eye-icon :size="16" />
                    </template>
                  </n-input>
                  <n-flex :justify="'end'">
                    <n-button :render-icon="() => h(ArrowBigRightIcon)" @click="handleKeyNext">
                        下一步
                    </n-button>
                  </n-flex>
                </n-flex>
              </n-tab-pane>
              <n-tab-pane name="config-setup" tab="配置文件设置" :disabled="isLoading">
                <n-flex vertical>
                  <n-card title="Port 模式配置" style="border: solid #999999">
                    <n-flex :justify="'space-between'" :align="'center'">
                      <n-text> 起始端口 (0~65535) </n-text>
                      <n-input-number
                        :button-placement="'both'"
                        :disabled="isLoading"
                        placeholder="请输入端口"
                        :min="0"
                        :max="65535"
                        v-model:value="portStartPort"
                      />
                    </n-flex>
                  </n-card>
                  <n-card title="Wildcard 模式配置" style="border: solid #999999">
                    <n-flex :justify="'space-between'" :align="'center'">
                      <n-text> 端口 (0~65535) </n-text>
                      <n-input-number
                        :button-placement="'both'"
                        :disabled="isLoading"
                        placeholder="请输入端口"
                        :min="0"
                        :max="65535"
                        v-model:value="wildPort"
                      />
                    </n-flex>
                    <n-flex
                      :justify="'space-between'"
                      :align="'center'"
                      style="gap: 32px; padding-top: 8px"
                    >
                      <n-text style="flex: 3"> 匹配域名 </n-text>
                      <n-input
                        style="flex: 17"
                        clearable
                        :disabled="isLoading"
                        placeholder="请输入域名"
                        v-model:value="wildDomain"
                      />
                    </n-flex>
                  </n-card>
                  <n-flex :justify="'space-between'" style="padding: 16px 0 0 0">
                    <n-button
                      :disabled="isLoading"
                      :render-icon="() => h(ArrowBigLeft)"
                      @click="handleConfigBack"
                      >上一步</n-button
                    >
                    <n-button
                      :loading="isLoading"
                      :render-icon="() => h(CheckIcon)"
                      @click="handleSetupCommit"
                      >提交</n-button
                    >
                  </n-flex>
                </n-flex>
              </n-tab-pane>
            </n-tabs>
          </n-flex>
        </template>
      </n-card>
    </n-flex>
  </n-layout>
</template>

<style scoped></style>
