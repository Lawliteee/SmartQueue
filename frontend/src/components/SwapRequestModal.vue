<template>
  <div class="swap-modal">
    <div class="swap-content">
      <div class="swap-header">
        <h3>Предложение обмена</h3>
      </div>

      <p class="swap-text">{{ fromName }}<span class="swap-pos"> ({{ fromPos }})</span>
         предлагает поменяться с вами местами</p>

      <div class="swap-actions">
        <button class="btn-decline" @click="$emit('decline')">Отказать</button>
        <button class="btn-accept" @click="$emit('accept')">Принять</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  fromName: { type: String, default: '' },
  fromPos: { type: Number, default: 0 }
})
defineEmits(['accept', 'decline'])
</script>

<style scoped>

.swap-modal {
  position: fixed;
  bottom: 32px;
  left: 24px;
  z-index: 20;
  pointer-events: auto;
  animation: slideInLeft 0.25s ease-out;
}

.swap-content {
  background: white;
  border-radius: 18px;
  width: 280px;
  box-sizing: border-box;
  border: 1px solid var(--divider);
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  overflow: hidden;
}

.swap-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 12px;
  border-bottom: 1px solid var(--divider);
}

h3 {
  font-size: 16px;
  font-weight: 500;
  color: var(--text);
  margin: 0;
}

.swap-text {
  font-size: 14px;
  padding: 14px 20px;
  line-height: 1.5;
  margin: 0;
}

.swap-actions {
  display: flex;
  gap: 8px;
  padding: 0 20px 16px;
}

.btn-decline {
  flex: 1;
  background: var(--panel);
  color: var(--text);
  border: none;
  border-radius: 10px;
  padding: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.20s;
}

.btn-decline:hover { background: #cfcfcf; }

.btn-accept {
  flex: 1;
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 10px;
  padding: 10px;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  font-family: 'Fira Sans', sans-serif;
  transition: background 0.20s;
}

.btn-accept:hover { background: var(--teal); }

@keyframes slideInLeft {
  from {
    transform: translateX(-110%);
    opacity: 0;
  } to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* Позиция того, кто предлагает обмен */
.swap-pos {
  font-size: 13px;
  color: var(--text-muted);
}
</style>
