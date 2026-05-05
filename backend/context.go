package main

import "context"

type contextKey string

const ctxUserID contextKey = "userID"

func withUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxUserID, id)
}

// UserIDFromContext возвращает userID из контекста запроса.
// Возвращает пустую строку, если пользователь не аутентифицирован.
func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserID).(string)
	return v
}
