import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import CreateQueueView from '../views/CreateQueueView.vue'
import QueueView from '../views/QueueView.vue'

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/create-queue', name: 'create-queue', component: CreateQueueView },
  { path: '/queue/:id', name: 'queue', component: QueueView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
