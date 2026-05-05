<template>
  <div class="account-modal" v-click-outside="() => $emit('close')">
    <div class="account-content">
      <p class="account-name">{{ user.displayName }}</p>
      <p class="account-email">{{ user.email }}</p>
      <hr class="divider" />
      <button class="btn-logout" @click="$emit('logout')">Выйти из аккаунта</button>
    </div>
  </div>
</template>

<script setup>
defineProps({ user: { type: Object, required: true } })
defineEmits(['close', 'logout'])
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
  min-width: 200px;
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

</style>
