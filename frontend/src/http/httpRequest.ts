import type { AxiosRequestConfig, AxiosResponse } from 'axios'
import { authService } from '@/services/authService'
import axios from 'axios'
import router from '@/router'

export const BASE_URL: string = '/api'

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('lp_token')
      router.push({ name: 'Login' })
    }
    return Promise.reject(error)
  },
)

export default class HttpRequest {
  private async _makeRequest(
    endpoint: string,
    method: 'get' | 'post' | 'delete' | 'put' = 'get',
    data?: unknown,
    responseType?: 'blob',
  ): Promise<AxiosResponse> {
    const token: string = authService.getToken() || ''
    const url: string = BASE_URL + endpoint

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: '*/*',
      Authorization: token ? `Bearer ${token}` : '',
    }

    const requestConfig: AxiosRequestConfig = {
      method,
      url,
      headers,
      data: data ? JSON.stringify(data) : undefined,
      responseType,
    }

    return await axios(requestConfig)
  }

  async status() {
    return this._makeRequest('/status', 'get')
  }

  async setup(
    key: string,
    port_start_at: number,
    wildcard_domain: string,
    wildcard_main_port: string,
  ) {
    return this._makeRequest('/setup', 'post', {
      key,
      port_start_at,
      wildcard_domain,
      wildcard_main_port,
    })
  }

  async login(key: string) {
    return this._makeRequest('/login', 'post', {
      key,
    })
  }

  async systemInfo() {
    return this._makeRequest('/v1/system/info', 'get')
  }

  async getConfig() {
    return this._makeRequest('/v1/config', 'get')
  }

  async updateConfig(data: unknown) {
    return this._makeRequest('/v1/config', 'post', data)
  }

  async listLogFiles() {
    return this._makeRequest('/v1/logs/files', 'get')
  }

  async downloadLogFile(name: string) {
    return this._makeRequest(`/v1/logs/files/${name}`, 'get')
  }

  async deleteLogFile(name: string) {
    return this._makeRequest(`/v1/logs/files/${name}`, 'delete')
  }

  async checkUpdate() {
    return this._makeRequest('/v1/update', 'get')
  }

  async getPreferences() {
    return this._makeRequest('/v1/preferences', 'get')
  }

  async updatePreferences(data: unknown) {
    return this._makeRequest('/v1/preferences', 'post', data)
  }

  async getServers() {
    return this._makeRequest('/v1/servers', 'get')
  }
}
