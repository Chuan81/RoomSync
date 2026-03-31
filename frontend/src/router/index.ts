import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import RoomList from '../views/RoomList.vue'
import BookingAdmin from '../views/BookingAdmin.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/login', component: Login },
    { path: '/register', component: Register },
    { 
      path: '/rooms', 
      component: RoomList,
      meta: { requiresAuth: true }
    },
    { 
      path: '/admin/bookings', 
      component: BookingAdmin,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
  ],
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const user = JSON.parse(localStorage.getItem('user') || '{}')

  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.meta.requiresAdmin && user.role !== 'admin') {
    next('/rooms')
  } else {
    next()
  }
})

export default router
