import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import RoomList from '../views/RoomList.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/login', component: Login },
    { 
      path: '/rooms', 
      component: RoomList,
      beforeEnter: (to, from, next) => {
        if (!localStorage.getItem('token')) next('/login')
        else next()
      }
    },
  ],
})

export default router
