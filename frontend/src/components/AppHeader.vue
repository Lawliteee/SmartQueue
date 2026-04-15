<template>
  <header>
    <router-link to="/" class="logo-link">
      <img class="logo" src="/logo.svg" alt="SQ logo" />
    </router-link>
    <span class="header-title">Smart Queue</span>
    <nav class="header-nav">
      <a href="#" @click.prevent="openLogin">Войти</a>
      <a href="#" @click.prevent="openRegister">Зарегистрироваться</a>
    </nav>
  </header>

  <LoginModal
    v-if="loginOpen"
    @close="loginOpen = false"
    @switch-to-register="switchToRegister"
  />
 
  <RegisterModal
    v-if="registerOpen"
    @close="registerOpen = false"
    @switch-to-login="switchToLogin"
  />
</template>

<script setup>
import { ref } from 'vue'
import LoginModal from './LoginModal.vue'
import RegisterModal from './RegisterModal.vue'
 
const loginOpen = ref(false)
const registerOpen = ref(false)
 
function openLogin() {
  loginOpen.value = true
  registerOpen.value = false
}
 
function openRegister() {
  registerOpen.value = true
  loginOpen.value = false
}
 
function switchToRegister() {
  loginOpen.value = false
  registerOpen.value = true
}
 
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
.header-nav a:hover {
  opacity: 0.7;
}
</style>