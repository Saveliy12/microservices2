package v1

import (
	"errors"
	"net/http"

	"github.com.Saveliy12/microservices2.git/order/internal/model"
)

func (a *API) CancelOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	orderUUID, err := extractActionOrderUUID(r.URL.Path, "/cancel")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order uuid")
		return
	}

	_, err = a.orderService.CancelOrder(r.Context(), orderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			writeError(w, http.StatusNotFound, "order not found")
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			writeError(w, http.StatusConflict, "paid order cannot be cancelled")
		case errors.Is(err, model.ErrOrderCancelled):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusBadGateway, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
