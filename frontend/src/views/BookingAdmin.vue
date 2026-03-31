<template>
  <div class="layout">
    <el-header class="header">
      <div class="logo">RoomSync - 审批中心</div>
      <div class="nav-links">
        <el-button link @click="$router.push('/rooms')">会议室大盘</el-button>
        <el-button link type="primary">审批管理</el-button>
      </div>
      <div class="user-info">
        <span>{{ userStore.user?.username }} ({{ userStore.user?.role }})</span>
        <el-button link @click="handleLogout">退出</el-button>
      </div>
    </el-header>

    <el-main>
      <el-card>
        <template #header>预约审批列表</template>
        <el-table :data="bookings" stripe style="width: 100%">
          <el-table-column prop="title" label="会议主题" />
          <el-table-column prop="room.name" label="会议室" />
          <el-table-column prop="user.username" label="申请人" />
          <el-table-column label="时间段" width="320">
            <template #default="scope">
              {{ formatTime(scope.row.start_time) }} - {{ formatTime(scope.row.end_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
            <template #default="scope">
              <el-tag :type="statusType(scope.row.status)">{{ scope.row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="scope">
              <div v-if="scope.row.status === 'pending'">
                <el-button size="small" type="success" @click="handleApprove(scope.row.id, 'approved')">通过</el-button>
                <el-button size="small" type="danger" @click="handleApprove(scope.row.id, 'rejected')">拒绝</el-button>
              </div>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import request from '@/api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const bookings = ref([])

const fetchBookings = async () => {
  const res = await request.get('/bookings')
  bookings.value = res.data
}

const handleApprove = async (id: number, status: string) => {
  await request.put(`/bookings/${id}/approve`, { status })
  ElMessage.success('操作成功')
  fetchBookings()
}

const formatTime = (time: string) => new Date(time).toLocaleString()

const statusType = (status: string) => {
  const map: any = { pending: 'warning', approved: 'success', rejected: 'danger', cancelled: 'info' }
  return map[status] || 'info'
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(() => {
  if (userStore.user?.role !== 'admin') {
    ElMessage.error('无权访问')
    router.push('/rooms')
    return
  }
  fetchBookings()
})
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; background: #001529; color: white; }
.logo { font-size: 18px; font-weight: bold; }
.nav-links { flex: 1; margin-left: 50px; }
.user-info { font-size: 14px; }
</style>
