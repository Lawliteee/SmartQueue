package main

// User — зарегистрированный пользователь.
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	PasswordHash string `json:"-"` // никогда не возвращаем клиенту
}

// RegisterRequest — тело запроса POST /api/auth/register
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

// LoginRequest — тело запроса POST /api/auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse — ответ при успешном входе/регистрации
type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
}
