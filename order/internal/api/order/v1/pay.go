package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com.Saveliy12/microservices2.git/order/internal/converter"
	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

func (a *API) PayOrder(w http.ResponseWriter, r *http.Request) {
	type request struct {
		PaymentMethod string `json:"payment_method"`
	}

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	orderUUID, err := extractActionOrderUUID(r.URL.Path, "/pay")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order uuid")
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.PaymentMethod == "" {
		writeError(w, http.StatusBadRequest, "payment_method is required")
		return
	}

	paymentMethod, err := converter.PaymentMethodFromString(req.PaymentMethod)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := a.orderService.PayOrder(r.Context(), orderUUID, paymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			writeError(w, http.StatusNotFound, "order not found")
		case errors.Is(err, model.ErrOrderAlreadyPaid), errors.Is(err, model.ErrOrderCancelled):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusBadGateway, err.Error())
		}
		return
	}

	resp, err := converter.ModelToPayResponse(order)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func extractActionOrderUUID(path, suffix string) (string, error) {
	if !strings.HasSuffix(path, suffix) {
		return "", errors.New("invalid path")
	}

	value := strings.TrimSuffix(path, suffix)
	value = strings.TrimSuffix(value, "/")
	return extractOrderUUID(value, "/api/v1/orders/")
}
