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
      <h3>Введите ссылку на очередь:</h3>
      <input v-model="queueLink" type="text" class="modal-input" />
      <button class="modal-btn" @click="submitQueue">Вступить</button>
    </div>
  </div>
</template>



<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const modalOpen = ref(false) // Открыто ли окно
const queueLink = ref('')    // Ссылка, котору. ввел пользователь

function openModal() {
  modalOpen.value = true
}

function submitQueue() {
  if (queueLink.value.trim()) {
    window.location.href = queueLink.value.trim()
  }
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
  max-width: 400px;
  text-align: center;
}

.modal-content h3 {
  font-size: 21px;
  font-weight: 400;
  margin-bottom: 18px;
  color: var(--text);
}

.modal-input {
  width: 100%;
  padding: 14px;
  margin-bottom: 30px;
  border: 1px solid var(--divider);
  border-color: var(--teal);
  border-radius: 12px;
  font-family: 'Fira Sans', sans-serif;
  font-size: 15px;
}

.modal-input:focus {
  outline: none;
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
}

.modal-btn:hover {
  background: var(--teal);
}
</style>