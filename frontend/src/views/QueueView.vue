<template>
  <div class="page-body">
    <!-- Кнопка чата -->
    <button class="chat-btn" @click="showChat = true" title="Чат">
      <img src="/icons/chat.png" alt="чат" width="30" height="30" />
    </button>

    <!-- Список участников -->
    <button class="burger-btn" @click="openParticipants" title="Список участников">
      <img src="/icons/burger.png" alt="меню" width="22" height="22" />
    </button>

    <!-- Заголовок и осн. инфа -->
    <section class="queue-info">
      <div class="queue-title-wrapper">
      <h2 class="queue-title">{{ queueTitle }}</h2>
      <div v-show="queueDescription" class="info-icon" data-tooltip="" @mouseenter="showTooltip = true" @mouseleave="showTooltip = false">
        i
      <div v-if="showTooltip" class="tooltip">{{ queueDescription }}</div>
    </div>
  </div>
    </section>
    <section class="content">
      <p v-if="waitTime === -1" class="your-turn">Ваша очередь!</p>
      <p v-else-if="currentNumber === 0" class="wait-time">Очередь ещё не началась</p>
      <p v-else-if="waitTime > 0" class="wait-time">
          Осталось ждать: {{ waitTime }} мин.
        </p>

        <p v-else class="wait-time">
          Ожидание рассчитывается...
        </p>
      <p v-if="currentNumber > 0" class="ahead-count">Перед вами: {{ peopleAhead }} чел.</p>
      <p v-else class="ahead-count" style="visibility: hidden">placeholder</p>
      <div class="btn-row">

        <!-- Кнопка пропуска меня -->
        <template v-if="skipFeature">
          <button v-if="!isSkipped && waitTime !== -1" class="btn-skip" @click="skipMe":disabled="currentNumber === 0 || skipCooldown">
            Пропустить меня
          </button>
          <button v-else-if="isSkipped" class="btn-return"@click="returnMe">
            Вернуться ({{ skipTimeLeft }}с)
           </button>
        </template>

        <!-- Обычная кнопка -->
        <button v-if="waitTime !== -1" class="btn-leave" @click="leaveQueue">
          Покинуть очередь
        </button>

        <!-- Кнопка для вызванного -->
        <button v-else-if="imFreeFeature" class="btn-leave" @click="imFree">
          Я освободился
        </button>

        <!-- Заблокированная версия -->
        <button v-else class="btn-leave btn-leave--disabled" disabled>
          Покинуть очередь
        </button>
      </div>
    </section>

    <!-- Модалки -->
    <QueueFinishedModal v-if="showFinished" @confirm="onFinishedConfirm" />
    <ParticipantsModal v-if="showParticipants":participants="participants":currentParticipant="currentParticipant":myId="getCookie('participantId')"
      :swapEnabled="swapEnabled" @close="showParticipants = false" @swap-requested="onSwapRequested":queueStarted="currentNumber > 0"/>
    <SwapRequestModal v-if="showSwapRequest":fromName="swapFromName":fromPos="swapFromPos" @accept="acceptSwap" @decline="declineSwap"/>
    <KickedModal v-if="showKicked" @confirm="router.push('/')"/>
    <ChatModal v-if="showChat" @close="showChat = false" />
    <SwapDeclinedModal v-if="showSwapDeclined":fromName="swapDeclinedName" @confirm="showSwapDeclined = false"/>
  </div>
</template>

<script setup>
import { getCookie, removeCookie } from '../utils/cookies.js'
import QueueFinishedModal from '../components/QueueFinishedModal.vue'
import ParticipantsModal from '../components/ParticipantsModal.vue'
import SwapRequestModal from '../components/SwapRequestModal.vue'
import ChatModal from '../components/ChatModal.vue'
import KickedModal from '../components/KickedModal.vue'
import SwapDeclinedModal from '../components/SwapDeclinedModal.vue'
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const showFinished = ref(false)
const showKicked = ref(false)
const showChat = ref(false)
const imFreeFeature = ref(false)

const showParticipants = ref(false)
const participants = ref([])

const route = useRoute()
const router = useRouter()

const showSwapRequest = ref(false)
const swapFromName = ref('')
const swapFromPos = ref(0)

const showSwapDeclined = ref(false)
const swapDeclinedName = ref('')

const queueTitle = ref('')
const waitTime = ref(0)
const peopleAhead = ref(0)
const currentNumber = ref(0)
const currentParticipant = ref(null)
const pendingSwapId = ref(null)
const swapEnabled = ref(false)
const queueDescription = ref('')
const showTooltip = ref(false)

const isSkipped = ref(false)
const skipUntil = ref(null)
const skipTimeLeft = ref(0)
let skipTimer = null
const skipFeature = ref(false)
const skipDurationMinutes = ref(10)
const skipCooldown = ref(false)


let ws = null
onMounted(async () => {
  await fetchQueue()   // первичная загрузка через HTTP
  connectWS()
})

onUnmounted(() => {
  ws?.close()
  window.removeEventListener('keydown', onKeydown)
  clearInterval(skipTimer)
})

function connectWS() {
  const queueId = route.params.id
  if (!queueId) return
  const participantId = getCookie('participantId') ?? ''
  ws = new WebSocket(
    `ws://localhost:8080/api/ws/queue/${queueId}?participantId=${participantId}`
  )

  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data)
    if (msg.type === 'queue_update') {
      handleUpdate(msg.data)
    } else if (msg.type === 'swap_request') {
      // Пришло предложение обмена
      pendingSwapId.value = msg.swapId
      swapFromName.value = msg.fromName
      swapFromPos.value = msg.fromPos
      showSwapRequest.value = true
    } else if (msg.type === 'swap_declined') {
      // Нам отказали
      swapDeclinedName.value = msg.fromName
      showSwapDeclined.value = true
    }
  }

  ws.onclose = () => {
    const queueId = route.params.id
    if (queueId) setTimeout(connectWS, 2000)
  }
}

async function fetchQueue() {
  const queueId = route.params.id
  console.log('queueId:', queueId)
  try {
    const response = await axios.get(`http://localhost:8080/api/queues/${queueId}`)
    console.log('response:', response.data)
    handleUpdate(response.data)
  } catch (error) {
    console.error('fetchQueue error:', error)
    router.push('/')
  }
}

function handleUpdate(data) {
  participants.value = data.participants || []
  currentNumber.value = data.currentNumber ?? 0
  currentParticipant.value = data.currentParticipant ?? null
  imFreeFeature.value = data.imFreeFeature ?? false
  swapEnabled.value = data.swapPositions ?? false
  queueDescription.value = data.description ?? ''
  queueTitle.value = data.name
  skipFeature.value = data.skipFeature ?? false
  skipDurationMinutes.value = data.skipDuration ?? 10

  if (data.finished) {
    ws?.close()
    showFinished.value = true
    return
  }

  const myId = getCookie('participantId')
  const cp = data.currentParticipant

  if (cp && cp.id === myId) {  // если нет в очереди значит выгнали
    queueTitle.value = data.name
    waitTime.value = -1 // флаг "вызван"
    peopleAhead.value = 0
    return
  }

  // Моя позиция в очереди
  const myPos = data.participants.findIndex(p => p.id === myId)

  if (myPos === -1) { // Если текущий вызванный - я
    removeCookie('participantId')
    ws?.close()
    showKicked.value = true
    return
  }

  peopleAhead.value = myPos
  if (myPos === -1) {
    waitTime.value = null
  } else {
    waitTime.value = data.etAs?.[myPos] ?? null
  }

  const myParticipant = data.participants.find(p => p.id === myId)
  if (myParticipant) {
    isSkipped.value = myParticipant.skipped
    skipUntil.value = myParticipant.skipUntil ? new Date(myParticipant.skipUntil) : null
    if (myParticipant.skipped && skipUntil.value) {
      startSkipCountdown()
    }
  } else {
    isSkipped.value = false
    skipUntil.value = null
  }
}

async function imFree() {
  await axios.post(
    `http://localhost:8080/api/queues/${route.params.id}/im-free`
  )
  removeCookie('participantId')
  router.push('/')
}

async function leaveQueue() {
  const participantId = getCookie('participantId')
  if (participantId) {
    await axios.delete(
      `http://localhost:8080/api/queues/${route.params.id}/participants/${participantId}`,
    )
    removeCookie('participantId')
  }
  router.push('/')
}

// Открывает список участников
function openParticipants() {
  showParticipants.value = true
}

function onFinishedConfirm() {
  showFinished.value = false
  router.push('/')
}

// Согласие на обмен местами
async function acceptSwap() {
  await axios.post(`http://localhost:8080/api/queues/${route.params.id}/swap/respond`, {
    swapId: pendingSwapId.value,
    accept: true,
  })
  showSwapRequest.value = false
  pendingSwapId.value = null
}

// Отказ от обемена местами
async function declineSwap() {
  await axios.post(`http://localhost:8080/api/queues/${route.params.id}/swap/respond`, {
    swapId: pendingSwapId.value,
    accept: false,
  })
  showSwapRequest.value = false
  pendingSwapId.value = null
}


// Тестирование открытия модалки на T
window.addEventListener('keydown', onKeydown)
function onKeydown(e) {
  if (e.key === 't') {
    swapFromName.value = 'Иван'
    showSwapRequest.value = true
  } else if (e.key === 'y') {
    swapDeclinedName.value = 'Иван'
    showSwapDeclined.value = true
  }
}

// Получение предложения об обмене
async function onSwapRequested(targetParticipant) {
  const myId = getCookie('participantId')
  const myName = participants.value.find(p => p.id === myId)?.name ?? ''

  await axios.post(
    `http://localhost:8080/api/queues/${route.params.id}/swap/request`,
    {
      fromId: myId,
      fromName: myName,
      toId: targetParticipant.id,
    }
  )
  showParticipants.value = false
}


// Пропуск участника
function startSkipCountdown() {
  clearInterval(skipTimer)
  skipTimer = setInterval(() => {
    if (!skipUntil.value) { clearInterval(skipTimer); return }
    const left = Math.max(0, Math.ceil((skipUntil.value - Date.now()) / 1000))
    skipTimeLeft.value = left
    if (left === 0) {
      clearInterval(skipTimer)
      isSkipped.value = false
    }
  }, 1000)
}

async function skipMe() {
  const participantId = getCookie('participantId')
  await axios.post(
    `http://localhost:8080/api/queues/${route.params.id}/skip`,
    { participantId }
  )
  skipCooldown.value = true
  setTimeout(() => { skipCooldown.value = false }, 20_000) // Отключение кнопки на минуту
}

async function returnMe() {
  const participantId = getCookie('participantId')
  await axios.post(
    `http://localhost:8080/api/queues/${route.params.id}/return`,
    { participantId }
  )
  skipCooldown.value = true
  setTimeout(() =>  { skipCooldown.value = false }, 20_000) // Отключение кнопки на минуту
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
  position: fixed;
  top: 90px;
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
  transition: background 0.2s;
  flex-shrink: 0;
}

.burger-btn:hover {
  background: #cfcfcf;
}

/* Чат */
.chat-btn {
  position: fixed;
  top: 90px;
  left: 24px;
  width: 46px;
  height: 46px;
  background: var(--panel);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
  flex-shrink: 0;
}

.chat-btn:hover {
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

.btn-row {
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: center;
}

.btn-skip {
  background: var(--panel);
  color: #2c2c2c;
  border: none;
  border-radius: 12px;
  padding: 16px 32px;
  font-size: 18px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.2s;
}

.btn-skip:hover {
  background: #cfcfcf;
}

.btn-skip:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

/* Информация об очереди */
.queue-title-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: center;
}

.info-icon {
  display: inline-flex;
  width: 17px;
  height: 17px;
  border-radius: 50%;
  border: 1.5px solid var(--text-muted);
  color: var(--text-muted);
  font-size: 10px;
  font-weight: 700;
  align-items: center;
  justify-content: center;
  cursor: help;
  position: relative;
  user-select: none;
}

.tooltip {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 50%;
  transform: translateX(-50%);
  background: #2c2c2c;
  color: white;
  font-size: 13px;
  font-weight: 400;
  padding: 8px 12px;
  border-radius: 8px;
  white-space: nowrap;
  z-index: 100;
}

/* Кнопка пропуска */
.btn-return {
  background: #de8900;
  color: white;
  border: none;
  border-radius: 12px;
  padding: 16px 32px;
  font-size: 18px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.2s;
}
.btn-return:hover { background: #c95e00; }

</style>
