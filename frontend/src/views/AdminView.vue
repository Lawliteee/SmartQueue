<template>
  <div class="admin-page">

    <!-- Кнопка чата -->
    <button class="chat-btn" @click="showChat = true" title="Чат">
      <img src="/icons/chat.png" alt="чат" width="30" height="30" />
    </button>

    <!-- Верхняя панель -->
    <div class="queue-bar">
      <span class="queue-name">{{ queue.name }}</span>
      <span class="queue-time">начало: {{ queue.startTime }}</span>
      <button class="btn-copy" :class="{ 'btn-copy--copied': copied }" @click="copyId">
        <span class="btn-text">{{ copied ? 'Скопировано' : 'Скопировать ID' }}</span>
      </button>
    </div>

    <!-- Текущий участник -->
    <div class="current-block" :class="{ 'current-block-ns': notStarted }">
      <div class="current-label">
        <span class="current-text">{{ notStarted ? 'Следующий:' : 'Текущий участник:' }}</span>
        <span class="current-number" v-if="!notStarted">#{{ queue.currentNumber }}</span>
      </div>
      <div class="current-participant-info">
        <span class="participant-name">{{ displayedParticipant?.name ?? '—' }}</span>
        <span class="participant-priority">Приоритет {{ displayedParticipant?.priority ?? '-' }}</span>
      </div>
    </div>

    <!-- Список участников -->
    <div class="participants-panel":class="{ 'participants-panel-ns': notStarted }">
      <div v-for="(p, idx) in displayedParticipants" :key="p.id" class="participant-row" :class="{ 'participant-row--first': idx === 0 }">
        <span class="p-number">{{ idx + 1 }}.</span>
        <span class="p-name">{{ p.name }}</span>
        <button v-if="idx !== 0"
          class="btn-remove" @click="removeParticipant(p.id)"
          title="Исключить">Х</button>
      </div>
     <div v-if="displayedParticipants.length === 0" class="empty-list">Участников пока нет</div>
    </div>

    <!-- кнопки внизу -->
    <div class="actions">
      <button class="btn-next" @click="callNext">
        {{ notStarted ? 'Начать очередь' : 'Позвать следующего' }}
      </button>
      <button class="btn-secondary" @click="finishQueue">Завершить очередь</button>
    </div>
  </div>

  <!-- Модалки -->
  <ChatModal v-if="showChat" :messages="chatMessages" myId="admin" :ws="ws" @close="showChat = false"/>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import api from '../utils/api.js'
import ChatModal from '../components/ChatModal.vue'

const showChat = ref(false)
const chatMessages = ref([])

const notStarted = computed(() => queue.value.currentNumber === 0)

const route = useRoute()
const router = useRouter()

const queue = ref({
  id: null,
  name: '',
  startTime: '',
  currentNumber: 1,
  participants: []
})

// Участник для отображения
const displayedParticipant = computed(() =>
  queue.value.currentParticipant ?? queue.value.participants[0] ?? null
)

// Список для отображения
const displayedParticipants = computed(() => {
  const cp = queue.value.currentParticipant
  if (!cp) {
    return queue.value.participants
  }
  return [cp, ...queue.value.participants] // Список с текущим и последующими участниками
})

let ws = null

onMounted(async () => {
  await fetchQueue()
  connectWS()
})

onUnmounted(() => {
  ws?.close()
})

function connectWS() {
  const queueId = route.params.id
  ws = new WebSocket(
    `ws://localhost:8080/api/ws/queue/${queueId}?participantId=admin&senderName=${encodeURIComponent('Администратор')}`
  )

  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data)
    if (msg.type === 'queue_update') {
      const data = msg.data
      queue.value = {
        id: data.id,
        name: data.name,
        startTime: data.startTime || '12:00',
        currentNumber: data.currentNumber ?? 0,
        participants: data.participants || [],
        currentParticipant: data.currentParticipant || null,
      }
    } else if (msg.type === 'chat_message') {
      chatMessages.value.push({ senderId: msg.senderId, senderName: msg.senderName, text: msg.text,})
      if (chatMessages.value.length > 20) chatMessages.value.shift()
    } else if (msg.type === 'chat_history') {
      chatMessages.value = msg.messages
    }
  }

  ws.onclose = () => {
    setTimeout(connectWS, 2000)
  }
}

async function fetchQueue() {
  const queueId = route.params.id
  try {
    const response = await api.get(`/queues/${queueId}`)
    const data = response.data
    queue.value = {
      id: data.id,
      name: data.name,
      startTime: data.startTime || '12:00',
      currentNumber: data.currentNumber ?? null,
      participants: data.participants || [],
      currentParticipant: data.currentParticipant || null,
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
    setTimeout(() => { copied.value = false }, 1000)
  }
}

// Вызываем следующего участника
async function callNext() {
  const response = await api.post(`/admin/queues/${route.params.id}/next`)
  queue.value.currentNumber = response.data.currentNumber
  queue.value.participants = response.data.participants
  queue.value.currentParticipant = response.data.currentParticipant ?? null
}

// Удаляет участника очереди
async function removeParticipant(participantId) {
  try {
    // Отправляем запрос на удаление
    await axios.delete(
      `http://localhost:8080/api/queues/${route.params.id}/participants/${participantId}`
    )
  } catch (error) {
    console.error('Ошибка при удалении участника:', error)
  }
}

// Завершаем очередь
async function finishQueue() {
  await api.post(`/admin/queues/${route.params.id}/finish`)
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
  position: relative;
}

/* Кнопка чата */
.chat-btn {
  position: fixed;
  top: 90px;
  left: 20px;
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
  z-index: 10;
}

.chat-btn:hover { background: #cfcfcf; }

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
  white-space: nowrap;
  min-width: 130px;
  text-align: center;
  transition: background 0.20s, color 0.20s;
}

.btn-copy:hover { background: #cfcfcf; }

.btn-copy--copied {
  background: var(--teal-dark);
  color: white;
  text-align: center;
}

.btn-copy--copied:hover {
  background: var(--teal);
}

.btn-text {
  display: inline-block;
  width: 100%;
}

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
  position: relative;
  flex-shrink: 0;
  margin-top: auto;
}

.btn-secondary {
  position: absolute;
  right: 0;
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

.participants-panel-ns { opacity: 0.4; }

.current-block-ns { opacity: 0.4; }

/* Кнопка удаления участника */
.btn-remove {
  margin-left: auto;
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 14px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 6px;
  line-height: 1;
  transition: color 0.20s;
}

.btn-remove:hover { color: var(--text); }
</style>
