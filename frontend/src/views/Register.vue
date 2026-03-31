<template>
  <div class="register-container">
    <el-card class="register-card">
      <h2>加入 RoomSync</h2>
      <el-form :model="registerForm" :rules="rules" ref="registerRef" label-width="0px">
        <el-form-item prop="username">
          <el-input v-model="registerForm.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="registerForm.email" placeholder="邮箱" prefix-icon="Message" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-button type="primary" style="width: 100%" @click="handleRegister" :loading="loading">立即注册</el-button>
        <div class="footer-links">
          已有账号？<el-button link type="primary" @click="$router.push('/login')">去登录</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const registerRef = ref()
const loading = ref(false)

const registerForm = ref({
  username: '',
  email: '',
  password: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, message: '密码至少6位', trigger: 'blur' }],
}

const handleRegister = async () => {
  if (!registerRef.value) return
  await registerRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await request.post('/users/register', registerForm.value)
        ElMessage.success('注册成功，请登录')
        router.push('/login')
      } catch (error) {
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.register-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}
.register-card {
  width: 400px;
  text-align: center;
}
.footer-links {
  margin-top: 15px;
  font-size: 14px;
  color: #666;
}
</style>
