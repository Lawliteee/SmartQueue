import { ref } from 'vue'
import { getUser } from './auth.js'

const currentUser = ref(getUser())

export function useCurrentUser() {
  return currentUser
}