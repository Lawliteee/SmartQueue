<template>
  <div class="page-body">
    <section class="queue-info">
      <h2 class="queue-title">{{ queueTitle }}</h2>
    </section>
    <section class="content">
      <p class="wait-time">Осталось ждать: {{ waitTime }} мин.</p>
      <p class="ahead-count">Перед вами: {{ peopleAhead }} чел.</p>
      <button class="btn-leave" @click="leaveQueue">Покинуть очередь</button>
    </section>
  </div>
</template>



<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()

const queueTitle = ref('')
const waitTime = ref(0)
const peopleAhead = ref(0)

onMounted(async () => {
  const queueId = route.params.id
  try {
    const response = await axios.get(`http://localhost:8080/api/queues/${queueId}`)
    const data = response.data
    queueTitle.value = data.name
    
    waitTime.value = data.waitTime || 0
    peopleAhead.value = data.peopleAhead || 0
  } catch (error) {
    alert('Очередь не найдена')
    router.push('/')
  }
})

const leaveQueue = () => {
  router.push('/')
}
</script>



<style scoped>
.page-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.queue-info {
  text-align: center;
  margin-bottom: 60px;
}

.queue-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--text);
}

.content {
  text-align: center;
}

.wait-time {
  font-size: 40px;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 8px;
}

.ahead-count {
  font-size: 16px;
  font-weight: 400;
  color: var(--text-muted);
  margin-bottom: 60px;
}

.btn-leave {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 16px 56px;
  font-size: 18px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.2s;
}

.btn-leave:hover {
  background: var(--teal);
}
</style>