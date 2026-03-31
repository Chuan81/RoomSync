<template>
  <div class="layout">
    <el-header class="header">
      <div class="logo">RoomSync</div>
      <div class="nav-links">
        <el-button link type="primary">会议室大盘</el-button>
        <el-button v-if="isAdmin" link @click="$router.push('/admin/bookings')">审批管理</el-button>
      </div>
      <div class="user-info">
        <span>{{ userStore.user?.username }} ({{ userStore.user?.role }})</span>
        <el-button link @click="handleLogout">退出</el-button>
      </div>
    </el-header>

    <el-main class="main">
      <div class="toolbar">
        <h3>会议室大盘</h3>
        <el-button v-if="isAdmin" type="primary" @click="showAddDialog = true">新增会议室</el-button>
      </div>

      <el-row :gutter="20">
        <el-col v-for="room in rooms" :key="room.id" :span="8">
          <el-card class="room-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span>{{ room.name }}</span>
                <el-tag :type="room.need_approval ? 'warning' : 'success'">
                  {{ room.need_approval ? '需审核' : '免审核' }}
                </el-tag>
              </div>
            </template>
            <p><strong>容量:</strong> {{ room.capacity }}人</p>
            <p><strong>地点:</strong> {{ room.location }}</p>
            <p><strong>设备:</strong> {{ room.equipment || '无' }}</p>
            <div class="actions">
              <el-button type="primary" @click="openBooking(room)">立即预约</el-button>
              <el-button v-if="isAdmin" type="danger" link @click="handleDelete(room.id)">删除</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-main>

    <!-- 预约对话框 -->
    <el-dialog v-model="showBookingDialog" title="发起预约" width="550px">
      <div v-if="busySlots.length > 0" class="busy-info">
        <p style="color: #f56c6c; font-size: 14px; margin-bottom: 10px;">
          <strong>注意：以下时间段已被占用（标红）：</strong>
        </p>
        <div class="busy-list">
          <el-tag v-for="slot in busySlots" :key="slot.id" type="danger" effect="dark" class="slot-tag">
            {{ formatTime(slot.start_time) }} 至 {{ formatTime(slot.end_time) }}
          </el-tag>
        </div>
      </div>
      <el-divider v-if="busySlots.length > 0" />
      
      <el-form :model="bookingForm" label-width="80px">
        <el-form-item label="会议名称">
          <el-input v-model="bookingForm.title" />
        </el-form-item>
        <el-form-item label="时间段">
          <el-date-picker
            v-model="bookingTime"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始"
            end-placeholder="结束"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBookingDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBooking">确认预约</el-button>
      </template>
    </el-dialog>

    <!-- 新增会议室对话框 -->
    <el-dialog v-model="showAddDialog" title="新增会议室" width="500px">
      <el-form :model="roomForm" label-width="100px">
        <el-form-item label="名称"><el-input v-model="roomForm.name" /></el-form-item>
        <el-form-item label="容量"><el-input-number v-model="roomForm.capacity" :min="1" /></el-form-item>
        <el-form-item label="地点"><el-input v-model="roomForm.location" /></el-form-item>
        <el-form-item label="设备"><el-input v-model="roomForm.equipment" /></el-form-item>
        <el-form-item label="需要审批"><el-switch v-model="roomForm.need_approval" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddRoom">创建</el-button>
      </template>
    </el-dialog>
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

const rooms = ref<any[]>([])
const busySlots = ref<any[]>([])
const showBookingDialog = ref(false)
const showAddDialog = ref(false)
const selectedRoom = ref<any>(null)

const bookingForm = ref({ title: '' })
const bookingTime = ref<[string, string]>(['', ''])

const roomForm = ref({
  name: '',
  capacity: 10,
  location: '',
  equipment: '',
  need_approval: false,
})

const fetchRooms = async () => {
  const res = await request.get('/rooms')
  rooms.value = res.data
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const openBooking = async (room: any) => {
  selectedRoom.value = room
  busySlots.value = [] // Reset
  try {
    const res = await request.get(`/rooms/${room.id}/bookings`)
    busySlots.value = res.data
  } catch (error) {}
  showBookingDialog.value = true
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN', {
    month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit'
  })
}

const handleBooking = async () => {
  if (!bookingTime.value || bookingTime.value.length < 2) {
    return ElMessage.warning('请选择时间段')
  }
  try {
    await request.post('/bookings', {
      room_id: selectedRoom.value.id,
      title: bookingForm.value.title,
      start_time: bookingTime.value[0],
      end_time: bookingTime.value[1],
    })
    ElMessage.success('预约成功')
    showBookingDialog.value = false
  } catch (error) {}
}

const handleAddRoom = async () => {
  await request.post('/rooms', roomForm.value)
  ElMessage.success('创建成功')
  showAddDialog.value = false
  fetchRooms()
}

const handleDelete = (id: number) => {
  ElMessageBox.confirm('确定删除该会议室吗？', '警告', { type: 'warning' }).then(async () => {
    await request.delete(`/rooms/${id}`)
    ElMessage.success('已删除')
    fetchRooms()
  })
}

onMounted(fetchRooms)
</script>

<style scoped>
.layout { min-height: 100vh; background-color: #f0f2f5; }
.header { 
  display: flex; justify-content: space-between; align-items: center; 
  background: #001529; color: white; padding: 0 20px; 
}
.logo { font-size: 20px; font-weight: bold; }
.nav-links { flex: 1; margin-left: 50px; }
.user-info { font-size: 14px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.room-card { margin-bottom: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.actions { margin-top: 20px; display: flex; justify-content: space-between; }
.busy-info { margin-bottom: 20px; }
.busy-list { display: flex; flex-wrap: wrap; gap: 8px; }
.slot-tag { margin-bottom: 5px; }
</style>
