<template>
  <div class="admin-page">

    <!-- Верхняя панель -->
    <div class="queue-bar">
      <span class="queue-name">{{ queue.name }}</span>
      <span class="queue-time">начало: {{ queue.startTime }}</span>
      <button class="btn-copy" :class="{ 'btn-copy--copied': copied }" @click="copyId">
        {{ copied ? 'Скопировано' : 'Скопировать ID' }}
      </button>
    </div>

    <!-- Текущий участник -->
    <div class="current-block">
      <div class="current-label">
        <span class="current-text">Текущий участник:</span>
        <span class="current-number">#{{ queue.currentNumber }}</span>
      </div>
      <div class="current-participant-info">
        <span class="participant-name">{{ currentParticipant?.name ?? '—' }}</span>
        <span class="participant-priority">Приоритет {{ currentParticipant?.priority ?? '-' }}</span>
      </div>
    </div>

    <!-- Список участников -->
    <div class="participants-panel">
      <div
        v-for="(p, idx) in queue.participants"
        :key="p.id"
        class="participant-row"
        :class="{ 'participant-row--first': idx === 0 }"
      >
        <span class="p-number">{{ idx + 1 }}.</span>
        <span class="p-name">{{ p.name }}</span>
      </div>
      <div v-if="queue.participants.length === 0" class="empty-list">
        Участников пока нет
      </div>
    </div>

    <!-- кнопки внизу -->
    <div class="actions">
      <button class="btn-secondary">Открыть чат</button>
      <button class="btn-next" @click="callNext">Позвать следующего</button>
      <button class="btn-secondary" @click="finishQueue">Завершить очередь</button>
    </div>
  </div>

</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

let pollTimer = null

const route = useRoute()
const router = useRouter()

const queue = ref({
  id: null,
  name: '',
  startTime: '',
  currentNumber: 1,
  participants: []
})

const currentParticipant = computed(() => {
  // Текущий участник - первый в списке
  if (queue.value.participants && queue.value.participants.length > 0) {
    return queue.value.participants[0]
  }
  return null
})

onMounted(async () => {
  await fetchQueue()
  pollTimer = setInterval(fetchQueue, 3000)
})

onUnmounted(() => clearInterval(pollTimer))

async function fetchQueue() {
  const queueId = route.params.id
  try {
    const response = await axios.get(`http://localhost:8080/api/queues/${queueId}`)
    const data = response.data
    queue.value = {
      id: data.id,
      name: data.name,
      startTime: data.startTime || '12:00',
      currentNumber: data.currentNumber || 1,
      participants: data.participants || []
    }
  } catch (error) {
    console.error('Ошибка загрузки:', error)
  }
}

const copied = ref(false) 
function copyId() {
  if (queue.value.id) {
    navigator.clipboard.writeText(queue.value.id)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}

function callNext() {
  // TODO: POST
  queue.value.currentNumber++
  if (queue.value.participants.length > 0) {
    queue.value.participants.shift()
  }
}

function finishQueue() {
  // TODO: POST
  router.push('/')
}

</script>

<style scoped>
.admin-page {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0px 48px 55px;
  gap: 20px;
  max-width: 900px;
  width: 100%;
  margin: 0 auto;
  box-sizing: border-box;
  height: 100vh;
  max-height: 100vh;
  overflow: hidden;
}

.queue-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  flex-shrink: 0;
}

.queue-name {
  font-size: 18px;
  font-weight: 500;
  color: var(--text);
}

.queue-time {
  font-size: 16px;
  font-weight: 400;
  color: var(--text-muted);
}

.btn-copy {
  margin-left: auto;
  background: var(--panel);
  color: var(--text);
  border: none;
  border-radius: 10px;
  padding: 8px 18px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.20s, color 0.20s;
  white-space: nowrap;
}
 
.btn-copy:hover { background: #cfcfcf; }
 
.btn-copy--copied {
  background: var(--teal-dark);
  color: white;
}
 
.btn-copy--copied:hover { background: var(--teal); }

.current-block {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 28px;
  flex-wrap: wrap;
  flex-shrink: 0;
}

.current-label {
  display: flex;
  align-items: baseline;
  gap: 16px;
}

.current-text {
  font-size: 44px;
  font-weight: 700;
  color: var(--text);
}

.current-number {
  font-size: 60px;
  font-weight: 700;
  color: var(--text);
  line-height: 1;
}

.current-participant-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.participant-name {
  font-size: 20px;
  font-weight: 500;
  color: var(--text);
  line-height: 1.2;
}

.participant-priority {
  font-size: 13px;
  font-weight: 400;
  color: var(--text-muted);
}

.participants-panel {
  background: var(--panel);
  border-radius: 18px;
  padding: 8px 0;
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.participant-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  border-top: 1px solid var(--divider);
}

.participant-row--first {
  border-top: none;
}

.participant-row--first .p-name {
  font-weight: 700;
}

.p-number {
  font-size: 14px;
  font-weight: 400;
  color: var(--text);
  min-width: 24px;
}

.p-name {
  font-size: 14px;
  font-weight: 400;
  color: var(--text);
}

.empty-list {
  padding: 40px 28px;
  text-align: center;
  color: var(--text-muted);
  font-size: 15px;
}

.actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  flex-shrink: 0;
  margin-top: auto;
}

.btn-secondary {
  background: var(--panel);
  color: var(--text);
  border: none;
  border-radius: 12px;
  padding: 12px 24px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.20s;
  white-space: nowrap;
}

.btn-secondary:hover {
  background: #cfcfcf;
}

.btn-next {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 12px 36px;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.18s;
  white-space: nowrap;
}

.btn-next:hover {
  background: var(--teal);
}
</style>