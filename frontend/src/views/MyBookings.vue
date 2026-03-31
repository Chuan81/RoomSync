<template>
  <div class="layout">
    <el-header class="header">
      <div class="logo">RoomSync</div>
      <div class="nav-links">
        <el-button link @click="$router.push('/rooms')">会议室大盘</el-button>
        <el-button link type="primary">我的预约</el-button>
        <el-button v-if="isAdmin" link @click="$router.push('/admin/bookings')">审批管理</el-button>
      </div>
      <div class="user-info">
        <span>{{ userStore.user?.username }}</span>
        <el-button link @click="handleLogout">退出</el-button>
      </div>
    </el-header>

    <el-main>
      <el-card>
        <template #header>我的预约请求</template>
        <el-table :data="bookings" stripe>
          <el-table-column prop="title" label="会议主题" />
          <el-table-column prop="room.name" label="会议室" />
          <el-table-column label="预约时间" width="300">
            <template #default="scope">
              {{ formatTime(scope.row.start_time) }} - {{ formatTime(scope.row.end_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
            <template #default="scope">
              <el-tag :type="statusType(scope.row.status)">{{ scope.row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button 
                v-if="['pending', 'approved'].includes(scope.row.status)"
                type="danger" 
                link 
                @click="handleCancel(scope.row.id)"
              >撤销预约</el-button>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import request from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const isAdmin = computed(() => userStore.user?.role === 'admin')
const bookings = ref([])

const fetchMyBookings = async () => {
  const res = await request.get('/bookings/my')
  bookings.value = res.data
}

const handleCancel = (id: number) => {
  ElMessageBox.confirm('确定要撤销这条预约吗？如果是已同意的预约，撤销后将释放会议室资源。', '二次确认', {
    confirmButtonText: '确定撤销',
    cancelButtonText: '点错了',
    type: 'warning'
  }).then(async () => {
    await request.put(`/bookings/${id}/cancel`)
    ElMessage.success('预约已撤销')
    fetchMyBookings()
  })
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

onMounted(fetchMyBookings)
</script>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; background: #001529; color: white; padding: 0 20px; }
.logo { font-size: 20px; font-weight: bold; }
.nav-links { flex: 1; margin-left: 50px; }
</style>
