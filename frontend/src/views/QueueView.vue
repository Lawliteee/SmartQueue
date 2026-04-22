<template>
  <div class="page-body">

    <button class="burger-btn" @click="openParticipants" title="Список участников">
      <img src="/icons/burger.png" alt="меню" width="22" height="22" />
    </button>

    <section class="queue-info">
      <h2 class="queue-title">{{ queueTitle }}</h2>
    </section>
    <section class="content">
      <p v-if="waitTime === -1" class="your-turn">Ваша очередь!</p>
      <p v-else class="wait-time">Осталось ждать: {{ waitTime }} мин.</p>
      <p class="ahead-count">Перед вами: {{ peopleAhead }} чел.</p>
      <button class="btn-leave" :disabled="waitTime === -1" :class="{ 'btn-leave--disabled': waitTime === -1 }"
      @click="leaveQueue">
  Покинуть очередь
</button>
    </section>
  </div>
</template>



<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

let pollTimer = null

const route = useRoute()
const router = useRouter()

const queueTitle = ref('')
const waitTime = ref(0)
const peopleAhead = ref(0)

onMounted(async () => {
  await fetchQueue()                        // загружаем данные
  pollTimer = setInterval(fetchQueue, 1500) // устанавливаем таймер на 1,5 секунды
})

onUnmounted(() => clearInterval(pollTimer)) // останавливаем таймер при уходе со страницы

async function fetchQueue() {
  const queueId = route.params.id
  try {
    const response = await axios.get(`http://localhost:8080/api/queues/${queueId}`)
    const data = response.data

    if (data.finished) {
      alert('Очередь завершена')
      router.push('/')
      return
    }

    const myId = sessionStorage.getItem('participantId')
    const cp = data.currentParticipant

    // Если текущий вызванный - я
    if (cp && cp.id === myId) {
      queueTitle.value = data.name
      waitTime.value = -1  // флаг "вызван"
      peopleAhead.value = 0
      return
    }

    // Моя позиция в очереди ожидания
    const myPos = data.participants.findIndex(p => p.id === myId)

    if (myPos === -1) {       // если нет в очереди значит выгнали
      sessionStorage.removeItem('participantId')
      router.push('/')
      return
    }
    queueTitle.value  = data.name
    peopleAhead.value = myPos === -1 ? 0 : myPos
    waitTime.value = myPos === -1 ? 0 : myPos * 5
  } catch (error) {
    clearInterval(pollTimer)
    router.push('/')
  }
}

async function leaveQueue() {
  const participantId = sessionStorage.getItem('participantId')
  if (participantId) {
    await axios.delete(
      `http://localhost:8080/api/queues/${route.params.id}/participants/${participantId}`
    )
    sessionStorage.removeItem('participantId')
  }
  router.push('/')
}

// TODO
function openParticipants() {
  // список участников
}
</script>



<style scoped>
.page-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

/* Список участников */
.burger-btn {
  position: absolute;
  top: -8px;
  right: 24px;
  width: 46px;
  height: 46px;
  background: var(--panel);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.20s;
  flex-shrink: 0;
}
 
.burger-btn:hover {
  background: #cfcfcf;
}


/* --------------------------- */

.queue-info {
  text-align: center;
  margin-bottom: 60px;
  margin-top: -60px;
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

.your-turn {
  font-size: 40px;
  font-weight: 700;
  color: var(--teal-dark);
  margin-bottom: 8px;  
}

.btn-leave--disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>