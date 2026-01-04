package handlers

import "net/http"

func InitHandlerHealth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Если запрос дошёл, то сервер работает - отсылаем OK
		w.WriteHeader(http.StatusOK)
	})
}
