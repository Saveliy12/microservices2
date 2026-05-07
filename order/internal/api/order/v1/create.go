package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

func (a *API) CreateOrder(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserUUID  string   `json:"user_uuid"`
		PartUUIDs []string `json:"part_uuids"`
	}

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.UserUUID == "" || len(req.PartUUIDs) == 0 {
		writeError(w, http.StatusBadRequest, "user_uuid and part_uuids are required")
		return
	}

	order, err := a.orderService.CreateOrder(r.Context(), req.UserUUID, req.PartUUIDs)
	if err != nil {
		if errors.Is(err, model.ErrInvalidOrderInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"uuid":        order.OrderUUID,
		"total_price": order.TotalPrice,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
