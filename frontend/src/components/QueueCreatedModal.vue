<template>
  <div class="modal" @click.self="$emit('close')">
    <div class="modal-content">
      <h3>Очередь создана!</h3>
      <p class="label">Идентификатор для участников:</p>
      <div class="id-row">
        <input class="id-input" :value="queueId" readonly />
        <button class="btn-copy" @click="copy">{{ copied ? 'Скопировано' : 'Копировать' }}</button>
      </div>
      <button class="modal-btn" @click="$emit('close')">Перейти к очереди</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const props = defineProps({ queueId: String })
defineEmits(['close'])

const copied = ref(false)

// Копирует id очереди
function copy() {
  navigator.clipboard.writeText(props.queueId)
  copied.value = true
  setTimeout(() => { copied.value = false }, 1500)
}
</script>

<style scoped>

.modal {
  position: fixed;
  top: 0;left: 0;
  width: 100%; height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 16px;
  width: 100%;
  max-width: 420px;
  text-align: center;
  box-sizing: border-box;
}

h3 {
  font-size: 21px;
  font-weight: 400;
  margin-bottom: 20px;
  color: var(--text);
}

.label {
  font-size: 14px;
  text-align: left;
  color: var(--text-muted);
  margin-bottom: 10px;
  padding-left: 4px;
}

.id-row {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}

.id-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid var(--teal);
  border-radius: 10px;
  font-family: 'Fira Sans', sans-serif;
  font-size: 14px;
  color: var(--text);
  background: white;
  outline: none;
}

.btn-copy {
  padding: 10px 16px;
  border: 1px solid var(--teal);
  border-radius: 10px;
  background: white;
  font-family: 'Fira Sans', sans-serif;
  font-size: 14px;
  cursor: pointer;
  color: var(--text);
  white-space: nowrap;
  transition: background 0.20s;
}

.btn-copy:hover { background: var(--panel); }

.modal-btn {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 14px 1px;
  font-size: 16px;
  font-weight: 700;
  width: 70%;
  font-family: 'Fira Sans', sans-serif;
  cursor: pointer;
  transition: background 0.20s;
}

.modal-btn:hover { background: var(--teal); }
</style>