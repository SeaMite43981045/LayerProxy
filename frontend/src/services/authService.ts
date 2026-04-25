const BASE_URL: string = 'http://localhost:23754/api'

export const authService = {
  setToken(token: string) {
    localStorage.setItem('lp_token', token)
  },

  getToken(): string | null {
    return localStorage.getItem('lp_token')
  }
}

export interface UnauthorizedResponse {
    error: string
}
