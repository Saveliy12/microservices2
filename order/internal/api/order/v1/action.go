package v1

import (
	"net/http"
	"strings"
)

func (a *API) OrderAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	switch {
	case strings.HasSuffix(r.URL.Path, "/pay"):
		a.PayOrder(w, r)
	case strings.HasSuffix(r.URL.Path, "/cancel"):
		a.CancelOrder(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}
