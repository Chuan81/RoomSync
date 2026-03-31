<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2>RoomSync 会议室预约</h2>
      <el-form :model="loginForm" label-width="0px">
        <el-form-item>
          <el-input v-model="loginForm.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="loginForm.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-button type="primary" style="width: 100%" @click="handleLogin">登录</el-button>
        <div class="footer-links">
          还没有账号？<el-button link type="primary" @click="$router.push('/register')">立即注册</el-button>
        </div>
        <div style="margin-top: 10px; text-align: center; color: #999; font-size: 12px;">
          首个注册用户将自动获得管理员权限
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/api'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loginForm = ref({
  username: '',
  password: '',
})

const handleLogin = async () => {
  try {
    const res = await request.post('/users/login', loginForm.value)
    userStore.setLogin(res.data.user, res.data.token)
    ElMessage.success('登录成功')
    router.push('/rooms')
  } catch (error) {}
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}
.login-card {
  width: 400px;
  text-align: center;
}
.footer-links {
  margin-top: 15px;
  font-size: 14px;
  color: #666;
}
</style>
