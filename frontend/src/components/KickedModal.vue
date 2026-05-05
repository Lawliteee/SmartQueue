<template>
  <div class="modal">
    <div class="modal-content">
      <h3>Вы покидаете очередь</h3>
      <p class="subtitle">Вы больше не числитесь в очереди</p>
      <button class="modal-btn" @click="$emit('confirm')">На главную</button>
    </div>
  </div>
</template>

<script setup>
    defineEmits(['confirm'])
</script>

<style scoped>

.modal {
  position: fixed;
  top: 0; left: 0;
  width: 100%; height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.modal-content {
  background: white;
  padding: 25px;
  border-radius: 18px;
  width: 100%;
  max-width: 380px;
  text-align: center;
  box-sizing: border-box;
}

h3 {
  font-size: 25px;
  font-weight: 400;
  margin-bottom: 12px;
  color: var(--text);
}

.subtitle {
  font-size: 15px;
  color: var(--text-muted);
  line-height: 1.6;
  margin-bottom: 28px;
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
  transition: background 0.20s;
}

.modal-btn:hover { background: var(--teal); }

</style>
