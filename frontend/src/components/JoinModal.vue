<template>
  <div class="modal" @click.self="$emit('close')">
    <div class="modal-content">
      <h3>Вступить в очередь</h3>
      <div class="input-group">
        <input v-model="userName" type="text" class="modal-input" placeholder="Ваше имя" />
        <input v-model="queueId" type="text" class="modal-input" placeholder="Идентификатор очереди" />
      </div>
      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
      <button
        class = "modal-btn"
        :disabled = "!canSubmit || loading"
        :class = "{ 'modal-btn--disabled': !canSubmit || loading }"
        @click = "submitQueue">
        Вступить
      </button>
    </div>
  </div>
</template>


<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
 
const emit = defineEmits(['close'])
const router = useRouter()
 
const userName  = ref('')
const queueId   = ref('')
const loading   = ref(false)
const errorMsg  = ref('')
 
// Проверяется возможность перехода в очередь (если поля заполнены)
const canSubmit = computed(() =>
  userName.value.trim() !== '' && queueId.value.trim() !== ''
)

// Переход в очередь
async function submitQueue() {
  if (!canSubmit.value || loading.value) return
 
  // Если вставили полную ссылку - извлекаем ID
  const raw = queueId.value.trim()
  const id  = raw.includes('/') ? raw.split('/').filter(Boolean).pop() : raw
 
  loading.value  = true
  errorMsg.value = ''
 
try {
  // Сначала проверяем статус очереди
  const check = await axios.get(`http://localhost:8080/api/queues/${id}`)
  if (check.data.finished) {
    errorMsg.value = 'Эта очередь уже завершена'
    loading.value = false
    return
  }

  const max = check.data.maxParticipants
  if (max && check.data.participants.length >= max) {
    errorMsg.value = `Очередь заполнена`
    loading.value = false
    return
  }

  // Вступаем и сохраняем response
  const response = await axios.post(`http://localhost:8080/api/queues/${id}/join`, {
    name: userName.value.trim(),
  })
  sessionStorage.setItem('participantId', response.data.participantId)
  router.push(`/queue/${id}`)
  emit('close')
} catch (err) {
  if (err.response?.status === 409) {
    errorMsg.value = 'Очередь заполнена'
  } else if (err.response?.status === 404) {
    errorMsg.value = 'Очередь не найдена. Проверьте идентификатор.'
  } else {
    errorMsg.value = 'Не удалось вступить в очередь.'
  }
} finally {
  loading.value = false
}
}
</script>


<style scoped>
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
 
.error-msg {
  font-size: 13px;
  color: #e85656;
  margin-bottom: 14px;
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
  transition: background 0.20s, opacity 0.20s;
}
 
.modal-btn:hover:not(.modal-btn--disabled) {
  background: var(--teal);
}
 
.modal-btn--disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>