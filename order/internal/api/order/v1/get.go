package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com.Saveliy12/microservices2.git/order/internal/converter"
	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

func (a *API) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	orderUUID, err := extractOrderUUID(r.URL.Path, "/api/v1/orders/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order uuid")
		return
	}

	order, err := a.orderService.GetOrder(r.Context(), orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			writeError(w, http.StatusNotFound, "order not found")
			return
		}
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	resp := converter.ModelToGetResponse(order)
	writeJSON(w, http.StatusOK, resp)
}

func extractOrderUUID(path, prefix string) (string, error) {
	if !strings.HasPrefix(path, prefix) {
		return "", errors.New("invalid path")
	}

	value := strings.TrimPrefix(path, prefix)
	value = strings.Trim(value, "/")
	if value == "" || strings.Contains(value, "/") {
		return "", errors.New("invalid order uuid")
	}

	return value, nil
}
