import type { AxiosRequestConfig, AxiosResponse } from 'axios'
import { authService } from '@/services/authService'
import axios from 'axios'

export const BASE_URL: string = 'http://localhost:23754/api'

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

  async setup(key: string) {
    return this._makeRequest('/setup', 'post', {
      key,
    })
  }

  async login(key: string) {
    return this._makeRequest('/login', 'post', {
      key,
    })
  }
}
