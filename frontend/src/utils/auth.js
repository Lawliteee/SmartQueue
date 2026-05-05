import { getCookie, setCookie, removeCookie } from './cookies.js'

export function getUser() {
  const raw = getCookie('user') // Получаем куки с пользователем
  if (!raw) {
    return null
  }

  try {
    return JSON.parse(raw)
  } catch { return null }
}

export function saveUser(user) {
  setCookie('user', JSON.stringify(user), 30) // Сохраняем в куки на 30 дней
}

export function clearUser() {
  removeCookie('user')
}
