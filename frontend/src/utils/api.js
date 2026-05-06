import axios from 'axios'
import { getToken } from './auth.js'

const api = axios.create({ // создаем экземпляр axios
  baseURL:'http://localhost:8080/api',
})

// Перехватывает запросы
api.interceptors.request.use(config => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}` // Добавляем в заголовки запроса токен
  }
  return config
})

export default api // Экспортируем настроенный экземпляр
