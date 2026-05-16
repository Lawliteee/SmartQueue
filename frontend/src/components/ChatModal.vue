<template>
  <div class="chat-modal">
    <div class="chat-content">
      <div class="chat-header">
        <h3>Чат</h3>
        <button class="btn-close" @click="$emit('close')">Х</button>
      </div>

      <div class="chat-body" ref="messagesEl">
        <div v-if="messages.length === 0" class="chat-empty">
          Сообщений пока нет
        </div>
        <div v-for="(msg, idx) in messages" :key="idx" class="msg-wrapper"
          :class="msg.senderId === myId ? 'msg-wrapper--mine' : 'msg-wrapper--theirs'"
        >
          <div class="msg-bubble" :class="msg.senderId === myId ? 'msg-bubble--mine' : ''">
            <span class="msg-sender">{{ msg.senderName }}</span>
            <p class="msg-text">{{ msg.text }}</p>
          </div>
        </div>
      </div>

      <div class="chat-input-row">
        <input v-model="inputText" class="chat-input" placeholder="Сообщение..." maxlength="100" @keydown.enter="send"/>
        <button class="btn-send" @click="send" :disabled="!inputText.trim()">→</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'

const props = defineProps({
  messages: { type: Array, default: () => [] },
  myId: { type: String, default: '' },
  ws: { type: Object, default: null },
})
defineEmits(['close'])

const inputText = ref('')
const messagesEl = ref(null)

function send() {
  const text = inputText.value.trim()
  if (!text || !props.ws) return
  props.ws.send(JSON.stringify({ type: 'chat_message', text }))
  inputText.value = ''
}

watch(() => props.messages.length, async () => {
  await nextTick()
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight
  }
})
</script>

<style scoped>
.chat-modal {
  position: fixed;
  top: 65px;
  left: 20px;
  z-index: 10;
  pointer-events: none;
}

.chat-content {
  background: white;
  border-radius: 18px;
  width: 320px;
  height: min(600px, 70vh);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  margin-top: 20px;
  border: 1px solid var(--divider);
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  pointer-events: auto;
  animation: slideInLeft 0.3s ease-out;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--divider);
  flex-shrink: 0;
}

h3 { font-size: 17px; font-weight: 500; color: var(--text); margin: 0; }

.btn-close {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: var(--text-muted);
  padding: 0;
}

.btn-close:hover { color: var(--text); }

.chat-body {
  flex: 1;
  overflow-y: auto;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.chat-empty { margin: auto; font-size: 14px; color: var(--text-muted); }

.msg-wrapper { display: flex; }
.msg-wrapper--mine { justify-content: flex-end; }

.msg-bubble {
  max-width: 75%;
  background: var(--panel);
  border-radius: 12px;
  padding: 7px 11px;
}
.msg-bubble--mine { background: var(--teal); }

.msg-sender { display: block; font-size: 11px; font-weight: 600; color: var(--text-muted); margin-bottom: 2px; }
.msg-bubble--mine .msg-sender { color: rgba(255,255,255,0.7); }

.msg-text { font-size: 14px; color: var(--text); margin: 0; word-break: break-word; line-height: 1.4; }
.msg-bubble--mine .msg-text { color: white; }

.chat-input-row {
  display: flex;
  gap: 8px;
  padding: 10px 14px;
  border-top: 1px solid var(--divider);
  flex-shrink: 0;
}

.chat-input {
  flex: 1;
  padding: 9px 12px;
  border: 1px solid var(--divider);
  border-radius: 10px;
  font-family: 'Fira Sans', sans-serif;
  font-size: 14px;
  color: var(--text);
  outline: none;
}
.chat-input:focus { border-color: var(--teal); }

.btn-send {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 10px;
  width: 38px;
  font-size: 17px;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.2s;
}
.btn-send:hover { background: var(--teal); }



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
