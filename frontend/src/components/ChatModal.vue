<template>
  <div class="chat-modal">
    <div class="chat-content">
      <div class="chat-header">
        <h3>Чат</h3>
        <button class="btn-close" @click="$emit('close')">Х</button>
      </div>
      <div class="chat-body">
        <p>Не реализовано 0_0</p>
      </div>
    </div>
  </div>
</template>

<script setup>
defineEmits(['close'])
</script>

<style scoped>

.chat-modal {
  position: fixed;
  top: 65px;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  padding-left: 20px;
  box-sizing: border-box;
  z-index: 10;
  pointer-events: none;
}

.chat-content {
  background: white;
  border-radius: 18px;
  width: 320px;
  height: 600px;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  overflow: hidden;
  margin: 20px 0;
  border: 1px solid var(--divider);
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  pointer-events: auto;
  animation: slideInLeft 0.3s ease-out;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--divider);
  flex-shrink: 0;
}

h3 {
  font-size: 18px;
  font-weight: 500;
  color: var(--text);
  margin: 0;
}

.btn-close {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: var(--text-muted);
  padding: 0;
  line-height: 1;
}

.btn-close:hover { color: var(--text); }

.chat-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

@keyframes slideInLeft {
  from {
    transform: translateX(-110%);
    opacity: 0;
  } to {
    transform: translateX(0);
    opacity: 1;
  }
}
</style>
