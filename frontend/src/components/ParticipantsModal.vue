<template>
  <div class="modal" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h3>Участники очереди</h3>
        <button class="btn-close" @click="$emit('close')">Х</button>
      </div>
      <div class="participants-list">
        <div v-for="(p, idx) in displayedParticipants":key="p.id" class="participant-row"
        :class="{'participant-row--first': idx === 0,'participant-row--skipped': p.skipped}">
          <span class="p-number">{{ idx + 1 }}.</span>
          <span class="p-name">{{ p.name }}</span>
          <button v-if="swapEnabled && p.id !== myId && !(p.id === currentParticipant?.id) && myId !== currentParticipant?.id"
            class="btn-swap" @click="requestSwap(p)" title="Предложить обмен">
            <img src="/icons/swap.png" alt="обмен" width="16" height="16" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const emit = defineEmits(['close', 'swap-requested'])

const props = defineProps({
  participants: { type: Array, default: () => [] },
  currentParticipant: { type: Object, default: null },
  myId: { type: String, default: null },
  queueStarted: { type: Boolean, default: false },
  swapEnabled: { type: Boolean, default: false }
})

const displayedParticipants = computed(() => {
  if (!props.currentParticipant) return props.participants
  return [props.currentParticipant, ...props.participants]
})

// Предлагает другому пользователю обмен местами
function requestSwap(targetParticipant) {
  emit('swap-requested', targetParticipant)
}
</script>

<style scoped>
.modal {
  position: fixed;
  top: 70px;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  z-index: 10;
}

.modal-content {
  background: white;
  border-radius: 18px;
  width: 320px;
  height: 600px;
  max-width: 380px;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  overflow: hidden;
  margin: 20px;
  animation: slideInRight 0.3s ease-out;
  pointer-events: auto;
  border: 1px solid var(--divider);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.modal-header {
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

.btn-close:hover {
  color: var(--text);
}

.participants-list {
  overflow-y: auto;
  padding: 8px 0;
}

.participant-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 24px;
  border-top: 1px solid var(--divider);
}

.participant-row--first {
  border-top: none;
}

.participant-row--first .p-name {
  font-weight: 700;
}

.p-number {
  font-size: 14px;
  color: var(--text);
  min-width: 24px;
}

.p-name {
  font-size: 14px;
  color: var(--text);
}

.empty {
  padding: 40px 24px;
  text-align: center;
  color: var(--text-muted);
  font-size: 15px;
}

@keyframes slideInRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.participant-row--skipped .p-name, .participant-row--skipped .p-number { opacity: 0.5; }

/* Обмен местами */
.btn-swap {
  margin-left: auto;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  opacity: 0.6;
  transition: opacity 0.18s;
}

.btn-swap:hover {
  opacity: 1;
}
</style>
