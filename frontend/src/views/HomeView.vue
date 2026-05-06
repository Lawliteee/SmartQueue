<template>
  <main>
    <h1>Управляйте очередями без хаоса и потери времени</h1>
    <p class="subtitle">Smart Queue — гибкая система виртуальных очередей с чатом, обменом местами и приоритетами</p>

    <section class="role-section">
      <h2>Выберите роль:</h2>
      <div class="role-cards">
        <div class="role-card">
          <span class="role-label">Участник</span>
          <button class="btn-primary1" @click="modalOpen = true">Вступить в очередь</button>
        </div>
        <div class="role-card">
          <span class="role-label">Организатор</span>
          <button class="btn-primary2" :class="{ 'btn-disabled': !currentUser }"
            :disabled="!currentUser" @click="handleCreateQueue">
            Создать очередь
          </button>
          <p v-if="!currentUser" class="auth-hint">Необходимо авторизоваться</p>
        </div>
      </div>
    </section>
  </main>

  <JoinQueueModal v-if="modalOpen" @close="modalOpen = false" />
</template>



<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import JoinQueueModal from '../components/JoinModal.vue'
import { getUser } from '../utils/auth.js'

const router = useRouter()
const modalOpen = ref(false)
const currentUser = ref(null)

onMounted(() => {
  currentUser.value = getUser()
})

function handleCreateQueue() {
  if (!currentUser.value) return
  router.push('/create-queue')
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

.btn-primary1:hover, .btn-primary2:hover { background: var(--teal); }


.btn-disabled {
  background: #6d6d6d;
  cursor: not-allowed;
  opacity: 0.7;
}

.btn-disabled:hover {
  background: #535353;
}

.auth-hint {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
  margin-bottom: 0;
}
</style>
