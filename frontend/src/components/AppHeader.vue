<template>
  <header>
    <router-link to="/" class="logo-link">
      <img class="logo" src="/logo.png" alt="SQ logo" />
    </router-link>
    <span class="header-title">Smart Queue</span>
    <nav class="header-nav">

      <!-- Незалогинен -->
      <template v-if="!currentUser">
        <a href="#" @click.prevent="openLogin">Войти</a>
        <a href="#" @click.prevent="openRegister">Зарегистрироваться</a>
      </template>

      <!-- Залогинен -->
      <div v-else class="avatar-wrap" ref="avatarRef">
        <button class="avatar-btn" @click="accountOpen = !accountOpen">
          <img src="/icons/avatar.png" alt="аккаунт" width="36" height="36" />
        </button>
        <AccountModal v-if="accountOpen":user="currentUser" @close="accountOpen = false" @logout="logout"/>
      </div>

    </nav>
  </header>

  <LoginModal
    v-if="loginOpen" @close="loginOpen = false"
    @switch-to-register="switchToRegister" @logged-in="onAuth"
  />
  <RegisterModal
    v-if="registerOpen" @close="registerOpen = false"
    @switch-to-login="switchToLogin" @registered="onAuth"
  />
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import LoginModal from './LoginModal.vue'
import RegisterModal from './RegisterModal.vue'
import AccountModal from './AccountModal.vue'
import { clearUser } from '../utils/auth.js'
import { useCurrentUser } from '../utils/useCurrentUser.js'
const currentUser = useCurrentUser()

const loginOpen = ref(false)
const registerOpen = ref(false)
const accountOpen = ref(false)
const avatarRef = ref(null)

onMounted(() => {
  document.addEventListener('click', onClickOutside) // Подписываемся на клики
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutside)  // Отписываемся от кликов
})

function onClickOutside(e) {
  if (avatarRef.value && !avatarRef.value.contains(e.target)) {
    accountOpen.value = false
  }
}

function onAuth(user) {
  currentUser.value = user
  loginOpen.value = false
  registerOpen.value = false
}

function logout() {
  clearUser()
  currentUser.value = null
  accountOpen.value = false
}

// Открываем модальное окно Логина
function openLogin() {
  loginOpen.value = true
  registerOpen.value = false
}

// Открываем модальное окно Регистрации
function openRegister() {
  registerOpen.value = true
  loginOpen.value = false
}

// Переключаемся с Логина на Регистрацию
function switchToRegister() {
  loginOpen.value = false
  registerOpen.value = true
}

// Переключаемся с Регистрации на Логин
function switchToLogin() {
  registerOpen.value = false
  loginOpen.value = true
}
</script>

<style scoped>
header {
  background: var(--teal);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  height: 70px;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.logo-link {
  display: block;
  line-height: 0;
  cursor: pointer;
  z-index: 1;
}

.logo {
  height: 50px;
  width: auto;
}

.header-title {
  font-size: 25px;
  font-weight: 800;
  color: var(--text);
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.header-nav {
  display: flex;
  gap: 20px;
  z-index: 1;
}

.header-nav a {
  color: white;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  opacity: 0.95;
}

.header-nav a:hover { opacity: 0.7; }

/* Аватарка */
.avatar-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  border-radius: 50%;
  overflow: hidden;
  transition: opacity 0.20s;
}

.avatar-btn:hover { opacity: 0.9; }
</style>
