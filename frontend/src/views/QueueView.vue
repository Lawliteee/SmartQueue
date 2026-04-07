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
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const queueTitle = ref('Тестовая очередь')
const waitTime = ref(0)
const peopleAhead = ref(0)

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