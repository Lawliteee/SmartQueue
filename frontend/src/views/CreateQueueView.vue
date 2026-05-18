<template>
  <main>
    <h1>Создание очереди</h1>

    <section class="form-panel">

      <!-- Название очереди -->
      <div class="form-row first-row">
      <label class="form-label" for="queueName">Название очереди</label>
        <div class="form-control">
          <input id="queueName" type="text" class="form-input" v-model="queueName"/>
        </div>
      </div>

      <!-- Описание -->
      <div class="form-row">
        <label class="form-label" for="description">Описание</label>
        <div class="form-control">
          <textarea id="description" class="form-input form-textarea" v-model="description"></textarea>
        </div>
      </div>

      <!-- Время начала и макс. учатников -->
      <div class="form-row">
        <div class="inline-fields">
          <div class="inline-field">
            <label class="inline-label">Время начала</label>
            <input type="time" class="form-input select-small" v-model="startTime" />
          </div>
          <div class="inline-field">
            <label class="inline-label">Макс. участников</label>
            <input type="text" class="form-input input-small" v-model="maxParticipants" />
          </div>
          <label class="checkbox-item">
            <input type="checkbox" v-model="swapPositions" />
            <span>Обмен позициями</span>
          </label>
        </div>
      </div>

      <!-- Приоритеты -->
      <div class="form-row">
        <div class="priority-row disabled">
          <label class="checkbox-item">
            <input type="checkbox" v-model="hasPriority" />
            <span>Приоритеты</span>
          </label>
          <label class="inline-label ml">Количество приоритетов</label>
          <input type="text" class="form-input input-small" v-model="priorityCount" />
          <label class="inline-label ml">Начальный приоритет</label>
          <input type="text" class="form-input input-small" v-model="initialPriority" />
        </div>
      </div>

      <!-- Уже не Ненужная фигня -->

      <div class="form-row">
        <div class="checkbox-row">
          <label class="checkbox-item">
            <input type="checkbox" v-model="skipFeature" />
            <span>Функция "Пропустить меня"</span>
          </label>
          <div class="inline-field">
            <input type="number" class="form-input new-input-small" v-model="skipDuration" :disabled="!skipFeature" min="1" max="99"/>
            <label class="inline-label">мин.</label>
          </div>
          <label class="checkbox-item">
            <input type="checkbox" v-model="imFreeFeature" />
            <span>Функция "Я освободился"</span>
          </label>

          <span class="info-icon" title="Пропустить меня - участник временно становится в конец очереди; Я освободился - участники могут без администратора двигать очередь">i</span>
        </div>
      </div>

      <!-- Администраторы -->
      <div class="form-row disabled">
        <span class="form-label">Администраторы</span>
        <div class="form-control">
          <div class="admin-row">
            <input type="text" class="form-input" placeholder="Добавьте администраторов..." v-model="adminInput" @keyup.enter="addAdmin" />
            <button class="btn-add" @click="addAdmin">+</button>
          </div>
          <div v-if="admins.length > 0" class="admin-list">
            <span v-for="(admin, idx) in admins" :key="idx" class="admin-tag"> {{ admin }}
              <button @click="removeAdmin(idx)" class="remove-admin">&times;</button>
            </span>
          </div>
        </div>
      </div>
    </section>

    <section class="actions">
      <button class="btn-back" @click="goBack">Назад</button>
      <button class="btn-create" @click="createQueue">Создать</button>
    </section>

    <QueueCreatedModal v-if="createdQueueId":queueId="createdQueueId" @close="router.push(`/admin/${createdQueueId}`)"/>
  </main>
</template>



<script setup>
import QueueCreatedModal from '../components/QueueCreatedModal.vue'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { getUser } from '../utils/auth.js'

const router = useRouter()

const createdQueueId = ref(null)

const queueName = ref('')
const description = ref('')
const startTime = ref('12:00')
const maxParticipants = ref('')
const hasPriority = ref(false)
const priorityCount = ref('')
const initialPriority = ref('')
const skipFeature = ref(false)
const skipDuration = ref(10)
const imFreeFeature = ref(false)
const swapPositions = ref(false)
const adminInput = ref('')
const admins = ref([])

// Возвращает на главную старницу
const goBack = () => {
  router.push('/')
}

// Создает очередь
const createQueue = async () => {
  try {
    if (!queueName.value.trim()) {
    alert('Введите название очереди')
    return
    }

    const payload = { // Объект с данными
      queueName: queueName.value,
      description: description.value,
      startTime: startTime.value,
      maxParticipants: parseInt(maxParticipants.value) || 0,
      hasPriority: hasPriority.value,
      priorityCount: parseInt(priorityCount.value) || 0,
      initialPriority: parseInt(initialPriority.value) || 0,
      skipFeature: skipFeature.value,
      skipDuration: parseInt(skipDuration.value) || 10,
      imFreeFeature: imFreeFeature.value,
      swapPositions: swapPositions.value,
      admins: [getUser()?.id, ...admins.value].filter(Boolean)
    }

    // Отправляем Post запрос
    const response = await axios.post('http://localhost:8080/api/queues', payload)
    const { link } = response.data

    createdQueueId.value = response.data.id // Модальное окно с идентификатором
  } catch (error) {
    alert('Не удалось создать очередь')
  }
}
</script>



<style scoped>
main {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 20px 60px;
}

h1 {
  font-size: 32px;
  font-weight: 400;
  color: var(--text);
  text-align: center;
  max-width: 700px;
  width: 100%;
  margin-bottom: 24px;
}

.form-panel {
  background: var(--panel);
  border-radius: 18px;
  padding: 6px 32px;
  width: 100%;
  max-width: 720px;
}

.form-row {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 18px 0;
  border-top: 1px solid var(--divider);
}

.form-row.first-row {
  border-top: none;
}

.form-label {
  font-size: 15px;
  font-weight: 400;
  color: var(--text);
  min-width: 170px;
  padding-top: 9px;
  flex-shrink: 0;
}

.form-control {
  flex: 1;
}

.form-input {
  font-family: 'Fira Sans', sans-serif;
  font-size: 14px;
  background: var(--input-bg);
  border: none;
  border-radius: 8px;
  padding: 9px 14px;
  width: 100%;
  color: var(--text);
  outline: none;
  transition: box-shadow 0.15s;
}

.form-input:focus {
  box-shadow: 0 0 0 2px var(--teal);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.inline-fields {
  display: flex;
  gap: 32px;
  flex: 1;
  align-items: center;
  flex-wrap: wrap;
}

.inline-field {
  display: flex;
  align-items: center;
  gap: 10px;
}

.inline-label {
  font-size: 14px;
  font-weight: 400;
  white-space: nowrap;
  color: var(--text);
}

.new-inline-label {
  font-size: 12px;
  font-weight: 400;
  white-space: nowrap;
  color: var(--text);
}

.select-small {
  width: auto;
  cursor: pointer;
  padding: 9px 10px;
}

.input-small {
  width: 70px;
}

.new-input-small {
  width: 60px;
  padding-right: 0px;
  margin-right: 0px;
}

.ml {
  margin-left: 16px;
}

.priority-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  flex: 1;
}

.checkbox-row {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
  flex: 1;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 14px;
  font-weight: 400;
  cursor: pointer;
  color: var(--text);
}

input[type="checkbox"] {
  width: 17px;
  height: 17px;
  accent-color: var(--teal-dark);
  cursor: pointer;
  flex-shrink: 0;
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
  cursor: default;
  user-select: none;
}

.admin-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.admin-row .form-input {
  flex: 1;
}

.btn-add {
  background: var(--input-bg);
  border: 1px solid var(--divider);
  border-radius: 8px;
  width: 38px;
  height: 38px;
  font-size: 22px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text);
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.15s;
  flex-shrink: 0;
}

.btn-add:hover {
  background: #f0f0f0;
}

.admin-list {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.actions {
  display: flex;
  gap: 16px;
  margin-top: 28px;
  width: 100%;
  max-width: 700px;
  justify-content: center;
}

.btn-back {
  background: var(--panel);
  color: var(--text);
  border: none;
  border-radius: 10px;
  padding: 14px 40px;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.15s;
}

.btn-back:hover {
  background: #cfcfcf;
}

.btn-create {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 10px;
  padding: 14px 60px;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.15s;
}

.btn-create:hover {
  background: var(--teal);
}

.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.form-input:disabled {
  opacity: 0.5;
  pointer-events: none;
}

</style>
