<template>
  <div class="modal" @click.self="$emit('close')">
    <div class="modal-content">
      <button class="modal-close" @click="$emit('close')">X</button>
      <h3>Регистрация аккаунта</h3>

      <div class="input-group">
        <input v-model="name" type="text" placeholder="Имя" class="modal-input" />
        <input v-model="email" type="email" placeholder="Email" class="modal-input" />
        <input v-model="password" type="password" placeholder="Пароль" class="modal-input" />
      </div>

      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

      <button class="modal-btn" @click="register">Зарегистрироваться</button>

      <p class="switch-text">Уже есть аккаунт?
        <a href="#" @click.prevent="$emit('switch-to-login')">Войти</a>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { saveUser } from '../utils/auth.js'

const emit = defineEmits(['close', 'switch-to-login', 'registered'])

const name = ref('')
const email = ref('')
const password = ref('')
const agreed = ref(false)
const errorMsg = ref('')

async function register() {
  try {
    const res = await axios.post('http://localhost:8080/api/auth/register', {
      email: email.value,
      password: password.value,
      displayName: name.value,
    })
    saveUser(res.data.user, res.data.token)
    emit('registered', res.data.user)
  } catch (err) {
    if (err.response?.status === 409) {
      errorMsg.value = 'Email уже зарегистрирован'
    } else if (err.response?.status === 400) {
      errorMsg.value = 'Пароль минимум 6 символов'
    } else {
      errorMsg.value = 'Ошибка регистрации'
    }
  }
}
</script>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.modal-content {
  background: white;
  padding: 40px 30px 32px;
  border-radius: 18px;
  width: 100%;
  max-width: 380px;
  text-align: center;
  position: relative;
}

.modal-close {
  position: absolute;
  top: 14px;
  right: 18px;
  background: none;
  border: none;
  font-size: 16px;
  color: var(--text-muted);
  cursor: pointer;
  line-height: 1;
  padding: 6px 6px 4px 6px;
  border-radius: 6px;
  transition: background 0.20s;
}

.modal-close:hover {
  background: var(--divider);
}

h3 {
  font-size: 22px;
  font-weight: 400;
  color: var(--text);
  margin-bottom: 24px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.modal-input {
  width: 100%;
  padding: 14px 16px;
  border: 1px solid var(--divider);
  border-radius: 12px;
  font-family: 'Fira Sans', sans-serif;
  font-size: 15px;
  color: var(--text);
  background: #f5f5f5;
  box-sizing: border-box;
  transition: border-color 0.18s;
}

.modal-input:focus {
  outline: none;
  border-color: var(--teal);
  background: white;
}

.modal-input::placeholder {
  color: #aaa;
}

.modal-btn {
  background: var(--teal-dark);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 14px 0;
  font-size: 17px;
  font-weight: 700;
  width: 80%;
  font-family: 'Fira Sans', sans-serif;
  cursor: pointer;
  transition: background 0.20s;
}

.modal-btn:hover {
  background: var(--teal);
}

.switch-text {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 18px;
  margin-bottom: 0;
}

.switch-text a {
  color: var(--teal-dark);
  text-decoration: none;
  font-weight: 500;
}

.switch-text a:hover {
  text-decoration: underline;
}

.error-msg {
  font-size: 13px;
  color: #e85656;
  margin-bottom: 12px;
}
</style>
