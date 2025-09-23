import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  school: string
  student_id: string
  role: 'student' | 'teacher' | 'admin'
  solved: number
  submitted: number
  rating: number
  is_active: boolean
  created_at: string
  last_login_at: string
}

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isLoggedIn = ref(false)

  // 初始化用户状态
  const initUser = () => {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    
    if (storedToken && storedUser) {
      token.value = storedToken
      user.value = JSON.parse(storedUser)
      isLoggedIn.value = true
    }
  }

  // 设置用户信息
  const setUser = (userData: User, userToken: string) => {
    user.value = userData
    token.value = userToken
    isLoggedIn.value = true
    
    localStorage.setItem('token', userToken)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  // 清除用户信息
  const clearUser = () => {
    user.value = null
    token.value = null
    isLoggedIn.value = false
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 初始化
  initUser()

  return {
    user,
    token,
    isLoggedIn,
    setUser,
    clearUser,
    initUser
  }
}) 