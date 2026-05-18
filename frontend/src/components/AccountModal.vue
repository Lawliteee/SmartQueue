<template>
  <div class="account-modal">
    <div class="account-content">
      <p class="account-name">{{ user.displayName }}</p>
      <p class="account-email">{{ user.email }}</p>

      <!-- Мои очереди -->
      <hr class="divider" />
      <p class="section-label">Мои очереди</p>
      <div v-if="adminQueues.length > 0" class="queue-list">
        <div v-for="q in adminQueues" :key="q.id" class="queue-item queue-item--admin" @click="$emit('go-to-admin', q.id)">
          {{ q.name }}
        </div>
      </div>
      <div v-if="participantQueues.length > 0" class="queue-list">
        <div v-for="q in participantQueues" :key="q.id" class="queue-item" @click="$emit('go-to-queue', q)">
          {{ q.name }}
        </div>
      </div>
      <p v-else class="queue-empty">Нет активных очередей</p>

      <hr class="divider" />
      <button class="btn-logout" @click="$emit('logout')">Выйти из аккаунта</button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  user: { type: Object, required: true },
  adminQueues: { type: Array, default: () => [] },
  participantQueues: { type: Array, default: () => [] },
})
defineEmits(['close', 'logout', 'go-to-admin', 'go-to-queue'])
</script>

<style scoped>
.account-modal {
  position: absolute;
  top: calc(100% + 8px);
  right: 6px;
  z-index: 200;
  animation: fadeIn 0.15s ease-out;
}

.account-content {
  background: white;
  border-radius: 14px;
  padding: 16px 20px;
  min-width: 220px;
  max-width: 280px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
  border: 1px solid var(--divider);
}

.account-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
  margin: 0 0 4px;
}

.account-email {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
}

.divider {
  border: none;
  border-top: 1px solid var(--divider);
  margin: 12px 0;
}

.section-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0 0 8px;
}

.queue-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.queue-item {
  font-size: 14px;
  color: var(--text);
  cursor: pointer;
  padding: 6px 8px;
  border-radius: 8px;
  transition: background 0.15s;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.queue-item:hover { background: var(--panel); }

.queue-item--admin {
  font-weight: 700;
}

.queue-empty {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
}

.btn-logout {
  background: none;
  border: none;
  font-family: 'Fira Sans', sans-serif;
  font-size: 14px;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0;
  transition: color 0.18s;
}
.btn-logout:hover { color: var(--text); }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-6px); }
  to   { opacity: 1; transform: translateY(0); }
}
</style>
