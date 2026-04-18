<template>
  <main>
    <h1>Управляйте очередями без хаоса и потери времени</h1>
    <p class="subtitle">Smart Queue — гибкая система виртуальных очередей с чатом, обменом местами и приоритетами</p>

    <section class="role-section">
      <h2>Выберите роль:</h2>
      <div class="role-cards">
        <div class="role-card">
          <span class="role-label">Участник</span>
          <button class="btn-primary1" @click="openModal">Вступить в очередь</button>
        </div>
        <div class="role-card">
          <span class="role-label">Организатор</span>
          <button class="btn-primary2" @click="router.push('/create-queue')">Создать очередь</button>
        </div>
      </div>
    </section>
  </main>

  <div class="modal" v-if="modalOpen" @click.self="modalOpen = false">
    <div class="modal-content">
      <h3>Вступить в очередь</h3>
      <div class="input-group">
        <input v-model="userName"  type="text" class="modal-input" placeholder="Ваше имя" />
        <input v-model="queueLink" type="text" class="modal-input" placeholder="Ссылка на очередь" />
      </div>
      <button class="modal-btn":disabled="!canSubmit":class="{ 'modal-btn--disabled': !canSubmit }" @click="submitQueue">
        Вступить
      </button>
    </div>
  </div>
</template>



<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const userName  = ref('')
const modalOpen = ref(false)
const queueLink = ref('')

const canSubmit = computed(() =>
  userName.value.trim() !== '' && queueLink.value.trim() !== ''
)

function openModal() {
  modalOpen.value = true
}

async function submitQueue() {
  if (!canSubmit.value) return

  const url = queueLink.value.trim()
  const parts = url.split('/')
  const queueId = parts[parts.length - 1]

  // TODO
  // await axios.post(`/api/queues/${queueId}/join`, { name: userName.value.trim() })

  router.push(`/queue/${queueId}`)
  modalOpen.value = false
}
</script>

<style scoped>
main {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 70px 20px;
  text-align: center;
}

h1 {
  font-size: 40px;
  font-weight: 400;
  color: var(--text);
  max-width: 1200px;
  line-height: 1.2;
  margin-bottom: 14px;
}

.subtitle {
  font-size: 16px;
  color: var(--text-muted);
  font-weight: 400;
  max-width: 520px;
  line-height: 1.6;
  margin-bottom: 92px;
}

.role-section h2 {
  font-size: 28px;
  font-weight: 400;
  margin-bottom: 28px;
  color: var(--text);
}

.role-cards {
  display: flex;
}

.role-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 0 40px;
}

.role-card + .role-card {
  border-left: 2px solid var(--divider);
}

.role-label {
  font-size: 18px;
  font-weight: 400;
  color: var(--text);
}

.btn-primary1, .btn-primary2 {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 10px;
  padding: 14px 32px;
  font-size: 18px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.18s;
  white-space: nowrap;
}

.btn-primary2 {
  padding: 14px 40px;
}

.btn-primary1:hover, .btn-primary2:hover {
  background: var(--teal);
}

/* Модальное окно */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.modal-content {
  background: white;
  padding: 35px;
  border-radius: 18px;
  width: 100%;
  max-width: 400px;
  text-align: center;
  box-sizing: border-box;
}

.modal-content h3 {
  font-size: 21px;
  font-weight: 400;
  margin-bottom: 18px;
  color: var(--text);
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 28px;
}

.modal-input {
  width: 100%;
  padding: 14px;
  border: 1px solid var(--teal);
  border-radius: 12px;
  font-family: 'Fira Sans', sans-serif;
  font-size: 15px;
  color: var(--text);
  box-sizing: border-box;
  transition: border-color 0.18s;
}

.modal-input::placeholder { color: #aaa; }

.modal-input:focus {
  outline: none;
  border-color: var(--teal-dark);
}

.modal-btn {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 14px 20px;
  font-size: 16px;
  font-weight: 700;
  width: 60%;
  font-family: 'Fira Sans', sans-serif;
  cursor: pointer;
  transition: background 0.18s, opacity 0.18s;
}

.modal-btn:hover:not(.modal-btn--disabled) {
  background: var(--teal);
}

.modal-btn--disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>