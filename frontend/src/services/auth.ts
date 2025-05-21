import { authApi } from '@/api/client'

export const useAuthApi = () => {
  return {
    login: (data: { email: string; password: string }) => 
      authApi.post('/login', data),
    register: (data: { email: string; password: string; full_name: string }) => 
      authApi.post('/register', data),
    logout: () => 
      authApi.post('/logout'),
    getCurrentUser: () => 
      authApi.get('/me'),
    updateProfile: (data: { full_name: string }) => 
      authApi.put('/profile', data),
    changePassword: (data: { current_password: string; new_password: string }) =>
      authApi.put('/password', data)
  }
} 