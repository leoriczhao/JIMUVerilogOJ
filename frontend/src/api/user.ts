import api from './index'
import type { User } from '../stores/user'

// 注册请求接口
export interface RegisterRequest {
  username: string
  email: string
  password: string
  nickname?: string
  school?: string
  student_id?: string
}

// 登录请求接口
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应接口
export interface LoginResponse {
  token: string
  user: User
  expires_in: number
}

// 更新个人信息请求接口
export interface UpdateProfileRequest {
  nickname?: string
  avatar?: string
  school?: string
  student_id?: string
}

// 用户注册
export const register = (data: RegisterRequest): Promise<User> => {
  return api.post('/users/register', data)
}

// 用户登录
export const login = (data: LoginRequest): Promise<LoginResponse> => {
  return api.post('/users/login', data)
}

// 获取个人信息
export const getProfile = (): Promise<User> => {
  return api.get('/users/profile')
}

// 更新个人信息
export const updateProfile = (data: UpdateProfileRequest): Promise<User> => {
  return api.put('/users/profile', data)
}