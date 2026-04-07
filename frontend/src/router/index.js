import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import CreateQueueView from '../views/CreateQueueView.vue'
import QueueView from '../views/QueueView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',              component: HomeView },
    { path: '/create-queue', component: CreateQueueView },
    { path: '/queue',        component: QueueView },
  ]
})

export default router
