import MessageComponent from '@/components/MessageComponent.vue'
import type { MessageType } from 'naive-ui'
import type { MessageApiInjection } from 'naive-ui/es/message/src/MessageProvider'
import { h } from 'vue'

export let globalMessage: MessageApiInjection

export const registerGlobalMessage = (instance: MessageApiInjection) => {
  globalMessage = instance
}

export const createMessage = (type: MessageType, title: string, content: string = '') => {
  if (!globalMessage) {
    console.warn('Message 实例尚未注册，请确保 App.vue 已正确配置')
    return
  }

  globalMessage.create('', {
    keepAliveOnHover: true,
    type: type,
    render: () =>
      h(MessageComponent, {
        type: type,
        title: title,
        content: content,
      }),
  })
}
