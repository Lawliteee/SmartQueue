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

export function saveUser(user, token) {
  setCookie('user', JSON.stringify(user), 30) // Сохраняем в куки на 30 дней
  setCookie('token', token, 3) // токен на 3 дня
}

export function getToken() {
  return getCookie('token')
}

export function clearUser() {
  removeCookie('user')
  removeCookie('token')
}
