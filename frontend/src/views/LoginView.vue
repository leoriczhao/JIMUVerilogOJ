<template>
  <div class="login">
    <!-- 主要内容区域 -->
    <div class="main-content">
      <div class="login-card">
          <div class="login-header">
            <h2>用户登录</h2>
            <p>欢迎回到 Verilog OJ</p>
          </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="80px"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名或邮箱"
            :prefix-icon="User"
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            @click="handleLogin"
            class="login-button"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
        
        <div class="register-link">
          <span>还没有账户？</span>
          <router-link to="/register" class="link">立即注册</router-link>
        </div>
      </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { login, type LoginRequest } from '../api/user'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

// 表单引用
const loginFormRef = ref<FormInstance>()

// 加载状态
const loading = ref(false)

// 表单数据
const loginForm = reactive<LoginRequest>({
  username: '',
  password: ''
})

// 表单验证规则
const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名或邮箱', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    // 表单验证
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 调用登录API
    const response = await login(loginForm)
    
    // 保存用户信息到store
    userStore.setUser(response.user, response.token)
    
    ElMessage.success('登录成功！')
    
    // 跳转到首页
    router.push('/')
    
  } catch (error: any) {
    console.error('登录失败:', error)
    
    // 处理错误信息
    if (error.response?.status === 401) {
      ElMessage.error('用户名或密码错误')
    } else if (error.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    } else if (error.message) {
      ElMessage.error(error.message)
    } else {
      ElMessage.error('登录失败，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login {
  min-height: calc(100vh - 120px); /* 减去Header和Footer的高度 */
}

.main-content {
  padding: 0 40px;
  background: #f8f9fa;
  min-height: calc(100vh - 120px);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  margin-top: 60px; /* 为固定导航栏留出空间 */
}

.login-card {
  max-width: 400px;
  margin: 0 auto;
  background: white;
  border-radius: 12px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  color: #333;
  margin-bottom: 8px;
  font-size: 28px;
  font-weight: 600;
}

.login-header p {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.login-form {
  margin-top: 20px;
}

.login-form :deep(.el-form-item__label) {
  color: #333;
  font-weight: 500;
}

.login-form :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6;
  transition: all 0.3s;
  width: 100%;
  min-width: 100%;
}

.login-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc;
}

.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff;
}

.login-form :deep(.el-input) {
  width: 100%;
}

.login-form :deep(.el-input__inner) {
  width: 100%;
  box-sizing: border-box;
}

.login-form :deep(.el-input__suffix) {
  position: absolute;
  right: 8px;
}

.login-button {
  width: 100%;
  height: 45px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.register-link {
  text-align: center;
  margin-top: 20px;
  color: #666;
  font-size: 14px;
}

.register-link .link {
  color: #409eff;
  text-decoration: none;
  font-weight: 500;
  margin-left: 5px;
}

.register-link .link:hover {
  text-decoration: underline;
}

/* 桌面端优化 */
@media (min-width: 1200px) {
  .main-content {
    padding: 60px 60px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .main-content {
    padding: 60px 20px;
  }
  
  .login-card {
    padding: 30px 20px;
  }
  
  .login-header h2 {
    font-size: 24px;
  }
}
</style>