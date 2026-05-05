export function setCookie(name, value, days = 7) {
    const expires = new Date(Date.now() + days * 864e5).toUTCString() // Дата истечения
    // Устанавливаем кику
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/; SameSite=Lax`
}

export function getCookie(name) {
    const match = document.cookie.split('; ').find(row => row.startsWith(name + '='))
    return match ? decodeURIComponent(match.split('=')[1]) : null
}

export function removeCookie(name) {
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/`
}
