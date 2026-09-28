package health

import (
	"encoding/json"
	"net/http"
)

// Handler HTTP health check эндпоинт сервиса
// Не зависит от внешних ресурсов: отвечает 200, пока жив процесс
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type response struct {
	Status string `json:"status"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response{Status: "SERVING"})
}
